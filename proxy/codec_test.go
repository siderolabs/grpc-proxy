package proxy_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/mem"

	"github.com/siderolabs/grpc-proxy/proxy"
	talos_testproto "github.com/siderolabs/grpc-proxy/testservice"
)

func TestCodec_ReadYourWrites(t *testing.T) {
	d := []byte{0xDE, 0xAD, 0xBE, 0xEF}

	for key, val := range map[string][]byte{
		"short message": d,
		"long message":  bytes.Repeat(d, 3072),
	} {
		t.Run(key, func(t *testing.T) {
			framePtr := proxy.NewFrame(nil)
			codec := proxy.Codec()

			buffer := mem.Copy(val, mem.DefaultBufferPool())
			defer func() { buffer.Free() }()

			require.NoError(t, codec.Unmarshal(mem.BufferSlice{buffer}, framePtr), "unmarshalling must go ok")
			out, err := codec.Marshal(framePtr)
			require.NoError(t, err, "no marshal error")
			require.Equal(t, val, out.Materialize(), "output and data must be the same")

			out.Free()
			buffer.Free()
			buffer = mem.Copy([]byte{0x55}, mem.DefaultBufferPool())

			// reuse
			require.NoError(t, codec.Unmarshal(mem.BufferSlice{buffer}, framePtr), "unmarshalling must go ok")
			out, err = codec.Marshal(framePtr)
			require.NoError(t, err, "no marshal error")
			require.Equal(t, []byte{0x55}, out.Materialize(), "output and data must be the same")

			out.Free()
		})
	}
}

func TestCodecUsualMessage(t *testing.T) {
	const msg = "short message"

	for key, val := range map[string]string{
		"short message": "edbca",
		"long message":  strings.Repeat(msg, 3072),
	} {
		t.Run(key, func(t *testing.T) {
			msg := &talos_testproto.PingRequest{Value: val}

			codec := proxy.Codec()

			out, err := codec.Marshal(msg)
			require.NoError(t, err, "no marshal error")

			defer out.Free()

			var dst talos_testproto.PingRequest

			require.NoError(t, codec.Unmarshal(out, &dst), "unmarshalling must go ok")
			require.NotZero(t, dst.Value, "output must not be zero")
			require.Equal(t, msg.Value, dst.Value, "output and data must be the same")
		})
	}
}
