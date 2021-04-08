package wasm

import(
	"testing"
    )


//
// Test stack usage
//
func TestStack(t *testing.T) {
	stack := Stack{}

	// Push 2 items of different types
	stack.Push(int(1))
	stack.Push(float32(2.5))

	// Stack is no longer empty
	if (stack.IsEmpty()) {
		t.Error("Stack is empty unexpectedly")
	}

	// Pop both items and validate
	data1, err := stack.Pop()
	if (err != nil) {
		t.Error("Unexpected error: ", err)
	}
	if (err == nil && data1.(float32) != 2.5) {
		t.Errorf("Unexpected data: %f", data1)
	}

	data2, err := stack.Pop()
	if (err != nil) {
		t.Error("Unexpected error: ", err)
	}
	if (err == nil && data2.(int) != 1) {
		t.Errorf("Unexpected data: %d", data2)
	}

	// Stack is empty now
	if (!stack.IsEmpty()) {
		t.Error("Stack still has data, unexpectedly")
	}
}
