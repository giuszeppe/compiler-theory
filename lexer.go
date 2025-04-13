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

const (
	Identifier TokenType = iota + 1
	Integer
	Whitespace
	Equals
	Semicolon
	Error
	End
	For
	If
	Else
	While
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
			"other",
		},
		StateList:  []int{0, 1, 2, 3, 4, 5},
		StatesAccp: []int{1, 2, 3, 4, 5},
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
	fmt.Println("TABLE")
	for i := 0; i < lexer.Rows; i++ {
		fmt.Println(lexer.Tx[i])
	}
	return lexer
}

func (l *Lexer) initialiseTable() {
	l.Tx[0][slices.Index(l.LexemeList, "letter")] = 1
	l.Tx[0][slices.Index(l.LexemeList, "_")] = 1
	l.Tx[1][slices.Index(l.LexemeList, "letter")] = 1
	l.Tx[1][slices.Index(l.LexemeList, "digit")] = 1

	//White Space
	l.Tx[0][slices.Index(l.LexemeList, "ws")] = 2
	l.Tx[2][slices.Index(l.LexemeList, "ws")] = 2

	//Eq sign (=)
	l.Tx[0][slices.Index(l.LexemeList, "eq")] = 3

	//Integers
	l.Tx[0][slices.Index(l.LexemeList, "digit")] = 4
	l.Tx[4][slices.Index(l.LexemeList, "digit")] = 4

	//Semicolon sign (;)
	l.Tx[0][slices.Index(l.LexemeList, "sc")] = 5
}

func (l *Lexer) isAcceptingState(state int) bool {
	return slices.Index(l.StatesAccp, state) != -1
}

func (l *Lexer) getTokenTypeByFinalState(state int, lexeme string) Token {
	if state == 1 {
		token, err := getKeywordTokenByLexeme(lexeme)
		if err == nil {
			return token
		}
		return Token{Type: Identifier, Lexeme: lexeme}
	}
	if state == 2 {
		return NewToken(Whitespace, lexeme)
	}
	if state == 3 {
		return NewToken(Equals, lexeme)
	}
	if state == 4 {
		return NewToken(Integer, lexeme)
	}
	if state == 5 {
		return NewToken(Semicolon, lexeme)
	}
	return Token{Type: Error, Lexeme: lexeme}
}

func getKeywordTokenByLexeme(lexeme string) (Token, error) {
	switch lexeme {
	case "for":
		return Token{For, lexeme}, nil
	case "if":
		return Token{If, lexeme}, nil
	case "else":
		return Token{Else, lexeme}, nil
	case "while":
		return Token{While, lexeme}, nil
	}

	return Token{}, fmt.Errorf("No specified keyword")
}

func (l *Lexer) isEndOfInput(srcProgramStr string, idx int) bool {
	return idx > len(strings.Split(srcProgramStr, ""))-1
}

func (l *Lexer) nextChar(srcProgramStr string, idx int) (bool, byte) {
	if !l.isEndOfInput(srcProgramStr, idx) {
		return true, srcProgramStr[idx]
	}
	return false, '.'
}

func (l *Lexer) NextToken(src string, idx int) (Token, string) {
	state := 0
	stack := []int{-2}
	lexeme := ""

	if l.isEndOfInput(src, idx) {
		return Token{Type: Error, Lexeme: "end"}, "end"
	}

	for state != -1 {
		if l.isAcceptingState(state) {
			stack = []int{}
		}
		stack = append(stack, state)

		ok, ch := l.nextChar(src, idx)
		lexeme += string(ch)
		if !ok {
			fmt.Printf("LAST LEXEME: %v\n", lexeme)
			break
		}
		idx++

		cat := l.catChar(ch)
		state = l.Tx[state][slices.Index(l.LexemeList, cat)]
		fmt.Printf("Lexeme: %s => NEXT STATE: %d => CAT: %s => CHAR: %c => STACK: %v\n",
			lexeme, state, cat, ch, stack)

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
			fmt.Printf("POPPED => %v\n", stack)
			if len(lexeme) > 0 {
				lexeme = lexeme[:len(lexeme)-1]
			}
		} else {
			state = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			break
		}
	}
	fmt.Printf("Lexeme: %v with state: %v\n", lexeme, state)

	if syntaxError {
		return Token{Type: Error, Lexeme: "error"}, "error"
	}

	if l.isAcceptingState(state) {
		return l.getTokenTypeByFinalState(state, lexeme), lexeme
	}
	return Token{Type: Error, Lexeme: "error"}, "error"
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
		fmt.Printf("Nxt TOKEN: %s (%s) => IDX: %d\n", token.Type, lexeme, idx)
		if token.Type == Error || token.Type == End {
			break
		}
	}

	if token.Type == End {
		fmt.Println("Encountered end of Input token!! Done")
	}

	return tokens
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isAlpha(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func (l *Lexer) catChar(ch byte) string {
	if isDigit(ch) {
		return "digit"
	}
	if isAlpha(ch) {
		return "letter"
	}
	if ch == '_' {
		return "_"
	}
	if ch == ';' {
		return "sc"
	}
	if ch == '=' {
		return "eq"
	}
	if ch == ' ' {
		return "ws"
	}
	return "other"
}
