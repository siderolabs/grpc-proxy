# gRPC Proxy

[![Go Report Card](https://goreportcard.com/badge/github.com/talos-systems/grpc-proxy)](https://goreportcard.com/report/github.com/talos-systems/grpc-proxy)
[![GoDoc](http://img.shields.io/badge/GoDoc-Reference-blue.svg)](https://godoc.org/github.com/talos-systems/grpc-proxy)
[![Apache 2.0 License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

[gRPC Go](https://github.com/grpc/grpc-go) Proxy server

## Project Goal

Build a transparent reverse proxy for gRPC targets that will make it easy to expose gRPC services
over the Internet.

This includes:

* no needed knowledge of the semantics of requests exchanged in the call (independent rollouts)
* easy declarative definition of backends and their mappings to frontends
* simple round-robin load balancing of inbound requests from a single connection to multiple backends

## Proxying Modes

There are two proxying modes supported:

* one to one: in this mode data passed back and forth is transmitted as is without any modifications;
* one to many: one client connection gets mapped into multiple upstream connections, results might be aggregated
  (for unary calls), errors translated into response messages; this mode requires a special layout of protobuf messages.

## Proxy Handler

The package [`proxy`](proxy/) contains a generic gRPC reverse proxy handler that allows a gRPC server not to
know about method names and their request/response data types.
Please consult the package documentation.
Here you can find an example usage.

First, define `Backend` implementation to identify specific upstream.
For one to one proxying, `SingleBackend` might be used:

```go
conn, err := grpc.NewClient(
    "api-service.staging.svc.local",
    grpc.WithDefaultCallOptions(grpc.ForceCodec(proxy.Codec())),
)
if err != nil {
    log.Fatal(err)
}

backend := &proxy.SingleBackend{
    GetConn: func(ctx context.Context) (context.Context, *grpc.ClientConn, error) {
        md, _ := metadata.FromIncomingContext(ctx)

        // Copy the inbound metadata explicitly.
        outCtx := metadata.NewOutgoingContext(ctx, md.Copy())

        return outCtx, conn, nil
    },
}
```

Defining a `StreamDirector` that decides where (if at all) to send the request

```go
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
```

Then you need to register it with a `grpc.Server`.
The server may have other handlers that will be served locally:

```go
server := grpc.NewServer(
    grpc.ForceServerCodec(proxy.Codec()),
    grpc.UnknownServiceHandler(
        proxy.TransparentHandler(director),
        proxy.WithMode(proxy.One2One),
    ))
pb_test.RegisterTestServiceServer(server, &testImpl{})
```

## One to Many Proxying

In one to many proxying mode, it's critical to identify source of each message proxied back from the upstreams.
Also upstream error shouldn't fail whole request and instead return errors as messages back.
In order to achieve this goal, protobuf response message should follow the same structure:

1. Every response should be `repeated` list of response messages so that responses from multiple upstreams might be
concatenated to build a combined response from all the upstreams.

2. Response should contain common metadata fields which allow grpc-proxy to inject source information and error information
into response.

## Talks

* "Transparent gRPC proxy in Go" at [GopherCon Russia 2021](https://www.gophercon-russia.ru/) [slides](https://speakerdeck.com/smira/transparent-grpc-gateway-in-go)

## History

This is a fork of awesome [mwitkow/grpc-proxy](https://github.com/mwitkow/grpc-proxy) package with added support
for one to many proxying.

## License

`grpc-proxy` is released under the Apache 2.0 license.
See [LICENSE.txt](LICENSE.txt).
