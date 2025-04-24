package main

import (
	"fmt"
	"testing"
)

func assertASTNodeEqual(t *testing.T, expected, actual ASTNode) {
	if expected == nil && actual == nil {
		return
	}
	if expected == nil || actual == nil {
		t.Fatalf("AST nodes are not equal: expected %v, got %v", expected, actual)
	}
	if fmt.Sprintf("%T", expected) != fmt.Sprintf("%T", actual) {
		t.Fatalf("AST node types are not equal: expected %T, got %T", expected, actual)
	}

	switch e := expected.(type) {
	case *ASTProgramNode:
		a := actual.(*ASTProgramNode)
		assertASTNodeEqual(t, &e.Block, &a.Block)
	case *ASTBlockNode:
		a := actual.(*ASTBlockNode)
		if len(e.Stmts) != len(a.Stmts) {
			t.Fatalf("AST block statements length mismatch: expected %d, got %d", len(e.Stmts), len(a.Stmts))
		}
		for i := range e.Stmts {
			assertASTNodeEqual(t, e.Stmts[i], a.Stmts[i])
		}
	case *ASTAssignmentNode:
		a := actual.(*ASTAssignmentNode)
		assertASTNodeEqual(t, &e.Id, &a.Id)
		assertASTNodeEqual(t, e.Expr, a.Expr)
	case *ASTVariableNode:
		a := actual.(*ASTVariableNode)
		if e.Token != a.Token {
			t.Fatalf("AST variable tokens are not equal: expected %v, got %v", e.Token, a.Token)
		}
	case *ASTIntegerNode:
		a := actual.(*ASTIntegerNode)
		if e.Value != a.Value {
			t.Fatalf("AST integer values are not equal: expected %d, got %d", e.Value, a.Value)
		}
	case *ASTVarDeclNode:
		a := actual.(*ASTVarDeclNode)
		if e.Name != a.Name {
			t.Fatalf("AST variable declaration names are not equal: expected %s, got %s", e.Name, a.Name)
		}
		if e.Type != a.Type {
			t.Fatalf("AST variable declaration types are not equal: expected %s, got %s", e.Type, a.Type)
		}
		assertASTNodeEqual(t, e.Expression, a.Expression)
	case *ASTTypeNode:
		a := actual.(*ASTTypeNode)
		if e.Name != a.Name {
			t.Fatalf("AST type names are not equal: expected %s, got %s", e.Name, a.Name)
		}
	case *ASTFuncDeclNode:
		a := actual.(*ASTFuncDeclNode)
		if e.Name != a.Name {
			t.Fatalf("AST function names are not equal: expected %s, got %s", e.Name, a.Name)
		}
		if e.ReturnType != a.ReturnType {
			t.Fatalf("AST function return types are not equal: expected %s, got %s", e.ReturnType, a.ReturnType)
		}

		assertASTNodeEqual(t, e.Params, a.Params)
		assertASTNodeEqual(t, e.Block, a.Block)
	case *ASTIfNode:
		a := actual.(*ASTIfNode)
		assertASTNodeEqual(t, e.Condition, a.Condition)
		assertASTNodeEqual(t, e.ThenBlock, a.ThenBlock)
		if e.ElseBlock != nil && a.ElseBlock != nil {
			assertASTNodeEqual(t, e.ElseBlock, a.ElseBlock)
		} else if e.ElseBlock != nil || a.ElseBlock != nil {
			t.Fatalf("AST if nodes else body mismatch: expected %v, got %v", e.ElseBlock, a.ElseBlock)
		}
	case *ASTFormalParamsNode:
		a := actual.(*ASTFormalParamsNode)
		if len(e.Params) != len(a.Params) {
			t.Fatalf("AST formal parameters length mismatch: expected %d, got %d", len(e.Params), len(a.Params))
		}
		for i := range e.Params {
			assertASTNodeEqual(t, e.Params[i], a.Params[i])
		}
	case *ASTActualParamsNode:
		a := actual.(*ASTActualParamsNode)
		if len(e.Params) != len(a.Params) {
			t.Fatalf("AST actual parameters length mismatch: expected %d, got %d", len(e.Params), len(a.Params))
		}
		for i := range e.Params {
			assertASTNodeEqual(t, e.Params[i], a.Params[i])
		}
	case *ASTFuncCallNode:
		a := actual.(*ASTFuncCallNode)
		if e.Name != a.Name {
			t.Fatalf("AST function call names are not equal: expected %s, got %s", e.Name, a.Name)
		}
		assertASTNodeEqual(t, e.Params, a.Params)
	case *ASTBinaryOpNode:
		a := actual.(*ASTBinaryOpNode)
		assertASTNodeEqual(t, e.Left, a.Left)
		if e.Operator != a.Operator {
			t.Fatalf("AST binary operator mismatch: expected %s, got %s", e.Operator, a.Operator)
		}
		assertASTNodeEqual(t, e.Right, a.Right)
	case *ASTFormalParamNode:
		a := actual.(*ASTFormalParamNode)
		if e.Name != a.Name {
			t.Fatalf("AST formal parameter names are not equal: expected %s, got %s", e.Name, a.Name)
		}
		if e.Type != a.Type {
			t.Fatalf("AST formal parameter types are not equal: expected %s, got %s", e.Type, a.Type)
		}
	case *ASTActualParamNode:
		a := actual.(*ASTActualParamNode)
		assertASTNodeEqual(t, e.Value, a.Value)
	case *ASTWhileNode:
		a := actual.(*ASTWhileNode)
		assertASTNodeEqual(t, e.Condition, a.Condition)
		assertASTNodeEqual(t, e.Block, a.Block)
	case *ASTUnaryOpNode:
		a := actual.(*ASTUnaryOpNode)
		if e.Operator != a.Operator {
			t.Fatalf("AST unary operator mismatch: expected %s, got %s", e.Operator, a.Operator)
		}
		assertASTNodeEqual(t, e.Operand, a.Operand)

	default:
		t.Fatalf("Unsupported AST node type: %T", expected)
	}
}
func TestParsingAssignment(t *testing.T) {
	program := "x = 2;"
	parser := NewParser(program)
	grammar := NewGrammar()
	node, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}

	expectedAST := &ASTProgramNode{
		Block: ASTBlockNode{Stmts: []ASTNode{
			&ASTAssignmentNode{
				Id:   ASTVariableNode{Token: Token{Type: Identifier, Lexeme: "x"}},
				Expr: &ASTIntegerNode{Value: 2},
			},
		}},
	}

	assertASTNodeEqual(t, expectedAST, node)
}

