package main

import (
	"fmt"
	"strconv"
)

func main() {
	program := `if (3) {
		let x:int = 3;
	} else {
	 	let y:int = 4;
	}`
	parser := NewParser(program)
	printVisitor := PrintNodesVisitor{}
	grammar := NewGrammar()
	node, err := parser.Parse(grammar)
	if err != nil {
		panic(err)
	}
	node.Accept(&printVisitor)
}

type OpList struct {
	Pairs []struct {
		Op    string
		Right ASTNode
	}
}

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
		RHS: []Symbol{Identifier, EqualsToken, "Expr", SemicolonToken},
		Action: func(ch []ASTNode) ASTNode {
			// ch[0] and ch[1] were terminals; ch[2] is *ASTExpressionNode
			varTok := ch[0].(*ASTSimpleExpression).Token
			exprN := ch[2]
			return &ASTAssignmentNode{
				Id:   ASTVariableNode{Token: varTok},
				Expr: exprN,
			}
		},
	})
	// — Statement → VariableDecl ';'
	g.Rules = append(g.Rules, Rule{
		LHS: "Statement",
		RHS: []Symbol{Let, Identifier, ColonToken, "TypeRule", EqualsToken, "Expr", SemicolonToken},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTVarDeclNode{
				Name:       ch[1].(*ASTSimpleExpression).Token.Lexeme,
				Type:       ch[3].(*ASTTypeNode).Name,
				Expression: ch[5],
			}
		},
	})

	// — TypeRule → 'float' | 'int' | 'color' | 'bool' |
	g.Rules = append(g.Rules, Rule{
		LHS: "TypeRule",
		RHS: []Symbol{FloatType},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTTypeNode{}
		},
	})
	// — TypeRule → 'float' | 'int' | 'color' | 'bool' |
	g.Rules = append(g.Rules, Rule{
		LHS: "TypeRule",
		RHS: []Symbol{IntType},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTTypeNode{}
		},
	})
	// — TypeRule → 'float' | 'int' | 'color' | 'bool' |
	g.Rules = append(g.Rules, Rule{
		LHS: "TypeRule",
		RHS: []Symbol{BoolType},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTTypeNode{}
		},
	})
	// — TypeRule → 'float' | 'int' | 'color' | 'bool' |
	g.Rules = append(g.Rules, Rule{
		LHS: "TypeRule",
		RHS: []Symbol{ColourType},
		Action: func(ch []ASTNode) ASTNode {
			return &ASTTypeNode{}
		},
	})

	// — Expr → SimpleExpr ExprPrime
	g.Rules = append(g.Rules, Rule{
		LHS: "Expr",
		RHS: []Symbol{"SimpleExpr", "ExprPrime"},
		Action: func(ch []ASTNode) ASTNode {
			// if ExprPrime is just ε, it returns SimpleExpr wrapped as ASTExpressionNode
			expr, ok := ch[1].(*ASTExpressionNode)
			if ok {
				_, ok := expr.Expr.(*ASTEpsilon)
				if ok {
					return ch[0]
				}
			}
			return ch[1]
		},
	})

	// — ExprPrime → RelOp SimpleExpr ExprPrime
	g.Rules = append(g.Rules, Rule{
		LHS: "ExprPrime",
		RHS: []Symbol{RelOpToken, "SimpleExpr", "ExprPrime"},
		Action: func(ch []ASTNode) ASTNode {
			// ch[0] is *ASTSimpleExpression wrapping the relop token:
			op := ch[0].(*ASTSimpleExpression).Token.Lexeme
			left := ch[1]
			rightTail := ch[2].(*ASTExpressionNode)
			// build a binary op and then possibly flatten with rightTail
			bin := &ASTBinaryOpNode{
				Operator: op,
				Left:     left,
				Right:    rightTail.Expr,
			}
			return &ASTExpressionNode{Expr: bin, Type: ""} // type‑checking later
		},
	})

	// — ExprPrime → ε
	g.Rules = append(g.Rules, Rule{
		LHS: "ExprPrime",
		RHS: []Symbol{},
		Action: func(ch []ASTNode) ASTNode {
			if len(ch) == 0 {
				return &ASTExpressionNode{Expr: &ASTEpsilon{}}
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
					Operator: pair.Op,
					Left:     node,
					Right:    pair.Right,
				}
			}

			return &ASTExpressionNode{
				Expr: node,
			}
		},
	})

	// — SimpleExprPrime → '+' Term SimpleExprPrime
	g.Rules = append(g.Rules, Rule{
		LHS: "SimpleExprPrime",
		RHS: []Symbol{PlusToken, "Term", "SimpleExprPrime"},
		Action: func(ch []ASTNode) ASTNode {
			op := ch[0].(*ASTSimpleExpression).Token.Lexeme
			term := ch[1]
			tailExpr := ch[2].(*ASTExpressionNode)
			tail := tailExpr.Expr.(*ASTOpList)

			pairs := append([]struct {
				Op    string
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
					Operator: pair.Op,
					Left:     node,
					Right:    pair.Right,
				}
			}

			return &ASTExpressionNode{
				Expr: node,
			}
		},
	})

	// — TermPrime → '*' Factor TermPrime
	g.Rules = append(g.Rules, Rule{
		LHS: "TermPrime",
		RHS: []Symbol{StarToken, "Factor", "TermPrime"},
		Action: func(ch []ASTNode) ASTNode {
			op := ch[0].(*ASTSimpleExpression).Token.Lexeme
			term := ch[1]
			tailExpr := ch[2].(*ASTExpressionNode)
			tail := tailExpr.Expr.(*ASTOpList)

			pairs := append([]struct {
				Op    string
				Right ASTNode
			}{{op, term}}, tail.Pairs...)

			return &ASTExpressionNode{
				Expr: &ASTOpList{Pairs: pairs},
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
			return &ASTExpressionNode{
				Expr: &ASTIntegerNode{Name: tok.Lexeme, Value: v},
				Type: "int",
			}
		},
	})

	// — Factor → Identifier
	g.Rules = append(g.Rules, Rule{
		LHS: "Factor",
		RHS: []Symbol{Identifier},
		Action: func(ch []ASTNode) ASTNode {
			tok := ch[0].(*ASTSimpleExpression).Token
			return &ASTExpressionNode{
				Expr: &ASTVariableNode{Token: tok},
				Type: "",
			}
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

	// … you can keep adding rules for
	//      Assignment, PrintStatement, VarDecl, If, While, For, etc.
	// in exactly the same pattern: LHS, then RHS slice mixing TokenType & string,
	// then an Action that builds one of your AST*Node structs.

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
