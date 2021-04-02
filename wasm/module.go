package wasm

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

var InvalidModule = errors.New("Invalid module")

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

func (module Module) String() string {
	var builder strings.Builder

	builder.WriteString("Module:\n")
	for _, section := range module.section {
		if (section != nil) {
			builder.WriteString(fmt.Sprintf("%s\n", section))
		}
	}

	return builder.String()
}

//
// Load and return an entire WASM module
//
func ReadModule(reader io.Reader) (Module, error) {
	module := Module{}

	//
	// Read through the preamble
	//
	var magic, version uint32
	err := binary.Read(reader, binary.LittleEndian, &magic)
	if (err != nil) {
		log.Printf("Unable to read module signature: %s\n", err)
		return module, InvalidModule
	}
	if (magic != MagicSignature) {
		log.Printf("Unexpected magic signature: %#x\n", magic)
		return module, InvalidModule
	}
	err = binary.Read(reader, binary.LittleEndian, &version)
	if (err != nil) {
		log.Printf("Unable to read module version: %s\n", err)
		return module, InvalidModule
	}
	if (version != Version1) {
		log.Printf("Unexpected module version: %#x\n", version)
		return module, InvalidModule
	}

	//
	// Parse the individual sections
	//
	sections := make([]Section, SectionCountMax)
	for {
		section, err := readSection(reader)
		if (err == io.EOF) {
			break
		}
		if (err != nil) {
			log.Printf("Invalid section")
			return module, err
		}

		// Each type of section can occur at most once, except custom sections,
		// so just track by section id/type
		sections[ section.id() ] = section
	}
	module.section = sections

	return module, nil
}
