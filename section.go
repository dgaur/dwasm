package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

// Section ids to identify the content of each specific Section block
const (
	CustomSection		= 0
	TypeSection			= 1
	ImportSection		= 2
	FunctionSection		= 3
	TableSection		= 4
	MemorySection		= 5
	GlobalSection		= 6
	ExportSection		= 7
	StartSection		= 8
	ElementSection		= 9
	CodeSection			= 10
	DataSection			= 11
	DataCountSection	= 12

	SectionCountMax		= DataCountSection
)

//
// Section structure.  Describes the content of one complete WASM Section
//
type Section struct {
	id		uint8
	size	uint32
	content	[]byte
}

func (section Section) String() string {
	return fmt.Sprintf("Section type %#x, len %d", section.id, section.size)
}

//
// Parse and return a single Section from a wasm byte sequence
//
func loadSection(reader *bufio.Reader) (Section, error) {
	section := Section{}

	//
	// Each section consists of:
	//	- section id
	//	- section size
	//	- section content
	//
	err := binary.Read(reader, binary.LittleEndian, &section.id)
	if (err == io.EOF) {
		return section, err
	}
	if (err != nil) {
		log.Fatalf("Unable to read section id: %s\n", err)
	}

	section.size, err = readULEB128(reader)
	if (err != nil) {
		log.Fatalf("Unable to read section size: %s\n", err)
	}

	section.content = make([]byte, section.size)
	err = binary.Read(reader, binary.LittleEndian, &section.content)
	if (err != nil) {
		log.Fatalf("Unable to read section content: %s\n", err)
	}

	return section, nil
}


