package packet

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO: Property based tests using quick.Check
func TestVarInt_ReadFrom(t *testing.T) {
	type varIntTestCase struct {
		Name   string
		Input  []byte
		Length int
		Value  int
		Output []byte
	}

	varIntTestCases := []varIntTestCase{
		{
			Name:   "Packet",
			Input:  []byte{16, 0, 255},
			Length: 1,
			Value:  16,
			Output: []byte{16},
		},
		{
			Name:   "Packet Small",
			Input:  []byte{16},
			Length: 1,
			Value:  16,
			Output: []byte{16},
		},
		{
			Name:   "Handshake Data",
			Input:  []byte{242, 5, 9, 108, 111, 99, 97, 108, 104, 111, 115, 116, 99, 221, 2},
			Length: 2,
			Value:  754,
			Output: []byte{242, 5},
		},
		{
			Name:   "Handshake Data Small",
			Input:  []byte{242, 5},
			Length: 2,
			Value:  754,
			Output: []byte{242, 5},
		},
	}

	for _, testCase := range varIntTestCases {
		var i VarInt
		n, err := i.ReadFrom(bytes.NewReader(testCase.Input))
		assert.NoError(t, err, "test case: %v", testCase.Name)
		assert.Equal(t, int64(testCase.Length), n, "test case: %v", testCase.Name)
		assert.Equal(t, VarInt(testCase.Value), i, "test case: %v", testCase.Name)

		buf := bytes.NewBuffer(nil)
		n, err = i.WriteTo(buf)
		assert.Equal(t, testCase.Output, buf.Bytes(), "test case: %v", testCase.Name)
	}
}

func TestVarInt_SimpleTable(t *testing.T) {
	type testCase struct {
		Value int32
		Bytes []byte
	}

	testCases := []testCase{
		{Value: 0, Bytes: []byte{0}},
		{Value: 1, Bytes: []byte{1}},
		{Value: 2, Bytes: []byte{2}},
		{Value: 127, Bytes: []byte{127}},
		{Value: 128, Bytes: []byte{128, 1}},
		{Value: 255, Bytes: []byte{255, 1}},
		{Value: 2097151, Bytes: []byte{255, 255, 127}},
		{Value: 2147483647, Bytes: []byte{255, 255, 255, 255, 7}},
		{Value: -1, Bytes: []byte{255, 255, 255, 255, 15}},
		{Value: -2147483648, Bytes: []byte{128, 128, 128, 128, 8}},
	}

	for _, testCase := range testCases {
		var i VarInt
		_, err := i.ReadFrom(bytes.NewReader(testCase.Bytes))
		assert.NoError(t, err)
		assert.Equal(t, testCase.Value, int32(i))
		buf := bytes.NewBuffer(nil)
		_, err = i.WriteTo(buf)
		assert.NoError(t, err)
		assert.Equal(t, testCase.Bytes, buf.Bytes())
	}
}

var result int32

func BenchmarkVarInt_ReadFrom(b *testing.B) {
	type testCase struct {
		Value int32
		Bytes []byte
	}
	for i := 0; i < b.N; i++ {
		testCases := []testCase{
			{Value: 0, Bytes: []byte{0}},
			{Value: 1, Bytes: []byte{1}},
			{Value: 2, Bytes: []byte{2}},
			{Value: 127, Bytes: []byte{127}},
			{Value: 128, Bytes: []byte{128, 1}},
			{Value: 255, Bytes: []byte{255, 1}},
			{Value: 2097151, Bytes: []byte{255, 255, 127}},
			{Value: 2147483647, Bytes: []byte{255, 255, 255, 255, 7}},
			{Value: -1, Bytes: []byte{255, 255, 255, 255, 15}},
			{Value: -2147483648, Bytes: []byte{128, 128, 128, 128, 8}},
		}

		for _, testCase := range testCases {
			var i VarInt
			_, _ = i.ReadFrom(bytes.NewReader(testCase.Bytes))
			result = int32(i)
		}
	}
}
