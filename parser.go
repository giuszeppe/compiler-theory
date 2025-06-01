package main

import (
	"fmt"
)

func NewParser(program string) Parser {
	lex := NewLexer()
	parser := Parser{
		Name:       "Parser",
		Lex:        &lex,
		Index:      -1,
		SrcProgram: program,
		Tokens:     lex.GenerateTokens(program),
		CrtToken:   NewToken(Error, ""),
		nextToken:  NewToken(Error, ""),
		ASTRoot:    ASTBlockNode{},
	}
	return parser

}

type Parser struct {
	Name       string
	Lex        *Lexer
	Index      int
	SrcProgram string
	Tokens     []Token
	CrtToken   Token
	nextToken  Token
	ASTRoot    ASTBlockNode
}
type Action func(children []ASTNode) ASTNode

type Rule struct {
	LHS    string   // name of the nonterminal
	RHS    []Symbol // sequence of terminals (TokenType) and nonterminals (string)
	Action Action   // builds the AST node when this rule is reduced
}

type Grammar struct {
	StartSymbol string
	Rules       []Rule
	Table       map[string]map[TokenType]int
}

// Symbol is either a terminal (TokenType) or a nonterminal (string).
type Symbol interface{}

// Parse performs an LL(1) parse of tokens against g, returning the root ASTNode.
func (p *Parser) Parse(g *Grammar) (ASTNode, error) {
	tokens := p.Tokens
	// parsing table: g.Table[nonterminal][lookahead] = ruleIndex
	table := g.Table

	// stack of “what symbols remain” for each active rule
	stack := make([][]Symbol, 0, 10)
	stack = append(stack, []Symbol{g.StartSymbol})

	// parallel stack of AST fragments.  Each frame’s first element is the Action,
	// the rest are its (yet‑to‑be‑filled) children.
	type frag struct {
		act      Action
		children []ASTNode
	}
	astStack := make([]frag, 0, 10)
	// dummy root action just returns its single child
	astStack = append(astStack, frag{
		act: func(ch []ASTNode) ASTNode { return ch[0] },
	})

	pos := 0
	for {
		// if the root frame is empty, we’re done
		if len(stack[0]) == 0 {
			if pos < len(tokens) && tokens[pos].Type != End {
				return nil, fmt.Errorf("extraneous input starting at token %d: Token: %v", pos, tokens[4])
			}
			break
		}
		// Skip whitespace, newlines, and comments
		if tokens[pos].Type == WhitespaceToken || tokens[pos].Type == NewLineToken || tokens[pos].Type == CommentSingleLine || tokens[pos].Type == CommentMultiLine {
			pos++
			continue
		}

		// look at innermost frame
		topFrame := &stack[len(stack)-1]

		// we’ve just completed a rule → reduce
		if len(*topFrame) == 0 {
			stack = stack[:len(stack)-1]
			completed := astStack[len(astStack)-1]
			astStack = astStack[:len(astStack)-1]

			node := completed.act(completed.children)

			parent := &astStack[len(astStack)-1]
			parent.children = append(parent.children, node)
			// pop the nonterminal in the parent’s symbol list
			parentFrame := &stack[len(stack)-1]
			*parentFrame = (*parentFrame)[1:]
			continue
		}

		// otherwise, we still have symbols to match
		if pos >= len(tokens) {
			return nil, fmt.Errorf("incomplete input; ran out of tokens")
		}
		tok := tokens[pos]
		sym := (*topFrame)[0]

		switch s := sym.(type) {
		case TokenType:
			// terminal: must match exactly
			if tok.Type != s {
				return nil, fmt.Errorf("mismatch: expected %v but saw %v at token %d",
					s, tok.Type, pos)
			}
			// consume it
			*topFrame = (*topFrame)[1:]
			// wrap token into a leaf AST node
			leaf := &ASTSimpleExpression{Token: tok}
			astStack[len(astStack)-1].children = append(astStack[len(astStack)-1].children, leaf)
			pos++

		case string:
			// nonterminal: consult table
			lookups, ok := table[s]
			if !ok {
				return nil, fmt.Errorf("no parsing table row for nonterminal %q", s)
			}
			ri, ok := lookups[tok.Type]
			if !ok {
				return nil, fmt.Errorf("no rule for nonterminal %q on lookahead %v, position %v lexeme %v", s, tok.Type, pos, tok.Lexeme)
			}
			rule := g.Rules[ri]

			// push a new frame for that rule’s RHS
			rhsCopy := make([]Symbol, len(rule.RHS))
			copy(rhsCopy, rule.RHS)
			stack = append(stack, rhsCopy)

			// push its action to the AST stack
			astStack = append(astStack, frag{
				act: rule.Action,
			})

		default:
			return nil, fmt.Errorf("invalid symbol on stack: %T %#v", s, s)
		}
	}

	// apply the dummy root action
	rootFrag := astStack[0]
	return rootFrag.act(rootFrag.children), nil
}
