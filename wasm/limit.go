package wasm

import (
	"encoding/binary"
	"io"
)

//
// Limit structure for describing resizeable storage (memory, tables, etc)
//
type Limit struct {
	min uint32
	max uint32
}

func readLimit(reader io.Reader)(Limit, error) {
	limit := Limit{}

	var flag uint8
	err := binary.Read(reader, binary.LittleEndian, &flag)
	if (err != nil) {
		return limit, err
	}

	// Min field is always present
	limit.min, err = readULEB128(reader)
	if (err != nil) {
		return limit, err
	}

	// Flag value determines whether max field is present
	if (flag != 0) {
		limit.max, err = readULEB128(reader)
	}

	return limit, err
}
