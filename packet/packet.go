package packet

import (
	"bytes"
	"compress/zlib"
	"github.com/rotisserie/eris"
	"io"
)

type (
	// TODO: These are exactly the same, do we just hold a bool if it's compressed?
	UncompressedPacket struct {
		packetID   VarInt
		dataLength VarInt
		readCloser io.ReadCloser
	}

	CompressedPacket struct {
		packetID   VarInt
		dataLength VarInt
		readCloser io.ReadCloser
	}

	Packet interface {
		ID() VarInt
		DataLength() VarInt
		DataReader() (io.ReadCloser, error)
	}
)

func (p UncompressedPacket) ID() VarInt {
	return p.packetID
}

func (p UncompressedPacket) DataLength() VarInt {
	return p.dataLength
}

func (p UncompressedPacket) DataReader() (io.ReadCloser, error) {
	return p.readCloser, nil
}

func (p CompressedPacket) ID() VarInt {
	return p.packetID
}

func (p CompressedPacket) DataReader() (io.ReadCloser, error) {
	return p.readCloser, nil
}

func (p CompressedPacket) DataLength() VarInt {
	return p.dataLength
}

func MakeUncompressedPacket(reader io.Reader) (Packet, error) {
	// Length of the PacketID + Data
	var pktLen VarInt
	var pktId VarInt

	lens, err := ReadFieldsWithLengths(reader, &pktLen, &pktId)
	if err != nil {
		return nil, err
	}

	pktIdLen := lens[1]

	// Length is ID + data
	dataLen := pktLen - VarInt(pktIdLen)
	var byteData ByteArray
	baReader, err := ByteArrayReader(dataLen, reader)
	if err != nil {
		return nil, eris.Wrap(err, "failed to create byte array reader")
	}
	_, err = byteData.ReadFrom(baReader)

	if err != nil {
		return nil, eris.Wrap(err, "failed to read packet byte data")
	}

	return &UncompressedPacket{
		packetID:   pktId,
		dataLength: dataLen,
		readCloser: io.NopCloser(bytes.NewReader(byteData)),
	}, nil
}

func MakeCompressedPacket(reader io.Reader) (Packet, error) {
	var pktLen VarInt
	var dataLen VarInt
	var compressedData ByteArray
	err := ReadFields(reader, &pktLen, &dataLen, &compressedData)
	if err != nil {
		return nil, err
	}

	// zlReader needs closing within the packet code
	zlReader, err := zlib.NewReader(bytes.NewReader(compressedData))
	if err != nil {
		return nil, err
	}

	var pktId VarInt
	if _, err = pktId.ReadFrom(zlReader); err != nil {
		return nil, err
	}

	return &CompressedPacket{
		packetID:   pktId,
		dataLength: dataLen,
		readCloser: zlReader,
	}, nil
}

func MakePacket(id VarInt, dataLength VarInt, payload io.Reader) Packet {
	return &UncompressedPacket{
		packetID:   id,
		dataLength: dataLength,
		readCloser: io.NopCloser(payload),
	}
}

func WriteTo(pkt Packet, writer io.Writer) (count int64, err error) {
	buf := bytes.NewBuffer(nil)
	pkt.ID().WriteTo(buf)
	reader, err := pkt.DataReader()
	if err != nil {
		return
	}
	count, err = io.Copy(buf, reader)
	if err != nil {
		return
	}
	VarInt(buf.Len()).WriteTo(writer)
	nn, err := io.Copy(writer, buf)
	if err != nil {
		return
	}
	count += nn
	return
}
