package main

import "fmt"

func main() {
	program := "let x:int = 2+2+3+4+5+6 as float;"
    lex := NewLexer()
    src := lex.GenerateTokens(program)
    fmt.Println("TOKENS", src)
	parser := NewParser(src)
    printVisitor := PrintNodesVisitor{}
	parser.Parse()[0].Accept(&printVisitor)
}

// func main() {
// 	lexer := NewLexer()
// 	source := `let x:int=3;
//     x += 3;
//
//     `
// 	tokens := lexer.GenerateTokens(source)
// 	for _, tok := range tokens {
// 		fmt.Printf("Token: %-10v Lexeme: %q\n", tok.Type.String(), tok.Lexeme)
// 	}
// }
