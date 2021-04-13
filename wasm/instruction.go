package wasm

import (
	"errors"
)


// Internal VM coordination
var EndOfBlock			= errors.New("End of VM block")
var ReloadBytecode		= errors.New("call/ret/jmp, must recache bytecode")

// Runtime errors
var InvalidOpcode		= errors.New("Invalid opcode")
var UnreachableCode		= errors.New("Unexpected/unreachable code (opcode 0)")



//
// Signature for all interpreted VM instructions:
//	(VM handle, raw binary code, offset) => (status, new offset)
//
type InstructionFunction func(*WASMInterpreterThread)(error)

type Instruction struct {
	name		string
	function	InstructionFunction
}


//
// Opcode map.  Top-level structure that drives the interpreter.  For
// each supported opcode => function pointer to the underlying implementation.
//
var Opcode = map[uint8]Instruction {
	// Control instructions
	0x00:	Instruction{"unreachable",	unreachable},
	0x01:	Instruction{"nop",			nop},
	0x0B:	Instruction{"end",			end},

	// Variable instructions
	0x20:	Instruction{"local.get",	localget},

	// Numeric instructions
	0x6A:	Instruction{"i32.add",		i32add},
}


func end(thread *WASMInterpreterThread) error {
	// End of block/function/execution
	//@how to distinguish between return vs end of block?
	stackFrame, err := thread.popFrame()
	if (err != nil) {
		return err
	}

	// Restore prior thread context
	//@clean up dataStack/locals?
	thread.jump(stackFrame.caller)

	return EndOfBlock
}

func i32add(thread *WASMInterpreterThread) error {
	// Consumed the opcode
	thread.current.ip += 1

	// i32.add takes two arguments
	value0, err := thread.dataStack.Pop()
	if (err != nil) {
		return err
	}
	value1, err := thread.dataStack.Pop()
	if (err != nil) {
		return err
	}

	thread.dataStack.Push(value0.(int32) + value1.(int32))
	return nil
}

func localget(thread *WASMInterpreterThread) error {
	// Consumed the opcode
	thread.current.ip += 1

	// Stack frame contains pointer to the local parameters
	value, err := thread.callStack.Peek(0)
	if (err != nil) {
		return err
	}
	stackFrame := value.(StackFrame)

	// "local.get" takes one argument: an index into the function local parms
	index := thread.current.bytecode[thread.current.ip]
	local, err := thread.dataStack.Peek( stackFrame.locals - int(index) )
	if (err != nil) {
		return err
	}

	// Consumed the argument
	thread.current.ip += 1

	// Push the local parameter onto the immediate stack for later consumption
	thread.dataStack.Push(local)

	return nil
}

func nop(thread *WASMInterpreterThread) error {
	// No-op.  Continue execution at the next instruction
	thread.current.ip += 1
	return nil
}

func unreachable(thread *WASMInterpreterThread) error {
	// Somehow reached unexpected/non-executable code
	return UnreachableCode
}
