package nbt

import (
	"compress/gzip"
	"compress/zlib"
	"io"
)

type (
	CompressionType byte
)

const (
	GzipCompressed = byte(0x1F)
	ZLibCompressed = byte(0x78)
)

func CompressWrapReader(compression byte, reader io.Reader) (io.ReadCloser, error) {
	switch compression {
	case GzipCompressed:
		return gzip.NewReader(reader)
	case ZLibCompressed:
		return zlib.NewReader(reader)
	}
	return io.NopCloser(reader), nil
}

func CompressWrapWriter(compression byte, writer io.Writer) io.WriteCloser {
	switch compression {
	case GzipCompressed:
		return gzip.NewWriter(writer)
	case ZLibCompressed:
		return zlib.NewWriter(writer)
	}
	return nopWriterCloser(writer)
}

func nopWriterCloser(writer io.Writer) io.WriteCloser {
	return &nopCloser{writer}
}

type nopCloser struct {
	io.Writer
}

func (n nopCloser) Close() error { return nil }
