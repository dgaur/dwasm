package wasm

import(
    "testing"
    )

//
// Test decoding of MemorySection blocks
//
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
        { "decodeABInvalid",
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


//
// Test decoding of TableSection blocks
//
func TestTableSection(t *testing.T) {
    testCases := []struct{
        name        string
        encoded     []byte
        decoded     TableSection
        status      error
    }{
        // Normal decode, min only
        { "decode1",
          []byte{ 0x01, 0x70, 0x00, 0x01 },
          TableSection{
            []Table{
                { Limit{ 0x01, 0 }, 0x70 },
            },
          },
          nil },

        // Normal decode, min + max
        { "decodeAB",
          []byte{ 0x01, 0x70, 0x01, 0x0A, 0x0B },
          TableSection{
            []Table{
                { Limit{ 0x0A, 0x0B }, 0x70 },
            },
          },
          nil },

        // Normal decode, min only, 2 tables
        { "decode2",
          []byte{ 0x02, 0x70, 0x00, 0x01, 0x6F, 0x00, 0x02 },
          TableSection{
            []Table{
                { Limit{ 0x01, 0 }, 0x70 },
                { Limit{ 0x02, 0 }, 0x6F },
            },
          },
          nil },
    }

    for _, test := range testCases {
        t.Run(test.name, func(t *testing.T) {
            section := TableSection{}
            err := section.read(test.encoded)
            if (err != test.status) {
                t.Error("Unexpected decoding status: ", err)
            }
            if (err == nil) {
                if (len(section.table) != len(test.decoded.table)) {
                    t.Error("Unexpected decoded length: ", section)
                }

                // Assume each successful decode has at least 1 table
                if (section.table[0] != test.decoded.table[0]) {
                    t.Error("Unexpected decoded table[0]: ", section)
                }
            }
        })
    }
}
