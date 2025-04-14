package main

import (
	"fmt"
	"slices"
	"strings"
)

type Token struct {
	Type   TokenType
	Lexeme string
}

type TokenType int

func (t TokenType) String() string {
	switch t {
	case Identifier:
		return "Identifier"
	case Integer:
		return "Integer"
	case Whitespace:
		return "Whitespace"
	case Equals:
		return "Equals"
	case Semicolon:
		return "Semicolon"
	case LeftParen:
		return "LeftParen"
	case RightParen:
		return "RightParen"
	case Return:
		return "Return"
	case Operator:
		return "Operator"
	case Let:
		return "Let"
	case Error:
		return "Error"
	case End:
		return "End"
	case As:
		return "As"
	case Colon:
		return "Colon"
	case IntType:
		return "IntType"
	case BoolType:
		return "BoolType"
	case ColourType:
		return "ColourType"
	case FloatType:
		return "FloatType"
	case LeftCurly:
		return "LeftCurly"
	case RightCurly:
		return "RightCurly"
	case RelOp:
		return "RelOp"
	case LeftArrow:
		return "LeftArrow"
	case Comma:
		return "Comma"
	case CommentSingleLine:
		return "CommentSingleLine"
	case CommentMultiLine:
		return "CommentMultiLine"
	default:
		return "Unknown"
	}
}

const (
	Identifier TokenType = iota + 1
	Integer
	WhitespaceToken
	// Op
	EqualsToken
	// Syntax
	SemicolonToken
	LeftParenToken
	RightParenToken

	OperatorToken
	ColonToken
	LeftArrowToken

	RightCurlyToken
	LeftCurlyToken

	RelOpToken
	CommaToken

	HexNumber
	Float

	// Comments
	CommentSingleLine
	CommentMultiLine
	NewLine
	// Keywords (not in the count)
	Return
	Let
	Error
	End

	// Type
	IntType
	FloatType
	BoolType
	ColourType
	As

	True
	False
)

// Lexeme constants
const (
	Underscore = iota
	Letter
	Digit
	Whitespace
	Equals
	Semicolon
	LeftParen
	RightParen
	Operator
	Other
	Colon
	LeftArrow
	LeftCurly
	RightCurly
	RelOp
	Comma
	Hash
	Dot
	Slash
	Star
	Newline
	LexemeCount // total count of lexeme types
)

// States constants
const (
	StateStart = iota
	StateIdent
	StateWhitespace
	StateEquals
	StateSemicolon
	StateLeftParen
	StateRightParen
	StateOperator
	StateInt
	StateColon
	StateRelOpExtended
	StateLeftCurly
	StateRightCurly
	StateRelOp
	StateComma
	StateHex
	StateFloat
	StateMultilineComment
	StateNewline
	StateCount // total number of states
)

var finalStateToTokenType = map[int]TokenType{
	StateWhitespace: Whitespace,
	StateEquals:     Equals,
	StateSemicolon:  Semicolon,
	StateLeftParen:  LeftParen,
	StateRightParen: RightParen,
	StateOperator:   Operator,
	StateInt:        Integer,
	StateColon:      Colon,
	StateLeftCurly:  LeftCurly,
	StateRightCurly: RightCurly,
	StateRelOp:      RelOp,
	StateComma:      Comma,
	StateHex:        HexNumber,
	StateFloat:      Float,
	StateNewline:    NewLine,
}
var charCategoryMap = map[byte]string{
	'_':  "_",
	' ':  "ws",
	'\t': "ws",
	'\n': "nl",
	'=':  "eq",
	';':  "sc",
	'(':  "lp",
	')':  "rp",
	'+':  "op",
	'-':  "op",
	'*':  "op",
	'/':  "op",
	':':  "colon",
	'{':  "lc",
	'}':  "rc",
	'<':  "relop",
	'>':  "relop",
	'!':  "relop",
	',':  "comma",
	'#':  "hash",
	'.':  "dot",
}

func NewToken(t TokenType, lexeme string) Token {
	return Token{t, lexeme}
}

type Lexer struct {
	LexemeMap  map[string]int
	StateList  []int
	StatesAccp []int
	Rows       int
	Cols       int
	Tx         [][]int
}

