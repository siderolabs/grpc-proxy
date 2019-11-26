# proxy
--
    import "."

Package proxy provides a reverse proxy handler for gRPC.

The implementation allows a `grpc.Server` to pass a received ServerStream to a
ClientStream without understanding the semantics of the messages exchanged. It
basically provides a transparent reverse-proxy.

This package is intentionally generic, exposing a `StreamDirector` function that
allows users of this package to implement whatever logic of backend-picking,
dialing and service verification to perform.

See examples on documented functions.

## Usage

#### func  Codec

```go
func Codec() grpc.Codec
```
Codec returns a proxying grpc.Codec with the default protobuf codec as parent.

See CodecWithParent.

nolint: staticcheck

#### func  CodecWithParent

```go
func CodecWithParent(fallback grpc.Codec) grpc.Codec
```
CodecWithParent returns a proxying grpc.Codec with a user provided codec as
parent.

This codec is *crucial* to the functioning of the proxy. It allows the proxy
server to be oblivious to the schema of the forwarded messages. It basically
treats a gRPC message frame as raw bytes. However, if the server handler, or the
client caller are not proxy-internal functions it will fall back to trying to
decode the message using a fallback codec.

nolint: staticcheck

#### func  RegisterService

```go
func RegisterService(server *grpc.Server, director StreamDirector, serviceName string, options ...Option)
```
RegisterService sets up a proxy handler for a particular gRPC service and
method. The behavior is the same as if you were registering a handler method,
e.g. from a codegenerated pb.go file.

This can *only* be used if the `server` also uses grpc.CustomCodec()
ServerOption.

#### func  TransparentHandler

```go
func TransparentHandler(director StreamDirector, options ...Option) grpc.StreamHandler
```
TransparentHandler returns a handler that attempts to proxy all requests that
are not registered in the server. The indented use here is as a transparent
proxy, where the server doesn't know about the services implemented by the
backends. It should be used as a `grpc.UnknownServiceHandler`.

This can *only* be used if the `server` also uses grpc.CustomCodec()
ServerOption.

#### type Backend

```go
type Backend interface {
	// String provides backend name for logging and errors.
	String() string

	// GetConnection returns a grpc connection to the backend.
	//
	// The context returned from this function should be the context for the *outgoing* (to backend) call. In case you want
	// to forward any Metadata between the inbound request and outbound requests, you should do it manually. However, you
	// *must* propagate the cancel function (`context.WithCancel`) of the inbound context to the one returned.
	GetConnection(ctx context.Context) (context.Context, *grpc.ClientConn, error)

	// AppendInfo is called to enhance response from the backend with additional data.
	//
	// Usecase might be appending backend endpoint (or name) to the protobuf serialized response, so that response is enhanced
	// with source information. This is particularly important for one to many calls, when it is required to identify
	// response from each of the backends participating in the proxying.
	//
	// If not additional proxying is required, simply returning the buffer without changes works fine.
	AppendInfo(resp []byte) ([]byte, error)

	// BuildError is called to convert error from upstream into response field.
	//
	// BuildError is never called for one to one proxying, in that case all the errors are returned back to the caller
	// as grpc errors.
	//
	// When proxying one to many, if one the requests fails or upstream returns an error, it is undesirable to fail the whole
	// request and discard responses from other backends. BuildError converts (marshals) error from backend into protobuf encoded
	// response which is analyzed by the caller, so that caller reaching out to N upstreams receives N1 successful responses and
	// N2 error responses so that N1 + N2 == N.
	//
	// If BuildError returns nil, error is returned as grpc error (failing whole request).
	BuildError(err error) ([]byte, error)
}
```

Backend wraps information about upstream connection.

For simple one-to-one proxying, not much should be done in the Backend, simply
providing a connection is enough.

When proxying one-to-many and aggregating results, Backend might be used to
append additional fields to upstream response to support more complicated
proxying.

#### type Mode

```go
type Mode int
```

Mode specifies proxying mode: one2one (transparent) or one2many (aggregation,
error wrapping).

```go
const (
	One2One Mode = iota
	One2Many
)
```
Mode constants.

#### type Option

```go
type Option func(*handlerOptions)
```

Option configures gRPC proxy

#### func  WithMethodNames

```go
func WithMethodNames(methodNames ...string) Option
```
WithMethodNames configures list of method names to proxy for non-transparent
handler.

#### func  WithMode

```go
func WithMode(mode Mode) Option
```
WithMode sets proxying mode: One2One or One2Many.

Default mode is One2One.

#### func  WithStreamedDetector

```go
func WithStreamedDetector(detector StreamedDetectorFunc) Option
```
WithStreamedDetector configures a function to detect streamed methods.

