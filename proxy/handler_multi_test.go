// Copyright 2017 Michal Witkowski. All Rights Reserved.
// Copyright 2019 Andrey Smirnov. All Rights Reserved.
// See LICENSE for licensing terms.

package proxy_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/smira/grpc-proxy/proxy"
	pb "github.com/smira/grpc-proxy/testservice"
)

const (
	numUpstreams = 3
)

// asserting service is implemented on the server side and serves as a handler for stuff
type assertingMultiService struct {
	t      *testing.T
	server string
}

func (s *assertingMultiService) PingEmpty(ctx context.Context, _ *pb.Empty) (*pb.MultiPingReply, error) {
	// Check that this call has client's metadata.
	md, ok := metadata.FromIncomingContext(ctx)
	assert.True(s.t, ok, "PingEmpty call must have metadata in context")
	_, ok = md[clientMdKey]
	assert.True(s.t, ok, "PingEmpty call must have clients's custom headers in metadata")
	return &pb.MultiPingReply{
		Response: []*pb.MultiPingResponse{
			{
				Value:   pingDefaultValue,
				Counter: 42,
				Server:  s.server,
			},
		},
	}, nil
}

func (s *assertingMultiService) Ping(ctx context.Context, ping *pb.PingRequest) (*pb.MultiPingReply, error) {
	// Send user trailers and headers.
	grpc.SendHeader(ctx, metadata.Pairs(serverHeaderMdKey, "I like turtles."))         //nolint: errcheck
	grpc.SetTrailer(ctx, metadata.Pairs(serverTrailerMdKey, "I like ending turtles.")) //nolint: errcheck
	return &pb.MultiPingReply{
		Response: []*pb.MultiPingResponse{
			{
				Value:   ping.Value,
				Counter: 42,
				Server:  s.server,
			},
		},
	}, nil
}

func (s *assertingMultiService) PingError(ctx context.Context, ping *pb.PingRequest) (*pb.EmptyReply, error) {
	return nil, status.Errorf(codes.FailedPrecondition, "Userspace error.")
}

func (s *assertingMultiService) PingList(ping *pb.PingRequest, stream pb.MultiService_PingListServer) error {
	// Send user trailers and headers.
	stream.SendHeader(metadata.Pairs(serverHeaderMdKey, "I like turtles.")) //nolint: errcheck
	for i := 0; i < countListResponses; i++ {
		stream.Send(&pb.MultiPingReply{ //nolint: errcheck
			Response: []*pb.MultiPingResponse{
				{
					Value:   ping.Value,
					Counter: int32(i),
					Server:  s.server,
				},
			},
		})
	}
	stream.SetTrailer(metadata.Pairs(serverTrailerMdKey, "I like ending turtles.")) //nolint: errcheck
	return nil
}

func (s *assertingMultiService) PingStream(stream pb.MultiService_PingStreamServer) error {
	stream.SendHeader(metadata.Pairs(serverHeaderMdKey, "I like turtles.")) // nolint: errcheck
	counter := int32(0)
	for {
		ping, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			require.NoError(s.t, err, "can't fail reading stream")
			return err
		}
		pong := &pb.MultiPingReply{
			Response: []*pb.MultiPingResponse{
				{
					Value:   ping.Value,
					Counter: counter,
					Server:  s.server,
				},
			},
		}
		if err := stream.Send(pong); err != nil {
			require.NoError(s.t, err, "can't fail sending back a pong")
		}
		counter += 1
	}
	stream.SetTrailer(metadata.Pairs(serverTrailerMdKey, "I like ending turtles."))
	return nil
}

type assertingBackend struct {
	addr string
	i    int
}

func (b *assertingBackend) String() string {
	return fmt.Sprintf("backend%d", b.i)
}

func (b *assertingBackend) GetConnection(ctx context.Context) (context.Context, *grpc.ClientConn, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	// Explicitly copy the metadata, otherwise the tests will fail.
	outCtx := metadata.NewOutgoingContext(ctx, md.Copy())

	conn, err := grpc.DialContext(ctx, b.addr, grpc.WithInsecure(), grpc.WithCodec(proxy.Codec())) // nolint: staticcheck

	return outCtx, conn, err
}

func (b *assertingBackend) AppendInfo(resp []byte) ([]byte, error) {
	// decode protobuf embedded header
	typ, n1 := proto.DecodeVarint(resp)
	_, n2 := proto.DecodeVarint(resp[n1:]) // length

	if typ != (1<<3)|2 { // type: 2, field_number: 1
		return nil, fmt.Errorf("unexpected message format: %d", typ)
	}

	payload, err := proto.Marshal(&pb.ResponseMetadataPrepender{
		Metadata: &pb.ResponseMetadata{
			Hostname: fmt.Sprintf("server%d", b.i),
		},
	})

	// cut off embedded message header
	resp = resp[n1+n2:]
	// build new embedded message header
	prefix := append(proto.EncodeVarint((1<<3)|2), proto.EncodeVarint(uint64(len(resp)+len(payload)))...)
	resp = append(prefix, resp...)

	return append(resp, payload...), err
}

