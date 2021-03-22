package main

import "io"


//
// Parse and return a single uint32 value from a LEB128 sequence
//
// See:
// - https://en.wikipedia.org/wiki/LEB128 and
// - https://en.wikipedia.org/wiki/Variable-length_quantity
//
func readULEB128(reader io.ByteReader)(uint32, error) {
	var shift uint32
	var value uint32

	for {
		// Consume the next byte
		b, err := reader.ReadByte()
		if (err != nil) {
			return uint32(0xFFFFFFFF), err
		}

		// This byte provides the next 7 bits of the result
		value += ( uint32(b & 0x7F) << shift )
		shift += 7

		// Last byte?
		if ( (b & 0x80) == 0 ) {
			break
		}
	}

	return value, nil
}

