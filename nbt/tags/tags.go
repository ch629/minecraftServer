package tags

import (
	"io"
	"reflect"
)

type Tag byte

const (
	// End - Signifies the end of a TAG_Compound. It is only ever used inside a TAG_Compound, and is not named despite being in a TAG_Compound
	End Tag = iota
	// Byte - A single signed byte
	Byte
	// Short - A single signed, big endian 16 bit integer
	Short
	// Int - A single signed, big endian 32 bit integer
	Int
	// Long - A single signed, big endian 64 bit integer
	Long
	// Float -	A single, big endian IEEE-754 single-precision floating point number (NaN possible)
	Float
	// Double -	A single, big endian IEEE-754 double-precision floating point number (NaN possible)
	Double
	// ByteArray - A length-prefixed array of signed bytes. The prefix is a signed integer (thus 4 bytes)
	ByteArray
	// String - A length-prefixed modified UTF-8 string. The prefix is an unsigned short (thus 2 bytes) signifying
	// the length of the string in bytes
	String
	// List - A list of nameless tags, all of the same type. The list is prefixed with the Type ID of the items it
	// contains (thus 1 byte), and the length of the list as a signed integer (a further 4 bytes). If the length of
	// the list is 0 or negative, the type may be 0 (TAG_End) but otherwise it must be any other type.
	// (The notchian implementation uses TAG_End in that situation, but another reference implementation by Mojang uses
	// 1 instead; parsers should accept any type if the length is <= 0).
	List
	// Compound - Effectively a list of a named tags. Order is not guaranteed.
	Compound
	// IntArray - A length-prefixed array of signed integers. The prefix is a signed integer (thus 4 bytes) and
	// indicates the number of 4 byte integers.
	IntArray
	// LongArray - A length-prefixed array of signed longs. The prefix is a signed integer (thus 4 bytes) and
	// indicates the number of 8 byte longs.
	LongArray
)

var TagMap = map[reflect.Kind]Tag{
	reflect.Bool:    Byte,
	reflect.Uint8:   Byte,
	reflect.Int16:   Short,
	reflect.Uint16:  Short,
	reflect.Int32:   Int,
	reflect.Uint32:  Int,
	reflect.Int64:   Long,
	reflect.Uint64:  Long,
	reflect.Float32: Float,
	reflect.Float64: Double,
	reflect.String:  String,
	reflect.Struct:  Compound,
}

func (t Tag) String() string {
	return []string{
		"TAG_End",
		"TAG_Byte",
		"TAG_Short",
		"TAG_Int",
		"TAG_Long",
		"TAG_Float",
		"TAG_Double",
		"TAG_Byte_Array",
		"TAG_String",
		"TAG_List",
		"TAG_Compound",
		"TAG_Int_Array",
		"TAG_Long_Array",
	}[t]
}

func (t Tag) WriteTo(writer io.Writer) (int64, error) {
	if _, err := writer.Write([]byte{byte(t)}); err != nil {
		return 0, err
	}
	return 1, nil
}

func (t *Tag) ReadFrom(reader io.Reader) (int64, error) {
	bs := make([]byte, 1)
	if _, err := reader.Read(bs); err != nil {
		return 0, err
	}
	*t = Tag(bs[0])
	return 1, nil
}
