package main

import "testing"

func TestParsing(t *testing.T) {
	program := "x = 2"
	parser := NewParser(program)
	parser.Parse()
}
