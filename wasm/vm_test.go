package wasm

import(
	"bytes"
	"testing"
    )


//
// Test VM cases
//
func TestVMExecution(t *testing.T) {
    testCases := []struct{
        name        string
        module      []byte
		startfn     string
        status      error
    }{
		// Empty/null module, see samples/empty.wat
        { "empty-module",
          []byte{ 0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00 },
		  "InvalidFunction",
          MissingFunction },

		// single "nop" function, see samples/fnop.wat
        { "fnop-invalid-entry",
          []byte{ 0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00,
				  0x01, 0x04, 0x01, 0x60, 0x00, 0x00, 0x03, 0x02,
				  0x01, 0x00, 0x07, 0x08, 0x01, 0x04, 0x66, 0x6e,
				  0x6f, 0x70, 0x00, 0x00, 0x0a, 0x06, 0x01, 0x04,
				  0x00, 0x01, 0x01, 0x0b },
		  "InvalidFunction",
          MissingFunction },

		// single "nop" function, see samples/fnop.wat
        { "fnop",
          []byte{ 0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00,
				  0x01, 0x04, 0x01, 0x60, 0x00, 0x00, 0x03, 0x02,
				  0x01, 0x00, 0x07, 0x08, 0x01, 0x04, 0x66, 0x6e,
				  0x6f, 0x70, 0x00, 0x00, 0x0a, 0x06, 0x01, 0x04,
				  0x00, 0x01, 0x01, 0x0b },
		  "fnop",
          nil },
	}

    for _, test := range testCases {
        t.Run(test.name, func(t *testing.T) {
			// Initialize module for the next test case
            reader := bytes.NewReader(test.module)
            module, err := ReadModule(reader)
            if (err != nil) {
                t.Error("Unexpected decoding status: ", err)
            }

			// Initialize a new VM for the next test case
			config := VMConfig{ StartFn: test.startfn }
			vm, err := CreateVM(config)
            if (err != nil) {
                t.Error("Unexpected VM creation error: ", err)
            }

			// Attempt to run the actual test code
			err = vm.Execute(module, config)
            if (err != test.status) {
                t.Error("Unexpected VM status: ", err)
            }
        })
    }
}
