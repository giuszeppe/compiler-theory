package main

import (
	"fmt"
	"strconv"
)

func NewGrammar() *Grammar {
	g := &Grammar{
		StartSymbol: "Program",
		Rules:       []Rule{},
		Table:       make(map[string]map[TokenType]int),
	}

	// — Program → StmtList
	g.Rules = append(g.Rules, Rule{
		LHS: "Program",
		RHS: []Symbol{"StmtList"},
		Action: func(ch []ASTNode) ASTNode {
			// ch[0] is *ASTBlockNode
			blk := ch[0].(*ASTBlockNode)
			return &ASTProgramNode{Block: *blk}
		},
	})

	// — Program → { StmtList }
	g.Rules = append(g.Rules, Rule{
		LHS: "Program",
		RHS: []Symbol{LeftCurlyToken, "StmtList", RightCurlyToken},
		Action: func(ch []ASTNode) ASTNode {
			// ch[0] is *ASTBlockNode
			blk := ch[1].(*ASTBlockNode)
			return &ASTProgramNode{Block: *blk}
		},
	})
	// — StmtList → Statement StmtList
	g.Rules = append(g.Rules, Rule{
		LHS: "StmtList",
		RHS: []Symbol{"Statement", "StmtList"},
		Action: func(ch []ASTNode) ASTNode {
			// ch[0] is one ASTNode stmt, ch[1] is *ASTBlockNode with a slice
			stmt := ch[0]
			tail := ch[1].(*ASTBlockNode)
			return &ASTBlockNode{
				Name:  "",
				Stmts: append([]ASTNode{stmt}, tail.Stmts...),
			}
		},
	})

	// — StmtList → ε
	g.Rules = append(g.Rules, Rule{
		LHS: "StmtList",
		RHS: []Symbol{}, // empty
		Action: func(ch []ASTNode) ASTNode {
			return &ASTBlockNode{Name: "", Stmts: []ASTNode{}}
		},
	})

	// — Statement → Assignment ';'
	g.Rules = append(g.Rules, Rule{
		LHS: "Statement",
		RHS: []Symbol{"Identifier", EqualsToken, "Expr", SemicolonToken},
		Action: func(ch []ASTNode) ASTNode {
			// ch[0] and ch[1] were terminals; ch[2] is *ASTExpressionNode
			exprN := ch[2]
			varNode, _ := ch[0].(*ASTVariableNode)
			return &ASTAssignmentNode{
				Id:   *varNode,
				Expr: exprN,
			}
		},
	})

	g.Rules = append(g.Rules, Rule{
		LHS: "Identifier",
		RHS: []Symbol{Identifier, "IdentifierOrArrayAccess"},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTVariableNode{
				Token:  ch[0].(*ASTSimpleExpression).Token,
				Offset: ch[1],
			}
		},
	})

	g.Rules = append(g.Rules, Rule{
		LHS: "IdentifierOrArrayAccess",
		RHS: []Symbol{LeftBracketToken, "Expr", RightBracketToken},
		Action: func(ch []ASTNode) ASTNode {
			return ch[1]
		},
	})

	g.Rules = append(g.Rules, Rule{
		LHS: "IdentifierOrArrayAccess",
		RHS: []Symbol{},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTEpsilon{}
		},
	})

	// — Statement → VariableDecl ';'
	g.Rules = append(g.Rules, Rule{
		LHS: "Statement",
		RHS: []Symbol{Let, Identifier, ColonToken, "TypeRule", "VarDeclSuffix", SemicolonToken},
		Action: func(ch []ASTNode) ASTNode {
			if arrNode, ok := ch[4].(*ASTArrayNode); ok {
				ch[3].(*ASTTypeNode).Name += "[" + strconv.Itoa(arrNode.Size) + "]"
				ch[4].(*ASTArrayNode).Type = ch[3].(*ASTTypeNode).Name
			}
			// if-else to match the VarDeclSuffix and behave differently if it's an array or a normal expression
			return &ASTVarDeclNode{
				Token:      ch[1].(*ASTSimpleExpression).Token,
				Type:       ch[3].(*ASTTypeNode).Name,
				Expression: ch[4],
			}
		},
	})

	// — VarDeclSuffix →  '=' Expr
	g.Rules = append(g.Rules, Rule{
		LHS: "VarDeclSuffix",
		RHS: []Symbol{EqualsToken, "Expr"},
		Action: func(ch []ASTNode) ASTNode {
			return ch[1]
		},
	})

	// - VarDeclSuffix → '[' VarDeclArray
	g.Rules = append(g.Rules, Rule{
		LHS: "VarDeclSuffix",
		RHS: []Symbol{LeftBracketToken, "VarDeclArray"},
		Action: func(ch []ASTNode) ASTNode {
			// ch[0] is a left bracket, ch[1] is *ASTVarDeclArrayNode
			return ch[1]
		},
	})

	// - VarDeclArray → Integer ']' = '[' Literal VarDeclArrayTail ]'
	g.Rules = append(g.Rules, Rule{
		LHS: "VarDeclArray",
		RHS: []Symbol{Integer, RightBracketToken, EqualsToken, LeftBracketToken, "Literal", "VarDeclArrayTail"},
		Action: func(ch []ASTNode) ASTNode {
			arrayNode := ch[5].(*ASTArrayNode)
			arraySize := ch[0].(*ASTSimpleExpression).Token.Lexeme
			v, _ := strconv.Atoi(arraySize)
			arrayNode.Size = v
			arrayNode.Items = append([]ASTNode{ch[4]}, arrayNode.Items...)
			arrayNode.Token = ch[3].(*ASTSimpleExpression).Token
			return arrayNode
		},
	})

	// - VarDeclArrayTail → ',' Literal VarDeclArrayTail
	g.Rules = append(g.Rules, Rule{
		LHS: "VarDeclArrayTail",
		RHS: []Symbol{CommaToken, "Literal", "VarDeclArrayTail"},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTArrayNode{
				Size:  ch[2].(*ASTArrayNode).Size + 1,
				Items: append([]ASTNode{ch[1]}, ch[2].(*ASTArrayNode).Items...),
			}
		},
	})
	// - VarDeclArrayTail → epsilon
	g.Rules = append(g.Rules, Rule{
		LHS: "VarDeclArrayTail",
		RHS: []Symbol{RightBracketToken},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTArrayNode{
				Items: []ASTNode{},
			}
		},
	})

	// - VarDeclArray →  ']' = '[' Literal VarDeclArrayTail ']'
	g.Rules = append(g.Rules, Rule{
		LHS: "VarDeclArray",
		RHS: []Symbol{RightBracketToken, EqualsToken, LeftBracketToken, "Literal", "VarDeclArrayTail"},
		Action: func(ch []ASTNode) ASTNode {
			arrayNode := ch[4].(*ASTArrayNode)
			arrayNode.Size = arrayNode.Size + 1
			arrayNode.Items = append([]ASTNode{ch[3]}, arrayNode.Items...)
			return arrayNode
		},
	})

	// - Literal → IntegerLiteral
	g.Rules = append(g.Rules, Rule{
		LHS: "Literal",
		RHS: []Symbol{Integer},
		Action: func(ch []ASTNode) ASTNode {
			tok := ch[0].(*ASTSimpleExpression).Token
			v, _ := strconv.Atoi(tok.Lexeme)
			return &ASTIntegerNode{Name: tok.Lexeme, Value: v}
		},
	})

	// - Literal → FloatLiteral
	g.Rules = append(g.Rules, Rule{
		LHS: "Literal",
		RHS: []Symbol{Float},
		Action: func(ch []ASTNode) ASTNode {
			tok := ch[0].(*ASTSimpleExpression).Token
			v, _ := strconv.ParseFloat(tok.Lexeme, 64)
			return &ASTFloatNode{Name: tok.Lexeme, Value: v}
		},
	})

	// - Literal → True
	g.Rules = append(g.Rules, Rule{
		LHS: "Literal",
		RHS: []Symbol{True},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTBooleanNode{Value: true}
		},
	})

	// - Literal → False
	g.Rules = append(g.Rules, Rule{
		LHS: "Literal",
		RHS: []Symbol{False},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTBooleanNode{Value: false}
		},
	})

	// - Literal → Color
	g.Rules = append(g.Rules, Rule{
		LHS: "Literal",
		RHS: []Symbol{HexNumber},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTColorNode{
				Value: ch[0].(*ASTSimpleExpression).Token.Lexeme,
			}
		},
	})

	// - VarDeclArray → ']' | ',' Expr VarDeclArray
	// g.Rules = append(g.Rules, Rule{
	// 	LHS: "VarDeclArray",
	// 	RHS: []Symbol{RightBracketToken},
	// 	Action: func(ch []ASTNode) ASTNode {
	// 		return &ASTEpsilon{}
	// 	},
	// })

	// — TypeRule → 'float' | 'int' | 'color' | 'bool' |
	g.Rules = append(g.Rules, Rule{
		LHS: "TypeRule",
		RHS: []Symbol{FloatType},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTTypeNode{
				Name: ch[0].(*ASTSimpleExpression).Token.Lexeme,
			}
		},
	})
	// — TypeRule → 'float' | 'int' | 'color' | 'bool' |
	g.Rules = append(g.Rules, Rule{
		LHS: "TypeRule",
		RHS: []Symbol{IntType},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTTypeNode{
				Name: ch[0].(*ASTSimpleExpression).Token.Lexeme,
			}
		},
	})
	// — TypeRule → 'float' | 'int' | 'color' | 'bool' |
	g.Rules = append(g.Rules, Rule{
		LHS: "TypeRule",
		RHS: []Symbol{BoolType},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTTypeNode{
				Name: ch[0].(*ASTSimpleExpression).Token.Lexeme,
			}

		},
	})
	// — TypeRule → 'float' | 'int' | 'color' | 'bool' |
	g.Rules = append(g.Rules, Rule{
		LHS: "TypeRule",
		RHS: []Symbol{ColourType},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTTypeNode{
				Name: ch[0].(*ASTSimpleExpression).Token.Lexeme,
			}
		},
	})

	// — Expr → SimpleExpr ExprPrime ExprTail
	g.Rules = append(g.Rules, Rule{
		LHS: "Expr",
		RHS: []Symbol{"SimpleExpr", "ExprPrime", "ExprTail"},
		Action: func(ch []ASTNode) ASTNode {
			node := ch[0] // the first term (e.g., "2")
			typeCastNode, isTypeCasted := ch[2].(*ASTTypeCastNode)

			opListExpr := ch[1].(*ASTExpressionNode)
			opList := opListExpr.Expr.(*ASTOpList)

			for _, pair := range opList.Pairs {
				node = &ASTBinaryOpNode{
					Token:    pair.Op,
					Operator: pair.Op.Lexeme,
					Left:     node,
					Right:    pair.Right,
				}
			}

			if isTypeCasted {
				return &ASTTypeCastNode{
					Type: typeCastNode.Type,
					Expr: node,
				}
			}
			return node
		},
	})
	// — ExprTail → "as" TypeRule
	g.Rules = append(g.Rules, Rule{
		LHS: "ExprTail",
		RHS: []Symbol{As, "TypeRule"},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTTypeCastNode{Type: ch[1].(*ASTTypeNode).Name}
		},
	})

	// — ExprTail → ε
	g.Rules = append(g.Rules, Rule{
		LHS: "ExprTail",
		RHS: []Symbol{},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTEpsilon{}
		},
	})

	// — ExprPrime → RelOp SimpleExpr ExprPrime
	g.Rules = append(g.Rules, Rule{
		LHS: "ExprPrime",
		RHS: []Symbol{RelOpToken, "SimpleExpr", "ExprPrime"},
		Action: func(ch []ASTNode) ASTNode {
			op := ch[0].(*ASTSimpleExpression).Token
			term := ch[1]
			tailExpr := ch[2].(*ASTExpressionNode)
			tail := tailExpr.Expr.(*ASTOpList)

			pairs := append([]struct {
				Op    Token
				Right ASTNode
			}{{op, term}}, tail.Pairs...)

			return &ASTExpressionNode{
				Expr: &ASTOpList{Pairs: pairs},
			}
		},
	})

	// — ExprPrime → ε
	g.Rules = append(g.Rules, Rule{
		LHS: "ExprPrime",
		RHS: []Symbol{},
		Action: func(ch []ASTNode) ASTNode {
			if len(ch) == 0 {
				return &ASTExpressionNode{
					Expr: &ASTOpList{},
				}
			}
			return ch[0]
		},
	})

	// — SimpleExpr → Term SimpleExprPrime
	g.Rules = append(g.Rules, Rule{
		LHS: "SimpleExpr",
		RHS: []Symbol{"Term", "SimpleExprPrime"},
		Action: func(ch []ASTNode) ASTNode {
			node := ch[0] // the first term (e.g., "2")

			opListExpr := ch[1].(*ASTExpressionNode)
			opList := opListExpr.Expr.(*ASTOpList)

			for _, pair := range opList.Pairs {
				node = &ASTBinaryOpNode{
					Token:    pair.Op,
					Operator: pair.Op.Lexeme,
					Left:     node,
					Right:    pair.Right,
				}
			}

			return node

		},
	})

	// — SimpleExprPrime → '+' Term SimpleExprPrime
	g.Rules = append(g.Rules, Rule{
		LHS: "SimpleExprPrime",
		RHS: []Symbol{"AdditiveOperator", "Term", "SimpleExprPrime"},
		Action: func(ch []ASTNode) ASTNode {
			op := ch[0].(*ASTSimpleExpression).Token
			term := ch[1]
			tailExpr := ch[2].(*ASTExpressionNode)
			tail := tailExpr.Expr.(*ASTOpList)

			pairs := append([]struct {
				Op    Token
				Right ASTNode
			}{{op, term}}, tail.Pairs...)

			return &ASTExpressionNode{
				Expr: &ASTOpList{Pairs: pairs},
			}
		},
	})

	// — SimpleExprPrime → ε
	g.Rules = append(g.Rules, Rule{
		LHS: "SimpleExprPrime",
		RHS: []Symbol{},
		Action: func(ch []ASTNode) ASTNode {
			if len(ch) == 0 {
				return &ASTExpressionNode{
					Expr: &ASTOpList{},
				}
			}
			return ch[0]
		},
	})

	// — Term → Factor TermPrime
	g.Rules = append(g.Rules, Rule{
		LHS: "Term",
		RHS: []Symbol{"Factor", "TermPrime"},
		Action: func(ch []ASTNode) ASTNode {
			node := ch[0] // the first term (e.g., "2")

			opListExpr := ch[1].(*ASTExpressionNode)
			opList := opListExpr.Expr.(*ASTOpList)

			for _, pair := range opList.Pairs {
				node = &ASTBinaryOpNode{
					Token:    pair.Op,
					Operator: pair.Op.Lexeme,
					Left:     node,
					Right:    pair.Right,
				}
			}

			return node
		},
	})

	// — TermPrime → MultiplicativeOperator Factor TermPrime
	g.Rules = append(g.Rules, Rule{
		LHS: "TermPrime",
		RHS: []Symbol{"MultiplicativeOperator", "Factor", "TermPrime"},
		Action: func(ch []ASTNode) ASTNode {
			op := ch[0].(*ASTSimpleExpression).Token
			term := ch[1]
			tailExpr := ch[2].(*ASTExpressionNode)
			tail := tailExpr.Expr.(*ASTOpList)

			pairs := append([]struct {
				Op    Token
				Right ASTNode
			}{{op, term}}, tail.Pairs...)

			return &ASTExpressionNode{
				Expr: &ASTOpList{Pairs: pairs},
			}
		},
	})

	// - MultiplicativeOperator → '*' | '/' | 'and'
	g.Rules = append(g.Rules, Rule{
		LHS: "MultiplicativeOperator",
		RHS: []Symbol{StarToken},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTSimpleExpression{
				Token: ch[0].(*ASTSimpleExpression).Token,
			}
		},
	})
	g.Rules = append(g.Rules, Rule{
		LHS: "MultiplicativeOperator",
		RHS: []Symbol{SlashToken},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTSimpleExpression{
				Token: ch[0].(*ASTSimpleExpression).Token,
			}
		},
	})
	g.Rules = append(g.Rules, Rule{
		LHS: "MultiplicativeOperator",
		RHS: []Symbol{AndToken},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTSimpleExpression{
				Token: ch[0].(*ASTSimpleExpression).Token,
			}
		},
	})

	// — AdditiveOperator → '+' | '-' | 'or'
	g.Rules = append(g.Rules, Rule{
		LHS: "AdditiveOperator",
		RHS: []Symbol{PlusToken},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTSimpleExpression{
				Token: ch[0].(*ASTSimpleExpression).Token,
			}
		},
	})
	g.Rules = append(g.Rules, Rule{
		LHS: "AdditiveOperator",
		RHS: []Symbol{MinusToken},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTSimpleExpression{
				Token: ch[0].(*ASTSimpleExpression).Token,
			}
		},
	})
	g.Rules = append(g.Rules, Rule{
		LHS: "AdditiveOperator",
		RHS: []Symbol{OrToken},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTSimpleExpression{
				Token: ch[0].(*ASTSimpleExpression).Token,
			}
		},
	})

	// — TermPrime → ε
	g.Rules = append(g.Rules, Rule{
		LHS: "TermPrime",
		RHS: []Symbol{},
		Action: func(ch []ASTNode) ASTNode {
			if len(ch) == 0 {
				return &ASTExpressionNode{
					Expr: &ASTOpList{},
				}
			}
			return ch[0]
		},
	})

	// — Factor → IntegerLiteral
	g.Rules = append(g.Rules, Rule{
		LHS: "Factor",
		RHS: []Symbol{Integer},
		Action: func(ch []ASTNode) ASTNode {
			tok := ch[0].(*ASTSimpleExpression).Token
			v, _ := strconv.Atoi(tok.Lexeme)
			return &ASTIntegerNode{Name: tok.Lexeme, Value: v}
		},
	})

	// — Factor → FloatLiteral
	g.Rules = append(g.Rules, Rule{
		LHS: "Factor",
		RHS: []Symbol{Float},
		Action: func(ch []ASTNode) ASTNode {
			tok := ch[0].(*ASTSimpleExpression).Token
			v, _ := strconv.ParseFloat(tok.Lexeme, 64)
			return &ASTFloatNode{Name: tok.Lexeme, Value: v}
		},
	})

	// — Factor → SubExpr
	g.Rules = append(g.Rules, Rule{
		LHS: "Factor",
		RHS: []Symbol{"SubExpr"},
		Action: func(ch []ASTNode) ASTNode {
			return ch[0]
		},
	})

	// — SubExpr → '(' Expr ')'
	g.Rules = append(g.Rules, Rule{
		LHS: "SubExpr",
		RHS: []Symbol{LeftParenToken, "Expr", RightParenToken},
		Action: func(ch []ASTNode) ASTNode {
			return ch[1] // the Expr inside
		},
	})

	// - Statement → 'if' '(' Expr ')' <Block> [ 'else' <Block> ]
	g.Rules = append(g.Rules, Rule{
		LHS: "Statement",
		RHS: []Symbol{If, "Expr", "Block", "IfTail"},
		Action: func(ch []ASTNode) ASTNode {
			ifNode := ASTIfNode{
				Condition: ch[1],
				ThenBlock: ch[2].(*ASTBlockNode),
				ElseBlock: ch[3],
			}
			return &ifNode
		},
	})
	// - IfTail → 'else' <Block>
	g.Rules = append(g.Rules, Rule{
		LHS: "IfTail",
		RHS: []Symbol{Else, "Block"},
		Action: func(ch []ASTNode) ASTNode {
			return ch[1]
		},
	})

	// - IfTail → ε
	g.Rules = append(g.Rules, Rule{
		LHS: "IfTail",
		RHS: []Symbol{},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTEpsilon{}
		},
	})
	// - Block → '{' StmtList '}'
	g.Rules = append(g.Rules, Rule{
		LHS: "Block",
		RHS: []Symbol{LeftCurlyToken, "StmtList", RightCurlyToken},
		Action: func(ch []ASTNode) ASTNode {
			// ch[1] is *ASTBlockNode
			blk := ch[1].(*ASTBlockNode)
			blk.Name = "Block"
			return blk
		},
	})

	// - Statement -> 'while'  Expr  <Block>
	g.Rules = append(g.Rules, Rule{
		LHS: "Statement",
		RHS: []Symbol{While, "Expr", "Block"},
		Action: func(ch []ASTNode) ASTNode {
			whileNode := ASTWhileNode{
				Condition: ch[1],
				Block:     ch[2].(*ASTBlockNode),
			}
			return &whileNode
		},
	})

	// - Statement -> 'for' '(' Assignment ';' Expr ';' Assignment ')' <Block>
	g.Rules = append(g.Rules, Rule{
		LHS: "Statement",
		RHS: []Symbol{For, LeftParenToken, "ForVarDecl", SemicolonToken, "Expr", SemicolonToken, "ForAssignment", RightParenToken, "Block"},
		Action: func(ch []ASTNode) ASTNode {
			forNode := ASTForNode{
				VarDecl:   ch[2],
				Condition: ch[4],
				Increment: ch[6],
				Block:     ch[8].(*ASTBlockNode),
			}
			return &forNode
		},
	})
	// - ForVarDecl → 'let' Identifier ':' TypeRule '=' Expr
	g.Rules = append(g.Rules, Rule{
		LHS: "ForVarDecl",
		RHS: []Symbol{Let, Identifier, ColonToken, "TypeRule", "VarDeclSuffix"},
		Action: func(ch []ASTNode) ASTNode {
			if arrNode, ok := ch[4].(*ASTArrayNode); ok {
				ch[3].(*ASTTypeNode).Name += "[" + strconv.Itoa(arrNode.Size) + "]"
				ch[4].(*ASTArrayNode).Type = ch[3].(*ASTTypeNode).Name
			}
			return &ASTVarDeclNode{
				Token:      ch[1].(*ASTSimpleExpression).Token,
				Type:       ch[3].(*ASTTypeNode).Name,
				Expression: ch[4],
			}
		},
	})
	// - ForVarDecl → ε
	g.Rules = append(g.Rules, Rule{
		LHS: "ForVarDecl",
		RHS: []Symbol{},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTEpsilon{}
		},
	})

	// - ForAssignment → Identifier '=' Expr
	g.Rules = append(g.Rules, Rule{
		LHS: "ForAssignment",
		RHS: []Symbol{"Identifier", EqualsToken, "Expr"},
		Action: func(ch []ASTNode) ASTNode {
			// ch[0] is *ASTSimpleExpression wrapping the var token
			varnode := ch[0].(*ASTVariableNode)
			exprN := ch[2]
			return &ASTAssignmentNode{
				Id:   *varnode,
				Expr: exprN,
			}
		},
	})
	// - ForAssignment → ε
	g.Rules = append(g.Rules, Rule{
		LHS: "ForAssignment",
		RHS: []Symbol{},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTEpsilon{}
		},
	})
	// - Statement → Fun Identifier '(' FormalParams ')' '->' TypeRule Block
	g.Rules = append(g.Rules, Rule{
		LHS: "Statement",
		RHS: []Symbol{Fun, Identifier, LeftParenToken, "FormalParams", RightParenToken, LeftArrowToken, "TypeRule", "ArrayTypeSignature", "Block"},
		Action: func(ch []ASTNode) ASTNode {
			retType := ch[6].(*ASTTypeNode).Name
			if _, ok := ch[7].(*ASTEpsilon); !ok {
				retType += "[" + ch[7].(*ASTSimpleExpression).Token.Lexeme + "]"
			}
			return &ASTFuncDeclNode{
				Token:      ch[1].(*ASTSimpleExpression).Token,
				Params:     ch[3],
				ReturnType: retType,
				Block:      ch[8].(*ASTBlockNode),
			}
		},
	})

	// - FormalParams → Identifier ':' TypeRule FormalParamsTail
	g.Rules = append(g.Rules, Rule{
		LHS: "FormalParams",
		RHS: []Symbol{Identifier, ColonToken, "TypeRule", "ArrayTypeSignature", "FormalParamsTail"},
		Action: func(ch []ASTNode) ASTNode {
			varType := ch[2].(*ASTTypeNode).Name
			if _, ok := ch[3].(*ASTEpsilon); !ok {
				varType += "[" + ch[3].(*ASTSimpleExpression).Token.Lexeme + "]"
			}
			param := ASTVarDeclNode{
				Token:      ch[0].(*ASTSimpleExpression).Token,
				Type:       varType,
				Expression: &ASTExpressionNode{Expr: &ASTEpsilon{}},
			}
			tail := ch[4].(*ASTFormalParamsNode)

			return &ASTFormalParamsNode{
				Params: append([]ASTNode{&param}, tail.Params...),
			}
		},
	})

	// - FormalParamsTail → ',' Identifier ':' TypeRule FormalParamsTail
	g.Rules = append(g.Rules, Rule{
		LHS: "FormalParamsTail",
		RHS: []Symbol{CommaToken, Identifier, ColonToken, "TypeRule", "ArrayTypeSignature", "FormalParamsTail"},
		Action: func(ch []ASTNode) ASTNode {
			varType := ch[3].(*ASTTypeNode).Name
			if _, ok := ch[4].(*ASTEpsilon); !ok {
				varType += "[" + ch[4].(*ASTSimpleExpression).Token.Lexeme + "]"
			}
			param := ASTVarDeclNode{
				Token:      ch[1].(*ASTSimpleExpression).Token,
				Type:       varType,
				Expression: &ASTExpressionNode{Expr: &ASTEpsilon{}},
			}
			tail := ch[5].(*ASTFormalParamsNode)

			return &ASTFormalParamsNode{
				Params: append([]ASTNode{&param}, tail.Params...),
			}
		},
	})

	// - ArrayTypeSignature → '[' Integer ']'
	g.Rules = append(g.Rules, Rule{
		LHS: "ArrayTypeSignature",
		RHS: []Symbol{LeftBracketToken, Integer, RightBracketToken},
		Action: func(ch []ASTNode) ASTNode {
			// ch[0] is a left bracket, ch[1] is *ASTVarDeclArrayNode
			return ch[1]
		},
	})

	// - ArrayTypeSignature → ε
	g.Rules = append(g.Rules, Rule{
		LHS: "ArrayTypeSignature",
		RHS: []Symbol{},
		Action: func(ch []ASTNode) ASTNode {
			// ch[0] is a left bracket, ch[1] is *ASTVarDeclArrayNode
			return &ASTEpsilon{}
		},
	})

	// - FormalParamsTail → ε
	g.Rules = append(g.Rules, Rule{
		LHS: "FormalParamsTail",
		RHS: []Symbol{},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTFormalParamsNode{
				Params: []ASTNode{},
			}
		},
	})

	// - FormalParams → ε
	g.Rules = append(g.Rules, Rule{
		LHS: "FormalParams",
		RHS: []Symbol{},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTFormalParamsNode{
				Params: []ASTNode{},
			}
		},
	})

	// - Statement → __print Expr ';'
	g.Rules = append(g.Rules, Rule{
		LHS: "Statement",
		RHS: []Symbol{Print, "Expr", SemicolonToken},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTBuiltinFuncNode{
				Token: ch[0].(*ASTSimpleExpression).Token,
				Args:  []ASTNode{ch[1]},
			}
		},
	})

	// - Statement → __delay Expr ';'
	g.Rules = append(g.Rules, Rule{
		LHS: "Statement",
		RHS: []Symbol{Delay, "Expr", SemicolonToken},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTBuiltinFuncNode{
				Token: ch[0].(*ASTSimpleExpression).Token,
				Args:  []ASTNode{ch[1]},
			}
		},
	})

	// - Statement → __write Expr ',' Expr ',' Expr ';'
	g.Rules = append(g.Rules, Rule{
		LHS: "Statement",
		RHS: []Symbol{Write, "Expr", CommaToken, "Expr", CommaToken, "Expr", SemicolonToken},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTBuiltinFuncNode{
				Token: ch[0].(*ASTSimpleExpression).Token,
				Args:  []ASTNode{ch[1], ch[3], ch[5]},
			}
		},
	})

	// - Statement → __write_box Expr ',' Expr ',' Expr ',' Expr ',' Expr ';'
	g.Rules = append(g.Rules, Rule{
		LHS: "Statement",
		RHS: []Symbol{WriteBox, "Expr", CommaToken, "Expr", CommaToken, "Expr", CommaToken, "Expr", CommaToken, "Expr", SemicolonToken},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTBuiltinFuncNode{
				Token: ch[0].(*ASTSimpleExpression).Token,
				Args:  []ASTNode{ch[1], ch[3], ch[5], ch[7], ch[9]},
			}
		},
	})

	g.Rules = append(g.Rules, Rule{
		LHS: "Statement",
		RHS: []Symbol{ClearToken, "Expr", SemicolonToken},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTBuiltinFuncNode{
				Token: ch[0].(*ASTSimpleExpression).Token,
				Args:  []ASTNode{ch[1]},
			}
		},
	})

	// - Factor → __width
	g.Rules = append(g.Rules, Rule{
		LHS: "Factor",
		RHS: []Symbol{PadWidth},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTBuiltinFuncNode{
				Token: ch[0].(*ASTSimpleExpression).Token,
				Args:  []ASTNode{},
			}
		},
	})

	// - Factor → __height
	g.Rules = append(g.Rules, Rule{
		LHS: "Factor",
		RHS: []Symbol{PadHeight},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTBuiltinFuncNode{
				Token: ch[0].(*ASTSimpleExpression).Token,
				Args:  []ASTNode{},
			}
		},
	})

	// - Factor → __read Expr ',' Expr
	g.Rules = append(g.Rules, Rule{
		LHS: "Factor",
		RHS: []Symbol{"ReadExpr"},
		Action: func(ch []ASTNode) ASTNode {
			return ch[0]
		},
	})

	g.Rules = append(g.Rules, Rule{
		LHS: "ReadExpr",
		RHS: []Symbol{PadRead, LeftParenToken, "Expr", CommaToken, "Expr", RightParenToken},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTBuiltinFuncNode{
				Token: ch[0].(*ASTSimpleExpression).Token,
				Args:  []ASTNode{ch[2], ch[4]},
			}
		},
	})

	//	- Factor → __random_int Expr
	g.Rules = append(g.Rules, Rule{
		LHS: "Factor",
		RHS: []Symbol{PadRandI, LeftParenToken, "Expr", RightParenToken},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTBuiltinFuncNode{
				Token: ch[0].(*ASTSimpleExpression).Token,
				Args:  []ASTNode{ch[2]},
			}
		},
	})

	// — Factor → Identifier IdentifierOrFunctionCall
	g.Rules = append(g.Rules, Rule{
		LHS: "Factor",
		RHS: []Symbol{Identifier, "IdentifierOrFunctionCall"},
		Action: func(ch []ASTNode) ASTNode {
			_, isFuncCall := ch[1].(*ASTFuncCallNode)
			if isFuncCall {
				funcCall := ch[1].(*ASTFuncCallNode)
				funcCall.Name = ch[0].(*ASTSimpleExpression).Token
				return funcCall
			}
			return &ASTVariableNode{Token: ch[0].(*ASTSimpleExpression).Token, Offset: ch[1]}
		},
	})

	// — IdentifierOrFunctionCall → [ Expr ]
	g.Rules = append(g.Rules, Rule{
		LHS: "IdentifierOrFunctionCall",
		RHS: []Symbol{LeftBracketToken, "Expr", RightBracketToken},
		Action: func(ch []ASTNode) ASTNode {
			return ch[1]
		},
	})
	// — IdentifierOrFunctionCall → ε
	g.Rules = append(g.Rules, Rule{
		LHS: "IdentifierOrFunctionCall",
		RHS: []Symbol{},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTEpsilon{}
		},
	})

	// - IdentifierOrFunctionCall →  '(' ActualParams ')'
	g.Rules = append(g.Rules, Rule{
		LHS: "IdentifierOrFunctionCall",
		RHS: []Symbol{LeftParenToken, "ActualParams", RightParenToken},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTFuncCallNode{
				Params: ch[1].(*ASTActualParamsNode),
			}
		},
	})

	// - ActualParams → Expr ActualParamsTail
	g.Rules = append(g.Rules, Rule{
		LHS: "ActualParams",
		RHS: []Symbol{"Expr", "ActualParamsTail"},
		Action: func(ch []ASTNode) ASTNode {
			param := ch[0]
			tail := ch[1].(*ASTActualParamsNode)
			return &ASTActualParamsNode{
				Params: append([]ASTNode{param}, tail.Params...),
			}
		},
	})

	// - ActualParams → epsilon
	g.Rules = append(g.Rules, Rule{
		LHS: "ActualParams",
		RHS: []Symbol{},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTActualParamsNode{
				Params: []ASTNode{},
			}
		},
	})

	// - ActualParamsTail → ',' Expr ActualParamsTail
	g.Rules = append(g.Rules, Rule{
		LHS: "ActualParamsTail",
		RHS: []Symbol{CommaToken, "Expr", "ActualParamsTail"},
		Action: func(ch []ASTNode) ASTNode {
			param := ch[1]
			tail := ch[2].(*ASTActualParamsNode)
			return &ASTActualParamsNode{
				Params: append([]ASTNode{param}, tail.Params...),
			}
		},
	})

	// - ActualParamsTail → ε
	g.Rules = append(g.Rules, Rule{
		LHS: "ActualParamsTail",
		RHS: []Symbol{},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTActualParamsNode{
				Params: []ASTNode{},
			}
		},
	})

	// - Statement → Block
	g.Rules = append(g.Rules, Rule{
		LHS: "Statement",
		RHS: []Symbol{"Block"},
		Action: func(ch []ASTNode) ASTNode {
			// ch[0] is *ASTBlockNode
			blk := ch[0].(*ASTBlockNode)
			return &ASTBlockNode{Name: "Statement", Stmts: blk.Stmts}
		},
	})

	// - Factor -> Unary
	g.Rules = append(g.Rules, Rule{
		LHS: "Factor",
		RHS: []Symbol{"Unary"},
		Action: func(ch []ASTNode) ASTNode {
			return ch[0]
		},
	})

	// - Unary -> '-' Factor
	g.Rules = append(g.Rules, Rule{
		LHS: "Unary",
		RHS: []Symbol{"UnaryOperator", "Factor"},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTUnaryOpNode{
				Operator: ch[0].(*ASTSimpleExpression).Token.Lexeme,
				Operand:  ch[1],
			}
		},
	})

	// - UnaryOperator -> '-' | 'not'
	g.Rules = append(g.Rules, Rule{
		LHS: "UnaryOperator",
		RHS: []Symbol{MinusToken},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTSimpleExpression{
				Token: ch[0].(*ASTSimpleExpression).Token,
			}
		},
	})
	// - UnaryOperator -> '-' | 'not'
	g.Rules = append(g.Rules, Rule{
		LHS: "UnaryOperator",
		RHS: []Symbol{NotToken},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTSimpleExpression{
				Token: ch[0].(*ASTSimpleExpression).Token,
			}
		},
	})

	// - Factor -> 'true'
	g.Rules = append(g.Rules, Rule{
		LHS: "Factor",
		RHS: []Symbol{True},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTBooleanNode{
				Value: true,
			}
		},
	})
	// - Factor -> 'false'
	g.Rules = append(g.Rules, Rule{
		LHS: "Factor",
		RHS: []Symbol{False},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTBooleanNode{
				Value: false,
			}
		},
	})

	// - Factor -> HexNumber
	g.Rules = append(g.Rules, Rule{
		LHS: "Factor",
		RHS: []Symbol{HexNumber},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTColorNode{
				Token: ch[0].(*ASTSimpleExpression).Token,
				Value: ch[0].(*ASTSimpleExpression).Token.Lexeme,
			}
		},
	})

	// - Statement -> 'return' Expr ';'
	g.Rules = append(g.Rules, Rule{
		LHS: "Statement",
		RHS: []Symbol{Return, "Expr", SemicolonToken},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTReturnNode{
				Token: ch[0].(*ASTSimpleExpression).Token,
				Expr:  ch[1],
			}
		},
	})

	// — finally, build the LL(1) table:
	g.Table = genTable(g)
	return g
}

