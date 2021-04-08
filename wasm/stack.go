package wasm

import (
	"container/list"
	"errors"
)


var StackUnderflow = errors.New("Stack underflow")


type Stack struct {
	data	list.List
}

func CreateStack() Stack {
	return Stack{}
}

func (stack Stack) IsEmpty() bool {
	return(stack.data.Len() == 0)
}

func (stack *Stack) Pop() (interface{}, error) {
	if (!stack.IsEmpty()) {
		element := stack.data.Front()
		stack.data.Remove(element)
		return element.Value, nil
	} else {
		return nil, StackUnderflow
	}
}

func (stack *Stack) Push(value interface{}) {
	stack.data.PushFront(value)
}

