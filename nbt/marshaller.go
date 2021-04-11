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

// TODO: Reader output instead?
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

func (e *encoder) EncodeValue(v reflect.Value) (err error) {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	rootField := v.FieldByName("Root")
	rootName := ""
	if rootField.IsValid() && rootField.Kind() == reflect.String {
		rootName = rootField.String()
	}
	namedTag := NamedTag{
		Tag:  tags.Compound,
		Name: rootName,
	}
	if _, err = namedTag.WriteTo(e.buf); err != nil {
		return
	}

	return e.EncodeInternalStruct(v)
}

func (e *encoder) EncodeInternalStruct(v reflect.Value) (err error) {
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if i == 0 && typ.Field(i).Name == "Root" {
			continue
		}
		if err = e.EncodeField(typ, i, v.Field(i)); err != nil {
			return
		}

	}
	return e.writeTag(tags.End)
}

func (e *encoder) EncodeField(typ reflect.Type, i int, field reflect.Value) (err error) {
	typeField := typ.Field(i)
	fieldTags := makeTags(typeField.Tag)

	if fieldTags.isOptional() {
		// TODO: This makes any optional field missing if it's 0 or false, will need to use pointers if this is not expected
		if field.IsZero() {
			return
		}
		if field.Kind() == reflect.Ptr {
			// Pull out from pointer
			field = field.Elem()
		}
	}

	fieldName := strings.ToLower(typeField.Name)

	makeNamedField := func(field Field) NamedField {
		return NamedField{
			Field:   field,
			Name:    fieldName,
			NameTag: fieldTags.Name,
		}
	}

	makeNamedTag := func(tag tags.Tag) NamedTag {
		return NamedTag{
			Tag:     tag,
			Name:    fieldName,
			NameTag: fieldTags.Name,
		}
	}

	var fieldEncoder FieldEncoder
	if fieldFunc, ok := fieldMap[field.Kind()]; ok {
		fieldEncoder = NamedField{
			Field:   fieldFunc(field),
			Name:    fieldName,
			NameTag: fieldTags.Name,
		}
	} else {
		switch field.Kind() {
		case reflect.Slice:
			sliceType := field.Type().Elem().Kind()
			if isList(sliceType, fieldTags) {
				l := field.Len()
				if _, err = writeAll(e.buf, makeNamedTag(tags.List), tags.TagMap[sliceType], Int(l)); err != nil {
					return
				}

				if fieldFunc, ok := fieldMap[sliceType]; ok {
					for i := 0; i < l; i++ {
						// TODO: Try to make this recursive if possible?
						if _, err = writeAll(e.buf, fieldFunc(field.Index(i))); err != nil {
							return
						}
					}
				} else {
					switch sliceType {
					// Compound List
					case reflect.Struct:
						for i := 0; i < field.Len(); i++ {
							// Compound lists just have an END tag between each
							if err = e.EncodeInternalStruct(field.Index(i)); err != nil {
								return
							}
						}
					}
				}
			} else {
				switch sliceType {
				// []byte
				case reflect.Uint8:
					ba := ByteArray(field.Bytes())
					fieldEncoder = makeNamedField(&ba)
				// IntArray
				case reflect.Int32:
					intArr := field.Interface().([]int32)
					if _, err = writeAll(e.buf, makeNamedTag(tags.IntArray), Int(len(intArr))); err != nil {
						return
					}
					for _, i := range intArr {
						if _, err = Int(i).WriteTo(e.buf); err != nil {
							return
						}
					}
				// LongArray
				case reflect.Int64:
					longArr := field.Interface().([]int64)
					if _, err = writeAll(e.buf, makeNamedTag(tags.LongArray), Int(len(longArr))); err != nil {
						return
					}
					for _, l := range longArr {
						if _, err = Long(l).WriteTo(e.buf); err != nil {
							return
						}
					}
				default:
					return eris.Errorf("unknown slice type %v", sliceType)
				}
			}
		// Compound
		case reflect.Struct:
			if _, err = makeNamedTag(tags.Compound).WriteTo(e.buf); err != nil {
				return
			}
			if err = e.EncodeValue(field); err != nil {
				return err
			}
		default:
			return eris.Errorf("unknown type '%v'", field.Kind())
		}
	}
	if fieldEncoder != nil {
		if _, err = fieldEncoder.WriteTo(e.buf); err != nil {
			return
		}
	}
	return
}

func isList(kind reflect.Kind, nbtTags nbtTags) bool {
	return strings.ToLower(nbtTags.Type) == "list" || (kind != reflect.Uint8 && kind != reflect.Int32 && kind != reflect.Int64)
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

func (e *encoder) writeTag(tag tags.Tag) error {
	_, err := tag.WriteTo(e.buf)
	return err
}
