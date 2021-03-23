package wasm

import(
	"bytes"
	"testing"
	)

// Test decoding of packed ULEB128 values
func TestULEB128Decoding(t *testing.T) {
	testCases := []struct{
		name		string
		encoded		[]byte
		decoded		uint32
	}{
		{ "0x00",		[]byte{ 0x00 },			0x00   },
		{ "0x01",		[]byte{ 0x01 },			0x01   },
		{ "0x7F00",		[]byte{ 0x7F },			0x7F   },
		{ "0xFF 0x00",	[]byte{ 0xFF, 0x00 },	0x7F   },
		{ "0xFF 0x01",	[]byte{ 0xFF, 0x01 },	0xFF   },
		{ "0xFF 0x03",	[]byte{ 0xFF, 0x03 },	0x01FF },
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			reader := bytes.NewReader(test.encoded)
			value, err := readULEB128(reader)
			if (err != nil) {
				t.Error("Unexpected decoding error: ", err)
			}
			if (value != test.decoded) {
				t.Error("Unexpected decoded value: ", value)
			}
		})
	}
}
