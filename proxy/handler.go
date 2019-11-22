// Copyright 2017 Michal Witkowski. All Rights Reserved.
// Copyright 2019 Andrey Smirnov. All Rights Reserved.
// See LICENSE for licensing terms.

package proxy

import (
	"context"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	clientStreamDescForProxying = &grpc.StreamDesc{
		ServerStreams: true,
		ClientStreams: true,
	}
)

// RegisterService sets up a proxy handler for a particular gRPC service and method.
// The behaviour is the same as if you were registering a handler method, e.g. from a codegenerated pb.go file.
//
// This can *only* be used if the `server` also uses grpcproxy.CodecForServer() ServerOption.
//
// streamedMethodNames is only important for one-2-many proxying.
func RegisterService(server *grpc.Server, director StreamDirector, serviceName string, methodNames []string, streamedMethodNames []string) {
	streamer := &handler{
		director:        director,
		streamedMethods: map[string]struct{}{},
	}

	for _, methodName := range streamedMethodNames {
		streamer.streamedMethods["/"+serviceName+"/"+methodName] = struct{}{}
	}

	fakeDesc := &grpc.ServiceDesc{
		ServiceName: serviceName,
		HandlerType: (*interface{})(nil),
	}
	for _, m := range methodNames {
		streamDesc := grpc.StreamDesc{
			StreamName:    m,
			Handler:       streamer.handler,
			ServerStreams: true,
			ClientStreams: true,
		}
		fakeDesc.Streams = append(fakeDesc.Streams, streamDesc)
	}
	server.RegisterService(fakeDesc, streamer)
}

// TransparentHandler returns a handler that attempts to proxy all requests that are not registered in the server.
// The indented use here is as a transparent proxy, where the server doesn't know about the services implemented by the
// backends. It should be used as a `grpc.UnknownServiceHandler`.
//
// This can *only* be used if the `server` also uses grpcproxy.CodecForServer() ServerOption.
func TransparentHandler(director StreamDirector) grpc.StreamHandler {
	streamer := &handler{
		director:        director,
		streamedMethods: map[string]struct{}{},
	}
	return streamer.handler
}

type handler struct {
	director        StreamDirector
	streamedMethods map[string]struct{}
}

type backendConnection struct {
	backend Backend

	backendConn *grpc.ClientConn
	connError   error

	clientStream grpc.ClientStream
}

// handler is where the real magic of proxying happens.
// It is invoked like any gRPC server stream and uses the gRPC server framing to get and receive bytes from the wire,
// forwarding it to a ClientStream established against the relevant ClientConn.
func (s *handler) handler(srv interface{}, serverStream grpc.ServerStream) error {
	// little bit of gRPC internals never hurt anyone
	fullMethodName, ok := grpc.MethodFromServerStream(serverStream)
	if !ok {
		return status.Errorf(codes.Internal, "lowLevelServerStream not exists in context")
	}

	backends, err := s.director(serverStream.Context(), fullMethodName)
	if err != nil {
		return err
	}

	if len(backends) == 0 {
		return status.Errorf(codes.Unavailable, "no backend connections for proxying")
	}

	var establishedConnections int
	backendConnections := make([]backendConnection, len(backends))

	clientCtx, clientCancel := context.WithCancel(serverStream.Context())
	defer clientCancel()

	for i := range backends {
		backendConnections[i].backend = backends[i]

		// We require that the backend's returned context inherits from the serverStream.Context().
		var outgoingCtx context.Context
		outgoingCtx, backendConnections[i].backendConn, backendConnections[i].connError = backends[i].GetConnection(clientCtx)

		if backendConnections[i].connError != nil {
			continue
		}

		backendConnections[i].clientStream, backendConnections[i].connError = grpc.NewClientStream(outgoingCtx, clientStreamDescForProxying,
			backendConnections[i].backendConn, fullMethodName)

		if backendConnections[i].connError != nil {
			continue
		}

		establishedConnections++
	}

	if len(backendConnections) != 1 {
		return s.handlerMulti(fullMethodName, serverStream, backendConnections)
	}

	// case of proxying one to one:
	if backendConnections[0].connError != nil {
		return backendConnections[0].connError
	}

	// Explicitly *do not close* s2cErrChan and c2sErrChan, otherwise the select below will not terminate.
	// Channels do not have to be closed, it is just a control flow mechanism, see
	// https://groups.google.com/forum/#!msg/golang-nuts/pZwdYRGxCIk/qpbHxRRPJdUJ
	s2cErrChan := s.forwardServerToClient(serverStream, &backendConnections[0])
	c2sErrChan := s.forwardClientToServer(&backendConnections[0], serverStream)
	// We don't know which side is going to stop sending first, so we need a select between the two.
	for i := 0; i < 2; i++ {
		select {
		case s2cErr := <-s2cErrChan:
			if s2cErr == io.EOF {
				// this is the happy case where the sender has encountered io.EOF, and won't be sending anymore./
				// the clientStream>serverStream may continue pumping though.
				//nolint: errcheck
				backendConnections[0].clientStream.CloseSend()
				break
			} else {
				// however, we may have gotten a receive error (stream disconnected, a read error etc) in which case we need
				// to cancel the clientStream to the backend, let all of its goroutines be freed up by the CancelFunc and
				// exit with an error to the stack
				return status.Errorf(codes.Internal, "failed proxying s2c: %v", s2cErr)
			}
		case c2sErr := <-c2sErrChan:
			// This happens when the clientStream has nothing else to offer (io.EOF), returned a gRPC error. In those two
			// cases we may have received Trailers as part of the call. In case of other errors (stream closed) the trailers
			// will be nil.
			serverStream.SetTrailer(backendConnections[0].clientStream.Trailer())
			// c2sErr will contain RPC error from client code. If not io.EOF return the RPC error as server stream error.
			if c2sErr != io.EOF {
				return c2sErr
			}
			return nil
		}
	}
	return status.Errorf(codes.Internal, "gRPC proxying should never reach this stage.")
}

func (s *handler) forwardClientToServer(src *backendConnection, dst grpc.ServerStream) chan error {
	ret := make(chan error, 1)
	go func() {
		f := &frame{}
		for i := 0; ; i++ {
			if err := src.clientStream.RecvMsg(f); err != nil {
				ret <- err // this can be io.EOF which is happy case
				break
			}

			var err error
			f.payload, err = src.backend.AppendInfo(f.payload)
			if err != nil {
				ret <- err
				break
			}

			if i == 0 {
				// This is a bit of a hack, but client to server headers are only readable after first client msg is
				// received but must be written to server stream before the first msg is flushed.
				// This is the only place to do it nicely.
				md, err := src.clientStream.Header()
				if err != nil {
					ret <- err
					break
				}
				if err := dst.SendHeader(md); err != nil {
					ret <- err
					break
				}
			}
			if err := dst.SendMsg(f); err != nil {
				ret <- err
				break
			}
		}
	}()
	return ret
}

func (s *handler) forwardServerToClient(src grpc.ServerStream, dst *backendConnection) chan error {
	ret := make(chan error, 1)
	go func() {
		f := &frame{}
		for i := 0; ; i++ {
			if err := src.RecvMsg(f); err != nil {
				ret <- err // this can be io.EOF which is happy case
				break
			}
			if err := dst.clientStream.SendMsg(f); err != nil {
				ret <- err
				break
			}
		}
	}()
	return ret
}