func NewLexer() Lexer {
	lexer := Lexer{
		LexemeMap: map[string]int{
			"_":         Underscore,
			"letter":    Letter,
			"digit":     Digit,
			"ws":        Whitespace,
			"eq":        Equals,
			"sc":        Semicolon,
			"lp":        LeftParen,
			"rp":        RightParen,
			"op":        Operator,
			"other":     Other,
			"colon":     Colon,
			"leftArrow": LeftArrow,
			"lc":        LeftCurly,
			"rc":        RightCurly,
			"relop":     RelOp,
			"comma":     Comma,
			"hash":      Hash,
			"dot":       Dot,
			"slash":     Slash,
			"star":      Star,
			"nl":        Newline,
		},

		StateList: make([]int, StateCount),
		StatesAccp: []int{
			StateIdent,
			StateWhitespace,
			StateEquals,
			StateSemicolon,
			StateLeftParen, StateRightParen, StateOperator, StateInt, StateColon, StateRelOpExtended, StateLeftCurly, StateRightCurly, StateRelOp, StateComma, StateHex, StateFloat, StateMultilineComment, StateNewline},

		// StateList:  []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
		// StatesAccp: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 18, 19},
	}
	lexer.Rows = StateCount
	lexer.Cols = LexemeCount

	lexer.Tx = make([][]int, lexer.Rows)
	for i := 0; i < lexer.Rows; i++ {
		lexer.Tx[i] = make([]int, lexer.Cols)
		for j := 0; j < lexer.Cols; j++ {
			lexer.Tx[i][j] = -1
		}
	}
	lexer.initialiseTable()
	return lexer
}

func (l *Lexer) initialiseTable() {
	l.Tx[StateStart][Letter] = StateIdent
	l.Tx[StateStart][Underscore] = StateIdent
	l.Tx[StateIdent][Letter] = StateIdent
	l.Tx[StateIdent][Digit] = StateIdent

	l.Tx[StateStart][Whitespace] = StateWhitespace
	l.Tx[StateWhitespace][Whitespace] = StateWhitespace

	l.Tx[StateStart][Equals] = StateEquals
	l.Tx[StateStart][Semicolon] = StateSemicolon
	l.Tx[StateStart][LeftParen] = StateLeftParen
	l.Tx[StateStart][RightParen] = StateRightParen

	l.Tx[StateStart][Operator] = StateOperator
	l.Tx[StateOperator][Equals] = StateOperator

	l.Tx[StateStart][Digit] = StateInt
	l.Tx[StateInt][Digit] = StateInt

	l.Tx[StateStart][Colon] = StateColon
	l.Tx[StateOperator][RelOp] = StateRelOpExtended

	l.Tx[StateStart][LeftCurly] = StateLeftCurly
	l.Tx[StateStart][RightCurly] = StateRightCurly

	l.Tx[StateStart][RelOp] = StateRelOp
	l.Tx[StateEquals][Equals] = StateRelOp
	l.Tx[StateRelOp][Equals] = StateRelOp

	l.Tx[StateStart][Comma] = StateComma

	l.Tx[StateStart][Hash] = StateHex
	l.Tx[StateHex][Digit] = StateHex
	l.Tx[StateHex][Letter] = StateHex

	l.Tx[StateInt][Dot] = StateFloat
	l.Tx[StateFloat][Digit] = StateFloat

	l.Tx[StateOperator][Operator] = StateMultilineComment

	l.Tx[StateStart][Newline] = StateNewline

	for idx := 0; idx < int(LexemeCount); idx++ {
		if idx != Newline {
			l.Tx[StateMultilineComment][idx] = StateMultilineComment
		}
	}
}

func (l *Lexer) isAcceptingState(state int) bool {
	return slices.Index(l.StatesAccp, state) != -1
}

func getKeywordTokenByLexeme(lexeme string) (Token, bool) {
	switch lexeme {
	case "return":
		return Token{Return, lexeme}, true
	case "let":
		return Token{Let, lexeme}, true
	case "as":
		return Token{As, lexeme}, true
	case "true":
		return Token{True, lexeme}, true
	case "false":
		return Token{False, lexeme}, true
	default:
		return Token{}, false
	}
}

