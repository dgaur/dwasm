package wasm

import (
	"io"
)

// Read a single name/string from a stream of bytes.  No side effects.
func readName(reader io.Reader) (string, error) {
	// Strings are encoded as a vector of characters
	nameLength, err := readVectorLength(reader)
	if (err != nil) {
		return "", err
	}
	name := make([]byte, nameLength)
	_, err = reader.Read(name)
	if (err != nil) {
		return "", err
	}

	return string(name), nil
}
