package packet

import (
	"github.com/rotisserie/eris"
	"io"
	"reflect"
)

type decoder struct {
	reader io.Reader
	bytes  int64
}

func Unmarshal(pkt Packet, i interface{}) error {
	reader, err := pkt.DataReader()
	if err != nil {
		return err
	}
	defer reader.Close()
	decoder := &decoder{
		reader: reader,
		bytes:  int64(pkt.DataLength()),
	}
	return decoder.Decode(i)
}

func (d *decoder) Decode(i interface{}) error {
	v := reflect.ValueOf(i)
	//if v.Kind() != reflect.Ptr {
	//	return eris.New("non-pointer passed to unmarshal")
	//}
	v = v.Elem()
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		typField := typ.Field(i)
		pktTag := typField.Tag.Get("pkt")
		bytesRead := int64(0)
		var err error
		// TODO: Split these out into separate functions
		switch field.Kind() {
		// Boolean
		case reflect.Bool:
			var b Boolean
			bytesRead, err = b.ReadFrom(d.reader)
			if err != nil {
				// TODO: Wrap
				return err
			}
			field.SetBool(bool(b))
		// Byte
		case reflect.Int8:
			var b Byte
			bytesRead, err = b.ReadFrom(d.reader)
			if err != nil {
				// TODO: Wrap
				return err
			}
			field.SetInt(int64(b))
		// UnsignedByte
		case reflect.Uint8:
			var b UnsignedByte
			bytesRead, err = b.ReadFrom(d.reader)
			if err != nil {
				// TODO: Wrap
				return err
			}
			field.SetUint(uint64(b))
		// Short
		case reflect.Int16:
			var s Short
			bytesRead, err = s.ReadFrom(d.reader)
			if err != nil {
				// TODO: Wrap
				return err
			}
			field.SetInt(int64(s))
		// UnsignedShort
		case reflect.Uint16:
			var s UnsignedShort
			bytesRead, err = s.ReadFrom(d.reader)
			if err != nil {
				// TODO: Wrap
				return err
			}
			field.SetUint(uint64(s))
		// Int
		case reflect.Int32:
			if pktTag == "VarInt" {
				var i VarInt
				bytesRead, err = i.ReadFrom(d.reader)
				if err != nil {
					// TODO: Wrap
					return err
				}
				field.SetInt(int64(i))
			} else {
				var i Int
				bytesRead, err = i.ReadFrom(d.reader)
				if err != nil {
					// TODO: Wrap
					return err
				}
				field.SetInt(int64(i))
			}
		// Long
		case reflect.Int64:
			if pktTag == "VarLong" {
				return eris.New("VarLong not implemented")
			} else {
				var l Long
				bytesRead, err = l.ReadFrom(d.reader)
				if err != nil {
					// TODO: Wrap
					return err
				}
				field.SetInt(int64(l))
			}
		// Float
		case reflect.Float32:
			var f Float
			bytesRead, err = f.ReadFrom(d.reader)
			if err != nil {
				// TODO: Wrap
				return err
			}
			field.SetFloat(float64(f))
		// Double
		case reflect.Float64:
			var do Double
			bytesRead, err = do.ReadFrom(d.reader)
			if err != nil {
				// TODO: Wrap
				return err
			}
			field.SetFloat(float64(do))
		// String
		case reflect.String:
			var s String
			bytesRead, err = s.ReadFrom(d.reader)
			if err != nil {
				// TODO: Wrap
				return err
			}
			field.SetString(string(s))
		case reflect.Slice:
		default:
			if field.CanInterface() {
				if err = d.Decode(field.Interface()); err != nil {
					return err
				}
			}
		}
		d.bytes -= bytesRead
	}
	return nil
}
