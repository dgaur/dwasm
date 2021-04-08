package wasm

import (
	"errors"
)


var StackUnderflow = errors.New("Stack underflow")


type Stack struct {
	data	[]interface{}
	top		int	// Next/unused index
}

func CreateStack(capacity int) Stack {
	return Stack{ data: make([]interface{}, capacity), top: 0 }
}

func (stack Stack) IsEmpty() bool {
	return(stack.top == 0)
}

// Peek at a specific item without modifying the stack.  No side effects.
// Useful for reading local variables (e.g., "local.get 0" instruction)
func (stack Stack) Peek(index int) (interface{}, error) {
	if (index < stack.top) {
		return stack.data[index], nil
	} else {
		return nil, StackUnderflow
	}
}

// Poke a specific item prior in the stack.  Useful for setting local variables.
// (e.g., "local.set 0" instruction)
func (stack *Stack) Poke(index int, value interface{}) {
	stack.data[index] = value
}

func (stack *Stack) Pop() (interface{}, error) {
	if (!stack.IsEmpty()) {
		stack.top--
		return stack.data[stack.top], nil
	} else {
		return nil, StackUnderflow
	}
}

func (stack *Stack) Push(value interface{}) {
	stack.data[stack.top] = value
	stack.top++
}

// Return current top-of-stack index.  No side effects.  Useful for determining
// base/frame pointer
func (stack Stack) Top() int {
	return (stack.top - 1)
}

