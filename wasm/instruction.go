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

func nop(thread *WASMInterpreterThread) error {
	// No-op.  Continue execution at the next instruction
	thread.current.ip += 1
	return nil
}

func unreachable(thread *WASMInterpreterThread) error {
	// Somehow reached unexpected/non-executable code
	return UnreachableCode
}
