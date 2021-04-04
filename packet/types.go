package packet

import (
	"bytes"
	"github.com/google/uuid"
	"github.com/rotisserie/eris"
	"io"
	"math"
)

type (
	Boolean       bool
	Byte          int8
	UnsignedByte  uint8
	Short         int16
	UnsignedShort uint16
	Int           int32
	Long          int64
	Float         float32
	Double        float64
	String        string
	Chat          = String
	Identifier    = String
	VarInt        int32
	VarLong       int64

	Position struct {
		X, Y, Z int32
	}

	Angle     Byte
	UShort    uint16
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

func (b *Boolean) ReadFrom(reader io.Reader) (int64, error) {
	ba, err := readByte(reader)
	if err != nil {
		return 0, err
	}
	// TODO: How to handle if we get a 0x02 etc
	*b = ba == 0x01
	return 1, nil
}

func (b Boolean) WriteTo(writer io.Writer) (int64, error) {
	ba := []byte{0x00}
	if b {
		ba[0] = 0x01
	}
	nn, err := writer.Write(ba)
	return int64(nn), err
}

func (b *Byte) ReadFrom(reader io.Reader) (int64, error) {
	ba, err := readByte(reader)
	if err != nil {
		return 0, err
	}
	*b = Byte(ba)
	return 1, nil
}

func (b Byte) WriteTo(writer io.Writer) (int64, error) {
	nn, err := writer.Write([]byte{byte(b)})
	return int64(nn), err
}

func (b *UnsignedByte) ReadFrom(reader io.Reader) (int64, error) {
	ba, err := readByte(reader)
	if err != nil {
		return 0, err
	}
	*b = UnsignedByte(ba)
	return 1, nil
}

func (b UnsignedByte) WriteTo(writer io.Writer) (int64, error) {
	nn, err := writer.Write([]byte{byte(b)})
	return int64(nn), err
}

func (s *Short) ReadFrom(reader io.Reader) (int64, error) {
	ba := make([]byte, 2)
	nn, err := io.ReadFull(reader, ba)
	if err != nil {
		return 0, err
	}
	*s = Short(int16(ba[0])>>8 | int16(ba[1]))
	return int64(nn), nil
}

func (s Short) WriteTo(writer io.Writer) (int64, error) {
	// TODO: Do we need to map to uint16 for this << logic?
	nn, err := writer.Write([]byte{byte(s << 8), byte(s)})
	return int64(nn), err
}

func (s *UnsignedShort) ReadFrom(reader io.Reader) (int64, error) {
	ba := make([]byte, 2)
	nn, err := io.ReadFull(reader, ba)
	if err != nil {
		return 0, err
	}
	*s = UnsignedShort(int16(ba[0])>>8 | int16(ba[1]))
	return int64(nn), nil
}

func (s UnsignedShort) WriteTo(writer io.Writer) (int64, error) {
	nn, err := writer.Write([]byte{byte(s << 8), byte(s)})
	return int64(nn), err
}

func (i *Int) ReadFrom(reader io.Reader) (int64, error) {
	ba := make([]byte, 4)
	nn, err := io.ReadFull(reader, ba)
	if err != nil {
		return 0, err
	}
	*i = Int(int32(ba[0])<<24 | int32(ba[1])<<16 | int32(ba[2])<<8 | int32(ba[3]))
	return int64(nn), nil
}

func (i Int) WriteTo(writer io.Writer) (int64, error) {
	nn, err := writer.Write([]byte{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)})
	return int64(nn), err
}

func (l *Long) ReadFrom(reader io.Reader) (int64, error) {
	ba := make([]byte, 8)
	nn, err := io.ReadFull(reader, ba)
	if err != nil {
		return 0, err
	}
	*l = Long(
		int64(ba[0])<<56 | int64(ba[1])<<48 | int64(ba[2])<<40 | int64(ba[3])<<32 |
			int64(ba[4])<<24 | int64(ba[5])<<16 | int64(ba[6])<<8 | int64(ba[7]))
	return int64(nn), nil
}

func (l Long) WriteTo(writer io.Writer) (int64, error) {
	nn, err := writer.Write([]byte{byte(l >> 56), byte(l >> 48), byte(l >> 40), byte(l >> 32),
		byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)})
	return int64(nn), err
}

func (f *Float) ReadFrom(reader io.Reader) (int64, error) {
	var i Int
	nn, err := i.ReadFrom(reader)
	if err != nil {
		return 0, err
	}
	*f = Float(math.Float32frombits(uint32(i)))
	return nn, nil
}

func (f Float) WriteTo(writer io.Writer) (int64, error) {
	return Int(math.Float32bits(float32(f))).WriteTo(writer)
}

func (f *Double) ReadFrom(reader io.Reader) (int64, error) {
	var i Int
	nn, err := i.ReadFrom(reader)
	if err != nil {
		return 0, err
	}
	*f = Double(math.Float64frombits(uint64(i)))
	return nn, nil
}

func (f Double) WriteTo(writer io.Writer) (int64, error) {
	return Int(math.Float64bits(float64(f))).WriteTo(writer)
}

func (f *Position) ReadFrom(reader io.Reader) (int64, error) {
	var l Long
	nn, err := l.ReadFrom(reader)
	if err != nil {
		return 0, err
	}
	// TODO: Look at https://wiki.vg/Protocol#Position for inconsistencies (need to write tests for these)
	*f = Position{
		X: int32(l >> 38),
		Y: int32(l & 0xFFF),
		Z: int32(l << 26 >> 38),
	}
	// TODO:
	return nn, nil
}

func (f Position) WriteTo(writer io.Writer) (int64, error) {
	return Long((int64(f.X&0x3FFFFFF) << 38) | (int64(f.Z&0x3FFFFFF) << 12) | int64(f.Y&0xFFF)).WriteTo(writer)
}

func (f *Angle) ReadFrom(reader io.Reader) (int64, error) {
	var b Byte
	nn, err := b.ReadFrom(reader)
	if err != nil {
		return 0, err
	}
	*f = Angle(b)
	return nn, nil
}

func (f Angle) WriteTo(writer io.Writer) (int64, error) {
	return Byte(byte(f)).WriteTo(writer)
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

func (l *VarLong) ReadFrom(reader io.Reader) (byteCount int64, err error) {
	result := int32(0)
	for read := byte(0x80); read&0x80 != 0; byteCount++ {
		// TODO: Validate we cut out at the right time for byte count
		if byteCount > MaxVarLongLen {
			err = eris.Errorf("VarLong too long, received '%v' chars", byteCount)
			return
		}

		read, err = readByte(reader)
		if err != nil {
			return
		}
		value := int32(read & 0x7F)
		result |= value << int32(7*byteCount)
	}
	*l = VarLong(result)
	return
}

// TODO: Max length
func (l VarLong) WriteTo(writer io.Writer) (int64, error) {
	buf := bytes.NewBuffer(nil)
	num := int32(l)
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
	nn, err := VarInt(len(s)).WriteTo(writer)
	if err != nil {
		return 0, err
	}

	bytesLen, err := ByteArray(s).WriteTo(writer)
	if err != nil {
		return nn, err
	}

	return nn + bytesLen, nil
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
	return int64(nn), err
}

func (b ByteArray) WriteTo(writer io.Writer) (int64, error) {
	baLen, err := writer.Write(b)
	if err != nil {
		return 0, err
	}
	return int64(baLen), nil
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
