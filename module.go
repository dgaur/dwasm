package main

import (
	"bufio"
	"encoding/binary"
	"log"
)

// Preamble/header
const MagicSignature	= 0x6d736100	// "\0asm"
const Version1			= 0x00000001

type Module struct {
}

func loadModule(reader *bufio.Reader) (Module, error) {
	module := Module{}
	
	//
	// Read through the preamble
	//
	var magic, version uint32
	err := binary.Read(reader, binary.LittleEndian, &magic)
	if (err != nil) {
		log.Fatalf("Unable to read module signature: %s\n", err)
	}
	if (magic != MagicSignature) {
		log.Fatalf("Unexpected magic signature: %#x\n", magic)
	}
	err = binary.Read(reader, binary.LittleEndian, &version)
	if (err != nil) {
		log.Fatalf("Unable to read module version: %s\n", err)
	}
	if (version != Version1) {
		log.Fatalf("Unexpected module version: %#x\n", version)
	}

	return module, nil
}
