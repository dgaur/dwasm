package wasm

import (
	"errors"
)


var MissingStartFunction = errors.New("Unable to find start/entry function")


//
// VM configuration
//
type VMConfig struct {
	StartFn	string
	//@JIT?
	//@resource allocation/sizing
}


//
// Abstracted interface for WASM virtual machine (interpreter, JIT, etc)
//
type WASMVM interface {
	//@id()
	Execute(Module, VMConfig) error
}


//
// WASM interpreter
//
type WASMInterpreter struct {
	//@some of these should be per-thread
	//@stack
	stats struct {
		//@start
		//@stop
		//@instructions
	}
}

//
// Run the actual interpreter
//
func (vm WASMInterpreter) Execute(module Module, config VMConfig) error {
	//
	// Locate the named start function / entry point
	//
	//@could be implicitly specified via the StartSection
	exportSection, ok := module.section[ExportSectionId].(ExportSection)
	if !ok {
		// No exported resources
		return MissingStartFunction
	}
	export, ok := exportSection.export[ config.StartFn ]
	if !ok {
		// No resource with this name
		return MissingStartFunction
	}
	if export.etype != ExportTypeFunction {
		// Wrong resource type
		return MissingStartFunction
	}

	codeSection, ok := module.section[CodeSectionId].(CodeSection)
	if !ok {
		// No code
		return MissingStartFunction
	}
	if (int(export.index) >= len(codeSection.function)) {
		// Function index is out of range
		return MissingStartFunction
	}

	// This is the actual start function, finally
	function := codeSection.function[ export.index ]
	_ = function.body //@and functions.local[]

	return nil
}


//
// Factory function for generating VM.  No side effects.
//
func CreateVM(config VMConfig) (WASMVM, error) {
	vm := WASMInterpreter{}

	//@initialize memory, etc
	//@link modules
	//@init + export WASI interfaces/hooks

	return vm, nil
}

