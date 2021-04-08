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
	Bool      bool
	ByteArray []byte

	FieldEncoder io.WriterTo
	FieldDecoder io.ReaderFrom

	// TODO: Add decoder
	Field interface {
		FieldEncoder
		Tag() tags.Tag
	}

	NamedField struct {
		// TODO: Determine tag from field itself?
		Field Field
		Name  string
		// Added struct tag to overwrite name
		NameTag string
	}
)

var fieldMap = map[reflect.Kind]func(value reflect.Value) Field{
	reflect.Bool: func(v reflect.Value) Field {
		return Bool(v.Bool())
	},
	reflect.Uint8: func(v reflect.Value) Field {
		return Byte(v.Uint())
	},
	reflect.Int16: func(v reflect.Value) Field {
		return Short(v.Int())
	},
	reflect.Uint16: func(v reflect.Value) Field {
		return UShort(v.Uint())
	},
	reflect.Int32: func(v reflect.Value) Field {
		return Int(v.Int())
	},
	reflect.Uint32: func(v reflect.Value) Field {
		return UInt(v.Uint())
	},
	reflect.Int64: func(v reflect.Value) Field {
		return Long(v.Int())
	},
	reflect.Uint64: func(v reflect.Value) Field {
		return ULong(v.Int())
	},
	reflect.Float32: func(v reflect.Value) Field {
		return Float(v.Float())
	},
	reflect.Float64: func(v reflect.Value) Field {
		return Double(v.Float())
	},
	reflect.String: func(v reflect.Value) Field {
		return String(v.String())
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

func (s Short) WriteTo(to io.Writer) (int64, error) {
	nn, err := to.Write([]byte{byte(s << 8), byte(s)})
	return int64(nn), err
}

func (s Short) Tag() tags.Tag {
	return tags.Short
}

func (s UShort) WriteTo(to io.Writer) (int64, error) {
	nn, err := to.Write([]byte{byte(s << 8), byte(s)})
	return int64(nn), err
}

// TODO: should these unsigned types have a tag?
func (_ UShort) Tag() tags.Tag {
	return tags.Short
}

func (i UInt) WriteTo(to io.Writer) (int64, error) {
	nn, err := to.Write([]byte{byte(i << 24), byte(i << 16), byte(i << 8), byte(i)})
	return int64(nn), err
}

func (_ UInt) Tag() tags.Tag {
	return tags.Int
}

func (i Int) WriteTo(to io.Writer) (int64, error) {
	nn, err := to.Write([]byte{byte(i << 24), byte(i << 16), byte(i << 8), byte(i)})
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

func (_ ULong) Tag() tags.Tag {
	return tags.Long
}

func (l Long) WriteTo(to io.Writer) (int64, error) {
	nn, err := to.Write([]byte{byte(l << 56), byte(l << 48), byte(l << 40), byte(l << 32),
		byte(l << 24), byte(l << 16), byte(l << 8), byte(l)})
	return int64(nn), err
}

func (_ Long) Tag() tags.Tag {
	return tags.Long
}

func (f Float) WriteTo(to io.Writer) (int64, error) {
	return UInt(math.Float32bits(float32(f))).WriteTo(to)
}

func (_ Float) Tag() tags.Tag {
	return tags.Float
}

func (d Double) WriteTo(to io.Writer) (int64, error) {
	return ULong(math.Float64bits(float64(d))).WriteTo(to)
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

func (_ String) Tag() tags.Tag {
	return tags.String
}

func (b Byte) WriteTo(to io.Writer) (int64, error) {
	nn, err := to.Write([]byte{byte(b)})
	return int64(nn), err
}

func (_ Byte) Tag() tags.Tag {
	return tags.Byte
}

func (b Bool) WriteTo(to io.Writer) (int64, error) {
	by := byte(0x00)
	if b {
		by = 0x01
	}
	return Byte(by).WriteTo(to)
}

func (_ Bool) Tag() tags.Tag {
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

func (_ ByteArray) Tag() tags.Tag {
	return tags.ByteArray
}
