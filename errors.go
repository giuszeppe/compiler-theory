package main

import "fmt"

func ErrVariableNotDeclared(tok Token) string {
	return fmt.Sprintf("Variable not declared: %s (at line %d, column %d)", tok.Lexeme, tok.Line, tok.Column)
}

func ErrVariableAlreadyDeclared(tok Token) string {
	return fmt.Sprintf("Variable already declared: %s (at line %d, column %d)", tok.Lexeme, tok.Line, tok.Column)
}

func ErrTypeMismatch(expected, got any, tok Token) string {
	return fmt.Sprintf("Type mismatch: expected %v, got %v (at line %d, column %d)", expected, got, tok.Line, tok.Column)
}

func ErrInvalidOffsetType(expected, got string, tok Token) string {
	return fmt.Sprintf("Invalid offset type: expected %s, got %s (at line %d, column %d)", expected, got, tok.Line, tok.Column)
}

func ErrNotVariableDeclaration(tok Token) string {
	return fmt.Sprintf("Not a variable declaration: %s (at line %d, column %d)", tok.Lexeme, tok.Line, tok.Column)
}

func ErrFunctionNotDeclared(tok Token) string {
	return fmt.Sprintf("Function not declared: %s (at line %d, column %d)", tok.Lexeme, tok.Line, tok.Column)
}

func ErrFunctionAlreadyDeclared(tok Token) string {
	return fmt.Sprintf("Function already declared: %s (at line %d, column %d)", tok.Lexeme, tok.Line, tok.Column)
}

func ErrParameterAlreadyDeclared(name string) string {
	return fmt.Sprintf("Parameter already declared: %s", name)
}

func ErrArgumentCountMismatch(expected, got int, tok Token) string {
	return fmt.Sprintf("Argument count mismatch: expected %d, got %d (at line %d, column %d)", expected, got, tok.Line, tok.Column)
}

func ErrInvalidColorValue(value string, tok Token) string {
	return fmt.Sprintf("Invalid color value: %s (at line %d, column %d)", value, tok.Line, tok.Column)
}

func ErrArraySize(size, items int, tok Token) string {
	return fmt.Sprintf("Array size must be greater than the number of items: %d < %d (at line %d, column %d)", size, items, tok.Line, tok.Column)
}

func ErrArraySizeNegative(size int, tok Token) string {
	return fmt.Sprintf("Array size must be greater than 0: %d (at line %d, column %d)", size, tok.Line, tok.Column)
}

func ErrUnknownExpressionType(t any) string {
	return fmt.Sprintf("Unknown expression type %T", t)
}

func ErrReturnTypeMismatch(expected, got any, tok Token) string {
	return fmt.Sprintf("Return type mismatch: expected %v, got %v (at line %d, column %d)", expected, got, tok.Line, tok.Column)
}

func ErrFunctionMustHaveReturn(tok Token) string {
	return fmt.Sprintf("Function must have a return statement (at line %d, column %d)", tok.Line, tok.Column)
}
