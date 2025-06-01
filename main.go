package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: program <source_file>")
		os.Exit(1)
	}

	filePath := os.Args[1]
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}
	program := string(content)

	parser := NewParser(program)
	printVisitor := NewPrintNodesVisitor()
	semanticVisitor := NewSemanticVisitor()
	generatorVisitor := NewGeneratorVisitor()
	grammar := NewGrammar()
	node, err := parser.Parse(grammar)
	if err != nil {
		panic(err)
	}

	node.Accept(printVisitor)
	node.Accept(semanticVisitor)
	node.Accept(generatorVisitor)
	//for _, instr := range generatorVisitor.Instructions {
	//	fmt.Println(instr)
	//}
}
