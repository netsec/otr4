package otr4

import (
	"bytes"
	"encoding/binary"
	"math/big"

	"github.com/twstrike/ed448"
	"golang.org/x/crypto/sha3"
)

func hashToScalar(in []byte) ed448.Scalar {
	hash := make([]byte, fieldBytes)
	sha3.ShakeSum256(hash, in)
	s := ed448.NewScalar(hash)
	return s
}

func appendBytes(bs ...interface{}) []byte {
	var b []byte

	if len(bs) < 2 {
		panic("programmer error: missing append arguments")
	}

	for _, e := range bs {
		switch i := e.(type) {
		case ed448.Point:
			b = append(b, i.Encode()...)
		case ed448.Scalar:
			b = append(b, i.Encode()...)
		case []byte:
			b = append(b, i...)
		default:
			panic("programmer error: invalid input")
		}
	}
	return b
}

func appendAndHash(bs ...interface{}) ed448.Scalar {
	return hashToScalar(appendBytes(bs...))
}

func appendWord32(b []byte, data uint32) []byte {
	return append(b, byte(data>>24), byte(data>>16), byte(data>>8), byte(data))
}

func appendWord64(b []byte, data int64) []byte {
	return append(b, byte(data>>56), byte(data>>48), byte(data>>40), byte(data>>32), byte(data>>24), byte(data>>16), byte(data>>8), byte(data))
}

func appendData(b, data []byte) []byte {
	return append(appendWord32(b, uint32(len(data))), data...)
}

func appendMPI(b []byte, data *big.Int) []byte {
	return appendData(b, data.Bytes())
}

func appendPoint(b []byte, p ed448.Point) []byte {
	return append(b, p.Encode()...)
}

func appendSignature(b []byte, data interface{}) []byte {
	var binBuf bytes.Buffer

	switch d := data.(type) {
	case *signature:
		binary.Write(&binBuf, binary.BigEndian, d)
		return append(b, binBuf.Bytes()...)
	case *dsaSignature:
		binary.Write(&binBuf, binary.BigEndian, d)
		return append(b, binBuf.Bytes()...)
	}
	return nil
}

func extractPoint(b []byte, cursor int) (ed448.Point, int, error) {
	if len(b) < 56 {
		return nil, 0, errInvalidLength
	}

	p := ed448.NewPointFromBytes()
	valid, err := p.Decode(b[cursor:cursor+fieldBytes], false)
	if !valid {
		return nil, 0, err
	}

	cursor += fieldBytes

	return p, cursor, err
}

func fromHexChar(c byte) (byte, bool) {
	switch {
	case '0' <= c && c <= '9':
		return c - '0', true
	case 'a' <= c && c <= 'f':
		return c - 'a' + 10, true
	case 'A' <= c && c <= 'F':
		return c - 'A' + 10, true
	}

	return 0, false
}

func parseToByte(str string) []byte {
	var bs []byte

	for _, s := range str {
		l, valid := fromHexChar(byte(s))
		if !valid {
			return nil
		}
		bs = append(bs, l)
	}

	return bs
}
