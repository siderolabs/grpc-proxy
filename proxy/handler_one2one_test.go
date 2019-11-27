// Copyright 2017 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package proxy_test

import (
	"context"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"fmt"

	"github.com/talos-systems/grpc-proxy/proxy"
	pb "github.com/talos-systems/grpc-proxy/testservice"
)

const (
	pingDefaultValue   = "I like kittens."
	clientMdKey        = "test-client-header"
	serverHeaderMdKey  = "test-client-header"
	serverTrailerMdKey = "test-client-trailer"

	rejectingMdKey = "test-reject-rpc-if-in-context"

	countListResponses = 20
)

// asserting service is implemented on the server side and serves as a handler for stuff
type assertingService struct {
	t *testing.T
}

func (s *assertingService) PingEmpty(ctx context.Context, _ *pb.Empty) (*pb.PingResponse, error) {
	// Check that this call has client's metadata.
	md, ok := metadata.FromIncomingContext(ctx)
	assert.True(s.t, ok, "PingEmpty call must have metadata in context")
	_, ok = md[clientMdKey]
	assert.True(s.t, ok, "PingEmpty call must have clients's custom headers in metadata")
	return &pb.PingResponse{Value: pingDefaultValue, Counter: 42}, nil
}

func (s *assertingService) Ping(ctx context.Context, ping *pb.PingRequest) (*pb.PingResponse, error) {
	// Send user trailers and headers.
	grpc.SendHeader(ctx, metadata.Pairs(serverHeaderMdKey, "I like turtles."))         //nolint: errcheck
	grpc.SetTrailer(ctx, metadata.Pairs(serverTrailerMdKey, "I like ending turtles.")) //nolint: errcheck
	return &pb.PingResponse{Value: ping.Value, Counter: 42}, nil
}

func (s *assertingService) PingError(ctx context.Context, ping *pb.PingRequest) (*pb.Empty, error) {
	return nil, status.Errorf(codes.FailedPrecondition, "Userspace error.")
}

func (s *assertingService) PingList(ping *pb.PingRequest, stream pb.TestService_PingListServer) error {
	// Send user trailers and headers.
	stream.SendHeader(metadata.Pairs(serverHeaderMdKey, "I like turtles.")) //nolint: errcheck
	for i := 0; i < countListResponses; i++ {
		stream.Send(&pb.PingResponse{Value: ping.Value, Counter: int32(i)}) //nolint: errcheck
	}
	stream.SetTrailer(metadata.Pairs(serverTrailerMdKey, "I like ending turtles.")) //nolint: errcheck
	return nil
}

func (s *assertingService) PingStream(stream pb.TestService_PingStreamServer) error {
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
		pong := &pb.PingResponse{Value: ping.Value, Counter: counter}
		if err := stream.Send(pong); err != nil {
			require.NoError(s.t, err, "can't fail sending back a pong")
		}
		counter += 1
	}
	stream.SetTrailer(metadata.Pairs(serverTrailerMdKey, "I like ending turtles."))
	return nil
}

// ProxyOne2OneSuite tests the "happy" path of handling: that everything works in absence of connection issues.
type ProxyOne2OneSuite struct {
	suite.Suite

	serverListener   net.Listener
	server           *grpc.Server
	proxyListener    net.Listener
	proxy            *grpc.Server
	serverClientConn *grpc.ClientConn

	client     *grpc.ClientConn
	testClient pb.TestServiceClient

	ctx       context.Context
	ctxCancel context.CancelFunc
}

func (s *ProxyOne2OneSuite) SetupTest() {
	s.ctx, s.ctxCancel = context.WithTimeout(context.TODO(), 120*time.Second)
}

func (s *ProxyOne2OneSuite) TearDownTest() {
	s.ctxCancel()
}

func (s *ProxyOne2OneSuite) TestPingEmptyCarriesClientMetadata() {
	ctx := metadata.NewOutgoingContext(s.ctx, metadata.Pairs(clientMdKey, "true"))
	out, err := s.testClient.PingEmpty(ctx, &pb.Empty{})
	require.NoError(s.T(), err, "PingEmpty should succeed without errors")
	require.Equal(s.T(), &pb.PingResponse{Value: pingDefaultValue, Counter: 42}, out)
}

func (s *ProxyOne2OneSuite) TestPingEmpty_StressTest() {
	for i := 0; i < 50; i++ {
		s.TestPingEmptyCarriesClientMetadata()
	}
}

func (s *ProxyOne2OneSuite) TestPingCarriesServerHeadersAndTrailers() {
	headerMd := make(metadata.MD)
	trailerMd := make(metadata.MD)
	// This is an awkward calling convention... but meh.
	out, err := s.testClient.Ping(s.ctx, &pb.PingRequest{Value: "foo"}, grpc.Header(&headerMd), grpc.Trailer(&trailerMd))
	require.NoError(s.T(), err, "Ping should succeed without errors")
	require.Equal(s.T(), &pb.PingResponse{Value: "foo", Counter: 42}, out)
	assert.Contains(s.T(), headerMd, serverHeaderMdKey, "server response headers must contain server data")
	assert.Len(s.T(), trailerMd, 1, "server response trailers must contain server data")
}

