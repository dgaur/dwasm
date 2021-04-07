package wasm

import (
	"errors"
)


var EndOfBlock			= errors.New("End of VM block") //@not really an error
var InvalidOpcode		= errors.New("Invalid opcode")
var UnreachableCode		= errors.New("Unexpected/unreachable code (opcode 0)")


//
// Signature for all interpreted VM instructions:
//	(VM handle, raw binary code, offset) => (status, new offset)
//
type InstructionFunction func(WASMVM, []byte, int)(error, int)

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

func end(vm WASMVM, code []byte, ip int) (error, int) {
	// End of block/function/execution
	return EndOfBlock, ip
}

func nop(vm WASMVM, code []byte, ip int) (error, int) {
	// Nop.  Continue execution at the next instruction
	return nil, ip+1
}

func unreachable(vm WASMVM, code []byte, ip int) (error, int) {
	// Somehow reached unexpected/non-executable code
	return UnreachableCode, ip
}
