package wasm



//
// Abstracted interface for WASM virtual machine (interpreter, JIT, etc)
//
type WASMVM interface {
	//@id()
	//@execute(module)
}


//
// WASM interpreter
//
type WASMInterpreter struct {
}


//
// Factory function for generating VM.  No side effects.
//
func CreateVM() WASMVM {
	vm := WASMInterpreter{}

	//@initialize memory, etc
	//@link modules
	//@init + export WASI interfaces/hooks

	return(vm)
}

