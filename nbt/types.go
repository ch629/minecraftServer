package nbt

import (
	"io"
	"math"
	"minecraftServer/nbt/tags"
	"reflect"
)

type (
	Short     int16
	UShort    uint16
	UInt      uint32
	Int       int32
	ULong     uint64
	Long      int64
	Float     float32
	Double    float64
	String    string
	Byte      uint8
	Boolean   bool
	ByteArray []byte

	FieldEncoder io.WriterTo
	FieldDecoder io.ReaderFrom

	Field interface {
		FieldEncoder
		FieldDecoder
		Tag() tags.Tag
	}

	NamedField struct {
		Field Field
		Name  string
		// Added struct tag to overwrite name
		NameTag string
	}

	NamedTag struct {
		Tag  tags.Tag
		Name string
		// Added struct tag to overwrite name
		NameTag string
	}
)

var fieldMap = map[reflect.Kind]func(value reflect.Value) Field{
	reflect.Bool: func(v reflect.Value) Field {
		b := Boolean(v.Bool())
		return &b
	},
	reflect.Uint8: func(v reflect.Value) Field {
		b := Byte(v.Uint())
		return &b
	},
	reflect.Int16: func(v reflect.Value) Field {
		s := Short(v.Int())
		return &s
	},
	reflect.Uint16: func(v reflect.Value) Field {
		s := UShort(v.Uint())
		return &s
	},
	reflect.Int32: func(v reflect.Value) Field {
		i := Int(v.Int())
		return &i
	},
	reflect.Uint32: func(v reflect.Value) Field {
		i := UInt(v.Uint())
		return &i
	},
	reflect.Int64: func(v reflect.Value) Field {
		l := Long(v.Int())
		return &l
	},
	reflect.Uint64: func(v reflect.Value) Field {
		l := ULong(v.Int())
		return &l
	},
	reflect.Float32: func(v reflect.Value) Field {
		f := Float(v.Float())
		return &f
	},
	reflect.Float64: func(v reflect.Value) Field {
		d := Double(v.Float())
		return &d
	},
	reflect.String: func(v reflect.Value) Field {
		s := String(v.String())
		return &s
	},
}

func writeAll(writer io.Writer, encoders ...FieldEncoder) (count int64, err error) {
	var nn int64
	for _, encoder := range encoders {
		nn, err = encoder.WriteTo(writer)
		if err != nil {
			// TODO: Give context to err
			return
		}
		count += nn
	}
	return
}

func (field NamedField) WriteTo(to io.Writer) (int64, error) {
	name := field.Name
	if len(field.NameTag) > 0 {
		name = field.NameTag
	}
	return writeAll(to, field.Field.Tag(), String(name), field.Field)
}

func (namedTag NamedTag) WriteTo(to io.Writer) (int64, error) {
	name := namedTag.Name
	if len(namedTag.NameTag) > 0 {
		name = namedTag.NameTag
	}
	return writeAll(to, namedTag.Tag, String(name))
}

func (s Short) WriteTo(to io.Writer) (int64, error) {
	nn, err := to.Write([]byte{byte(s << 8), byte(s)})
	return int64(nn), err
}

func (s *Short) ReadFrom(from io.Reader) (int64, error) {
	by := make([]byte, 2)
	nn, err := io.ReadFull(from, by)
	if err != nil {
		return 0, err
	}
	*s = Short(int16(by[0])>>8 | int16(by[1]))
	return int64(nn), nil
}

func (s Short) Tag() tags.Tag {
	return tags.Short
}

func (s UShort) WriteTo(to io.Writer) (int64, error) {
	nn, err := to.Write([]byte{byte(s << 8), byte(s)})
	return int64(nn), err
}

func (s *UShort) ReadFrom(from io.Reader) (int64, error) {
	by := make([]byte, 2)
	nn, err := io.ReadFull(from, by)
	if err != nil {
		return 0, err
	}
	*s = UShort(uint16(by[0])>>8 | uint16(by[1]))
	return int64(nn), nil
}

// TODO: should these unsigned types have a tag?
func (_ UShort) Tag() tags.Tag {
	return tags.Short
}

func (i UInt) WriteTo(to io.Writer) (int64, error) {
	nn, err := to.Write([]byte{byte(i << 24), byte(i << 16), byte(i << 8), byte(i)})
	return int64(nn), err
}

func (i *UInt) ReadFrom(from io.Reader) (int64, error) {
	by := make([]byte, 4)
	nn, err := io.ReadFull(from, by)
	if err != nil {
		return 0, err
	}
	*i = UInt(uint32(by[0])>>24 | uint32(by[1])>>16 | uint32(by[2])>>8 | uint32(by[3]))
	return int64(nn), err
}

func (_ UInt) Tag() tags.Tag {
	return tags.Int
}

func (i Int) WriteTo(to io.Writer) (int64, error) {
	nn, err := to.Write([]byte{byte(i << 24), byte(i << 16), byte(i << 8), byte(i)})
	return int64(nn), err
}

func (i *Int) ReadFrom(from io.Reader) (int64, error) {
	by := make([]byte, 4)
	nn, err := io.ReadFull(from, by)
	if err != nil {
		return 0, err
	}
	*i = Int(int32(by[0])>>24 | int32(by[1])>>16 | int32(by[2])>>8 | int32(by[3]))
	return int64(nn), err
}

