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
	case SemicolonToken:
		return "Semicolon"
	case LeftParenToken:
		return "LeftParen"
	case RightParenToken:
		return "RightParen"
	case Return:
		return "Return"
	case PlusToken:
		return "Plus"
	case Let:
		return "Let"
	case Error:
		return "Error"
	case End:
		return "End"
	case As:
		return "As"
	case ColonToken:
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
	case RelOpToken:
		return "RelOp"
	case LeftArrowToken:
		return "LeftArrow"
	case CommaToken:
		return "Comma"
	case CommentSingleLine:
		return "CommentSingleLine"
	case CommentMultiLine:
		return "CommentMultiLine"
	case NewLineToken:
		return "Newline"
	case If:
		return "If"
	case Else:
		return "Else"
	case While:
		return "While"
	case For:
		return "For"
	case Fun:
		return "Fun"
	case Print:
		return "Print"
	case Delay:
		return "Delay"
	case WriteBox:
		return "WriteBox"
	case Write:
		return "Write"
	case PadWidth:
		return "PadWidth"
	case PadHeight:
		return "PadHeight"
	case PadRead:
		return "PadRead"
	case PadRandI:
		return "PadRandI"
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
	PlusToken
	StarToken
	MinusToken
	SlashToken
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
	NewLineToken
	// Keywords (not in the count)
	Return
	Let
	Error
	End
	If
	Else
	While
	For
	Fun

	// Builtins
	PadWidth
	PadHeight
	PadRead
	PadRandI

	Print
	Delay
	WriteBox
	Write

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
	Plus
	Minus
	Slash
	Star
	Other
	Colon
	LeftArrow
	LeftCurly
	RightCurly
	RelOp
	Comma
	Hash
	Dot
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
	StatePlus
	StateStar
	StateSlash
	StateMinus
	StateInt
	StateColon
	StateRelOpExtended
	StateLeftCurly
	StateRightCurly
	StateRelOp
	StateComma
	StateHex
	StateFloat

	StateMultilineCommentOpen
	StateMultilineAlmostClosed
	StateMultilineClosed

	StateSinglelineComment
	StateNewline

	StateCount // total number of states
)

var finalStateToTokenType = map[int]TokenType{
	StateWhitespace: WhitespaceToken,
	StateEquals:     EqualsToken,
	StateSemicolon:  SemicolonToken,
	StateLeftParen:  LeftParenToken,
	StateRightParen: RightParenToken,

	StateMinus: MinusToken,
	StateStar:  StarToken,
	StateSlash: SlashToken,
	StatePlus:  PlusToken,

	StateInt:                   Integer,
	StateColon:                 ColonToken,
	StateLeftCurly:             LeftCurlyToken,
	StateRightCurly:            RightCurlyToken,
	StateRelOp:                 RelOpToken,
	StateComma:                 CommaToken,
	StateHex:                   HexNumber,
	StateFloat:                 Float,
	StateNewline:               NewLineToken,
	StateSinglelineComment:     CommentSingleLine,
	StateMultilineClosed:       CommentMultiLine,
	StateMultilineAlmostClosed: CommentMultiLine,
	StateMultilineCommentOpen:  CommentMultiLine,
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
	'+':  "plus",
	'-':  "minus",
	'*':  "star",
	'/':  "slash",
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
			"plus":      Plus,
			"minus":     Minus,
			"slash":     Slash,
			"star":      Star,
			"other":     Other,
			"colon":     Colon,
			"leftArrow": LeftArrow,
			"lc":        LeftCurly,
			"rc":        RightCurly,
			"relop":     RelOp,
			"comma":     Comma,
			"hash":      Hash,
			"dot":       Dot,
			"nl":        Newline,
		},

		StateList: make([]int, StateCount),
		StatesAccp: []int{
			StateIdent,
			StateWhitespace,
			StateEquals,
			StateSemicolon,
			StateLeftParen,
			StateRightParen,
			StatePlus,
			StateMinus,
			StateStar,
			StateSlash,
			StateInt,
			StateColon,
			StateRelOpExtended,
			StateLeftCurly,
			StateRightCurly,
			StateRelOp,
			StateComma,
			StateHex,
			StateFloat,
			StateSinglelineComment,
			StateMultilineCommentOpen,
			StateMultilineAlmostClosed,
			StateMultilineClosed,
			StateNewline,
		},
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
	l.Tx[StateIdent][Underscore] = StateIdent

	l.Tx[StateStart][Whitespace] = StateWhitespace
	l.Tx[StateWhitespace][Whitespace] = StateWhitespace

	l.Tx[StateStart][Equals] = StateEquals
	l.Tx[StateStart][Semicolon] = StateSemicolon
	l.Tx[StateStart][LeftParen] = StateLeftParen
	l.Tx[StateStart][RightParen] = StateRightParen

	l.Tx[StateStart][Plus] = StatePlus
	l.Tx[StateStart][Minus] = StateMinus
	l.Tx[StateStart][Star] = StateStar
	l.Tx[StateStart][Slash] = StateSlash

	l.Tx[StatePlus][Equals] = StatePlus
	l.Tx[StateStar][Equals] = StateStar
	l.Tx[StateSlash][Equals] = StateSlash
	l.Tx[StateMinus][Equals] = StateMinus

	l.Tx[StateStart][Digit] = StateInt
	l.Tx[StateInt][Digit] = StateInt

	l.Tx[StateStart][Colon] = StateColon

	l.Tx[StateMinus][RelOp] = StateRelOpExtended

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

	l.Tx[StateSlash][Slash] = StateSinglelineComment

	l.Tx[StateStart][Newline] = StateNewline

	for idx := 0; idx < int(LexemeCount); idx++ {
		if idx != Newline {
			l.Tx[StateSinglelineComment][idx] = StateSinglelineComment
		}
	}

	l.Tx[StateSlash][Star] = StateMultilineCommentOpen

	l.Tx[StateMultilineCommentOpen][Star] = StateMultilineAlmostClosed
	l.Tx[StateMultilineAlmostClosed][Slash] = StateMultilineClosed
	l.Tx[StateMultilineAlmostClosed][Star] = StateMultilineAlmostClosed

	for idx := 0; idx < int(LexemeCount); idx++ {
		if idx != Star {
			l.Tx[StateMultilineAlmostClosed][idx] = StateMultilineCommentOpen
			l.Tx[StateMultilineCommentOpen][idx] = StateMultilineCommentOpen
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
	case "if":
		return Token{If, lexeme}, true
	case "else":
		return Token{Else, lexeme}, true
	case "while":
		return Token{While, lexeme}, true
	case "for":
		return Token{For, lexeme}, true
	case "fun":
		return Token{Fun, lexeme}, true
	case "__print":
		return Token{Print, lexeme}, true
	case "__delay":
		return Token{Delay, lexeme}, true
	case "__write":
		return Token{Write, lexeme}, true
	case "__write_box":
		return Token{WriteBox, lexeme}, true
	case "__width":
		return Token{PadWidth, lexeme}, true
	case "__height":
		return Token{PadHeight, lexeme}, true
	case "__read":
		return Token{PadRead, lexeme}, true
	case "__random_int":
		return Token{PadRandI, lexeme}, true
	default:
		return Token{}, false
	}
}

func getTypeTokenByLexeme(lexeme string) (Token, bool) {
	switch lexeme {
	case "color":
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