func TestParsingIntVariableDeclaration(t *testing.T) {
	program := "let x:int = 2;"
	parser := NewParser(program)
	grammar := NewGrammar()
	node, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}

	expectedAST := &ASTProgramNode{
		Block: ASTBlockNode{Stmts: []ASTNode{
			&ASTVarDeclNode{
				Name:       "x",
				Type:       "int",
				Expression: &ASTIntegerNode{Value: 2},
			},
		}},
	}

	assertASTNodeEqual(t, expectedAST, node)
}

func TestParsingFunctionDeclaration(t *testing.T) {
	program := "fun main(a:int, b:int) -> int { a = a + b; }"
	parser := NewParser(program)
	grammar := NewGrammar()
	node, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}

	expectedAST := &ASTProgramNode{
		Block: ASTBlockNode{Stmts: []ASTNode{
			&ASTFuncDeclNode{
				Name:       "main",
				ReturnType: "int",
				Params: &ASTFormalParamsNode{
					Params: []ASTNode{
						&ASTFormalParamNode{Name: "a", Type: "int"},
						&ASTFormalParamNode{Name: "b", Type: "int"},
					},
				},
				Block: &ASTBlockNode{
					Stmts: []ASTNode{
						&ASTAssignmentNode{
							Id: ASTVariableNode{Token: Token{Identifier, "a"}},

							Expr: &ASTBinaryOpNode{
								Left:     &ASTVariableNode{Token: Token{Identifier, "a"}},
								Operator: "+",
								Right:    &ASTVariableNode{Token: Token{Identifier, "b"}},
							},
						},
					},
				},
			},
		}},
	}

	assertASTNodeEqual(t, expectedAST, node)
}
func TestParsingIfStatement(t *testing.T) {
	program := "if (x > 0) { y = 1; } else { y = 1; }"
	parser := NewParser(program)
	grammar := NewGrammar()
	node, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}

	expectedAST := &ASTProgramNode{
		Block: ASTBlockNode{Stmts: []ASTNode{
			&ASTIfNode{
				Condition: &ASTBinaryOpNode{
					Left:     &ASTVariableNode{Token: Token{Type: Identifier, Lexeme: "x"}},
					Operator: ">",
					Right:    &ASTIntegerNode{Value: 0},
				},
				ThenBlock: &ASTBlockNode{Stmts: []ASTNode{
					&ASTAssignmentNode{
						Id:   ASTVariableNode{Token: Token{Type: Identifier, Lexeme: "y"}},
						Expr: &ASTIntegerNode{Value: 1},
					},
				}},
				ElseBlock: &ASTBlockNode{Stmts: []ASTNode{
					&ASTAssignmentNode{
						Id:   ASTVariableNode{Token: Token{Type: Identifier, Lexeme: "y"}},
						Expr: &ASTIntegerNode{Value: 1},
					},
				}},
			},
		}},
	}

	assertASTNodeEqual(t, expectedAST, node)
}

