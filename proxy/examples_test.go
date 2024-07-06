// Copyright 2017 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package proxy_test

import (
	"context"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/siderolabs/grpc-proxy/proxy"
)

var director proxy.StreamDirector

// ExampleRegisterService is a simple example of registering a service with the proxy.
func ExampleRegisterService() {
	// A gRPC server with the proxying codec enabled.
	server := grpc.NewServer(grpc.ForceServerCodec(proxy.Codec()))

	// Register a TestService with 4 of its methods explicitly.
	proxy.RegisterService(server, director,
		"talos.testproto.TestService",
		proxy.WithMethodNames("PingEmpty", "Ping", "PingError", "PingList"),
		proxy.WithStreamedMethodNames("PingList"),
	)

	// Output:
}

// ExampleTransparentHandler is an example of redirecting all requests to the proxy.
func ExampleTransparentHandler() {
	grpc.NewServer(
		grpc.ForceServerCodec(proxy.Codec()),
		grpc.UnknownServiceHandler(proxy.TransparentHandler(director)),
	)

	// Output:
}

// Provide sa simple example of a director that shields internal services and dials a staging or production backend.
// This is a *very naive* implementation that creates a new connection on every request. Consider using pooling.
func ExampleStreamDirector() {
	simpleBackendGen := func(hostname string) (proxy.Backend, error) {
		conn, err := grpc.NewClient(
			hostname,
			grpc.WithDefaultCallOptions(grpc.ForceCodec(proxy.Codec())),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, err
		}

		return &proxy.SingleBackend{
			GetConn: func(ctx context.Context) (context.Context, *grpc.ClientConn, error) {
				md, _ := metadata.FromIncomingContext(ctx)

				// Copy the inbound metadata explicitly.
				outCtx := metadata.NewOutgoingContext(ctx, md.Copy())

				return outCtx, conn, nil
			},
		}, nil
	}

	stagingBackend, err := simpleBackendGen("api-service.staging.svc.local")
	if err != nil {
		log.Fatal("failed to create staging backend:", err)
	}

	prodBackend, err := simpleBackendGen("api-service.prod.svc.local")
	if err != nil {
		log.Fatal("failed to create production backend:", err)
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
				return proxy.One2One, []proxy.Backend{stagingBackend}, nil
			} else if val, exists := md[":authority"]; exists && val[0] == "api.example.com" {
				return proxy.One2One, []proxy.Backend{prodBackend}, nil
			}
		}

		return proxy.One2One, nil, status.Errorf(codes.Unimplemented, "Unknown method")
	}

	// Output:
}
