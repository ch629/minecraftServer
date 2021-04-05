package nbt

import (
	"bytes"
	"github.com/rotisserie/eris"
	"math"
	"reflect"
	"strings"
)

type (
	encoder struct {
		buf *bytes.Buffer
	}

	nbtTags struct {
		Type     string
		Name     string
		Optional string
	}
)

func MarshalToNBT(i interface{}) ([]byte, error) {
	enc := &encoder{
		buf: bytes.NewBuffer(nil),
	}
	if err := enc.Encode(i); err != nil {
		return nil, err
	}
	return enc.buf.Bytes(), nil
}

func MarshalValueToNBT(v reflect.Value) ([]byte, error) {
	enc := &encoder{
		buf: bytes.NewBuffer(nil),
	}
	if err := enc.EncodeValue(v); err != nil {
		return nil, err
	}
	return enc.buf.Bytes(), nil
}

func (e *encoder) Encode(i interface{}) error {
	return e.EncodeValue(reflect.ValueOf(i))
}

// TODO: Optionals -> Should these be pointers so can be nil?
// TODO: List
func (e *encoder) EncodeValue(v reflect.Value) (err error) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tags := makeTags(typ.Field(i).Tag)
		fieldName := typ.Field(i).Name

		if i == 0 {
			if err = e.writeTag(Compound); err != nil {
				return
			}
			if typ.Field(i).Name == "Root" {
				// TODO: Validate it's a string first
				rootValue := field.String()
				if err = e.writeString(rootValue); err != nil {
					return
				}
				continue
			} else {
				if err = e.writeShort(0); err != nil {
					return
				}
			}
		}

		switch field.Kind() {
		case reflect.Bool:
			if err = e.writeNamedTag(Byte, tags.Name, fieldName); err != nil {
				return
			}
			boolByte := byte(0x00)
			if field.Bool() {
				boolByte = 0x01
			}
			if err = e.buf.WriteByte(boolByte); err != nil {
				return
			}
		case reflect.Uint8:
			if err = e.writeNamedTag(Byte, tags.Name, fieldName); err != nil {
				return
			}
			if err = e.buf.WriteByte(uint8(field.Uint())); err != nil {
				return
			}
		case reflect.Int16:
			if err = e.writeNamedTag(Short, tags.Name, fieldName); err != nil {
				return
			}
			if err = e.writeShort(int16(field.Int())); err != nil {
				return
			}
		case reflect.Int32:
			if err = e.writeNamedTag(Int, tags.Name, fieldName); err != nil {
				return
			}
			if err = e.writeInt(int32(field.Int())); err != nil {
				return
			}
		case reflect.Int64:
			if err = e.writeNamedTag(Long, tags.Name, fieldName); err != nil {
				return
			}
			if err = e.writeLong(field.Int()); err != nil {
				return
			}
		case reflect.Float32:
			if err = e.writeNamedTag(Float, tags.Name, fieldName); err != nil {
				return
			}
			if err = e.writeFloat(float32(field.Float())); err != nil {
				return
			}
		case reflect.Float64:
			if err = e.writeNamedTag(Double, tags.Name, fieldName); err != nil {
				return
			}
			if err = e.writeDouble(field.Float()); err != nil {
				return
			}
		case reflect.Slice:
			sliceType := field.Type().Elem().Kind()
			switch sliceType {
			// []byte
			case reflect.Uint8:
				if err = e.writeNamedTag(ByteArray, tags.Name, fieldName); err != nil {
					return
				}
				bs := field.Bytes()
				if err = e.writeInt(int32(len(bs))); err != nil {
					return
				}
				if _, err = e.buf.Write(bs); err != nil {
					return
				}
			// IntArray
			case reflect.Int32:
				if err = e.writeNamedTag(IntArray, tags.Name, fieldName); err != nil {
					return
				}
				intArr := field.Interface().([]int32)
				if err = e.writeInt(int32(len(intArr))); err != nil {
					return
				}
				for _, i := range intArr {
					if err = e.writeInt(i); err != nil {
						return
					}
				}
			// LongArray
			case reflect.Int64:
				if err = e.writeNamedTag(LongArray, tags.Name, fieldName); err != nil {
					return
				}
				longArr := field.Interface().([]int64)
				if err = e.writeInt(int32(len(longArr))); err != nil {
					return
				}
				for _, i := range longArr {
					if err = e.writeLong(i); err != nil {
						return
					}
				}
			// Array of Compounds
			case reflect.Struct:
				for i := 0; i < field.NumField(); i++ {
					if err = e.EncodeValue(field.Index(i)); err != nil {
						return
					}
				}
			default:
				return eris.Errorf("unknown slice type %v", sliceType)
			}
		case reflect.String:
			if err = e.writeNamedTag(String, tags.Name, fieldName); err != nil {
				return
			}
			if err = e.writeString(field.String()); err != nil {
				return
			}
		// Compound
		case reflect.Struct:
			if err = e.EncodeValue(field); err != nil {
				return err
			}
		default:
			return eris.Errorf("unknown type '%v'", field.Kind())
		}
	}
	return e.writeTag(End)
}

func makeTags(tag reflect.StructTag) nbtTags {
	return nbtTags{
		Name:     tag.Get("nbt"),
		Optional: tag.Get("nbt_opt"),
	}
}

func (tags nbtTags) isOptional() bool {
	return tags.Optional == "true"
}

func (e *encoder) writeNamedTag(tag Tag, name string, fieldName string) error {
	if err := e.buf.WriteByte(byte(tag)); err != nil {
		return err
	}
	if len(name) == 0 {
		name = strings.ToLower(fieldName)
	}
	return e.writeString(name)
}

func (e *encoder) writeTag(tag Tag) error {
	return e.buf.WriteByte(byte(tag))
}

func (e *encoder) writeShort(s int16) error {
	_, err := e.buf.Write([]byte{byte(s << 8), byte(s)})
	return err
}

func (e *encoder) writeUnsignedShort(s uint16) error {
	_, err := e.buf.Write([]byte{byte(s << 8), byte(s)})
	return err
}

func (e *encoder) writeInt(i int32) error {
	_, err := e.buf.Write([]byte{byte(i << 24), byte(i << 16), byte(i << 8), byte(i)})
	return err
}

func (e *encoder) writeLong(l int64) error {
	_, err := e.buf.Write([]byte{byte(l << 56), byte(l << 48), byte(l << 40), byte(l << 32),
		byte(l << 24), byte(l << 16), byte(l << 8), byte(l)})
	return err
}

func (e *encoder) writeFloat(f float32) error {
	// TODO: Check if signs matter here
	return e.writeInt(int32(math.Float32bits(f)))
}

func (e *encoder) writeDouble(d float64) error {
	// TODO: Check if signs matter here
	return e.writeLong(int64(math.Float64bits(d)))
}

func (e *encoder) writeString(s string) error {
	// TODO: Docs say unsigned short
	if err := e.writeUnsignedShort(uint16(len(s))); err != nil {
		return err
	}
	_, err := e.buf.Write([]byte(s))
	return err
}