// genTable builds the LL(1) parsing table for g.
// It returns table[A][a] = index of the rule in g.Rules to apply when
// the current nonterminal is A and the lookahead token is a.
func genTable(g *Grammar) map[string]map[TokenType]int {
	// 1) collect all nonterminals
	nonterms := make(map[string]struct{})
	for _, r := range g.Rules {
		nonterms[r.LHS] = struct{}{}
	}

	// 2) prepare FIRST sets and nullable map
	first := make(map[string]map[TokenType]struct{})
	nullable := make(map[string]bool)
	for A := range nonterms {
		first[A] = make(map[TokenType]struct{})
		nullable[A] = false
	}

	// helper: add a terminal to first[A], return true if it was new
	addFirst := func(A TokenType, set map[TokenType]struct{}) bool {
		if _, ok := set[A]; !ok {
			set[A] = struct{}{}
			return true
		}
		return false
	}

	// 3) iterate to compute FIRST and nullable
	changed := true
	for changed {
		changed = false
		for _, rule := range g.Rules {
			A := rule.LHS
			// compute FIRST(RHS) and whether RHS is nullable
			rhsNullable := true
			for _, sym := range rule.RHS {
				switch s := sym.(type) {
				case TokenType:
					// terminal: FIRST(RHS) includes itself, and RHS isn't nullable past here
					if addFirst(s, first[A]) {
						changed = true
					}
					rhsNullable = false
					break
				case string:
					// nonterminal: inherit FIRST(s) minus ε
					for t := range first[s] {
						if addFirst(t, first[A]) {
							changed = true
						}
					}
					if !nullable[s] {
						rhsNullable = false
						break
					}
				}
				if !rhsNullable {
					// if this symbol blocks ε, stop scanning RHS
					break
				}
			}
			// if all symbols were nullable (or RHS empty), A is nullable
			if rhsNullable && !nullable[A] {
				nullable[A] = true
				changed = true
			}
		}
	}

	// 4) compute FOLLOW sets
	follow := make(map[string]map[TokenType]struct{})
	for A := range nonterms {
		follow[A] = make(map[TokenType]struct{})
	}
	// start symbol gets End-of-input
	follow[g.StartSymbol][End] = struct{}{}

	// helper to compute FIRST of a sequence of symbols
	firstOfSeq := func(seq []Symbol) (map[TokenType]struct{}, bool) {
		res := make(map[TokenType]struct{})
		seqNullable := true
		for _, sym := range seq {
			if !seqNullable {
				break
			}
			switch s := sym.(type) {
			case TokenType:
				res[s] = struct{}{}
				seqNullable = false
			case string:
				for t := range first[s] {
					res[t] = struct{}{}
				}
				if !nullable[s] {
					seqNullable = false
				}
			}
		}
		return res, seqNullable
	}

	// iterate to fixpoint on FOLLOW
	changed = true
	for changed {
		changed = false
		for _, rule := range g.Rules {
			A := rule.LHS
			rhs := rule.RHS
			for i, sym := range rhs {
				B, ok := sym.(string)
				if !ok {
					continue // skip terminals
				}
				// compute FIRST of β = rhs[i+1:]
				beta := rhs[i+1:]
				firstBeta, betaNullable := firstOfSeq(beta)
				// add FIRST(β) minus ε to FOLLOW(B)
				for t := range firstBeta {
					if t == End {
						continue
					}
					if _, seen := follow[B][t]; !seen {
						follow[B][t] = struct{}{}
						changed = true
					}
				}
				// if β nullable, add FOLLOW(A) to FOLLOW(B)
				if betaNullable {
					for t := range follow[A] {
						if _, seen := follow[B][t]; !seen {
							follow[B][t] = struct{}{}
							changed = true
						}
					}
				}
			}
		}
	}

	// 5) build the parsing table
	table := make(map[string]map[TokenType]int)
	for A := range nonterms {
		table[A] = make(map[TokenType]int)
	}

	for i, rule := range g.Rules {
		A := rule.LHS
		firstRHS, rhsNullable := firstOfSeq(rule.RHS)
		// for each terminal in FIRST(RHS), assign rule i
		for t := range firstRHS {
			table[A][t] = i
		}
		// if RHS nullable, for each b in FOLLOW(A), assign rule i
		if rhsNullable {
			for b := range follow[A] {
				table[A][b] = i
			}
		}
	}

	return table
}
func printParsingTable(g *Grammar) {
	fmt.Println("LL(1) Parsing Table:")
	fmt.Println("---------------------------")

	// Collect all terminal symbols used in the table
	terminalsSet := make(map[TokenType]struct{})
	for _, row := range g.Table {
		for tok := range row {
			terminalsSet[tok] = struct{}{}
		}
	}

	// Convert set to slice for consistent order
	var terminals []TokenType
	for tok := range terminalsSet {
		terminals = append(terminals, tok)
	}

	// Print header
	fmt.Printf("%-10s", "")
	for _, t := range terminals {
		fmt.Printf("%-15s", t.String())
	}
	fmt.Println()

	// Print rows for each nonterminal
	for nonterm, row := range g.Table {
		fmt.Printf("%-10s", nonterm)
		for _, t := range terminals {
			if ruleIndex, ok := row[t]; ok {
				rule := g.Rules[ruleIndex]
				fmt.Printf("%-15s", fmtRule(rule))
			} else {
				fmt.Printf("%-15s", "")
			}
		}
		fmt.Println()
	}
	fmt.Println("---------------------------")
}

// fmtRule formats a rule into a string like "A → B c"
func fmtRule(r Rule) string {
	result := r.LHS + " →"
	for _, sym := range r.RHS {
		switch s := sym.(type) {
		case TokenType:
			result += " " + s.String()
		case string:
			result += " " + s
		}
	}
	if len(r.RHS) == 0 {
		result += " ε"
	}
	return result
}
