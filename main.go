package main

import "fmt"

func main(){
    lexer := NewLexer()
    fmt.Println(lexer)
    srcProgram := "for; if; else; while;  "
    fmt.Printf("Tokens: %v", lexer.GenerateTokens(srcProgram))
}