func (b *assertingBackend) BuildError(err error) ([]byte, error) {
	return proto.Marshal(&pb.EmptyReply{
		Response: []*pb.EmptyResponse{
			{
				Metadata: &pb.ResponseMetadata{
					Hostname:      fmt.Sprintf("server%d", b.i),
					UpstreamError: err.Error(),
				},
			},
		},
	})
}

type MultiServiceSuite struct {
	suite.Suite

	serverListeners  []net.Listener
	servers          []*grpc.Server
	proxyListener    net.Listener
	proxy            *grpc.Server
	serverClientConn *grpc.ClientConn

	client     *grpc.ClientConn
	testClient pb.MultiServiceClient

	ctx       context.Context
	ctxCancel context.CancelFunc
}

func (s *MultiServiceSuite) TestPingEmptyCarriesClientMetadata() {
	ctx := metadata.NewOutgoingContext(s.ctx, metadata.Pairs(clientMdKey, "true"))
	out, err := s.testClient.PingEmpty(ctx, &pb.Empty{})
	require.NoError(s.T(), err, "PingEmpty should succeed without errors")

	expectedUpstreams := map[string]struct{}{}
	for i := 0; i < numUpstreams; i++ {
		expectedUpstreams[fmt.Sprintf("server%d", i)] = struct{}{}
	}

	s.Require().Len(out.Response, numUpstreams)
	for _, resp := range out.Response {
		s.Require().Equal(pingDefaultValue, resp.Value)
		s.Require().EqualValues(42, resp.Counter)

		// equal metadata set by proxy and server
		s.Require().Equal(resp.Metadata.Hostname, resp.Server)

		delete(expectedUpstreams, resp.Metadata.Hostname)
	}

	s.Require().Empty(expectedUpstreams)
}

func (s *MultiServiceSuite) TestPingEmpty_StressTest() {
	for i := 0; i < 50; i++ {
		s.TestPingEmptyCarriesClientMetadata()
	}
}

func (s *MultiServiceSuite) TestPingCarriesServerHeadersAndTrailers() {
	headerMd := make(metadata.MD)
	trailerMd := make(metadata.MD)
	// This is an awkward calling convention... but meh.
	out, err := s.testClient.Ping(s.ctx, &pb.PingRequest{Value: "foo"}, grpc.Header(&headerMd), grpc.Trailer(&trailerMd))
	require.NoError(s.T(), err, "Ping should succeed without errors")

	s.Require().Len(out.Response, numUpstreams)
	for _, resp := range out.Response {
		s.Require().Equal("foo", resp.Value)
		s.Require().EqualValues(42, resp.Counter)

		// equal metadata set by proxy and server
		s.Require().Equal(resp.Metadata.Hostname, resp.Server)
	}

	assert.Contains(s.T(), headerMd, serverHeaderMdKey, "server response headers must contain server data")
	assert.Len(s.T(), trailerMd, 1, "server response trailers must contain server data")
}

func (s *MultiServiceSuite) TestPingErrorPropagatesAppError() {
	out, err := s.testClient.PingError(s.ctx, &pb.PingRequest{Value: "foo"})
	s.Require().NoError(err, "error should be encapsulated in the response")

	s.Require().Len(out.Response, numUpstreams)
	for _, resp := range out.Response {
		s.Require().NotEmpty(resp.Metadata.UpstreamError)
		s.Require().NotEmpty(resp.Metadata.Hostname)
		s.Assert().Equal("rpc error: code = FailedPrecondition desc = Userspace error.", resp.Metadata.UpstreamError)
	}
}

func (s *MultiServiceSuite) TestDirectorErrorIsPropagated() {
	// See SetupSuite where the StreamDirector has a special case.
	ctx := metadata.NewOutgoingContext(s.ctx, metadata.Pairs(rejectingMdKey, "true"))
	_, err := s.testClient.Ping(ctx, &pb.PingRequest{Value: "foo"})
	require.Error(s.T(), err, "Director should reject this RPC")
	assert.Equal(s.T(), codes.PermissionDenied, status.Code(err))
	assert.Equal(s.T(), "testing rejection", status.Convert(err).Message())
}

