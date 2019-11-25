// Copyright 2019 Andrey Smirnov. All Rights Reserved.
// See LICENSE for licensing terms.

package proxy

import (
	"errors"
	"fmt"
	"io"

	"github.com/hashicorp/go-multierror"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *handler) handlerMulti(fullMethodName string, serverStream grpc.ServerStream, backendConnections []backendConnection) error {
	// wrap the stream for safe concurrent access
	serverStream = &ServerStreamWrapper{ServerStream: serverStream}

	s2cErrChan := s.forwardServerToClientsMulti(serverStream, backendConnections)
	var c2sErrChan chan error

	if _, streamed := s.streamedMethods[fullMethodName]; !streamed {
		c2sErrChan = s.forwardClientsToServerMultiUnary(backendConnections, serverStream)
	} else {
		c2sErrChan = s.forwardClientsToServerMultiStreaming(backendConnections, serverStream)
	}

	for i := 0; i < 2; i++ {
		select {
		case s2cErr := <-s2cErrChan:
			if s2cErr == io.EOF {
				// this is the happy case where the sender has encountered io.EOF, and won't be sending anymore./
				// the clientStream>serverStream may continue pumping though.
				for i := range backendConnections {
					if backendConnections[i].clientStream != nil {
						backendConnections[i].clientStream.CloseSend() //nolint: errcheck
					}
				}
				break
			} else {
				// however, we may have gotten a receive error (stream disconnected, a read error etc) in which case we need
				// to cancel the clientStream to the backend, let all of its goroutines be freed up by the CancelFunc and
				// exit with an error to the stack
				return status.Errorf(codes.Internal, "failed proxying s2c: %v", s2cErr)
			}
		case c2sErr := <-c2sErrChan:
			// c2sErr will contain RPC error from client code. If not io.EOF return the RPC error as server stream error.
			if c2sErr != io.EOF {
				return c2sErr
			}
			return nil
		}
	}
	return status.Errorf(codes.Internal, "gRPC proxying should never reach this stage.")
}

// formatError tries to format error from upstream as message to the client
func (s *handler) formatError(src *backendConnection, backendErr error) ([]byte, error) {
	payload, err := src.backend.BuildError(backendErr)
	if err != nil {
		return nil, fmt.Errorf("error building error for %s: %w", src.backend, err)
	}

	if payload == nil {
		err = backendErr
	}

	return payload, err
}

// sendError tries to deliver error back to the client via dst
//
// if sendError fails to deliver the error, error is returned
// if sendError successfully delivers the error, nil is returned
func (s *handler) sendError(src *backendConnection, dst grpc.ServerStream, backendErr error) error {
	payload, err := s.formatError(src, backendErr)
	if err != nil {
		return err
	}

	f := &frame{payload: payload}
	if err = dst.SendMsg(f); err != nil {
		return fmt.Errorf("error sending error back: %w", err)
	}

	return nil
}

