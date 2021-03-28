package wasm

import(
    "bytes"
    "testing"
    )


//
// Test decoding of Module blocks
//
func TestModule(t *testing.T) {
    testCases := []struct{
        name        string
        encoded     []byte
        decoded     Module
        status      error
    }{
        { "empty-module",
          []byte{ 0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00 },
          Module{ []Section{ nil } },
          nil },

        { "single-custom-section",
          []byte{ 0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00, 0x00, 0x05, 0x04, 't', 'e', 's', 't' },
          Module{ []Section{ CustomSection{ []byte("test"), "test" } } },
          nil },

        { "single-unknown-section",
          []byte{ 0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00, 0xEF, 0x01, 0xAB },
          Module{ []Section{ UnknownSection{ []byte{0xAB}, uint8(0xEF) } } },
          nil },

        { "bad-preamble",
          []byte{ 0x00 },
          Module{ []Section{} },
          InvalidModule },
    }

    for _, test := range testCases {
        t.Run(test.name, func(t *testing.T) {
            reader := bytes.NewReader(test.encoded)
            module, err := ReadModule(reader)
            if (err != test.status) {
                t.Error("Unexpected decoding status: ", err)
            }
            if (err == nil && module.section[0] != nil) {
                //@deeper comparison would be more meaningful here
                if (module.section[0].id() != test.decoded.section[0].id()) {
                    t.Error("Unexpected section[0]: ", module)
                }
            }
        })
    }
}

