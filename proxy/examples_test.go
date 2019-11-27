// Copyright 2017 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package proxy_test

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/talos-systems/grpc-proxy/proxy"
)

var (
	director proxy.StreamDirector
)

func ExampleRegisterService() {
	// A gRPC server with the proxying codec enabled.
	server := grpc.NewServer(grpc.CustomCodec(proxy.Codec()))
	// Register a TestService with 4 of its methods explicitly.
	proxy.RegisterService(server, director,
		"talos.testproto.TestService",
		proxy.WithMethodNames("PingEmpty", "Ping", "PingError", "PingList"),
		proxy.WithStreamedMethodNames("PingList"),
	)
}

func ExampleTransparentHandler() {
	grpc.NewServer(
		grpc.CustomCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(proxy.TransparentHandler(director)))
}

// Provide sa simple example of a director that shields internal services and dials a staging or production backend.
// This is a *very naive* implementation that creates a new connection on every request. Consider using pooling.
func ExampleStreamDirector() {
	simpleBackendGen := func(hostname string) proxy.Backend {
		return &proxy.SingleBackend{
			GetConn: func(ctx context.Context) (context.Context, *grpc.ClientConn, error) {
				md, _ := metadata.FromIncomingContext(ctx)

				// Copy the inbound metadata explicitly.
				outCtx := metadata.NewOutgoingContext(ctx, md.Copy())
				// Make sure we use DialContext so the dialing can be cancelled/time out together with the context.
				conn, err := grpc.DialContext(ctx, hostname, grpc.WithCodec(proxy.Codec())) // nolint: staticcheck

				return outCtx, conn, err
			},
		}
	}

	director = func(ctx context.Context, fullMethodName string) (proxy.Mode, []proxy.Backend, error) {
		// Make sure we never forward internal services.
		if strings.HasPrefix(fullMethodName, "/com.example.internal.") {
			return proxy.One2One, nil, status.Errorf(codes.Unimplemented, "Unknown method")
		}
		md, ok := metadata.FromIncomingContext(ctx)

		if ok {
			// Decide on which backend to dial
			if val, exists := md[":authority"]; exists && val[0] == "staging.api.example.com" {
				return proxy.One2One, []proxy.Backend{simpleBackendGen("api-service.staging.svc.local")}, nil
			} else if val, exists := md[":authority"]; exists && val[0] == "api.example.com" {
				return proxy.One2One, []proxy.Backend{simpleBackendGen("api-service.prod.svc.local")}, nil
			}
		}
		return proxy.One2One, nil, status.Errorf(codes.Unimplemented, "Unknown method")
	}
}
