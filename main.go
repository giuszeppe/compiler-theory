package main

import "fmt"

func main() {
	program := `
	let c:color = 0 as color;

	for (let i:int = 0; i < 64; i = i + 1) {
		c = __random_int (1677216) as color;	
		__clear c;

		__delay 16;
	}
	`

	parser := NewParser(program)
	printVisitor := NewSemanticVisitor()
	generatorVisitor := NewGeneratorVisitor()
	grammar := NewGrammar()
	node, err := parser.Parse(grammar)
	if err != nil {
		panic(err)
	}

	node.Accept(printVisitor)
	node.Accept(generatorVisitor)
	for _, instr := range generatorVisitor.Instructions {
		// fmt.Print(fmt.Sprint(idx) + " ")
		fmt.Println(instr)
	}
	// for idx, instr := range generatorVisitor.Instructions {
	// 	fmt.Print(fmt.Sprint(idx) + " ")
	// 	fmt.Println(instr)
	// }

}