func (s *ProxyOne2OneSuite) TestPingErrorPropagatesAppError() {
	_, err := s.testClient.PingError(s.ctx, &pb.PingRequest{Value: "foo"})
	require.Error(s.T(), err, "PingError should never succeed")
	assert.Equal(s.T(), codes.FailedPrecondition, status.Code(err))
	assert.Equal(s.T(), "Userspace error.", status.Convert(err).Message())
}

func (s *ProxyOne2OneSuite) TestDirectorErrorIsPropagated() {
	// See SetupSuite where the StreamDirector has a special case.
	ctx := metadata.NewOutgoingContext(s.ctx, metadata.Pairs(rejectingMdKey, "true"))
	_, err := s.testClient.Ping(ctx, &pb.PingRequest{Value: "foo"})
	require.Error(s.T(), err, "Director should reject this RPC")
	assert.Equal(s.T(), codes.PermissionDenied, status.Code(err))
	assert.Equal(s.T(), "testing rejection", status.Convert(err).Message())
}

func (s *ProxyOne2OneSuite) TestPingStream_FullDuplexWorks() {
	stream, err := s.testClient.PingStream(s.ctx)
	require.NoError(s.T(), err, "PingStream request should be successful.")

	for i := 0; i < countListResponses; i++ {
		ping := &pb.PingRequest{Value: fmt.Sprintf("foo:%d", i)}
		require.NoError(s.T(), stream.Send(ping), "sending to PingStream must not fail")
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if i == 0 {
			// Check that the header arrives before all entries.
			headerMd, err := stream.Header()
			require.NoError(s.T(), err, "PingStream headers should not error.")
			assert.Contains(s.T(), headerMd, serverHeaderMdKey, "PingStream response headers user contain metadata")
		}
		assert.EqualValues(s.T(), i, resp.Counter, "ping roundtrip must succeed with the correct id")
	}
	require.NoError(s.T(), stream.CloseSend(), "no error on close send")
	_, err = stream.Recv()
	require.Equal(s.T(), io.EOF, err, "stream should close with io.EOF, meaining OK")
	// Check that the trailer headers are here.
	trailerMd := stream.Trailer()
	assert.Len(s.T(), trailerMd, 1, "PingList trailer headers user contain metadata")
}

func (s *ProxyOne2OneSuite) TestPingStream_StressTest() {
	for i := 0; i < 50; i++ {
		s.TestPingStream_FullDuplexWorks()
	}
}

func (s *ProxyOne2OneSuite) SetupSuite() {
	var err error

	s.proxyListener, err = net.Listen("tcp", "127.0.0.1:0")
	require.NoError(s.T(), err, "must be able to allocate a port for proxyListener")
	s.serverListener, err = net.Listen("tcp", "127.0.0.1:0")
	require.NoError(s.T(), err, "must be able to allocate a port for serverListener")

	s.server = grpc.NewServer()
	pb.RegisterTestServiceServer(s.server, &assertingService{t: s.T()})

	// Setup of the proxy's Director.
	s.serverClientConn, err = grpc.Dial(s.serverListener.Addr().String(), grpc.WithInsecure(), grpc.WithCodec(proxy.Codec())) // nolint: staticcheck
	require.NoError(s.T(), err, "must not error on deferred client Dial")
	director := func(ctx context.Context, fullName string) (proxy.Mode, []proxy.Backend, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			if _, exists := md[rejectingMdKey]; exists {
				return proxy.One2One, nil, status.Errorf(codes.PermissionDenied, "testing rejection")
			}
		}

		return proxy.One2One, []proxy.Backend{
			&proxy.SingleBackend{
				GetConn: func(ctx context.Context) (context.Context, *grpc.ClientConn, error) {
					md, _ := metadata.FromIncomingContext(ctx)
					// Explicitly copy the metadata, otherwise the tests will fail.
					outCtx := metadata.NewOutgoingContext(ctx, md.Copy())
					return outCtx, s.serverClientConn, nil
				},
			},
		}, nil
	}
	s.proxy = grpc.NewServer(
		grpc.CustomCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(proxy.TransparentHandler(director)),
	)
	// Ping handler is handled as an explicit registration and not as a TransparentHandler.
	proxy.RegisterService(s.proxy, director,
		"talos.testproto.TestService",
		proxy.WithMethodNames("Ping"),
	)

	// Start the serving loops.
	s.T().Logf("starting grpc.Server at: %v", s.serverListener.Addr().String())
	go func() {
		s.server.Serve(s.serverListener) // nolint: errcheck
	}()
	s.T().Logf("starting grpc.Proxy at: %v", s.proxyListener.Addr().String())
	go func() {
		s.proxy.Serve(s.proxyListener) // nolint: errcheck
	}()

	ctx, ctxCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer ctxCancel()
	clientConn, err := grpc.DialContext(ctx, strings.Replace(s.proxyListener.Addr().String(), "127.0.0.1", "localhost", 1), grpc.WithInsecure())
	require.NoError(s.T(), err, "must not error on deferred client Dial")
	s.testClient = pb.NewTestServiceClient(clientConn)
}

func (s *ProxyOne2OneSuite) TearDownSuite() {
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
	if s.serverListener != nil {
		s.server.Stop()
		s.serverListener.Close()
	}
}

func TestProxyOne2OneSuite(t *testing.T) {
	suite.Run(t, &ProxyOne2OneSuite{})
}
