package wasm

import(
	"bytes"
	"io"
	"testing"
	)

// Test decoding of packed ULEB128 values
func TestULEB128Decoding(t *testing.T) {
	testCases := []struct{
		name		string
		encoded		[]byte
		decoded		uint32
		status		error
	}{
		{ "0x00",		[]byte{ 0x00 },			0x00,	nil },
		{ "0x01",		[]byte{ 0x01 },			0x01,	nil },
		{ "0x01",		[]byte{ 0x01 },			0x01,	nil },
		{ "0x7F",		[]byte{ 0x7F },			0x7F,	nil },
		{ "0xFF 0x00",	[]byte{ 0xFF, 0x00 },	0x7F,	nil },
		{ "0xFF 0x01",	[]byte{ 0xFF, 0x01 },	0xFF,	nil },
		{ "0xFF 0x03",	[]byte{ 0xFF, 0x03 },	0x01FF,	nil },

		// Weird cases
		{ "0x7F-extra",	[]byte{ 0x7F, 0xAA },	0x7F,	nil },
		{ "0xFF-eof",	[]byte{ 0xFF },			0,		io.EOF },

	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			reader := bytes.NewReader(test.encoded)
			value, err := readULEB128(reader)
			if (err != test.status) {
				t.Error("Unexpected decoding status: ", err)
			}
			if (err == nil && value != test.decoded) {
				t.Error("Unexpected decoded value: ", value)
			}
		})
	}
}
