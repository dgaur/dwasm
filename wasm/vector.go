package wasm

import "io"

func readVectorLength(reader io.Reader) (uint32, error) {
	return(readULEB128(reader))
}

