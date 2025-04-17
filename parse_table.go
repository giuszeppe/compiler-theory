package main

// import "fmt"
//
// type NonTerminal string
//
// const (
// 	VarDecl    NonTerminal = "VarDecl"
// 	TypeRule   NonTerminal = "Type"
// 	Expression NonTerminal = "Expression"
// )
//
// type Symbol interface {
// 	isSymbol()
// }
//
// type TerminalSymbol struct {
// 	Type TokenType
// }
//
// type NonTerminalSymbol struct {
// 	Name NonTerminal
// }
// type SemanticMarker struct {
// 	Builder func([]ASTNode) (ASTNode, int)
// 	Arity   int // how many children to pop
// }
//
// func (TerminalSymbol) isSymbol()    {}
// func (NonTerminalSymbol) isSymbol() {}
// func (SemanticMarker) isSymbol()    {}
//
// // Production represents a grammar production (e.g., E → T E')
// type Production struct {
// 	Expansion []Symbol
// 	Builder   func(children []ASTNode) (ASTNode, int)
// }
//
// // ParseTable is a map of non-terminal → terminal → production
// type ParseTable struct {
// 	table map[NonTerminal]map[TokenType]Production
// }
//
// // NewParseTable creates an empty parse table
// func NewParseTable() ParseTable {
// 	return ParseTable{
// 		table: make(map[NonTerminal]map[TokenType]Production),
// 	}
// }
//
// // AddEntry inserts a production into the parse table
// func (pt *ParseTable) AddEntry(nonTerminal NonTerminal, terminal TokenType, ss Production) {
// 	if pt.table[nonTerminal] == nil {
// 		pt.table[nonTerminal] = make(map[TokenType]Production)
// 	}
// 	pt.table[nonTerminal][terminal] = ss
// }
//
// // GetEntry retrieves a production based on non-terminal and terminal
// func (pt *ParseTable) GetEntry(nonTerminal NonTerminal, terminal TokenType) (Production, bool) {
// 	row, exists := pt.table[nonTerminal]
// 	if !exists {
// 		return Production{}, false
// 	}
// 	prod, found := row[terminal]
// 	return prod, found
// }
//
// // Print prints the entire parse table
// func (pt *ParseTable) Print() {
// 	for nonTerminal, row := range pt.table {
// 		for terminal, prod := range row {
// 			fmt.Printf("M[%s, %s] = %s → %v\n", nonTerminal, terminal, prod)
// 		}
// 	}
// }
//
// func (p *ParseTable) PopulateParseTable() {
//
// 	// ------------------------------
// 	// VarDecl ::= let Identifier : TypeRule = Expression ;
// 	// ------------------------------
// 	p.AddEntry("VarDecl", Let, Production{
// 		Expansion: []Symbol{
// 			TerminalSymbol{Type: Let},
// 			TerminalSymbol{Type: Identifier},
// 			TerminalSymbol{Type: ColonToken},
// 			NonTerminalSymbol{Name: "TypeRule"},
// 			TerminalSymbol{Type: EqualsToken},
// 			NonTerminalSymbol{Name: "Expression"},
// 			TerminalSymbol{Type: SemicolonToken},
// 		},
// 		Builder: func(children []ASTNode) (ASTNode, int) {
// 			return &ASTVarDeclNode{
// 				Name:       children[1].(*ASTLiteralNode).Token.Lexeme,
// 				Type:       children[3].(*ASTTypeNode).Name,
// 				Expression: ASTExpressionNode{},
// 			}, 3
// 		},
// 	})
//
// 	// ------------------------------
// 	// TypeRule ::= int | float | bool | colour
// 	// ------------------------------
// 	p.AddEntry("TypeRule", IntType, Production{
// 		Expansion: []Symbol{TerminalSymbol{Type: IntType}},
// 		Builder: func(children []ASTNode) (ASTNode, int) {
// 			return &ASTTypeNode{Name: "int"}, 0
// 		},
// 	})
//
// 	p.AddEntry("TypeRule", FloatType, Production{
// 		Expansion: []Symbol{TerminalSymbol{Type: FloatType}},
// 		Builder: func(children []ASTNode) (ASTNode, int) {
// 			return &ASTTypeNode{Name: "float"}, 0
// 		},
// 	})
//
// 	p.AddEntry("TypeRule", BoolType, Production{
// 		Expansion: []Symbol{TerminalSymbol{Type: BoolType}},
// 		Builder: func(children []ASTNode) (ASTNode, int) {
// 			return &ASTTypeNode{Name: "bool"}, 0
// 		},
// 	})
//
// 	p.AddEntry("TypeRule", ColourType, Production{
// 		Expansion: []Symbol{TerminalSymbol{Type: ColourType}},
// 		Builder: func(children []ASTNode) (ASTNode, int) {
// 			return &ASTTypeNode{Name: "colour"}, 0
// 		},
// 	})
//
// 	// ------------------------------
// 	// Factor ::= Integer | Identifier
// 	// ------------------------------
// 	p.AddEntry("Factor", Integer, Production{
// 		Expansion: []Symbol{TerminalSymbol{Type: Integer}},
// 		Builder: func(children []ASTNode) (ASTNode, int) {
// 			return &ASTLiteralNode{Token: children[0].(*ASTLiteralNode).Token}, 0
// 		},
// 	})
//
// 	p.AddEntry("Factor", Identifier, Production{
// 		Expansion: []Symbol{TerminalSymbol{Type: Identifier}},
// 		Builder: func(children []ASTNode) (ASTNode, int) {
// 			return &ASTVariableNode{Token: children[0].(*ASTVariableNode).Token}, 0
// 		},
// 	})
//
// 	// ------------------------------
// 	// Term ::= Factor
// 	// ------------------------------
// 	p.AddEntry("Term", Integer, Production{
// 		Expansion: []Symbol{NonTerminalSymbol{Name: "Factor"}},
// 		Builder: func(children []ASTNode) (ASTNode, int) {
// 			return children[0], 0
// 		},
// 	})
//
// 	p.AddEntry("Term", Identifier, Production{
// 		Expansion: []Symbol{NonTerminalSymbol{Name: "Factor"}},
// 		Builder: func(children []ASTNode) (ASTNode, int) {
// 			return children[0], 0
// 		},
// 	})
//
// 	// ------------------------------
// 	// SimpleExpr ::= Term
// 	// ------------------------------
// 	p.AddEntry("SimpleExpr", Integer, Production{
// 		Expansion: []Symbol{NonTerminalSymbol{Name: "Term"}},
// 		Builder: func(children []ASTNode) (ASTNode, int) {
// 			return children[0], 0
// 		},
// 	})
//
// 	p.AddEntry("SimpleExpr", Identifier, Production{
// 		Expansion: []Symbol{NonTerminalSymbol{Name: "Term"}},
// 		Builder: func(children []ASTNode) (ASTNode, int) {
// 			return children[0], 0
// 		},
// 	})
//
// 	// ------------------------------
// 	// Expression ::= SimpleExpr RelOp SimpleExpr
// 	// ------------------------------
// 	// p.AddEntry("Expression", Integer, Production{
// 	// 	Expansion: []Symbol{
// 	// 		NonTerminalSymbol{Name: "SimpleExpr"},
// 	// 		TerminalSymbol{Type: RelOpToken},
// 	// 		NonTerminalSymbol{Name: "SimpleExpr"},
// 	// 		TerminalSymbol{Type: As},
// 	// 	},
// 	// 	Builder: func(children []ASTNode) (ASTNode, int) {
// 	// 		fmt.Println(children[3])
// 	// 		return &ASTExpressionNode{
// 	// 			Left:     children[0],
// 	// 			Operator: children[1].(*ASTLiteralNode).Token.Lexeme,
// 	// 			Right:    children[2],
// 	// 		}, 2
// 	// 	},
// 	// })
//
// 	// Expression
// 	p.AddEntry("Expression", Integer, Production{
// 		Expansion: []Symbol{
// 			NonTerminalSymbol{Name: "SimpleExpr"},
// 			NonTerminalSymbol{Name: "ExpressionRest"},
// 			NonTerminalSymbol{Name: "ExpressionCast"},
// 		},
// 		Builder: func(children []ASTNode) (ASTNode, int) {
// 			expr := children[0]
//
// 			// Chain binary expressions
// 			if children[1] != nil {
// 				for _, part := range children[1].([]ASTNode) {
// 					b := part.(*ASTExpressionNode)
// 					b.Left = expr
// 					expr = b
// 				}
// 			}
//
// 			// Optional casting
// 			if children[2] != nil {
// 				return &ASTCastExprNode{
// 					Expr: expr,
// 					Type: children[2].(*ASTTypeNode).Name,
// 				}, 2
// 			}
//
// 			return expr, 2
// 		},
// 	})
//
// 	// ExpressionRest ::= RelOp SimpleExpr ExpressionRest | ε
// 	p.AddEntry("ExpressionRest", RelOpToken, Production{
// 		Expansion: []Symbol{
// 			TerminalSymbol{Type: RelOpToken},
// 			NonTerminalSymbol{Name: "SimpleExpr"},
// 			NonTerminalSymbol{Name: "ExpressionRest"},
// 		},
// 		Builder: func(children []ASTNode) (ASTNode, int) {
// 			current := &ASTExpressionNode{
// 				Operator: children[0].(*ASTLiteralNode).Token.Lexeme,
// 				Right:    children[1],
// 			}
// 			if children[2] != nil {
// 				rest := children[2].([]ASTNode)
// 				return append([]ASTNode{current}, rest...), 2
// 			}
// 			return []ASTNode{current}, 2
// 		},
// 	})
//
// 	p.AddEntry("ExpressionRest", AsToken, Production{
// 		Expansion: []Symbol{},
// 		Builder: func([]ASTNode) (ASTNode, int) {
// 			return nil, -1
// 		},
// 	})
//
// 	p.AddEntry("ExpressionRest", SemicolonToken, Production{
// 		Expansion: []Symbol{},
// 		Builder: func([]ASTNode) (ASTNode, int) {
// 			return nil, -1
// 		},
// 	})
//
// 	// ExpressionCast ::= 'as' Type
// 	p.AddEntry("ExpressionCast", AsToken, Production{
// 		Expansion: []Symbol{
// 			TerminalSymbol{Type: AsToken},
// 			NonTerminalSymbol{Name: "TypeRule"},
// 		},
// 		Builder: func(children []ASTNode) (ASTNode, int) {
// 			return children[1], 1
// 		},
// 	})
//
// 	// ExpressionCast ::= ε
// 	p.AddEntry("ExpressionCast", SemicolonToken, Production{
// 		Expansion: []Symbol{},
// 		Builder: func([]ASTNode) (ASTNode, int) {
// 			return nil, -1
// 		},
// 	})
//
// }
