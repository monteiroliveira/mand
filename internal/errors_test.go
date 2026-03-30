package internal

import (
	"testing"
)

func TestSyntaxError(t *testing.T) {
	err := SetSyntaxError()
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if err.Error() != "Syntax Error" {
		t.Errorf("expected 'Syntax Error', got %q", err.Error())
	}
}

func TestSemanticError(t *testing.T) {
	err := SetSemanticError()
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if err.Error() != "Semantic Error" {
		t.Errorf("expected 'Semantic Error', got %q", err.Error())
	}
}
