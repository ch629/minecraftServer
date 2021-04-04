package packet

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	type test struct {
		Int      int32  `pkt_type:"VarInt"`
		BALen    int32  `pkt_type:"VarInt"`
		BA       []byte `pkt_len:"BALen"`
		OptValue bool
		OptField int32 `pkt_opt:"OptValue"`
	}
	var te test
	buf := bytes.NewBuffer(nil)
	err := WriteFields(buf, VarInt(123456), VarInt(5), ByteArray([]byte{1, 2, 3, 4, 5}), Boolean(false))
	assert.NoError(t, err)
	assert.NoError(t, Unmarshal(&UncompressedPacket{
		readCloser: io.NopCloser(buf),
		dataLength: VarInt(buf.Len()),
	}, &te))
	assert.Equal(t, test{
		Int:      123456,
		BALen:    5,
		BA:       []byte{1, 2, 3, 4, 5},
		OptValue: false,
	}, te)
}
