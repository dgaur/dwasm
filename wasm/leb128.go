package wasm

import (
	"encoding/binary"
	"io"
)

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
		var b uint8
		err := binary.Read(reader, binary.LittleEndian, &b)
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

