package wasm

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
)

// Preamble/header constants
const (
	MagicSignature	= 0x6d736100	// "\0asm"
	Version1		= 0x00000001
)

//
// Module structure.  Describes the contents of one complete WASM module.
//
type Module struct {
	section []Section
}

//
// Load and return an entire WASM module
//
func LoadModule(reader *bufio.Reader) (Module, error) {
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

	//
	// Parse the individual sections
	//
	sections := make([]Section, SectionCountMax)
	for {
		section, err := loadSection(reader)
		if (err == io.EOF) {
			break
		}
		if (err != nil) {
			log.Fatalf("Invalid section")
		}
		log.Println(section)

		// Each type of section can occur at most once, except custom sections,
		// so just track by section id/type
		sections[ section.id ] = section
	}
	module.section = sections

	return module, nil
}
