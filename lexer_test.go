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
		{"abc", Token{Type: Identifier, Lexeme: "abc"}},
		{"_temp1", Token{Type: Identifier, Lexeme: "_temp1"}},

		// Numbers
		{"123", Token{Type: Integer, Lexeme: "123"}},
		{"#ab12cf", Token{Type: HexNumber, Lexeme: "#ab12cf"}},
		{"12.3", Token{Type: Float, Lexeme: "12.3"}},

		// Whitespace
		{" ", Token{Type: Whitespace, Lexeme: " "}},
		{"\n", Token{Type: NewLineToken, Lexeme: "\n"}},
		{"\t", Token{Type: Whitespace, Lexeme: "\t"}},

		// Symbols
		{"=", Token{Type: Equals, Lexeme: "="}},
		{";", Token{Type: SemicolonToken, Lexeme: ";"}},
		{"(", Token{Type: LeftParenToken, Lexeme: "("}},
		{")", Token{Type: RightParenToken, Lexeme: ")"}},
		// Operators
		{"+", Token{Type: PlusToken, Lexeme: "+"}},
		{"-", Token{Type: MinusToken, Lexeme: "-"}},
		{"*", Token{Type: StarToken, Lexeme: "*"}},
		{"/", Token{Type: SlashToken, Lexeme: "/"}},
		{"and", Token{Type: AndToken, Lexeme: "and"}},
		{"or", Token{Type: OrToken, Lexeme: "or"}},
		{"not", Token{Type: NotToken, Lexeme: "not"}},

		{":", Token{Type: ColonToken, Lexeme: ":"}},
		{",", Token{Type: CommaToken, Lexeme: ","}},
		{"{", Token{Type: LeftCurlyToken, Lexeme: "{"}},
		{"}", Token{Type: RightCurlyToken, Lexeme: "}"}},

		// Relational Operators
		{"<", Token{Type: RelOpToken, Lexeme: "<"}},
		{">", Token{Type: RelOpToken, Lexeme: ">"}},
		{"!=", Token{Type: RelOpToken, Lexeme: "!="}},
		{"==", Token{Type: RelOpToken, Lexeme: "=="}},
		{">=", Token{Type: RelOpToken, Lexeme: ">="}},
		{"<=", Token{Type: RelOpToken, Lexeme: "<="}},

		// Arrow
		{"->", Token{Type: LeftArrowToken, Lexeme: "->"}},

		// Keywords
		{"let", Token{Type: Let, Lexeme: "let"}},
		{"return", Token{Type: Return, Lexeme: "return"}},
		{"as", Token{Type: As, Lexeme: "as"}},
		{"true", Token{Type: True, Lexeme: "true"}},
		{"false", Token{Type: False, Lexeme: "false"}},

		// Types
		{"int", Token{Type: IntType, Lexeme: "int"}},
		{"float", Token{Type: FloatType, Lexeme: "float"}},
		{"bool", Token{Type: BoolType, Lexeme: "bool"}},
		{"colour", Token{Type: ColourType, Lexeme: "colour"}},

		{"//comment and dw", Token{Type: CommentSingleLine, Lexeme: "//comment and dw"}},
		{"/* comment and dw */", Token{Type: CommentMultiLine, Lexeme: "/* comment and dw */"}},
		{"/* comment and dw *", Token{Type: CommentMultiLine, Lexeme: "/* comment and dw *"}},
		{"/* comment and dw ", Token{Type: CommentMultiLine, Lexeme: "/* comment and dw "}},
		{`/* comment and dw
                       ez ez ez
                       `, Token{Type: CommentMultiLine, Lexeme: `/* comment and dw
                       ez ez ez
                       `}},
		{"[", Token{Type: LeftBracketToken, Lexeme: "["}},
		{"]", Token{Type: RightBracketToken, Lexeme: "]"}},
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
