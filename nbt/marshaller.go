package nbt

import (
	"bytes"
	"github.com/rotisserie/eris"
	"minecraftServer/nbt/tags"
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

// TODO: Optionals -> Should these be pointers so can be nil? -> reflect.Zero & Value.IsZero
// TODO: List
// TODO: Make types for all of these (WriteTo, ReadFrom)
func (e *encoder) EncodeValue(v reflect.Value) (err error) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if _, err = tags.Compound.WriteTo(e.buf); err != nil {
		return
	}

	rootField := v.FieldByName("Root")
	if !rootField.IsValid() || rootField.IsZero() {
		if err = e.writeShort(0); err != nil {
			return
		}
	} else {
		rootName := rootField.String()
		if err = e.writeString(rootName); err != nil {
			return
		}
	}

	return e.EncodeInternalStruct(v)
}

func (e *encoder) EncodeInternalStruct(v reflect.Value) (err error) {
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		typeField := typ.Field(i)
		field := v.Field(i)
		fieldTags := makeTags(typeField.Tag)
		fieldName := typeField.Name

		makeNamedField := func(field Field) NamedField {
			return NamedField{
				Field:   field,
				Name:    fieldName,
				NameTag: fieldTags.Name,
			}
		}

		if i == 0 && typeField.Name == "Root" {
			continue
		}

		var fieldEncoder FieldEncoder
		switch field.Kind() {
		case reflect.Bool:
			fieldEncoder = makeNamedField(Bool(field.Bool()))
		case reflect.Uint8:
			fieldEncoder = makeNamedField(Byte(field.Uint()))
		case reflect.Int16:
			fieldEncoder = makeNamedField(Short(field.Int()))
		case reflect.Int32:
			fieldEncoder = makeNamedField(Int(field.Int()))
		case reflect.Int64:
			fieldEncoder = makeNamedField(Long(field.Int()))
		case reflect.Float32:
			fieldEncoder = makeNamedField(Float(field.Float()))
		case reflect.Float64:
			fieldEncoder = makeNamedField(Double(field.Float()))
		case reflect.Slice:
			sliceType := field.Type().Elem().Kind()
			if fieldTags.Type == "List" {
				if err = e.writeNamedTag(tags.List, fieldTags.Name, fieldName); err != nil {
					return
				}
				switch sliceType {
				// List
				case reflect.Struct:
					if err = e.writeTag(tags.Compound); err != nil {
						return
					}
					if err = e.writeInt(int32(field.Len())); err != nil {
						return
					}
					for i := 0; i < field.Len(); i++ {
						if err = e.writeTag(tags.Compound); err != nil {
							return
						}
						if err = e.EncodeInternalStruct(field.Index(i)); err != nil {
							return
						}
					}
				}
			} else {
				switch sliceType {
				// []byte
				case reflect.Uint8:
					fieldEncoder = makeNamedField(ByteArray(field.Bytes()))
				// IntArray
				case reflect.Int32:
					if err = e.writeNamedTag(tags.IntArray, fieldTags.Name, fieldName); err != nil {
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
					if err = e.writeNamedTag(tags.LongArray, fieldTags.Name, fieldName); err != nil {
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
				default:
					return eris.Errorf("unknown slice type %v", sliceType)
				}
			}
		case reflect.String:
			fieldEncoder = makeNamedField(String(field.String()))
		// Compound
		case reflect.Struct:
			if err = e.writeNamedTag(tags.Compound, fieldTags.Name, fieldName); err != nil {
				return
			}
			if err = e.EncodeValue(field); err != nil {
				return err
			}
		default:
			return eris.Errorf("unknown type '%v'", field.Kind())
		}
		if fieldEncoder != nil {
			if _, err = fieldEncoder.WriteTo(e.buf); err != nil {
				return
			}
		}
	}
	return e.writeTag(tags.End)
}

func makeTags(tag reflect.StructTag) nbtTags {
	return nbtTags{
		Name:     tag.Get("nbt"),
		Optional: tag.Get("nbt_opt"),
		Type:     tag.Get("nbt_type"),
	}
}

func (tags nbtTags) isOptional() bool {
	return tags.Optional == "true"
}

func (e *encoder) writeNamedTag(tag tags.Tag, name string, fieldName string) error {
	if err := e.buf.WriteByte(byte(tag)); err != nil {
		return err
	}
	if len(name) == 0 {
		name = strings.ToLower(fieldName)
	}
	return e.writeString(name)
}

func (e *encoder) writeTag(tag tags.Tag) error {
	_, err := tag.WriteTo(e.buf)
	return err
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

func (e *encoder) writeString(s string) error {
	if err := e.writeUnsignedShort(uint16(len(s))); err != nil {
		return err
	}
	_, err := e.buf.Write([]byte(s))
	return err
}
