package testcover

import "testing"

func TestAdd(t *testing.T) {
	if Add(2, 3) != 5 {
		t.Error("Add(2, 3) should be 5")
	}
}

func TestSubtract(t *testing.T) {
	if Subtract(5, 3) != 2 {
		t.Error("Subtract(5, 3) should be 2")
	}
}

func TestMultiply(t *testing.T) {
	if Multiply(2, 3) != 6 {
		t.Error("Multiply(2, 3) should be 6")
	}
}

func TestDivide(t *testing.T) {
	if Divide(6, 2) != 3 {
		t.Error("Divide(6, 2) should be 3")
	}
	if Divide(5, 0) != 0 {
		t.Error("Divide(5, 0) should be 0")
	}
}