func (_ Int) Tag() tags.Tag {
	return tags.Int
}

func (l ULong) WriteTo(to io.Writer) (int64, error) {
	nn, err := to.Write([]byte{byte(l << 56), byte(l << 48), byte(l << 40), byte(l << 32),
		byte(l << 24), byte(l << 16), byte(l << 8), byte(l)})
	return int64(nn), err
}

func (l *ULong) ReadFrom(from io.Reader) (int64, error) {
	by := make([]byte, 8)
	nn, err := io.ReadFull(from, by)
	if err != nil {
		return 0, err
	}
	*l = ULong(uint64(by[0])>>56 | uint64(by[1])>>48 | uint64(by[2])>>40 | uint64(by[3])>>32 |
		uint64(by[4])>>24 | uint64(by[5])>>16 | uint64(by[6])>>8 | uint64(by[7]))
	return int64(nn), nil
}

func (_ ULong) Tag() tags.Tag {
	return tags.Long
}

func (l Long) WriteTo(to io.Writer) (int64, error) {
	nn, err := to.Write([]byte{byte(l << 56), byte(l << 48), byte(l << 40), byte(l << 32),
		byte(l << 24), byte(l << 16), byte(l << 8), byte(l)})
	return int64(nn), err
}

func (l *Long) ReadFrom(from io.Reader) (int64, error) {
	by := make([]byte, 8)
	nn, err := io.ReadFull(from, by)
	if err != nil {
		return 0, err
	}
	*l = Long(int64(by[0])>>56 | int64(by[1])>>48 | int64(by[2])>>40 | int64(by[3])>>32 |
		int64(by[4])>>24 | int64(by[5])>>16 | int64(by[6])>>8 | int64(by[7]))
	return int64(nn), nil
}

func (_ Long) Tag() tags.Tag {
	return tags.Long
}

func (f Float) WriteTo(to io.Writer) (int64, error) {
	return UInt(math.Float32bits(float32(f))).WriteTo(to)
}

func (f *Float) ReadFrom(from io.Reader) (int64, error) {
	var i Int
	nn, err := i.ReadFrom(from)
	if err != nil {
		return 0, err
	}
	*f = Float(math.Float32frombits(uint32(i)))
	return nn, nil
}

func (_ Float) Tag() tags.Tag {
	return tags.Float
}

func (d Double) WriteTo(to io.Writer) (int64, error) {
	return ULong(math.Float64bits(float64(d))).WriteTo(to)
}

func (d *Double) ReadFrom(from io.Reader) (int64, error) {
	var l Long
	nn, err := l.ReadFrom(from)
	if err != nil {
		return 0, err
	}
	*d = Double(math.Float64frombits(uint64(l)))
	return nn, nil
}

func (_ Double) Tag() tags.Tag {
	return tags.Double
}

func (s String) WriteTo(to io.Writer) (byteCount int64, err error) {
	byteCount, err = UShort(len(s)).WriteTo(to)
	if err != nil {
		return
	}

	nn, err := to.Write([]byte(s))
	byteCount += int64(nn)
	return
}

func (s *String) ReadFrom(from io.Reader) (int64, error) {
	var l UShort
	nn, err := l.ReadFrom(from)
	if err != nil {
		return 0, err
	}
	by := make([]byte, l)
	count, err := io.ReadFull(from, by)
	if err != nil {
		return nn, err
	}
	*s = String(by)
	return nn + int64(count), nil
}

func (_ String) Tag() tags.Tag {
	return tags.String
}

func (b Byte) WriteTo(to io.Writer) (int64, error) {
	nn, err := to.Write([]byte{byte(b)})
	return int64(nn), err
}

func (b *Byte) ReadFrom(from io.Reader) (int64, error) {
	by := make([]byte, 1)
	nn, err := io.ReadFull(from, by)
	if err != nil {
		return 0, err
	}
	*b = Byte(by[0])
	return int64(nn), nil
}

func (_ Byte) Tag() tags.Tag {
	return tags.Byte
}

func (b Boolean) WriteTo(to io.Writer) (int64, error) {
	by := byte(0x00)
	if b {
		by = 0x01
	}
	return Byte(by).WriteTo(to)
}

func (b *Boolean) ReadFrom(from io.Reader) (int64, error) {
	var by Byte
	nn, err := by.ReadFrom(from)
	if err != nil {
		return 0, err
	}
	*b = by == 0x01
	return nn, nil
}

func (_ Boolean) Tag() tags.Tag {
	return tags.Byte
}

func (ba ByteArray) WriteTo(to io.Writer) (int64, error) {
	count, err := Int(len(ba)).WriteTo(to)
	if err != nil {
		return 0, err
	}
	nn, err := to.Write(ba)
	return count + int64(nn), err
}

func (ba *ByteArray) ReadFrom(from io.Reader) (int64, error) {
	var l Int
	nn, err := l.ReadFrom(from)
	if err != nil {
		return 0, err
	}
	by := make([]byte, l)
	count, err := io.ReadFull(from, by)
	if err != nil {
		return nn, err
	}
	*ba = by
	return nn + int64(count), nil
}

func (_ ByteArray) Tag() tags.Tag {
	return tags.ByteArray
}
