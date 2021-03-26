package wasm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

var InvalidSection = errors.New("Invalid section")


//
// Section ids to identify the content of each specific Section block
//
const (
	CustomSectionId		= 0
	TypeSectionId		= 1
	ImportSectionId		= 2
	FunctionSectionId	= 3
	TableSectionId		= 4
	MemorySectionId		= 5
	GlobalSectionId		= 6
	ExportSectionId		= 7
	StartSectionId		= 8
	ElementSectionId	= 9
	CodeSectionId		= 10
	DataSectionId		= 11
	DataCountSectionId	= 12

	SectionCountMax		= DataCountSectionId
)

//
// Section interface.  All section objects present this interface, for
// parsing, validation, etc
//
type Section interface {
	id() uint32
	read([]byte) error
	validate() error
}


//
// Custom section.  These are opaque/proprietary, so cannot be parsed
// or decomposed any further
//
type CustomSection struct {
	content []byte
}

func (section CustomSection) id() uint32 {
	return CustomSectionId
}

func (section *CustomSection) read(content []byte) error {
	// Cannot be parsed, just consume the entire block and continue
	section.content = content
	return nil
}

func (section CustomSection) validate() error {
	// Cannot be validated, so assume always valid
	return nil
}

func (section CustomSection) String() string {
	return fmt.Sprintf("Custom/unknown section (size %d)", len(section.content))
}


//
// Memory section
//
type MemorySection struct {
	limit Limit
}

func (section MemorySection) id() uint32 {
	return MemorySectionId
}

func (section *MemorySection) read(content []byte) error {
	r := bytes.NewReader(content)

	// Memory section is encoded as a vector of memory limits
	mems, err := readVectorLength(r)
	if (mems > 1 || err != nil) {
		// At most 1 memory is allowed.  See section 2.5.5 of WASM spec 1.1
		return InvalidSection
	}
	section.limit, err = readLimit(r)
	return err
}

func (section MemorySection) validate() error {
	if (section.limit.min > section.limit.max) && (section.limit.max != 0) {
		return InvalidSection
	}
	return nil
}

func (section MemorySection) String() string {
	return fmt.Sprintf("Memory section: min %#x, max %#x",
		section.limit.min, section.limit.max)
}


//
// Table section
//
const (
	FuncRefType		= 0x70
	ExternRefType	= 0x6F
)

type Table struct {
	limit	Limit
	reftype	uint8
}

func readTable(reader *bytes.Reader) (Table, error) {
	table := Table{}

	reftype, err := reader.ReadByte()
	if (err != nil) {
		return table, err
	}
	table.reftype = reftype

	table.limit, err = readLimit(reader)
	return table, err
}

func (table Table) String() string {
	return fmt.Sprintf("min %#x, max %#x, type %#x",
		table.limit.min, table.limit.max, table.reftype)
}

type TableSection struct {
	table []Table
}

func (section TableSection) id() uint32 {
	return TableSectionId
}

func (section *TableSection) read(content []byte) error {
	r := bytes.NewReader(content)

	// Table section is encoded as a vector of tables limits + types
	count, err := readVectorLength(r)
	if (err != nil) {
		return InvalidSection
	}

	// Parse the individual tables
	table := make([]Table, count)
	for i := uint32(0); i < count; i++ {
		table[i], err = readTable(r)
		if (err != nil) {
			return err
		}
	}
	section.table = table

	return nil
}

func (section TableSection) validate() error {
	//@
	return nil
}

func (section TableSection) String() string {
	var builder strings.Builder

	builder.WriteString("Table section:\n")
	for _, table := range section.table {
		builder.WriteString(fmt.Sprintf("    table: %s\n", table))
	}
	return builder.String()
}


//
// Parse and return a single Section from a wasm byte sequence.  Each
// Section is basically encoded as a TLV structure
//
func readSection(reader io.Reader) (Section, error) {
	var section Section

	// Read the section id.  This id determines the type of section (code,
	// data, memory, etc)
	var id uint8
	err := binary.Read(reader, binary.LittleEndian, &id)
	if (err == io.EOF) {
		return section, err
	}
	if (err != nil) {
		log.Fatalf("Unable to read section id: %s\n", err)
	}

	// Read the section size. This determines the length of the remaining
	// section content, if any
	var size uint32
	size, err = readULEB128(reader)
	if (err != nil) {
		log.Fatalf("Unable to read section size: %s\n", err)
	}

	// Read the actual content bytes.  The format will vary depending on the
	// exact section id
	content := make([]byte, size)
	err = binary.Read(reader, binary.LittleEndian, &content)
	if (err != nil) {
		log.Fatalf("Unable to read section content: %s\n", err)
	}

	// Delegate the remaining parsing to the Section itself, based on the
	// Section id above
	switch(id) {
		case MemorySectionId:	section = &MemorySection{}
		case TableSectionId:	section = &TableSection{}
		default:				section = &CustomSection{}
	}
	err = section.read(content)
	if (err != nil) {
		log.Fatalf("Unable to parse section %#x content: %s\n", id, err)
	}

	return section, nil
}


