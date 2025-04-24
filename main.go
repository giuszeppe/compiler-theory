package main

func main() {
	program := `
	a = main(a,b);
	let a:int = 2;
	__print main(a);
	__print main();
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
