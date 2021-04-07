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


// Top-level CLI configuration
type CLIConfig struct {
	//@logging level
	//@disassemble?

	execute			bool
	filename		string //@list of files/modules
	showSections	bool
	validate		bool
	vm				wasm.VMConfig
}


// Parse + return any CLI configuration.  No side effects.
func initialize() CLIConfig {
	var config = CLIConfig{}

	// Describe all flags
	flag.StringVar(&config.vm.StartFn, "f", "",    "Start/entry function")
	flag.BoolVar(&config.showSections, "s", false, "Dump .wasm sections")
	flag.BoolVar(&config.validate,     "v", false, "Validate .wasm sections")
	flag.BoolVar(&config.execute,      "x", false, "Start VM + execute")
	

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
	//
	// Parse any CLI options
	//
	config := initialize()

	//
	// Load any modules
	//
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

	//
	// Dispatch any CLI options
	//
	if (config.showSections) {
		log.Println(module)
	}
	if (config.validate) {
		err = module.Validate()
		if (err != nil) {
			log.Fatalf("Module validation failed: %s\n", err)
		}
	}
	//@disassemble
	if (config.execute) {
		vm, err := wasm.CreateVM(config.vm)
		if (err != nil) {
			log.Fatalf("Unable to initialize VM: %s\n", err)
		}
		err = vm.Execute(module, config.vm)
		if (err == nil) {
			log.Printf("VM exited cleanly")
		} else {
			log.Fatalf("VM error: %s\n", err)
		}
	}
	
	return
}
