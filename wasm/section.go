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

	UnknownSectionId	= 16
	SectionCountMax		= UnknownSectionId + 1
)

//
// Section interface.  All section objects present this interface, for
// parsing, validation, etc
//
type Section interface {
	id() uint32
	validate() error
}


//
// Custom section.  These are mostly opaque/proprietary, so cannot be parsed
// or decomposed beyond the section name
//
type CustomSection struct {
	content []byte
	name string
}

// Factory function for decoding and generating a CustomSection from a stream
// of bytes.  No side effects.
func readCustomSection(content []byte) (CustomSection, error) {
	section := CustomSection{}

	// Extract the section name, for debugging
	reader := bytes.NewReader(content)
	name, err := readName(reader)
	if (err != nil) {
		return section, err
	}

	// Rest of the section is opaque data and cannot be parsed.  Just consume
	// the entire block and continue

	return CustomSection{ content, string(name) }, nil
}

func (section CustomSection) id() uint32 {
	return CustomSectionId
}

func (section CustomSection) validate() error {
	// Cannot be validated, so assume always valid
	return nil
}

func (section CustomSection) String() string {
	return fmt.Sprintf("Custom section '%s', size %d",
		section.name, len(section.content))
}


//
// Exports section
//
const (
	ExportTypeFunction	= 0x00
	ExportTypeTable		= 0x01
	ExportTypeMemory	= 0x02
	ExportTypeGlobal	= 0x03
)

var ExportTypeMap = map[int]string {
	ExportTypeFunction:	"function",
	ExportTypeTable:	"table",
	ExportTypeMemory:	"memory",
	ExportTypeGlobal:	"global",
}

// A descriptor for a single exported symbol/reference
type Export struct {
	name	string
	etype	uint8
	index	uint32
}

// Factory function for decoding + returning a single Export descriptor.  No
// side effects.
func readExport(reader *bytes.Reader) (Export, error) {
	export := Export{}

	// Symbol/reference name
	name, err := readName(reader)
	if (err != nil) {
		return export, err
	}
	export.name = name

	// Reference target (function, table, etc)
	etype, err := reader.ReadByte()
	if (err != nil) {
		return export, err
	}
	export.etype = etype

	// Exported index
	index, err := readULEB128(reader)
	if (err != nil) {
		return export, err
	}
	export.index = index

	return export, err
}

func (export Export) String() string {
	return fmt.Sprintf("export: %s, type %s, index %#x",
		export.name, ExportTypeMap[ int(export.etype) ], export.index)
}

// Top-level section for declaring Exported symbols/references
type ExportSection struct {
	export []Export
}

func (section ExportSection) id() uint32 {
	return ExportSectionId
}

// Factory function for decoding and generating an ExportSection from a stream
// of bytes.  No side effects.
func readExportSection(content []byte) (ExportSection, error) {
	section := ExportSection{}
	reader  := bytes.NewReader(content)

	// Export section is encoded as a vector of Export descriptors
	count, err := readVectorLength(reader)
	if (err != nil) {
		return section, err
	}

	// Parse the individual export descriptors
	export := make([]Export, count)
	for i := uint32(0); i < count; i++ {
		export[i], err = readExport(reader)
		if (err != nil) {
			return section, err
		}
	}
	section.export = export

	return section, nil
}


func (section ExportSection) validate() error {
	//@
	return nil
}

func (section ExportSection) String() string {
	var builder strings.Builder

	builder.WriteString("Export section:\n")
	for _, export := range section.export {
		builder.WriteString(fmt.Sprintf("    %s\n", export))
	}
	return builder.String()
}


//
// Memory section
//

// Single memory descriptor
type Memory struct {
	limit Limit
}

// Factory function for decoding + returning a single Memory descriptor.  No
// side effects.
func readMemory(reader *bytes.Reader) (Memory, error) {
	memory := Memory{}
	limit, err := readLimit(reader)
	memory.limit = limit
	return memory, err
}

func (memory Memory) String() string {
	return fmt.Sprintf("memory: min %#x, max %#x",
		memory.limit.min, memory.limit.max)
}

