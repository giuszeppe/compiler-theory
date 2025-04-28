package main

import "fmt"

func main() {
	program := `
	for (let i:int=0;i<20;i=i+1){
	for (let j:int=0;j<10;j=j+1){
			__write j,i, 1000 * i * j;	
			__delay 100;
	}
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
	for idx, instr := range generatorVisitor.Instructions {
		fmt.Print(fmt.Sprint(idx) + " ")
		fmt.Println(instr)
	}

}
