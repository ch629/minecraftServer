package packet

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	type test struct {
		Int int32 `pkt:"VarInt"`
	}
	var te test
	l := VarInt(1234567890)
	buf := bytes.NewBuffer(nil)
	_, err := l.WriteTo(buf)
	assert.NoError(t, err)
	assert.NoError(t, Unmarshal(&UncompressedPacket{readCloser: io.NopCloser(buf)}, &te))
	assert.Equal(t, test{Int: 1234567890}, te)
}
