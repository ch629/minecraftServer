package nbt

import (
	"fmt"
	"github.com/ch629/minecraftServer/nbt/tags"
	"github.com/rotisserie/eris"
	"io"
)

// TODO: Read compound as a map & then read through the interface to grab names & set from map, throw err if not optional & missing

type (
	nbtCompound struct {
		Name   string
		Values map[string]interface{}
	}

	decoder struct {
		reader    io.Reader
		bytesRead int64
	}
)

func Unmarshal(reader io.Reader, i interface{}) error {
	dec := decoder{
		reader: reader,
	}
	return dec.Decode(i)
}

func (d *decoder) Decode(i interface{}) error {
	return nil
}

func (d *decoder) ReadAsMap() (*nbtCompound, error) {
	var tag tags.Tag
	var tagName String
	count, err := readAll(d.reader, &tag, &tagName)
	if err != nil {
		return nil, err
	}
	d.bytesRead += count
	if tag != tags.Compound {
		return nil, eris.New("expected NBT to start with a compound")
	}
	baseCompound := &nbtCompound{
		Name:   string(tagName),
		Values: make(map[string]interface{}),
	}

loop:
	for {
		tag, err = d.ReadTag()
		if err != nil {
			// TODO: EOF
			return nil, err
		}

		switch tag {
		case tags.End:
			break loop
		case tags.Byte:
			var name String
			var b Byte
			count, err := readAll(d.reader, &name, &b)
			fmt.Println(count, err)
		}
	}

	// TODO: Validate last tag is end
	return baseCompound, nil
}

func (d *decoder) ReadTag() (tags.Tag, error) {
	var tag tags.Tag
	_, err := tag.ReadFrom(d.reader)
	return tag, err
}