func getTypeTokenByLexeme(lexeme string) (Token, bool) {
	switch lexeme {
	case "colour":
		return Token{ColourType, lexeme}, true
	case "int":
		return Token{IntType, lexeme}, true
	case "bool":
		return Token{BoolType, lexeme}, true
	case "float":
		return Token{FloatType, lexeme}, true
	default:
		return Token{}, false
	}
}

func (l *Lexer) getTokenTypeByFinalState(state int, lexeme string) Token {
	switch state {
	case StateIdent:
		if tok, ok := getKeywordTokenByLexeme(lexeme); ok {
			return tok
		}
		if tok, ok := getTypeTokenByLexeme(lexeme); ok {
			return tok
		}
		return Token{Identifier, lexeme}

	case StateRelOpExtended:
		if lexeme == "->" {
			return Token{LeftArrow, lexeme}
		}

	case StateMultilineComment:
		if strings.HasPrefix(lexeme, "//") {
			return Token{CommentSingleLine, lexeme}
		}
		if strings.HasPrefix(lexeme, "/*") && strings.HasSuffix(lexeme, "*/") {
			return Token{CommentMultiLine, lexeme}
		}
	}

	if tokenType, ok := finalStateToTokenType[state]; ok {
		return Token{tokenType, lexeme}
	}
	return Token{Error, lexeme}
}

func (l *Lexer) isEndOfInput(src string, idx int) bool {
	return idx >= len(src)
}

func (l *Lexer) nextChar(src string, idx int) (bool, byte) {
	if !l.isEndOfInput(src, idx) {
		return true, src[idx]
	}
	return false, '.'
}

func (l *Lexer) catChar(ch byte) string {
	// Basic classification via direct map
	if class, ok := charCategoryMap[ch]; ok {
		return class
	}

	// Fallback to functional classification
	switch {
	case isAlpha(ch):
		return "letter"
	case isDigit(ch):
		return "digit"
	}
	return "other"
}

func (l *Lexer) NextToken(src string, idx int) (Token, string) {
	state := 0
	stack := []int{-2}
	lexeme := ""

	if l.isEndOfInput(src, idx) {
		return Token{Type: End, Lexeme: "end"}, "end"
	}

	for state != -1 {
		if l.isAcceptingState(state) {
			stack = []int{}
		}
		stack = append(stack, state)

		exists, ch := l.nextChar(src, idx)
		lexeme += string(ch)
		if !exists {
			break
		}
		idx++

		cat := l.catChar(ch)
		idx, _ := l.LexemeMap[cat]
		state = l.Tx[state][idx]
	}

	if len(lexeme) > 0 {
		lexeme = lexeme[:len(lexeme)-1]
	}

	syntaxError := false
	for len(stack) > 0 {
		if stack[len(stack)-1] == -2 {
			syntaxError = true
			break
		}
		if !l.isAcceptingState(stack[len(stack)-1]) {
			stack = stack[:len(stack)-1]
			if len(lexeme) > 0 {
				lexeme = lexeme[:len(lexeme)-1]
			}
		} else {
			state = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			break
		}
	}

	if syntaxError {
		return Token{Type: Error, Lexeme: lexeme}, lexeme
	}
	if l.isAcceptingState(state) {
		return l.getTokenTypeByFinalState(state, lexeme), lexeme
	}
	return Token{Type: Error, Lexeme: lexeme}, lexeme
}

func (l *Lexer) GenerateTokens(src string) []Token {
	fmt.Println("INPUT::", src)
	tokens := []Token{}
	idx := 0
	token, lexeme := l.NextToken(src, idx)
	tokens = append(tokens, token)

	for token.Type != End {
		idx += len(lexeme)
		token, lexeme = l.NextToken(src, idx)
		tokens = append(tokens, token)
		if token.Type == Error {
			break
		}
	}

	return tokens
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func main() {
	lexer := NewLexer()
	source := `let x:int=3;
    x += 3;

    `
	tokens := lexer.GenerateTokens(source)
	for _, tok := range tokens {
		fmt.Printf("Token: %-10v Lexeme: %q\n", tok.Type.String(), tok.Lexeme)
	}
}