func TestParsingWhileLoop(t *testing.T) {
	program := "while x < 10 { x = x + 1; }"
	parser := NewParser(program)
	grammar := NewGrammar()
	node, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}

	expectedAST := &ASTProgramNode{
		Block: ASTBlockNode{Stmts: []ASTNode{
			&ASTWhileNode{
				Condition: &ASTBinaryOpNode{
					Left:     &ASTVariableNode{Token: Token{Type: Identifier, Lexeme: "x"}},
					Operator: "<",
					Right:    &ASTIntegerNode{Value: 10},
				},
				Block: &ASTBlockNode{Stmts: []ASTNode{
					&ASTAssignmentNode{
						Id: ASTVariableNode{Token: Token{Type: Identifier, Lexeme: "x"}},
						Expr: &ASTBinaryOpNode{
							Left:     &ASTVariableNode{Token: Token{Type: Identifier, Lexeme: "x"}},
							Operator: "+",
							Right:    &ASTIntegerNode{Value: 1},
						},
					},
				}},
			},
		}},
	}

	assertASTNodeEqual(t, expectedAST, node)
}

func TestParsingNestedBlocks(t *testing.T) {
	program := "{ let x:int = 5; { let y:int = 10; } }"
	parser := NewParser(program)
	grammar := NewGrammar()
	node, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}

	expectedAST := &ASTProgramNode{
		Block: ASTBlockNode{Stmts: []ASTNode{
			&ASTVarDeclNode{
				Name:       "x",
				Type:       "int",
				Expression: &ASTIntegerNode{Value: 5},
			},
			&ASTBlockNode{Stmts: []ASTNode{
				&ASTVarDeclNode{
					Name:       "y",
					Type:       "int",
					Expression: &ASTIntegerNode{Value: 10},
				},
			}},
		}},
	}

	assertASTNodeEqual(t, expectedAST, node)
}

func TestParsingFunctionCall(t *testing.T) {
	program := "result = add(3, 4);"
	parser := NewParser(program)
	grammar := NewGrammar()
	node, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}

	expectedAST := &ASTProgramNode{
		Block: ASTBlockNode{Stmts: []ASTNode{
			&ASTAssignmentNode{
				Id: ASTVariableNode{Token: Token{Type: Identifier, Lexeme: "result"}},
				Expr: &ASTFuncCallNode{
					Name: "add",
					Params: &ASTActualParamsNode{
						Params: []ASTNode{
							&ASTIntegerNode{Value: 3},
							&ASTIntegerNode{Value: 4},
						},
					},
				},
			},
		}},
	}

	assertASTNodeEqual(t, expectedAST, node)
}
func TestParsingInvalidInput(t *testing.T) {
	program := "x +"
	parser := NewParser(program)
	grammar := NewGrammar()
	_, err := parser.Parse(grammar)
	if err == nil {
		t.Fatalf("Expected parsing to fail for invalid input, but it succeeded")
	}
}

func TestParsingUnaryInput(t *testing.T) {
	program := "a = -x;"
	parser := NewParser(program)
	grammar := NewGrammar()
	node, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}

	expectedAST := &ASTProgramNode{
		Block: ASTBlockNode{Stmts: []ASTNode{
			&ASTAssignmentNode{
				Id: ASTVariableNode{Token: Token{Type: Identifier, Lexeme: "a"}},
				Expr: &ASTUnaryOpNode{
					Operator: "-",
					Operand:  &ASTVariableNode{Token: Token{Type: Identifier, Lexeme: "x"}},
				},
			},
		}},
	}

	assertASTNodeEqual(t, expectedAST, node)
}

func TestParsingUnaryInputOnFunctionCall(t *testing.T) {
	program := "a = -add(3, 4);"
	parser := NewParser(program)
	grammar := NewGrammar()
	node, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}

	expectedAST := &ASTProgramNode{
		Block: ASTBlockNode{Stmts: []ASTNode{
			&ASTAssignmentNode{
				Id: ASTVariableNode{Token: Token{Type: Identifier, Lexeme: "a"}},
				Expr: &ASTUnaryOpNode{
					Operator: "-",
					Operand: &ASTFuncCallNode{
						Name: "add",
						Params: &ASTActualParamsNode{
							Params: []ASTNode{
								&ASTIntegerNode{Value: 3},
								&ASTIntegerNode{Value: 4},
							},
						},
					},
				},
			},
		}},
	}

	assertASTNodeEqual(t, expectedAST, node)
}

