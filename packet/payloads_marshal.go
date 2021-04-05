package packet

import (
	"bytes"
	"github.com/google/uuid"
	"io"
	"minecraftServer/nbt"
	"reflect"
)

type (
	encoder struct {
		buf *bytes.Buffer
	}
)

func Marshal(i interface{}) ([]byte, error) {
	encoder := &encoder{
		buf: bytes.NewBuffer(nil),
	}
	return encoder.Encode(i)
}

func MarshalReader(i interface{}) (io.Reader, int, error) {
	encoder := &encoder{
		buf: bytes.NewBuffer(nil),
	}
	return encoder.EncodeReader(i)
}

func (e *encoder) Encode(i interface{}) ([]byte, error) {
	if err := e.EncodeValue(reflect.ValueOf(i)); err != nil {
		return nil, err
	}

	return e.buf.Bytes(), nil
}

func (e *encoder) EncodeReader(i interface{}) (io.Reader, int, error) {
	if err := e.EncodeValue(reflect.ValueOf(i)); err != nil {
		return nil, 0, err
	}

	return e.buf, e.buf.Len(), nil
}

func (e *encoder) EncodeValue(v reflect.Value) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tags := makeTags(typ.Field(i).Tag)

		// Handle optional
		if len(tags.PktOpt) > 0 && !v.FieldByName(tags.PktOpt).Bool() {
			continue
		}

		var encoder FieldEncoder
		switch field.Kind() {
		case reflect.Bool:
			encoder = Boolean(field.Bool())
		case reflect.Int8:
			encoder = Byte(field.Int())
		case reflect.Uint8:
			encoder = UnsignedByte(field.Int())
		case reflect.Int16:
			encoder = Short(field.Int())
		case reflect.Uint16:
			encoder = UnsignedShort(field.Int())
		case reflect.Int32:
			if tags.PktType == "VarInt" {
				encoder = VarInt(field.Int())
			} else {
				encoder = Int(field.Int())
			}
		case reflect.Int64:
			if tags.PktType == "VarLong" {
				encoder = VarLong(field.Int())
			} else {
				encoder = Long(field.Int())
			}
		case reflect.Float32:
			encoder = Float(field.Float())
		case reflect.Float64:
			encoder = Double(field.Float())
		case reflect.String:
			encoder = String(field.String())
		case reflect.Slice:
			sliceType := field.Type().Elem().Kind()
			if sliceType == reflect.Uint8 {
				encoder = ByteArray(field.Bytes())
			} else {
				l := field.Len()

				for i := 0; i < l; i++ {
					if err := e.EncodeValue(field.Index(i)); err != nil {
						return err
					}
				}
				continue
			}
		case reflect.Array:
			l := field.Len()
			sliceType := field.Type().Elem().Kind()

			// UUID
			if l == 16 && sliceType == reflect.Uint8 {
				u := field.Interface().(uuid.UUID)
				encoder = UUID(u)
			}
		default:
			if tags.PktType == "nbt" {
				bs, err := nbt.MarshalValueToNBT(field)
				if err != nil {
					return err
				}
				if _, err := e.buf.Write(bs); err != nil {
					return err
				}
				continue
			}
			if field.CanInterface() {
				if _, err := e.Encode(field.Interface()); err != nil {
					return err
				}
			}
		}
		if encoder != nil {
			_, err := encoder.WriteTo(e.buf)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
