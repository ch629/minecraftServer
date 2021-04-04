package packet

import (
	"bytes"
	"reflect"
)

type (
	encoder struct {
		buf *bytes.Buffer
	}
)

// TODO: []byte or writer?
func Marshal(i interface{}) ([]byte, error) {
	encoder := &encoder{
		buf: bytes.NewBuffer(nil),
	}
	return encoder.Encode(i)
}

func (e *encoder) Encode(i interface{}) ([]byte, error) {
	v := reflect.ValueOf(i)
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
			}
		default:
			// TODO:
		}
		if encoder != nil {
			_, err := encoder.WriteTo(e.buf)
			if err != nil {
				return nil, err
			}
		}
	}
	return e.buf.Bytes(), nil
}