type MemorySection struct {
	memory []Memory
}

// Factory function for decoding and generating a MemorySection from a stream
// of bytes.  No side effects.
func readMemorySection(content []byte) (MemorySection, error) {
	section := MemorySection{}
	reader  := bytes.NewReader(content)

	// Memory section is encoded as a vector of memory limits
	count, err := readVectorLength(reader)
	if (err != nil) {
		return section, err
	}

	// Parse the individual memories
	memory := make([]Memory, count)
	for i := uint32(0); i < count; i++ {
		memory[i], err = readMemory(reader)
		if (err != nil) {
			return section, err
		}
	}
	section.memory = memory

	return section, nil
}


func (section MemorySection) id() uint32 {
	return MemorySectionId
}

func (section MemorySection) validate() error {
	if (len(section.memory) > 1) {
		// At most 1 memory is allowed.  See section 2.5.5 of WASM spec 1.1
		return InvalidSection
	}
	return nil
}

func (section MemorySection) String() string {
	var builder strings.Builder

	builder.WriteString("Memory section:\n")
	for _, memory := range section.memory {
		builder.WriteString(fmt.Sprintf("    %s\n", memory))
	}
	return builder.String()
}


//
// Table section
//
const (
	FuncRefType		= 0x70
	ExternRefType	= 0x6F
)

var TableRefTypeMap = map[int]string {
	FuncRefType:	"function",
	ExternRefType:	"extern",
}

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
	return fmt.Sprintf("table: min %#x, max %#x, type %s",
		table.limit.min, table.limit.max, TableRefTypeMap[ int(table.reftype) ])
}

type TableSection struct {
	table []Table
}

func (section TableSection) id() uint32 {
	return TableSectionId
}

// Factory function for decoding and generating a TableSection from a stream
// of bytes.  No side effects.
func readTableSection(content []byte) (TableSection, error) {
	section := TableSection{}
	reader  := bytes.NewReader(content)

	// Table section is encoded as a vector of tables limits + types
	count, err := readVectorLength(reader)
	if (err != nil) {
		return section, err
	}

	// Parse the individual tables
	table := make([]Table, count)
	for i := uint32(0); i < count; i++ {
		table[i], err = readTable(reader)
		if (err != nil) {
			return section, err
		}
	}
	section.table = table

	return section, nil
}


func (section TableSection) validate() error {
	//@
	return nil
}

func (section TableSection) String() string {
	var builder strings.Builder

	builder.WriteString("Table section:\n")
	for _, table := range section.table {
		builder.WriteString(fmt.Sprintf("    %s\n", table))
	}
	return builder.String()
}


//
// Unknown section.  This is not part of the WASM spec, but useful for dealing
// with unexpected/unrecognized/unsupported sections
//
type UnknownSection struct {
	content []byte
	unknownId uint8
}

// Factory function for decoding and generating an UnknownSection from a stream
// of bytes.  No side effects.
func readUnknownSection(id uint8, content []byte) (UnknownSection, error) {
	return UnknownSection{ content, id }, nil
}

func (section UnknownSection) id() uint32 {
	return UnknownSectionId
}

func (section UnknownSection) validate() error {
	// Cannot be validated, so assume always valid
	return nil
}

func (section UnknownSection) String() string {
	// Include the first few bytes of the payload
	contentLength := len(section.content)
	previewLength := 4
	suffix := " ..."
	if (previewLength > contentLength) {
		previewLength = contentLength
		suffix = ""
	}

	return fmt.Sprintf("Unknown section %#x, size %d: % x%s",
		section.unknownId,
		contentLength,
		section.content[:previewLength],
		suffix)
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
		case CustomSectionId:	section, err = readCustomSection(content)
		case ExportSectionId:	section, err = readExportSection(content)
		case MemorySectionId:	section, err = readMemorySection(content)
		case TableSectionId:	section, err = readTableSection(content)
		default:				section, err = readUnknownSection(id, content)
	}

	return section, nil
}


