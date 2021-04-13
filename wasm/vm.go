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
	StartFn		string
	StartStack	[]int32
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
	bytecode	[]byte		// Cached slice of code[function]
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
	caller	InstructionPointer
	locals	int
}

// Jump to new function/instruction, as a result of call or return
func (thread *WASMInterpreterThread) jump(ip InstructionPointer) {
	thread.current = ip
}

// Save the current stack frame in preparation for a function call
func (thread *WASMInterpreterThread) pushFrame() {
	stackFrame := StackFrame{}

	// Save the current bytecode context
	stackFrame.caller.bytecode	= thread.current.bytecode
	stackFrame.caller.function	= thread.current.function
	stackFrame.caller.ip		= thread.current.ip //@plus calling instruction

	// Save the current stack location (i.e., the stack base pointer), since
	// locals are relative to this offset
	stackFrame.locals = thread.dataStack.Top()

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

	// Initialize the initial VM thread context.  Preload the data stack if
	// necessary
	thread := WASMInterpreterThread{
		callStack: CreateStack(32),
		dataStack: CreateStack(256),
	}
	for _, value := range config.StartStack {
		thread.dataStack.Push(value)
	}
	
	// Simulate a function call to the entry function, so that exit/unwinding
	// behaves properly
	thread.pushFrame()
	entryfn	:= codeSection.function[ int(export.index) ].body[:]
	thread.jump( InstructionPointer{ entryfn, int(export.index), 0 } )
	//@handle functions.local[]


	//
	// Main execution loop
	//
	for {
		// (Re)locate the next opcode in the bytecode, based on prior jumps, etc
		opcode := thread.current.bytecode[ thread.current.ip ]

		// Execute the actual bytecode instruction
		instruction, ok := Opcode[ opcode ]
		if (!ok) {
			log.Printf("VM invalid opcode %#x at IP %#x\n",
				opcode, thread.current.ip)
			return InvalidOpcode
		}
		err = instruction.function(&thread)

		// Deal with errors, branches, etc
		if (err == EndOfBlock && thread.callStack.IsEmpty()) {
			// Entry point returned, so exit here
			err = nil
			break
		} else if (err == ReloadBytecode) {
			// Recache a new bytecode block after a call/ret/jump
			thread.current.bytecode =
				codeSection.function[ thread.current.function ].body[:]
		} else if (err != nil) {
			log.Printf("VM runtime error at IP %#x: %s\n", thread.current.ip, err)
			break
		}
		// else, no error.  Continue executing at next linear IP
	}

	// Dump any data left on the stack, in the assumption that these are
	// the result(s) of some function/calculation
	for {
		if (thread.dataStack.IsEmpty()) {
			break
		}
		value, err := thread.dataStack.Pop()
		if (err != nil) {
			log.Printf("Thread stack error on exit: %s\n", err)
			break
		}
		log.Printf("Thread stack: %d\n", value.(int32)) //@assume int32 results
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