func (s *MultiServiceSuite) TestPingStream_FullDuplexWorks() {
	stream, err := s.testClient.PingStream(s.ctx)
	require.NoError(s.T(), err, "PingStream request should be successful.")

	for i := 0; i < countListResponses; i++ {
		ping := &pb.PingRequest{Value: fmt.Sprintf("foo:%d", i)}
		require.NoError(s.T(), stream.Send(ping), "sending to PingStream must not fail")

		expectedUpstreams := map[string]struct{}{}
		for j := 0; j < numUpstreams; j++ {
			expectedUpstreams[fmt.Sprintf("server%d", j)] = struct{}{}
		}

		// each upstream should send back response
		for j := 0; j < numUpstreams; j++ {
			resp, err := stream.Recv()
			s.Require().NoError(err)

			s.Assert().Len(resp.Response, 1)
			s.Assert().EqualValues(i, resp.Response[0].Counter, "ping roundtrip must succeed with the correct id")
			s.Assert().EqualValues(resp.Response[0].Metadata.Hostname, resp.Response[0].Server)

			delete(expectedUpstreams, resp.Response[0].Metadata.Hostname)
		}

		s.Require().Empty(expectedUpstreams)

		if i == 0 {
			// Check that the header arrives before all entries.
			headerMd, err := stream.Header()
			require.NoError(s.T(), err, "PingStream headers should not error.")
			assert.Contains(s.T(), headerMd, serverHeaderMdKey, "PingStream response headers user contain metadata")
		}
	}
	require.NoError(s.T(), stream.CloseSend(), "no error on close send")
	_, err = stream.Recv()
	require.Equal(s.T(), io.EOF, err, "stream should close with io.EOF, meaning OK")
	// Check that the trailer headers are here.
	trailerMd := stream.Trailer()
	assert.Len(s.T(), trailerMd, 1, "PingList trailer headers user contain metadata")
}

func (s *MultiServiceSuite) SetupTest() {
	s.ctx, s.ctxCancel = context.WithTimeout(context.TODO(), 120*time.Second)
}

func (s *MultiServiceSuite) TearDownTest() {
	s.ctxCancel()
}

func (s *MultiServiceSuite) SetupSuite() {
	var err error

	s.proxyListener, err = net.Listen("tcp", "127.0.0.1:0")
	require.NoError(s.T(), err, "must be able to allocate a port for proxyListener")

	s.serverListeners = make([]net.Listener, numUpstreams)

	for i := range s.serverListeners {
		s.serverListeners[i], err = net.Listen("tcp", "127.0.0.1:0")
		require.NoError(s.T(), err, "must be able to allocate a port for serverListener")
	}

	s.servers = make([]*grpc.Server, numUpstreams)

	for i := range s.servers {
		s.servers[i] = grpc.NewServer()
		pb.RegisterMultiServiceServer(s.servers[i],
			&assertingMultiService{
				t:      s.T(),
				server: fmt.Sprintf("server%d", i),
			})
	}

	backends := make([]*assertingBackend, numUpstreams)

	for i := range backends {
		backends[i] = &assertingBackend{
			i:    i,
			addr: s.serverListeners[i].Addr().String(),
		}
	}

	// Setup of the proxy's Director.
	director := func(ctx context.Context, fullName string) ([]proxy.Backend, error) {
		var targets []int

		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			if _, exists := md[rejectingMdKey]; exists {
				return nil, status.Errorf(codes.PermissionDenied, "testing rejection")
			}

			if mdTargets, exists := md["targets"]; exists {
				for _, strTarget := range mdTargets {
					t, err := strconv.Atoi(strTarget)
					if err != nil {
						return nil, err
					}

					targets = append(targets, t)
				}
			}
		}

		var result []proxy.Backend

		if targets == nil {
			for i := range backends {
				targets = append(targets, i)
			}
		}

		for _, t := range targets {
			result = append(result, backends[t])
		}

		return result, nil
	}

	s.proxy = grpc.NewServer(
		grpc.CustomCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(proxy.TransparentHandler(director)),
	)
	// Ping handler is handled as an explicit registration and not as a TransparentHandler.
	proxy.RegisterService(s.proxy, director,
		"smira.testproto.MultiService",
		[]string{"Ping", "PingStream"},
		[]string{"PingStream"})

	// Start the serving loops.
	for i := range s.servers {
		s.T().Logf("starting grpc.Server at: %v", s.serverListeners[i].Addr().String())
		go func(i int) {
			s.servers[i].Serve(s.serverListeners[i]) // nolint: errcheck
		}(i)
	}
	s.T().Logf("starting grpc.Proxy at: %v", s.proxyListener.Addr().String())
	go func() {
		s.proxy.Serve(s.proxyListener) // nolint: errcheck
	}()

	ctx, ctxCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer ctxCancel()
	clientConn, err := grpc.DialContext(ctx, strings.Replace(s.proxyListener.Addr().String(), "127.0.0.1", "localhost", 1), grpc.WithInsecure())
	require.NoError(s.T(), err, "must not error on deferred client Dial")
	s.testClient = pb.NewMultiServiceClient(clientConn)
}

func (s *MultiServiceSuite) TearDownSuite() {
	if s.client != nil {
		s.client.Close()
	}
	if s.serverClientConn != nil {
		s.serverClientConn.Close()
	}
	// Close all transports so the logs don't get spammy.
	time.Sleep(10 * time.Millisecond)

	if s.proxy != nil {
		s.proxy.Stop()
		s.proxyListener.Close()
	}

	for _, server := range s.servers {
		if server != nil {
			server.Stop()
		}
	}

	for _, serverListener := range s.serverListeners {
		if serverListener != nil {
			serverListener.Close()
		}
	}
}
func TestMultiServiceSuite(t *testing.T) {
	suite.Run(t, &MultiServiceSuite{})
}

func init() {
	grpclog.SetLogger(log.New(os.Stderr, "grpc: ", log.LstdFlags)) // nolint: staticcheck
}
