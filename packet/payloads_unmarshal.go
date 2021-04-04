package packet

import (
	"github.com/rotisserie/eris"
	"io"
	"reflect"
)

var ErrNotPointer = eris.New("non-pointer passed to unmarshal")

const (
	// pkt_type is used to differentiate between int32 and VarInt/int64 and VarLong
	tagPktType = "pkt_type"
	// pkt_len is used for array sizing, references another field if specified, otherwise uses the remaining data length
	tagPktLength = "pkt_len"
	// pkt_opt is used for optional fields, references a boolean field if specified
	tagPktOptional = "pkt_opt"
)

type (
	decoder struct {
		reader io.Reader
		bytes  int64
	}

	pktTags struct {
		PktType string
		PktLen  string
		PktOpt  string
	}
)

func makeTags(tag reflect.StructTag) pktTags {
	return pktTags{
		PktType: tag.Get(tagPktType),
		PktLen:  tag.Get(tagPktLength),
		PktOpt:  tag.Get(tagPktOptional),
	}
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
	if v.Kind() != reflect.Ptr {
		return ErrNotPointer
	}
	v = v.Elem()
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tags := makeTags(typ.Field(i).Tag)

		// Handle optional
		if len(tags.PktOpt) > 0 && !v.FieldByName(tags.PktOpt).Bool() {
			continue
		}

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
		// Int/VarInt
		case reflect.Int32:
			if tags.PktType == "VarInt" {
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
		// Long/VarLong
		case reflect.Int64:
			if tags.PktType == "VarLong" {
				var l VarLong
				bytesRead, err = l.ReadFrom(d.reader)
				if err != nil {
					// TODO: Wrap
					return err
				}
				field.SetInt(int64(l))
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
		// Arrays
		case reflect.Slice:
			bytesRead, err = d.handleSlice(v, field, tags)
			if err != nil {
				return err
			}
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

func (d *decoder) handleSlice(v reflect.Value, field reflect.Value, tags pktTags) (bytesRead int64, err error) {
	sliceType := field.Type().Elem().Kind()
	if sliceType == reflect.Uint8 {
		length := d.bytes
		if len(tags.PktLen) > 0 {
			lenField := v.FieldByName(tags.PktLen)
			length = lenField.Int()
		}

		var reader io.Reader
		reader, err = ByteArrayReader(VarInt(length), d.reader)
		if err != nil {
			return
		}
		var ba ByteArray
		bytesRead, err = ba.ReadFrom(reader)
		if err != nil {
			return
		}
		bs := []byte(ba)
		field.Set(reflect.ValueOf(bs))
		return
	}

	return 0, eris.Errorf("unknown slice type %v", sliceType)
}
