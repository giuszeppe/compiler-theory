package main

import (
	"fmt"
	"slices"
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
	default:
		return "Unknown"
	}
}

const (
	Identifier TokenType = iota + 1
	Integer
	Whitespace
	// Op
	Equals
	// Syntax
	Semicolon
	LeftParen
	RightParen

	Operator
	Colon
	LeftArrow

	RightCurly
	LeftCurly

	RelOp

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

func NewToken(t TokenType, lexeme string) Token {
	return Token{t, lexeme}
}

type Lexer struct {
	LexemeList []string
	StateList  []int
	StatesAccp []int
	Rows       int
	Cols       int
	Tx         [][]int
}

func NewLexer() Lexer {
	lexer := Lexer{
		LexemeList: []string{
			"_",
			"letter",
			"digit",
			"ws",
			"eq",
			"sc",
			"lp",
			"rp",
			"op",
			"other",
			"colon",
			"leftArrow",
			"lc",
			"rc",
			"relop",
		},
		StateList:  []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
		StatesAccp: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
	}
	lexer.Rows = len(lexer.StateList)
	lexer.Cols = len(lexer.LexemeList)

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
	l.Tx[0][slices.Index(l.LexemeList, "letter")] = 1
	l.Tx[0][slices.Index(l.LexemeList, "_")] = 1
	l.Tx[1][slices.Index(l.LexemeList, "letter")] = 1
	l.Tx[1][slices.Index(l.LexemeList, "digit")] = 1

	l.Tx[0][slices.Index(l.LexemeList, "ws")] = 2
	l.Tx[2][slices.Index(l.LexemeList, "ws")] = 2

	l.Tx[0][slices.Index(l.LexemeList, "eq")] = 3
	l.Tx[0][slices.Index(l.LexemeList, "sc")] = 4
	l.Tx[0][slices.Index(l.LexemeList, "lp")] = 5
	l.Tx[0][slices.Index(l.LexemeList, "rp")] = 6

	l.Tx[0][slices.Index(l.LexemeList, "op")] = 7
	l.Tx[7][slices.Index(l.LexemeList, "eq")] = 7


	l.Tx[0][slices.Index(l.LexemeList, "digit")] = 8
	l.Tx[8][slices.Index(l.LexemeList, "digit")] = 8

	l.Tx[0][slices.Index(l.LexemeList, "colon")] = 9

	l.Tx[0][slices.Index(l.LexemeList, "leftArrow")] = 10

	l.Tx[0][slices.Index(l.LexemeList, "lc")] = 11
	l.Tx[0][slices.Index(l.LexemeList, "rc")] = 12

	l.Tx[0][slices.Index(l.LexemeList, "relop")] = 13
	l.Tx[3][slices.Index(l.LexemeList, "eq")] = 13
	l.Tx[13][slices.Index(l.LexemeList, "eq")] = 13
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
	if state == 1 {
		tok, ok := getKeywordTokenByLexeme(lexeme)
		if ok {
			return tok
		}
		tok, ok = getTypeTokenByLexeme(lexeme)
		if ok {
			return tok
		}
		return Token{Identifier, lexeme}
	} else if state == 2 {
		return Token{Whitespace, lexeme}
	} else if state == 3 {
		return Token{Equals, lexeme}
	} else if state == 4 {
		return Token{Semicolon, lexeme}
	} else if state == 5 {
		return Token{LeftParen, lexeme}
	} else if state == 6 {
		return Token{RightParen, lexeme}
	} else if state == 7 {
		return Token{Operator, lexeme}
	} else if state == 8 {
		return Token{Integer, lexeme}
	} else if state == 9 {
		return Token{Colon, lexeme}
	} else if state == 10 {
		return Token{LeftArrow, lexeme}
	} else if state == 11 {
		return Token{LeftCurly, lexeme}
	} else if state == 12 {
		return Token{RightCurly, lexeme}
	} else if state == 13 {
		return Token{RelOp, lexeme}
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
	switch {
	case isAlpha(ch):
		return "letter"
	case isDigit(ch):
		return "digit"
	case ch == '_':
		return "_"
	case ch == ' ' || ch == '\t' || ch == '\n':
		return "ws"
	case ch == '=':
		return "eq"
	case ch == ';':
		return "sc"
	case ch == '(':
		return "lp"
	case ch == ')':
		return "rp"
	case ch == '+' || ch == '-' || ch == '/' || ch == '*':
		return "op"
	case ch == ':':
		return "colon"
	case ch == '{':
		return "lc"
	case ch == '}':
		return "rc"
	case ch == '<' || ch == '>' || ch == '!':
		return "relop"
	default:
		return "other"
	}
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
		state = l.Tx[state][slices.Index(l.LexemeList, cat)]
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
    source := "+= -= *= /="
	tokens := lexer.GenerateTokens(source)
	for _, tok := range tokens {
		fmt.Printf("Token: %-10v Lexeme: %q\n", tok.Type.String(), tok.Lexeme)
	}

}
