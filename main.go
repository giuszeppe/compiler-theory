package main

func main() {
	program := `fun foo(x:int, y:int) -> int { return x + y; } let z:int = foo(5);`

	parser := NewParser(program)
	printVisitor := NewSemanticVisitor()
	grammar := NewGrammar()
	node, err := parser.Parse(grammar)
	if err != nil {
		panic(err)
	}

	// Va rivisto un po' tutto perche' dare error quando nella symbol table non trovo una variable declaration e' molto errato, potrebbe difatti essere una funzione
	node.Accept(printVisitor)
}
