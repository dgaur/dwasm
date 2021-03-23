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

func initialize() string {
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

	return flag.Args()[0]
}


func main() {
	filename := initialize()
	wasmfile, err := os.Open(filename)
	if (err != nil) {
		log.Fatalf("Unable to open %s: %s\n", filename, err)
	}
	defer wasmfile.Close()

	reader := bufio.NewReader(wasmfile)
	_, _ = wasm.ReadModule(reader)
	
	return
}