func TestParsingDoubleUnaryInput(t *testing.T) {
	program := "a = --x;"
	parser := NewParser(program)
	grammar := NewGrammar()
	node, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}

	expectedAST := &ASTProgramNode{
		Block: ASTBlockNode{Stmts: []ASTNode{
			&ASTAssignmentNode{
				Id: ASTVariableNode{Token: Token{Type: Identifier, Lexeme: "a"}},
				Expr: &ASTUnaryOpNode{
					Operator: "-",
					Operand: &ASTUnaryOpNode{
						Operator: "-",
						Operand:  &ASTVariableNode{Token: Token{Type: Identifier, Lexeme: "x"}},
					},
				},
			},
		}},
	}

	assertASTNodeEqual(t, expectedAST, node)
}

func TestParsingNegativeInteger(t *testing.T) {
	program := "a = -5;"
	parser := NewParser(program)
	grammar := NewGrammar()
	node, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}

	expectedAST := &ASTProgramNode{
		Block: ASTBlockNode{Stmts: []ASTNode{
			&ASTAssignmentNode{
				Id: ASTVariableNode{Token: Token{Type: Identifier, Lexeme: "a"}},
				Expr: &ASTUnaryOpNode{
					Operator: "-",
					Operand:  &ASTIntegerNode{Value: 5},
				},
			},
		}},
	}

	assertASTNodeEqual(t, expectedAST, node)
}

func TestParsingMultiplicativeOperators(t *testing.T) {
	program := "a = 3 * 4 / 2 and 5;"
	parser := NewParser(program)
	grammar := NewGrammar()
	node, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}

	expectedAST := &ASTProgramNode{
		Block: ASTBlockNode{Stmts: []ASTNode{
			&ASTAssignmentNode{
				Id: ASTVariableNode{Token: Token{Type: Identifier, Lexeme: "a"}},
				Expr: &ASTBinaryOpNode{
					Left: &ASTBinaryOpNode{
						Left: &ASTBinaryOpNode{
							Left:     &ASTIntegerNode{Value: 3},
							Operator: "*",
							Right:    &ASTIntegerNode{Value: 4},
						},
						Operator: "/",
						Right:    &ASTIntegerNode{Value: 2},
					},
					Operator: "and",
					Right:    &ASTIntegerNode{Value: 5},
				},
			},
		}},
	}

	assertASTNodeEqual(t, expectedAST, node)
}

func TestParsingAdditiveOperators(t *testing.T) {
	program := "a = 3 + 4 - 2 or 5;"
	parser := NewParser(program)
	grammar := NewGrammar()
	node, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}

	expectedAST := &ASTProgramNode{
		Block: ASTBlockNode{Stmts: []ASTNode{
			&ASTAssignmentNode{
				Id: ASTVariableNode{Token: Token{Type: Identifier, Lexeme: "a"}},
				Expr: &ASTBinaryOpNode{

					Left: &ASTBinaryOpNode{
						Left: &ASTBinaryOpNode{
							Left:     &ASTIntegerNode{Value: 3},
							Operator: "+",
							Right:    &ASTIntegerNode{Value: 4},
						},
						Operator: "-",
						Right:    &ASTIntegerNode{Value: 2},
					},
					Operator: "or",
					Right:    &ASTIntegerNode{Value: 5},
				},
			},
		}},
	}

	assertASTNodeEqual(t, expectedAST, node)
}

func TestParsingUnaryOperators(t *testing.T) {
	program := "a = not x;"
	parser := NewParser(program)
	grammar := NewGrammar()
	node, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}

	expectedAST := &ASTProgramNode{
		Block: ASTBlockNode{Stmts: []ASTNode{
			&ASTAssignmentNode{
				Id: ASTVariableNode{Token: Token{Type: Identifier, Lexeme: "a"}},
				Expr: &ASTUnaryOpNode{
					Operator: "not",
					Operand:  &ASTVariableNode{Token: Token{Type: Identifier, Lexeme: "x"}},
				},
			},
		}},
	}

	assertASTNodeEqual(t, expectedAST, node)
}

func TestParsingEmptyProgram(t *testing.T) {
	program := ""
	parser := NewParser(program)
	grammar := NewGrammar()
	_, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Expected parsing to succeed for empty program, but it failed: %v", err)
	}
}
