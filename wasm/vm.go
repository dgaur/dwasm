package wasm


//
// VM configuration
//
type VMConfig struct {
	//@JIT?
	//@stats
	//@resource allocation/sizing
}


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
	//@stack
	stats struct {
		//@start
		//@stop
		//@instructions
	}
}


//
// Factory function for generating VM.  No side effects.
//
func CreateVM(config VMConfig) WASMVM {
	vm := WASMInterpreter{}

	//@initialize memory, etc
	//@link modules
	//@init + export WASI interfaces/hooks

	return(vm)
}

