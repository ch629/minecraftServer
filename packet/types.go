package packet

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
	"io"
)

type (
	VarInt    int32
	VarLong   int64
	UShort    uint16
	String    string
	ByteArray []byte
	UUID      uuid.UUID

	FieldEncoder io.WriterTo
	FieldDecoder io.ReaderFrom
)

const (
	MaxVarIntLen  = 5
	MaxVarLongLen = 10
)

func ReadFields(reader io.Reader, fields ...FieldDecoder) error {
	for i, field := range fields {
		if _, err := field.ReadFrom(reader); err != nil {
			return eris.Wrapf(err, "failed to read value at index %v", i)
		}
	}
	return nil
}

func ReadFieldsWithLengths(reader io.Reader, fields ...FieldDecoder) ([]int64, error) {
	lengths := make([]int64, len(fields))
	for i, readerFrom := range fields {
		if l, err := readerFrom.ReadFrom(reader); err != nil {
			return []int64{}, eris.Wrapf(err, "failed to read value at index %v", i)
		} else {
			lengths[i] = l
		}
	}
	return lengths, nil
}

func WriteFields(writer io.Writer, fields ...FieldEncoder) error {
	for i, field := range fields {
		if _, err := field.WriteTo(writer); err != nil {
			return eris.Wrapf(err, "failed to write value at index %v", i)
		}
	}
	return nil
}

func (i *VarInt) ReadFrom(reader io.Reader) (byteCount int64, err error) {
	result := int32(0)
	for read := byte(0x80); read&0x80 != 0; byteCount++ {
		// TODO: Validate we cut out at the right time for byte count
		if byteCount > MaxVarIntLen {
			err = eris.Errorf("VarInt too long, received '%v' chars", byteCount)
			return
		}

		read, err = readByte(reader)
		if err != nil {
			return
		}
		value := int32(read & 0x7F)
		result |= value << int32(7*byteCount)
	}
	*i = VarInt(result)
	return
}

// TODO: Max length
func (i VarInt) WriteTo(writer io.Writer) (int64, error) {
	buf := bytes.NewBuffer(nil)
	num := int32(i)
	count := int64(0)
	for {
		b := num & 0x7F
		num = int32(uint32(num) >> 7)
		if num != 0 {
			b |= 0x80
		}
		buf.WriteByte(byte(b))
		count++
		if num == 0 {
			break
		}
	}
	return io.CopyN(writer, buf, count)
}

func (s *UShort) ReadFrom(reader io.Reader) (byteCount int64, err error) {
	b := make([]byte, 2)
	byteCountInt, err := io.ReadFull(reader, b)
	byteCount = int64(byteCountInt)
	if err != nil {
		return
	}
	*s = UShort(int16(b[0])<<8 | int16(b[1]))
	return
}

func (s *String) ReadFrom(reader io.Reader) (int64, error) {
	var ba ByteArray
	byteCount, err := ba.ReadFrom(reader)
	if err != nil {
		return 0, err
	}

	*s = String(ba)
	return byteCount, nil
}

func (s String) WriteTo(writer io.Writer) (int64, error) {
	return ByteArray(s).WriteTo(writer)
}

// ByteArrayReader creates a io.Reader with a VarInt length prefixing the data
func ByteArrayReader(l VarInt, reader io.Reader) (io.Reader, error) {
	buf := bytes.NewBuffer(nil)
	_, err := l.WriteTo(buf)
	if err != nil {
		return nil, err
	}

	return io.MultiReader(buf, reader), nil
}

func ByteArrayFromLenReader(l VarInt, reader io.Reader) (ByteArray, error) {
	reader, err := ByteArrayReader(l, reader)

	if err != nil {
		return nil, err
	}

	var ba ByteArray
	// TODO: Return len too?
	_, err = ba.ReadFrom(reader)
	return ba, err
}

// ReadFrom creates a []byte, io.Reader needs to have a VarInt prefixing the byte data
func (b *ByteArray) ReadFrom(reader io.Reader) (int64, error) {
	var l VarInt
	_, err := l.ReadFrom(reader)
	if err != nil {
		return 0, err
	}
	bs := make([]byte, l)
	nn, err := io.ReadFull(reader, bs)
	*b = bs
	return int64(l) + int64(nn), err
}

func (b ByteArray) WriteTo(writer io.Writer) (int64, error) {
	l := VarInt(len(b))
	nn, err := l.WriteTo(writer)
	if err != nil {
		return 0, err
	}

	baLen, err := writer.Write(b)
	if err != nil {
		return 0, err
	}
	nn += int64(baLen)
	return nn, nil
}

func (u UUID) WriteTo(writer io.Writer) (int64, error) {
	nn, err := writer.Write(u[:])
	return int64(nn), err
}

func (u *UUID) ReadFrom(reader io.Reader) (int64, error) {
	// TODO: Do we need the * below?
	nn, err := io.ReadFull(reader, (*u)[:])
	return int64(nn), err
}

func readByte(r io.Reader) (byte, error) {
	if r, ok := r.(io.ByteReader); ok {
		return r.ReadByte()
	}
	v := make([]byte, 1)
	_, err := io.ReadFull(r, v)
	return v[0], err
}
