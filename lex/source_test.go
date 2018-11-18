package lex

import "testing"

func TestFileSource_ReadChar(t *testing.T) {
	source := FileSource{ path: "", source: "hello there"}

	char := source.ReadChar()
	if char != 'h' {
		t.Error(char, "is not h")
	}

	char = source.ReadChar()
	if char != 'e' {
		t.Error(char, "is not e")
	}

	char = source.ReadChar()
	if char != 'l' {
		t.Error(char, "is not l")
	}
}

func TestFileSource_ReadChar2(t *testing.T) {
	source := FileSource{ path: "", source: "日本語" }
	char := source.ReadChar()
	if char != '日' {
		t.Error(char, "is not valid")
	}
	char = source.ReadChar()
	if char != '本' {
		t.Error(char, "is not valid")
	}
	char = source.ReadChar()
	if char != '語' {
		t.Error(char, "is not valid")
	}
}