This is only important for one2many proxying.

#### func  WithStreamedMethodNames

```go
func WithStreamedMethodNames(streamedMethodNames ...string) Option
```
WithStreamedMethodNames configures list of streamed method names.

This is only important for one2many proxying. This option can't be used with
TransparentHandler.

#### type ServerStreamWrapper

```go
type ServerStreamWrapper struct {
	grpc.ServerStream
}
```

ServerStreamWrapper wraps grpc.ServerStream and adds locking to send path

#### func (*ServerStreamWrapper) SendHeader

```go
func (wrapper *ServerStreamWrapper) SendHeader(md metadata.MD) error
```
SendHeader sends the header metadata. The provided md and headers set by
SetHeader() will be sent. It fails if called multiple times.

#### func (*ServerStreamWrapper) SendMsg

```go
func (wrapper *ServerStreamWrapper) SendMsg(m interface{}) error
```
SendMsg sends a message. On error, SendMsg aborts the stream and the error is
returned directly.

SendMsg blocks until:

    - There is sufficient flow control to schedule m with the transport, or
    - The stream is done, or
    - The stream breaks.

SendMsg does not wait until the message is received by the client. An untimely
stream closure may result in lost messages.

It is safe to have a goroutine calling SendMsg and another goroutine calling
RecvMsg on the same stream at the same time, but it is not safe to call SendMsg
on the same stream in different goroutines.

#### func (*ServerStreamWrapper) SetHeader

```go
func (wrapper *ServerStreamWrapper) SetHeader(md metadata.MD) error
```
SetHeader sets the header metadata. It may be called multiple times. When call
multiple times, all the provided metadata will be merged. All the metadata will
be sent out when one of the following happens:

    - ServerStream.SendHeader() is called;
    - The first response is sent out;
    - An RPC status is sent out (error or success).

#### func (*ServerStreamWrapper) SetTrailer

```go
func (wrapper *ServerStreamWrapper) SetTrailer(md metadata.MD)
```
SetTrailer sets the trailer metadata which will be sent with the RPC status.
When called more than once, all the provided metadata will be merged.

#### type SingleBackend

```go
type SingleBackend struct {
	// GetConn returns a grpc connection to the backend.
	//
	// The context returned from this function should be the context for the *outgoing* (to backend) call. In case you want
	// to forward any Metadata between the inbound request and outbound requests, you should do it manually. However, you
	// *must* propagate the cancel function (`context.WithCancel`) of the inbound context to the one returned.
	GetConn func(ctx context.Context) (context.Context, *grpc.ClientConn, error)
}
```

SingleBackend implements a simple wrapper around get connection function of one
to one proxying.

SingleBackend implements Backend interface and might be used as an easy wrapper
for one to one proxying.

#### func (*SingleBackend) AppendInfo

```go
func (sb *SingleBackend) AppendInfo(resp []byte) ([]byte, error)
```
AppendInfo is called to enhance response from the backend with additional data.

#### func (*SingleBackend) BuildError

```go
func (sb *SingleBackend) BuildError(err error) ([]byte, error)
```
BuildError is called to convert error from upstream into response field.

#### func (*SingleBackend) GetConnection

```go
func (sb *SingleBackend) GetConnection(ctx context.Context) (context.Context, *grpc.ClientConn, error)
```
GetConnection returns a grpc connection to the backend.

#### func (*SingleBackend) String

```go
func (sb *SingleBackend) String() string
```

#### type StreamDirector

```go
type StreamDirector func(ctx context.Context, fullMethodName string) ([]Backend, error)
```

StreamDirector returns a list of Backend objects to forward the call to.

There are two proxying modes:

    1. one to one: StreamDirector returns a single Backend object - proxying is done verbatim, Backend.AppendInfo might
       be used to enhance response with source information (or it might be skipped).
    2. one to many: StreamDirector returns more than one Backend object - for unary calls responses from Backend objects
       are aggregated by concatenating protobuf responses (requires top-level `repeated` protobuf definition) and errors
       are wrapped as responses via BuildError. Responses are potentially enhanced via AppendInfo.

The presence of the `Context` allows for rich filtering, e.g. based on Metadata
(headers). If no handling is meant to be done, a `codes.NotImplemented` gRPC
error should be returned.

It is worth noting that the StreamDirector will be fired *after* all server-side
stream interceptors are invoked. So decisions around authorization, monitoring
etc. are better to be handled there.

See the rather rich example.

#### type StreamedDetectorFunc

```go
type StreamedDetectorFunc func(fullMethodName string) bool
```

StreamedDetectorFunc reports is gRPC is doing streaming (only for one2many
proxying).
