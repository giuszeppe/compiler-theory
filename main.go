package main

func main() {
	program := `
	{ let x:int = 5; { let y:int = 10; } }
	`

	parser := NewParser(program)
	printVisitor := PrintNodesVisitor{}
	grammar := NewGrammar()
	node, err := parser.Parse(grammar)
	if err != nil {
		panic(err)
	}
	node.Accept(&printVisitor)
}
