package wasm

import(
    "testing"
    )

// Test decoding of packed ULEB128 values
func TestMemorySection(t *testing.T) {
    testCases := []struct{
        name        string
        encoded     []byte
        decoded     MemorySection
        status      error
    }{
        // Normal decode, min only
        { "decode1",
          []byte{ 0x01, 0x00, 0x01 },
          MemorySection{ Limit{ 0x01, 0 } },
          nil },
        { "decodeF",
          []byte{ 0x01, 0x00, 0x0F },
          MemorySection{ Limit{ 0x0F, 0 } },
          nil },

        // Normal decode, min + max
        { "decodeAB",
          []byte{ 0x01, 0x01, 0x0A, 0x0B },
          MemorySection{ Limit{ 0x0A, 0x0B } },
          nil },

        // Invalid decode, only 1 Memory allowed
        { "decodeAB",
          []byte{ 0x02, 0x00, 0x01, 0x00, 0x02 },
          MemorySection{ Limit{ 0x01, 0 } },
          InvalidSection },
    }

    for _, test := range testCases {
        t.Run(test.name, func(t *testing.T) {
            mem := MemorySection{}
            err := mem.read(test.encoded)
            if (err != test.status) {
                t.Error("Unexpected decoding status: ", err)
            }
            if (err == nil && mem != test.decoded) {
                t.Error("Unexpected decoded memory: ", mem)
            }
        })
    }
}

