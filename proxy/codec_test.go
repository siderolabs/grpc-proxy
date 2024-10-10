package proxy_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/mem"

	"github.com/siderolabs/grpc-proxy/proxy"
)

func TestCodec_ReadYourWrites(t *testing.T) {
	framePtr := proxy.NewFrame(nil)
	data := []byte{0xDE, 0xAD, 0xBE, 0xEF}
	codec := proxy.Codec()

	buffer := mem.Copy(data, mem.DefaultBufferPool())
	defer buffer.Free()

	require.NoError(t, codec.Unmarshal(mem.BufferSlice{buffer}, framePtr), "unmarshalling must go ok")
	out, err := codec.Marshal(framePtr)
	require.NoError(t, err, "no marshal error")
	require.Equal(t, data, out.Materialize(), "output and data must be the same")

	out.Free()
	buffer.Free()
	buffer = mem.Copy([]byte{0x55}, mem.DefaultBufferPool())

	// reuse
	require.NoError(t, codec.Unmarshal(mem.BufferSlice{buffer}, framePtr), "unmarshalling must go ok")
	out, err = codec.Marshal(framePtr)
	require.NoError(t, err, "no marshal error")
	require.Equal(t, []byte{0x55}, out.Materialize(), "output and data must be the same")

	out.Free()
}
