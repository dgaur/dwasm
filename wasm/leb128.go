package wasm

import (
	"io"
)


//
// Consume exactly one byte from the reader.  Do *not* buffer or readahead,
// since the underlying Reader semantics are unknown.  Effectively equivalent
// to adding io.ByteReader interface, when the underlying Reader may not
// support it.
//
func read1(reader io.Reader)(byte, error) {
	buffer := make([]byte, 1)
	_, err := io.ReadFull(reader, buffer)
	return buffer[0], err
}


//
// Parse and return a single uint32 value from a LEB128 sequence
//
// See:
// - https://en.wikipedia.org/wiki/LEB128 and
// - https://en.wikipedia.org/wiki/Variable-length_quantity
//
func readULEB128(reader io.Reader)(uint32, error) {
	var shift uint32
	var value uint32

	for {
		// Consume the next byte
		b, err := read1(reader)
		if (err != nil) {
			return uint32(0xFFFFFFFF), err
		}

		// This byte provides the next 7 bits of the result
		value += ( uint32(b & 0x7F) << shift )
		shift += 7

		// The high-order bit determines whether this is the last byte
		if ( (b & 0x80) == 0 ) {
			break
		}
	}

	return value, nil
}

