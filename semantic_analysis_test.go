package main

import "testing"

func expectPanic(t *testing.T, f func(), msg string) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic: %s", msg)
		} else {
			if r != msg {
				t.Errorf("Expected panic message: %s, got: %s", msg, r)
			}
		}
	}()
	f()
}

func TestDoubleVariableDeclaration(t *testing.T) {
	program := `let x:int = 5; let x:float = 10;
	`
	parser := NewParser(program)
	grammar := NewGrammar()
	rootAST, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}
	visitor := NewSemanticVisitor()

	expectPanic(t, func() { rootAST.Accept(visitor) }, "Variable already declared: x")
}

func TestUndeclaredVariable(t *testing.T) {
	program := `let x:int = 5; let y:int = x + z;
	`
	parser := NewParser(program)
	grammar := NewGrammar()
	rootAST, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}
	visitor := NewSemanticVisitor()

	expectPanic(t, func() { rootAST.Accept(visitor) }, "Variable not declared: z")
}

func TestValidVariableDeclaration(t *testing.T) {
	program := `let x:int = 5; let y:float = 10.0;
	`
	parser := NewParser(program)
	grammar := NewGrammar()
	rootAST, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}
	visitor := NewSemanticVisitor()
	rootAST.Accept(visitor)
}

func TestValidVariableAssignment(t *testing.T) {
	program := `let x:int = 5; x = 10;
	`
	parser := NewParser(program)
	grammar := NewGrammar()
	rootAST, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}
	visitor := NewSemanticVisitor()
	rootAST.Accept(visitor)
}

func TestInvalidVariableAssignment(t *testing.T) {
	program := `let x:int = 5; y = 10;
	`
	parser := NewParser(program)
	grammar := NewGrammar()
	rootAST, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}
	visitor := NewSemanticVisitor()

	expectPanic(t, func() { rootAST.Accept(visitor) }, "Variable not declared: y")
}

func TestValidVariableUsage(t *testing.T) {
	program := `let x:int = 5; let y:int = x + 10;
	`
	parser := NewParser(program)
	grammar := NewGrammar()
	rootAST, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}
	visitor := NewSemanticVisitor()
	rootAST.Accept(visitor)
}

func TestDoubleFuncDeclaration(t *testing.T) {
	program := `fun foo() -> int { return 1; } fun foo() -> float { return 1.0; }
	`
	parser := NewParser(program)
	grammar := NewGrammar()
	rootAST, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}
	visitor := NewSemanticVisitor()

	expectPanic(t, func() { rootAST.Accept(visitor) }, "Function already declared: foo")
}

func TestUndeclaredFunc(t *testing.T) {
	program := `fun foo() -> int { return 1; } let x:int = bar();
	`
	parser := NewParser(program)
	grammar := NewGrammar()
	rootAST, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}
	visitor := NewSemanticVisitor()

	expectPanic(t, func() { rootAST.Accept(visitor) }, "Function not declared: bar")
}
func TestValidFuncDeclaration(t *testing.T) {
	program := `fun foo() -> int { return 1; } fun bar() -> float { return 1.0; }
	`
	parser := NewParser(program)
	grammar := NewGrammar()
	rootAST, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}
	visitor := NewSemanticVisitor()
	rootAST.Accept(visitor)
}

func TestOuterVariableAreSeenByInnerScopes(t *testing.T) {
	program := `let x:int = 5; fun foo() -> int { let y:int = x + 1; return y; } let z:int = foo();
	`
	parser := NewParser(program)
	grammar := NewGrammar()
	rootAST, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}
	visitor := NewSemanticVisitor()
	rootAST.Accept(visitor)
}

func TestInnerVariableAreNotSeenByOuterScopes(t *testing.T) {
	program := `let x:int = 5; fun foo() -> int { let y:int = 10; return y; } let z:int = x + y;
	`
	parser := NewParser(program)
	grammar := NewGrammar()
	rootAST, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}
	visitor := NewSemanticVisitor()

	expectPanic(t, func() { rootAST.Accept(visitor) }, "Variable not declared: y")
}

func TestValidBlock(t *testing.T) {
	program := `let x:int = 5; { let y:int = x + 1; }
	`
	parser := NewParser(program)
	grammar := NewGrammar()
	rootAST, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}
	visitor := NewSemanticVisitor()
	rootAST.Accept(visitor)
}

func TestInvalidBlock(t *testing.T) {
	program := `let x:int = 5; { let y:int = x + 1; } { let z:int = y + 1; }
	`
	parser := NewParser(program)
	grammar := NewGrammar()
	rootAST, err := parser.Parse(grammar)
	if err != nil {
		t.Fatalf("Failed to parse program: %v", err)
	}
	visitor := NewSemanticVisitor()

	expectPanic(t, func() { rootAST.Accept(visitor) }, "Variable not declared: y")
}