// one:many proxying, unary call version (merging results)
func (s *handler) forwardClientsToServerMultiUnary(sources []backendConnection, dst grpc.ServerStream) chan error {
	ret := make(chan error, 1)

	payloadCh := make(chan []byte, len(sources))
	errCh := make(chan error, len(sources))

	for i := 0; i < len(sources); i++ {
		go func(src *backendConnection) {
			errCh <- func() error {
				if src.connError != nil {
					payload, err := s.formatError(src, src.connError)
					if err != nil {
						return err
					}

					payloadCh <- payload
					return nil
				}

				f := &frame{}
				for j := 0; ; j++ {
					if err := src.clientStream.RecvMsg(f); err != nil {
						if err == io.EOF {
							// This happens when the clientStream has nothing else to offer (io.EOF), returned a gRPC error. In those two
							// cases we may have received Trailers as part of the call. In case of other errors (stream closed) the trailers
							// will be nil.
							dst.SetTrailer(src.clientStream.Trailer())
							return nil
						}

						payload, err := s.formatError(src, err)
						if err != nil {
							return err
						}

						payloadCh <- payload
						return nil
					}
					if j == 0 {
						// This is a bit of a hack, but client to server headers are only readable after first client msg is
						// received but must be written to server stream before the first msg is flushed.
						// This is the only place to do it nicely.
						md, err := src.clientStream.Header()
						if err != nil {
							payload, err := s.formatError(src, err)
							if err != nil {
								return err
							}

							payloadCh <- payload
							return nil
						}

						if err := dst.SetHeader(md); err != nil {
							return fmt.Errorf("error setting headers from client %s: %w", src.backend, err)
						}
					}

					var err error
					f.payload, err = src.backend.AppendInfo(f.payload)
					if err != nil {
						return fmt.Errorf("error appending info for %s: %w", src.backend, err)
					}

					payloadCh <- f.payload
				}
			}()
		}(&sources[i])
	}

	go func() {
		var multiErr *multierror.Error

		for range sources {
			multiErr = multierror.Append(multiErr, <-errCh)
		}

		if multiErr.ErrorOrNil() != nil {
			ret <- multiErr.ErrorOrNil()
			return
		}

		close(payloadCh)

		var merged []byte
		for b := range payloadCh {
			merged = append(merged, b...)
		}

		ret <- dst.SendMsg(&frame{payload: merged})
	}()

	return ret
}

// one:many proxying, streaming version (no merge)
func (s *handler) forwardClientsToServerMultiStreaming(sources []backendConnection, dst grpc.ServerStream) chan error {
	ret := make(chan error, 1)

	errCh := make(chan error, len(sources))

	for i := range sources {
		go func(src *backendConnection) {
			errCh <- func() error {
				if src.connError != nil {
					return s.sendError(src, dst, src.connError)
				}

				f := &frame{}
				for j := 0; ; j++ {
					if err := src.clientStream.RecvMsg(f); err != nil {
						if err == io.EOF {
							// This happens when the clientStream has nothing else to offer (io.EOF), returned a gRPC error. In those two
							// cases we may have received Trailers as part of the call. In case of other errors (stream closed) the trailers
							// will be nil.
							dst.SetTrailer(src.clientStream.Trailer())
							return nil
						}
						return s.sendError(src, dst, err)
					}
					if j == 0 {
						// This is a bit of a hack, but client to server headers are only readable after first client msg is
						// received but must be written to server stream before the first msg is flushed.
						// This is the only place to do it nicely.
						md, err := src.clientStream.Header()
						if err != nil {
							return s.sendError(src, dst, err)
						}
						if err := dst.SetHeader(md); err != nil {
							return fmt.Errorf("error setting headers from client %s: %w", src.backend, err)
						}
					}

					var err error
					f.payload, err = src.backend.AppendInfo(f.payload)
					if err != nil {
						return fmt.Errorf("error appending info for %s: %w", src.backend, err)
					}

					if err = dst.SendMsg(f); err != nil {
						return fmt.Errorf("error sending back to server from %s: %w", src.backend, err)
					}
				}

			}()
		}(&sources[i])
	}

	go func() {
		var multiErr *multierror.Error

		for range sources {
			multiErr = multierror.Append(multiErr, <-errCh)
		}

		ret <- multiErr.ErrorOrNil()
	}()

	return ret
}

func (s *handler) forwardServerToClientsMulti(src grpc.ServerStream, destinations []backendConnection) chan error {
	ret := make(chan error, 1)
	go func() {
		f := &frame{}
		for {
			if err := src.RecvMsg(f); err != nil {
				ret <- err
				return
			}

			errCh := make(chan error)

			for i := range destinations {
				go func(dst *backendConnection) {
					errCh <- func() error {
						if dst.clientStream == nil || dst.connError != nil {
							return nil // skip it
						}

						return dst.clientStream.SendMsg(f)
					}()
				}(&destinations[i])
			}

			liveDestinations := 0
			for range destinations {
				if err := <-errCh; err == nil {
					liveDestinations++
				}
			}

			if liveDestinations == 0 {
				ret <- errors.New("no backend connections to forward to are available")
				return
			}
		}
	}()
	return ret
}
