package main

import "fmt"

func ErrVariableNotDeclared(name string) string {
	return fmt.Sprintf("Variable not declared: %s", name)
}

func ErrVariableAlreadyDeclared(name string) string {
	return fmt.Sprintf("Variable already declared: %s", name)
}

func ErrTypeMismatch(expected, got any) string {
	return fmt.Sprintf("Type mismatch: expected %v, got %v", expected, got)
}

func ErrInvalidOffsetType(expected, got string) string {
	return fmt.Sprintf("Invalid offset type: expected %s, got %s", expected, got)
}

func ErrNotVariableDeclaration(name string) string {
	return fmt.Sprintf("Not a variable declaration: %s", name)
}

func ErrFunctionNotDeclared(name string) string {
	return fmt.Sprintf("Function not declared: %s", name)
}

func ErrFunctionAlreadyDeclared(name string) string {
	return fmt.Sprintf("Function already declared: %s", name)
}

func ErrParameterAlreadyDeclared(name string) string {
	return fmt.Sprintf("Parameter already declared: %s", name)
}

func ErrArgumentCountMismatch(expected, got int) string {
	return fmt.Sprintf("Argument count mismatch: expected %d, got %d", expected, got)
}

func ErrInvalidColorValue(value string) string {
	return fmt.Sprintf("Invalid color value: %s", value)
}

func ErrArraySize(size, items int) string {
	return fmt.Sprintf("Array size must be greater than the number of items: %d < %d", size, items)
}

func ErrArraySizeNegative(size int) string {
	return fmt.Sprintf("Array size must be greater than 0: %d", size)
}

func ErrUnknownExpressionType(t any) string {
	return fmt.Sprintf("Unknown expression type %T", t)
}

func ErrReturnTypeMismatch(expected, got any) string {
	return fmt.Sprintf("Return type mismatch: expected %v, got %v", expected, got)
}

func ErrFunctionMustHaveReturn() string {
	return "Function must have a return statement"
}
