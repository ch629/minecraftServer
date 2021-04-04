package packet

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestMarshal(t *testing.T) {
	type testStruct struct {
		Int      int32  `pkt_type:"VarInt"`
		BALen    int32  `pkt_type:"VarInt"`
		BA       []byte `pkt_len:"BALen"`
		OptValue bool
		OptField int32 `pkt_opt:"OptValue"`
		Pos      Position
	}

	test := testStruct{
		Int:      123456,
		BALen:    5,
		BA:       []byte{1, 2, 3, 4, 5},
		OptValue: false,
		Pos: Position{
			X: 10,
			Y: 20,
			Z: 30,
		},
	}

	bs, err := Marshal(&test)
	assert.NoError(t, err)

	buf := bytes.NewBuffer(bs)
	var newTest testStruct
	assert.NoError(t, Unmarshal(&UncompressedPacket{
		readCloser: io.NopCloser(buf),
		dataLength: VarInt(buf.Len()),
	}, &newTest))
	assert.Equal(t, test, newTest)
}
