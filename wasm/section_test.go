package wasm

import(
    "testing"
    )


//
// Test decoding of ExportSection blocks
//
func TestExportSection(t *testing.T) {
    testCases := []struct{
        name        string
        encoded     []byte
        decoded     ExportSection
        status      error
    }{
		// 1 export
        { "export1",
          []byte{ 1,
				  7, 'e', 'x', 'p', 'o', 'r', 't', '1', 0, 0 },
          ExportSection{
			[]Export{
				{ "export1", 0, 0 },
			},
		  },
          nil },

		// 2 exports
        { "export2",
          []byte{ 2,
				  7, 'e', 'x', 'p', 'o', 'r', 't', '1', 0, 0,
                  7, 'e', 'x', 'p', 'o', 'r', 't', '2', 2, 0xF },
          ExportSection{
			[]Export{
				{ "export1", 0, 0 },
				{ "export2", 2, 0xF },
			},
		  },
          nil },
    }

    for _, test := range testCases {
        t.Run(test.name, func(t *testing.T) {
            section, err := readExportSection(test.encoded)
            if (err != test.status) {
                t.Error("Unexpected decoding status: ", err)
            }
            if (err == nil) {
                if (len(section.export) != len(test.decoded.export)) {
                    t.Error("Unexpected decoded length: ", section)
                }

                // Assume each successful decode has at least 1 export
                if (section.export[0] != test.decoded.export[0]) {
                    t.Error("Unexpected decoded export[0]: ", section)
                }
			}
        })
    }
}


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
          MemorySection{
			[]Memory{
				Memory { Limit{ 0x01, 0 } },
			},
		  },
          nil },
        { "decodeF",
          []byte{ 0x01, 0x00, 0x0F },
          MemorySection{
			[]Memory{
				Memory { Limit{ 0x0F, 0 } },
			},
		  },
          nil },

        // Normal decode, min + max
        { "decodeAB",
          []byte{ 0x01, 0x01, 0x0A, 0x0B },
          MemorySection{
			[]Memory{
				Memory { Limit{ 0x0A, 0x0B } },
			},
		  },
          nil },
    }

    for _, test := range testCases {
        t.Run(test.name, func(t *testing.T) {
            section, err := readMemorySection(test.encoded)
            if (err != test.status) {
                t.Error("Unexpected decoding status: ", err)
            }
            if (err == nil) {
                if (len(section.memory) != len(test.decoded.memory)) {
                    t.Error("Unexpected decoded length: ", section)
                }

                // Assume each successful decode has at least 1 memory
                if (section.memory[0] != test.decoded.memory[0]) {
                    t.Error("Unexpected decoded memory[0]: ", section)
                }
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
            section, err := readTableSection(test.encoded)
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


//
// Test decoding of TypeSection blocks
//
func TestTypeSection(t *testing.T) {
    testCases := []struct{
        name        string
        encoded     []byte
        decoded     TypeSection
        status      error
    }{
        // Normal decode, no params or results
        { "decode-no-param-no-result",
          []byte{ 0x01, 0x60, 0x00, 0x00 },
          TypeSection{
            []FunctionType{
                { ResultType{}, ResultType{} },
            },
          },
          nil },

        // Normal decode, multiple params, multiple results
        { "decode-multiple-param-multiple-result",
          []byte{ 0x02,
		          0x60, 0x00, 0x00,
				  0x60, 0x02, 0x7F, 0x7E, 0x02, 0x7D, 0x7C },
          TypeSection{
            []FunctionType{
                { ResultType{},             ResultType{} },
                { ResultType{ 0x7F, 0x7E }, ResultType{ 0x7D, 0x7C } },
            },
          },
          nil },
		  
        // Bad delimiter (0xAA instead of 0x60)
        { "bad-ftype-delimiter-0xAA",
          []byte{ 0x01, 0xAA, 0x00, 0x00 },
          TypeSection{
            []FunctionType{
                { ResultType{}, ResultType{} },
            },
          },
          InvalidSection },
	}

    for _, test := range testCases {
        t.Run(test.name, func(t *testing.T) {
            section, err := readTypeSection(test.encoded)
            if (err != test.status) {
                t.Error("Unexpected decoding status: ", err)
            }
            if (err == nil) {
                if (len(section.ftype) != len(test.decoded.ftype)) {
                    t.Error("Unexpected decoded ftype length: ", section)
                }

                // Assume each successful decode has at least 1 function-type
                if (len(section.ftype[0].parameter) != len(test.decoded.ftype[0].parameter)) {
                    t.Error("Unexpected decoded param length: ", section)
                }
                if (len(section.ftype[0].result) != len(test.decoded.ftype[0].result)) {
                    t.Error("Unexpected decoded result length: ", section)
                }

            }
        })
    }
}
