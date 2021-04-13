package wasm

import(
	"testing"
    )


//
// Test stack usage
//
func TestStack(t *testing.T) {
	stack := CreateStack(32)

	// Push 2 items of different types
	stack.Push(int(1))
	stack.Push(float32(2.5))

	// Stack is no longer empty
	if (stack.IsEmpty()) {
		t.Error("Stack is empty unexpectedly")
	}

	// Pop and validate
	data1, err := stack.Pop()
	if (err != nil) {
		t.Error("Unexpected error: ", err)
	}
	if (err == nil && data1.(float32) != 2.5) {
		t.Errorf("Unexpected data: %f", data1)
	}

	// Peek and validate
	top := stack.Top()
	data2, err := stack.Peek(top)
	if (err != nil) {
		t.Error("Unexpected error: ", err)
	}
	if (err == nil && data2.(int) != 1) {
		t.Errorf("Unexpected data: %d", data2)
	}

	// Poke a new value
	stack.Poke(top, 10)

	// Pop and validate
	data3, err := stack.Pop()
	if (err != nil) {
		t.Error("Unexpected error: ", err)
	}
	if (err == nil && data3.(int) != 10) {
		t.Errorf("Unexpected data: %d", data3)
	}

	// Stack should be empty now
	if (!stack.IsEmpty()) {
		t.Error("Stack still has data, unexpectedly")
	}
}
