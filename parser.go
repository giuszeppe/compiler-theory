package main

import (
	"fmt"
	"strconv"
)

type Parser struct {
	tokens       []Token
	index        int
	stack        []interface{}
	astNodes     []ASTNode // your base AST interface
	currentToken Token
}

func NewParser(tokens []Token) *Parser {
	// Remove whitespace tokens
	filtered := []Token{}
	for _, t := range tokens {
		if t.Type != Whitespace {
			filtered = append(filtered, t)
		}
	}
	p := &Parser{
		tokens:   filtered,
		index:    0,
		stack:    []interface{}{"Program"},
		astNodes: []ASTNode{},
	}
	if len(p.tokens) > 0 {
		p.currentToken = p.tokens[0]
	} else {
		p.currentToken = Token{Type: End, Lexeme: "END"}
	}
	return p
}

func (p *Parser) advance() {
	p.index++
	if p.index < len(p.tokens) {
		p.currentToken = p.tokens[p.index]
	} else {
		p.currentToken = Token{Type: End, Lexeme: "END"}
	}
}

func (p *Parser) Parse() []ASTNode {
	for len(p.stack) > 0 {
		top := p.stack[len(p.stack)-1]
		p.stack = p.stack[:len(p.stack)-1]

		switch symbol := top.(type) {
		case string:
			switch symbol {
			case "Program", "Block", "Statement", "VariableDecl", "Assignment", "Expression":
				p.applyRule(symbol)
			default:
				fmt.Printf("Unknown non-terminal: %v\n", symbol)
			}
		case TokenType:
			if p.currentToken.Type == symbol {
				p.advance()
			} else {
				fmt.Printf("Syntax Error: Expected %v, got %v\n", symbol, p.currentToken.Type)
				return nil
			}
		}
	}
	return p.astNodes
}

func (p *Parser) applyRule(nonTerminal string) {
	tt := p.currentToken.Type

	switch nonTerminal {
	case "Program":
		p.stack = append(p.stack, "Block")

	case "Block":
		if tt == Let || tt == Identifier {
			p.stack = append(p.stack, "Block", "Statement")
		} else if tt == End {
			// epsilon
		}

	case "Statement":
		if tt == Let {
			p.stack = append(p.stack, SemicolonToken, "VariableDecl")
		} else if tt == Identifier {
			p.stack = append(p.stack, SemicolonToken, "Assignment")
		} else {
			panic(fmt.Sprintf("Syntax Error: Invalid start of statement: %v", tt))
		}

	case "VariableDecl":
		p.expect(Let)
		idToken := p.expect(Identifier)
		p.expect(ColonToken)
		typeToken := p.expectOne([]TokenType{IntType, FloatType, BoolType, ColourType})
		p.expect(EqualsToken)
		expr := p.parseExpression()
		temp := &ASTVarDeclNode{idToken.Lexeme, typeToken.Lexeme, expr}
		p.astNodes = append(p.astNodes, temp)

	case "Assignment":
		idToken := p.expect(Identifier)
		p.expect(EqualsToken)
		expr := p.parseExpression()
		p.astNodes = append(p.astNodes, &ASTAssignmentNode{ASTVariableNode{idToken}, expr})
	}
}

func (p *Parser) expect(t TokenType) Token {
	if p.currentToken.Type == t {
		token := p.currentToken
		p.advance()
		return token
	}
	panic(fmt.Sprintf("Expected %v, got %v", t, p.currentToken.Type))
}

func (p *Parser) expectOne(types []TokenType) Token {
	for _, t := range types {
		if p.currentToken.Type == t {
			token := p.currentToken
			p.advance()
			return token
		}
	}
	panic(fmt.Sprintf("Expected one of %v, got %v", types, p.currentToken.Type))
}

// Expression parsing with binary ops (+ and -)
func (p *Parser) parseExpression() ASTExpressionNode {
	simpleExpr := p.parseSimpleExpr()
	if p.currentToken.Type == As {
		p.advance()

        castType := p.expectOne([]TokenType{IntType, FloatType, ColourType, BoolType}).Lexeme

		return ASTExpressionNode{Expr: simpleExpr, Type: castType}
	}
	return ASTExpressionNode{Expr: simpleExpr}
}

func (p *Parser) parseSimpleExpr() ASTNode {
	node := p.parseTerm()
	for p.currentToken.Type == OperatorToken && (p.currentToken.Lexeme == "+" || p.currentToken.Lexeme == "-") {
		opToken := p.currentToken
		p.advance()
		right := p.parseTerm()
		node = &ASTBinaryOpNode{opToken.Lexeme, node, right}
	}
	return node
}

func (p *Parser) parseTerm() ASTNode {
	if p.currentToken.Type == Integer {
		val := p.currentToken.Lexeme
		p.advance()
		i, _ := strconv.Atoi(val)
		return &ASTIntegerNode{Value: i}
	}
	panic(fmt.Sprintf("Expected literal or factor, got %v", p.currentToken.Type))
}
