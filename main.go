package main

func main() {
	program := `
	let x:int = 45; //this is fine
while (x < 50) {
__print MoreThan50(x); //"false" *5 since bool operator is <
x = x + 1;
}

let x:int = 45; //re-declaration in the same scope ... not allowed!!
while (MoreThan50(x)) {
__print MoreThan50(x); //"false" x5 since bool operator is <=
x = x + 1;
 }
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
