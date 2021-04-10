package wasm

import (
	"log"
	"errors"
)


var MissingFunction = errors.New("Unable to find function")

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
// Bytecode instruction pointer (IP).  Wrapper for current function index
// and offset into the block of bytecode for that function.
//
type InstructionPointer struct {
	function	int
	ip			int
}

//
// Context for a single interpreter thread: stacks, current IP, etc
//
type WASMInterpreterThread struct {
	callStack	Stack
	current		InstructionPointer
	dataStack	Stack
	//@logger + unique id

	stats struct {
		//@start
		//@stop
		//@instructions
	}
}

type StackFrame struct {
	caller InstructionPointer
}

// Jump to new function/instruction, as a result of call or return
func (thread *WASMInterpreterThread) jump(ip InstructionPointer) {
	thread.current = ip
}

// Save the current stack frame in preparation for a function call
func (thread *WASMInterpreterThread) pushFrame() {
	stackFrame := StackFrame{}
	stackFrame.caller.function	= thread.current.function
	stackFrame.caller.ip		= thread.current.ip //@plus calling instruction
	thread.callStack.Push(stackFrame)
}

// Unwind the stack frame created by pushFrame()
func (thread *WASMInterpreterThread) popFrame() (StackFrame, error) {
	stackFrame, err := thread.callStack.Pop()
	if (err != nil) {
		return StackFrame{}, err
	}
	return stackFrame.(StackFrame), nil
}


//
// WASM interpreter
//
type WASMInterpreter struct {
	//@threads?
}


//
// Run the actual interpreter
//
func (vm WASMInterpreter) Execute(module Module, config VMConfig) error {
	var err error

	//
	// Locate the named start function / entry point
	//
	exportSection, ok := module.section[ExportSectionId].(ExportSection)
	if !ok {
		// No exported resources
		return MissingFunction
	}
	export, ok := exportSection.export[ config.StartFn ]
	if !ok {
		// No resource with this name
		return MissingFunction
	}
	if export.etype != ExportTypeFunction {
		// Wrong resource type
		return MissingFunction
	}

	codeSection, ok := module.section[CodeSectionId].(CodeSection)
	if !ok {
		// No code
		return MissingFunction
	}
	if (int(export.index) >= len(codeSection.function)) {
		// Function index is out of range
		return MissingFunction
	}

	// Initialize the initial VM thread context
	thread := WASMInterpreterThread{
		callStack: CreateStack(32),
		dataStack: CreateStack(256),
	}

	// Simulate a function call to the entry function, so that exit/unwinding
	// behaves properly
	thread.pushFrame()
	thread.jump( InstructionPointer{ int(export.index), 0 } )
	//@handle functions.local[]


	//
	// Main execution loop
	//
	for {
		function	:= codeSection.function[ thread.current.function ]
		opcode		:= function.body[ thread.current.ip ]

		instruction, ok := Opcode[ opcode ]
		if (!ok) {
			log.Printf("VM invalid opcode %#x at IP %#x\n", opcode, thread.current.ip)
			return InvalidOpcode
		}

		// Execute the actual bytecode instruction
		err = instruction.function(&thread)
		if (err == EndOfBlock && thread.callStack.IsEmpty()) {
			err = nil
			break
		} else if (err != nil) {
			log.Printf("VM runtime error at IP %#x: %s\n", thread.current.ip, err)
			break
		}
		// else, no error.  Continue executing at new IP
	}

	return err
}


//
// Factory function for generating VM.  No side effects.
//
func CreateVM(config VMConfig) (WASMVM, error) {
	vm := WASMInterpreter{}

	//@run module.start() function, if any
	//@initialize memory, etc
	//@link modules
	//@init + export WASI interfaces/hooks

	return vm, nil
}

