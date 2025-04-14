package main

import (
	"testing"
)

func TestAllTokenTypes(t *testing.T) {
	lexer := NewLexer()

	tests := []struct {
		input    string
		expected Token
	}{
		// Identifiers
		{"abc", Token{Identifier, "abc"}},
		{"_temp1", Token{Identifier, "_temp1"}},

		// Numbers
		{"123", Token{Integer, "123"}},
		{"#ab12cf", Token{HexNumber, "#ab12cf"}},
		{"12.3", Token{Float, "12.3"}},


		// Whitespace
		{" ", Token{Whitespace, " "}},
		{"\n", Token{Whitespace, "\n"}},
		{"\t", Token{Whitespace, "\t"}},

		// Symbols
		{"=", Token{Equals, "="}},
		{";", Token{Semicolon, ";"}},
		{"(", Token{LeftParen, "("}},
		{")", Token{RightParen, ")"}},
		{"+", Token{Operator, "+"}},
		{"-", Token{Operator, "-"}},
		{"*", Token{Operator, "*"}},
		{"/", Token{Operator, "/"}},
		{":", Token{Colon, ":"}},
		{",", Token{Comma, ","}},
		{"{", Token{LeftCurly, "{"}},
		{"}", Token{RightCurly, "}"}},

		// Relational Operators
		{"<", Token{RelOp, "<"}},
		{">", Token{RelOp, ">"}},
		{"!=", Token{RelOp, "!="}},
		{"==", Token{RelOp, "=="}},
		{">=", Token{RelOp, ">="}},
		{"<=", Token{RelOp, "<="}},

		// Arrow
		{"->", Token{LeftArrow, "->"}},

		// Keywords
		{"let", Token{Let, "let"}},
		{"return", Token{Return, "return"}},
		{"as", Token{As, "as"}},
		{"true", Token{True, "true"}},
		{"false", Token{False, "false"}},

		// Types
		{"int", Token{IntType, "int"}},
		{"float", Token{FloatType, "float"}},
		{"bool", Token{BoolType, "bool"}},
		{"colour", Token{ColourType, "colour"}},


		{"//comment and dw", Token{CommentSingleLine, "//comment and dw"}},
	}

	for _, test := range tests {
		tokens := lexer.GenerateTokens(test.input)
		if len(tokens) == 0 {
			t.Fatalf("No token returned for input %q", test.input)
		}
		got := tokens[0]
		if got.Type != test.expected.Type || got.Lexeme != test.expected.Lexeme {
			t.Errorf("Input: %q — Expected token {Type: %v, Lexeme: %q}, got {Type: %v, Lexeme: %q}",
				test.input, test.expected.Type, test.expected.Lexeme, got.Type, got.Lexeme)
		}
	}
}
