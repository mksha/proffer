package common

import (
	"testing"
)

func TestDummyFunction(t *testing.T) {
	expected := "Dummy function"
	got := dummyFunction()

	if got != expected {
		t.Errorf("Expected: %s, Got: %s", expected, got)
	}
}
