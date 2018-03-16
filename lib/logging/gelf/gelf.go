package gelf

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"crypto/rand"
	"encoding/binary"
	"io"
	"math"
	"net"
	"strconv"
)

// ZapWriteSyncer mirrors zap.WriteSyncer
type ZapWriteSyncer interface {
	io.Writer
	Sync() error
}

// New returns an implementation of ZapWriteSyncer which should be compatible with zap.WriteSyncer
func New(config Config) ZapWriteSyncer {
	return &gelf{Config: config}
}

// Gelf is an operational structure which holds the implementation of WriteSyncer
type gelf struct {
	Config
}

// Sync implements the WriteSyncer Sync method, nothing to do here so just noop
func (g *gelf) Sync() error {
	return nil
}

// Write implements a Writer Write method
func (g *gelf) Write(p []byte) (int, error) {
	var (
		buf bytes.Buffer
		err error
	)

	switch g.Config.Compression {
	case CompressionGZip:
		buf, err = g.compressGZip(p)

	case CompressionZLib:
		buf, err = g.compressZLib(p)

	default:
		buf, err = g.compressNone(p)
	}

	if err != nil {
		return 0, err
	}

	chunksize := g.Config.MaxChunkSize
	length := buf.Len()

	if length > chunksize {
		chunkCountInt := int(math.Ceil(float64(length) / float64(chunksize)))

		id := make([]byte, 8)
		rand.Read(id)

		for i, index := 0, 0; i < length; i, index = i+chunksize, index+1 {
			packet := g.createChunkedMessage(index, chunkCountInt, id, &buf)
			_, err := g.send(packet.Bytes())
			if err != nil {
				return 0, err
			}
		}
	} else {
		_, err := g.send(buf.Bytes())
		if err != nil {
			return 0, err
		}
	}

	return len(p), nil
}

// createChunkedMessage creates UDP packets for transmission using the chunk size
func (g *gelf) createChunkedMessage(index int, chunkCountInt int, id []byte, compressed *bytes.Buffer) bytes.Buffer {
	var packet bytes.Buffer

	chunksize := g.Config.MaxChunkSize

	packet.Write(g.intToBytes(30))
	packet.Write(g.intToBytes(15))
	packet.Write(id)

	packet.Write(g.intToBytes(index))
	packet.Write(g.intToBytes(chunkCountInt))

	packet.Write(compressed.Next(chunksize))

	return packet
}

// intToBytes writes numbers into a byte message using LittleEndian
func (g *gelf) intToBytes(i int) []byte {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, int8(i))

	return buf.Bytes()
}

// compressNone just returns the message as a bytes.Buffer
func (g *gelf) compressNone(b []byte) (bytes.Buffer, error) {
	var buf bytes.Buffer

	_, err := buf.Write(b)
	if err != nil {
		return buf, err
	}

	return buf, nil
}

// compressGZip squashes the buffer using gzip compression
func (g *gelf) compressGZip(b []byte) (bytes.Buffer, error) {
	var buf bytes.Buffer
	comp := gzip.NewWriter(&buf)

	_, err := comp.Write(b)
	if err != nil {
		return buf, err
	}

	if err := comp.Close(); err != nil {
		return buf, err
	}

	return buf, nil
}

// compressZLib squashes the buffer using zlib compression
func (g *gelf) compressZLib(b []byte) (bytes.Buffer, error) {
	var buf bytes.Buffer
	comp := zlib.NewWriter(&buf)

	_, err := comp.Write(b)
	if err != nil {
		return buf, err
	}

	if err := comp.Close(); err != nil {
		return buf, err
	}

	return buf, nil
}

// end connects to the configured GELF server and transmits a packet
func (g *gelf) send(b []byte) (int, error) {
	var addr = g.Config.Host + ":" + strconv.Itoa(g.Config.Port)

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return 0, err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return 0, err
	}

	return conn.Write(b)
}

// EOF
