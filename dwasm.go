//
// WASM main CLI
//
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"wasm"
)


// VM/CLI configuration
type WASMConfig struct {
	//@logging level
	//@validate?
	//@disassemble?

	filename		string //@list of files/modules
	showSections	bool
}


// Parse + return any CLI configuration.  No side effects.
func initialize() WASMConfig {
	var config = WASMConfig{}

	// Describe all flags
	flag.BoolVar(&config.showSections, "s", false, "Dump .wasm sections")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: %s [options] input.wasm\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Parse + validate any command-line arguments
	flag.Parse()
	if (flag.NArg() != 1) {
		flag.Usage()
	}
	config.filename = flag.Args()[0]

	return config
}


func main() {
	// Parse any CLI options
	config := initialize()

	// Load any modules
	wasmfile, err := os.Open(config.filename)
	if (err != nil) {
		log.Fatalf("Unable to open %s: %s\n", config.filename, err)
	}
	defer wasmfile.Close()

	reader := bufio.NewReader(wasmfile)
	module, err := wasm.ReadModule(reader)
	if (err != nil) {
		log.Fatalf("Unable to load module %s: %s\n", config.filename, err)
	}

	//@link modules if necessary

	// Dispatch any CLI options
	if (config.showSections) {
		log.Println(module)
	}
	//@disassemble, execute, etc

	return
}
