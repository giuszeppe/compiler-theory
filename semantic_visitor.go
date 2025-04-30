package main

import (
	"fmt"
	"strings"
)

type Scope map[string]ASTNode

type SymbolTable struct {
	Scopes Stack[Scope]
}

func (st *SymbolTable) Push() {
	if !st.Scopes.IsEmpty() {
		// Copy the current scope to the new scope
		currentScope, err := st.Scopes.Peek()
		if err != nil {
			return
		}
		newScope := make(Scope)
		for k, v := range currentScope {
			newScope[k] = v
		}
		st.Scopes.Push(newScope)
		return
	}
	st.Scopes.Push(make(Scope))
}
func (st *SymbolTable) Pop() {
	st.Scopes.Pop()
}

func (st *SymbolTable) Lookup(name string) (ASTNode, bool) {
	currentScope, err := st.Scopes.Peek()
	if err != nil {
		return nil, false
	}
	if node, ok := currentScope[name]; ok {
		return node, true
	}

	return nil, false
}

func (st *SymbolTable) Insert(name string, node ASTNode) {
	currentScope, err := st.Scopes.Peek()
	if err != nil {
		return
	}
	currentScope[name] = node
}

type SemanticVisitor struct {
	SymbolTable *SymbolTable
}

func NewSemanticVisitor() *SemanticVisitor {
	return &SemanticVisitor{
		SymbolTable: &SymbolTable{
			Scopes: Stack[Scope]{},
		},
	}
}
func (v *SemanticVisitor) VisitIntegerNode(node *ASTIntegerNode) {
	// Do nothing
}
func (v *SemanticVisitor) VisitVariableNode(node *ASTVariableNode) {
	_, ok := v.SymbolTable.Lookup(node.Token.Lexeme)
	if !ok {
		panic("Variable not declared: " + node.Token.Lexeme)
	}
	if _, isEpsilon := node.Offset.(*ASTEpsilon); !isEpsilon {
		offsetType := getExpressionType(node.Offset, *v.SymbolTable)
		if offsetType != "int" {
			panic("Invalid offset type: expected int, got " + offsetType)
		}
	}

}

func getExpressionType(node ASTNode, symbolTable SymbolTable) string {
	switch n := node.(type) {
	case *ASTIntegerNode:
		return "int"
	case *ASTFloatNode:
		return "float"
	case *ASTBooleanNode:
		return "bool"
	case *ASTColorNode:
		return "color"
	case *ASTVariableNode:
		variableNode, ok := node.(*ASTVariableNode)
		if !ok {
			panic("Not a variable node: ")
		}
		// Check if the variable is declared
		val, ok := symbolTable.Lookup(variableNode.Token.Lexeme)
		if !ok {
			panic("Variable not declared: " + variableNode.Token.Lexeme)
		}
		varDeclNode, ok := val.(*ASTVarDeclNode)
		if !ok {
			panic("Not a variable declaration: ")
		}
		if _, isEpsilon := variableNode.Offset.(*ASTEpsilon); !isEpsilon {
			itemType := varDeclNode.Type[:strings.Index(varDeclNode.Type, "[")]
			return itemType
		}

		return varDeclNode.Type

	case *ASTBinaryOpNode:
		leftType := getExpressionType(n.Left, symbolTable)
		rightType := getExpressionType(n.Right, symbolTable)
		if leftType != rightType {
			panic("Type mismatch: expected " + leftType + ", got " + rightType)
		}
		binaryOpNode := node.(*ASTBinaryOpNode)
		if binaryOpNode.Operator == "<" || binaryOpNode.Operator == ">" || binaryOpNode.Operator == "<=" || binaryOpNode.Operator == ">=" || binaryOpNode.Operator == "==" {
			return "bool"
		}
		return leftType
	case *ASTUnaryOpNode:
		return getExpressionType(n.Operand, symbolTable)
	case *ASTAssignmentNode:
		return getExpressionType(n.Expr, symbolTable)
	case *ASTFuncCallNode:
		// Check if the function is declared
		val, ok := symbolTable.Lookup(n.Name)
		if !ok {
			panic("Function not declared: " + n.Name)
		}
		funcDeclNode, ok := val.(*ASTFuncDeclNode)
		if !ok {
			panic("Not a function declaration: " + n.Name)
		}
		formalParamsNode, _ := funcDeclNode.Params.(*ASTFormalParamsNode)
		// Check if the types of the arguments match
		actualParamsNode, _ := n.Params.(*ASTActualParamsNode)

		// Check if the number of arguments matches
		if len(actualParamsNode.Params) != len(formalParamsNode.Params) {
			panic("Argument count mismatch: expected " + fmt.Sprint(len(formalParamsNode.Params)) + ", got " + fmt.Sprint(len(actualParamsNode.Params)))
		}
		for i, param := range actualParamsNode.Params {
			paramType := getExpressionType(param, symbolTable)
			formParamNode := formalParamsNode.Params[i].(*ASTVarDeclNode)
			funcParamType := formParamNode.Type

			if paramType != funcParamType {
				panic("Type mismatch: expected " + funcParamType + ", got " + paramType)
			}
		}
		return funcDeclNode.ReturnType
	case *ASTReturnNode:
		return getExpressionType(n.Expr, symbolTable)
	case *ASTExpressionNode:
		return getExpressionType(n.Expr, symbolTable)
	case *ASTTypeCastNode:
		n, _ = node.(*ASTTypeCastNode)

		return n.Type
	case *ASTArrayNode:
		// Check if the array is declared
		arrNode := node.(*ASTArrayNode)
		return arrNode.Type
	case *ASTEpsilon:
		return ""
	case *ASTBuiltinFuncNode:
		switch n.Name {
		case "__random_int":
			if len(n.Args) != 1 {
				panic("Invalid number of arguments for __random_int: expected 1, got " + fmt.Sprint(len(n.Args)))
			}
			argType := getExpressionType(n.Args[0], symbolTable)
			if argType != "int" {
				panic("Invalid argument type for __random_int: expected int, got " + argType)
			}
			return "int"
		case "__delay":
			if len(n.Args) != 1 {
				panic("Invalid number of arguments for __delay: expected 1, got " + fmt.Sprint(len(n.Args)))
			}
			argType := getExpressionType(n.Args[0], symbolTable)
			if argType != "int" {
				panic("Invalid argument type for __delay: expected int, got " + argType)
			}
			return ""
		case "__height":
			return "int"
		case "__width":
			return "int"
		case "__write":
			if len(n.Args) != 3 {
				panic("Invalid number of arguments for __write: expected 3, got " + fmt.Sprint(len(n.Args)))
			}
			arg1Type := getExpressionType(n.Args[0], symbolTable)
			arg2Type := getExpressionType(n.Args[1], symbolTable)
			arg3Type := getExpressionType(n.Args[2], symbolTable)
			if arg1Type != "int" {
				panic("Invalid argument type for __write: expected int, got " + arg1Type)
			}
			if arg2Type != "int" {
				panic("Invalid argument type for __write: expected int, got " + arg2Type)
			}
			if arg3Type != "color" {
				panic("Invalid argument type for __write: expected color, got " + arg3Type)
			}
			return ""
		case "__print":
			if len(n.Args) != 1 {
				panic("Invalid number of arguments for __print: expected 1, got " + fmt.Sprint(len(n.Args)))
			}
		case "__write_box":
			if len(n.Args) != 5 {
				panic("Invalid number of arguments for __random_int: expected 5 got " + fmt.Sprint(len(n.Args)))
			}
			arg1Type := getExpressionType(n.Args[0], symbolTable)
			arg2Type := getExpressionType(n.Args[1], symbolTable)
			arg3Type := getExpressionType(n.Args[2], symbolTable)
			arg4Type := getExpressionType(n.Args[3], symbolTable)
			arg5Type := getExpressionType(n.Args[4], symbolTable)
			if arg1Type != "int" {
				panic("Invalid argument type for __write_box: expected int, got " + arg1Type)
			}
			if arg2Type != "int" {
				panic("Invalid argument type for __write_box: expected int, got " + arg2Type)
			}
			if arg3Type != "int" {
				panic("Invalid argument type for __write_box: expected int, got " + arg3Type)
			}
			if arg4Type != "int" {
				panic("Invalid argument type for __write_box: expected int, got " + arg4Type)
			}
			if arg5Type != "color" {
				panic("Invalid argument type for __write_box: expected color, got " + arg5Type)
			}
			return "int"
		case "__read":
			if len(n.Args) != 2 {
				panic("Invalid number of arguments for __read: expected 2, got " + fmt.Sprint(len(n.Args)))
			}
			arg1Type := getExpressionType(n.Args[0], symbolTable)
			arg2Type := getExpressionType(n.Args[1], symbolTable)
			if arg1Type != "int" {
				panic("Invalid argument type for __read: expected int, got " + arg1Type)
			}
			if arg2Type != "int" {
				panic("Invalid argument type for __read: expected int, got " + arg2Type)
			}
			return "int"
		default:
			panic("Unknown builtin function: " + n.Name)
		}
		return ""

	default:
		panic("Unknown expression type" + fmt.Sprintf("%T", node))
	}
}

func (v *SemanticVisitor) VisitAssignmentNode(node *ASTAssignmentNode) {
	// Check if the variable is declared
	val, ok := v.SymbolTable.Lookup(node.Id.Token.Lexeme)
	if !ok {
		panic("Variable not declared: " + node.Id.Token.Lexeme)
	}

	varDeclNode, ok := val.(*ASTVarDeclNode)
	if !ok {
		panic("Not a variable declaration: " + node.Id.Token.Lexeme)
	}

	if _, isEpsilon := node.Id.Offset.(*ASTEpsilon); !isEpsilon {
		offsetType := getExpressionType(node.Id.Offset, *v.SymbolTable)
		if offsetType != "int" {
			panic("Invalid offset type: expected int, got " + offsetType)
		}

		if getExpressionType(node.Expr, *v.SymbolTable) != varDeclNode.Type[:strings.Index(varDeclNode.Type, "[")] {
			panic(fmt.Sprintf("Type mismatch: expected %v, got %v", varDeclNode.Type[:strings.Index(varDeclNode.Type, "[")], getExpressionType(node.Expr, *v.SymbolTable)))
		}
	} else {
		// Check if the type of the expression matches the variable type
		if getExpressionType(node.Expr, *v.SymbolTable) != varDeclNode.Type {
			panic(fmt.Sprintf("Type mismatch: expected %v, got %v", varDeclNode.Type, getExpressionType(node.Expr, *v.SymbolTable)))
		}

	}
	node.Expr.Accept(v)
}

func (v *SemanticVisitor) VisitVarDeclNode(node *ASTVarDeclNode) {
	// Check if the variable is already declared in the current scope
	if _, ok := v.SymbolTable.Lookup(node.Name); ok {
		panic("Variable already declared: " + node.Name)
	}

	nodeType := getExpressionType(node.Expression, *v.SymbolTable)
	// Check if the type is valid
	if nodeType != "" && nodeType != node.Type {
		panic(fmt.Sprintf("Type mismatch: expected %v, got %v", node.Type, getExpressionType(node.Expression, *v.SymbolTable)))
	}

	v.SymbolTable.Insert(node.Name, node)
	node.Expression.Accept(v)
}
func (v *SemanticVisitor) VisitBlockNode(node *ASTBlockNode) {
	for _, stmt := range node.Stmts {
		pushAndPopIfBlock(v, stmt)
	}
}
func (v *SemanticVisitor) VisitTypeNode(node *ASTTypeNode) {
	// Do nothing
}
func (v *SemanticVisitor) VisitFunctionNode(node *ASTFuncDeclNode) {
	// Check if the function is already declared in the current scope
}

func (v *SemanticVisitor) VisitProgramNode(node *ASTProgramNode) {
	// Visit the block node
	pushAndPopIfBlock(v, &node.Block)
}
func (v *SemanticVisitor) VisitIfNode(node *ASTIfNode) {
	// Visit the condition and the block
	node.Condition.Accept(v)
	pushAndPopIfBlock(v, node.ThenBlock)
	if node.ElseBlock != nil {
		pushAndPopIfBlock(v, node.ElseBlock)
	}
}

func pushAndPopIfBlock(v *SemanticVisitor, block ASTNode) {
	// Check if the block is a block node
	if _, ok := block.(*ASTBlockNode); ok {
		v.SymbolTable.Push()
		block.Accept(v)
		v.SymbolTable.Pop()
	} else {
		block.Accept(v)
	}

}

func (v *SemanticVisitor) VisitWhileNode(node *ASTWhileNode) {
	// Visit the condition and the block
	node.Condition.Accept(v)
	pushAndPopIfBlock(v, node.Block)
}

func (v *SemanticVisitor) VisitForNode(node *ASTForNode) {
	// Visit the initialization, condition, and block
	v.SymbolTable.Push()
	node.VarDecl.Accept(v)
	node.Condition.Accept(v)
	node.Increment.Accept(v)
	node.Block.Accept(v)
	v.SymbolTable.Pop()

}

func (v *SemanticVisitor) VisitTypeCastNode(node *ASTTypeCastNode) {
	// Visit the expression
	node.Expr.Accept(v)
}

func (v *SemanticVisitor) VisitFormalParamsNode(node *ASTFormalParamsNode) {
	// Visit each parameter
	for _, param := range node.Params {
		param.Accept(v)
	}
}
func (v *SemanticVisitor) VisitEpsilon(node *ASTEpsilon) {
	// Do nothing
}
func (v *SemanticVisitor) VisitActualParamNode(node *ASTActualParamNode) {
	// Visit the expression
	node.Value.Accept(v)
}

func (v *SemanticVisitor) VisitFuncCallNode(node *ASTFuncCallNode) {
	// Check if the function is declared
	_, ok := v.SymbolTable.Lookup(node.Name)
	if !ok {
		panic("Function not declared: " + node.Name)
	}

	node.Params.Accept(v)
}

func (v *SemanticVisitor) VisitPrintNode(node *ASTPrintNode) {
	// Visit the expression
	node.Expr.Accept(v)
}
func (v *SemanticVisitor) VisitBinaryOpNode(node *ASTBinaryOpNode) {
	// Visit the left and right operands
	node.Left.Accept(v)
	node.Right.Accept(v)
}
func (v *SemanticVisitor) VisitUnaryOpNode(node *ASTUnaryOpNode) {
	// Visit the operand
	node.Operand.Accept(v)
}

func (v *SemanticVisitor) VisitBooleanNode(node *ASTBooleanNode) {
	// Do nothing
}
func (v *SemanticVisitor) VisitColorNode(node *ASTColorNode) {
	hexValue := node.Value
	if (len(hexValue) != 7 && len(hexValue) != 4) || hexValue[0] != '#' {
		panic("Invalid color value: " + hexValue)
	}
	// Check if the color value is valid
	if _, err := fmt.Sscanf(hexValue, "#%x", new(int)); err != nil {
		panic("Invalid color value: " + hexValue)
	}
	// Do nothing
}
func (v *SemanticVisitor) VisitBuiltinFuncNode(node *ASTBuiltinFuncNode) {
	// Visit the arguments
	for _, arg := range node.Args {
		arg.Accept(v)
	}
}
func (v *SemanticVisitor) VisitReturnNode(node *ASTReturnNode) {
	// Visit the expression
	node.Expr.Accept(v)
}

func (v *SemanticVisitor) VisitActualParamsNode(node *ASTActualParamsNode) {
	// Visit each actual parameter
	for _, param := range node.Params {
		param.Accept(v)
	}
}

func (v *SemanticVisitor) VisitExpressionNode(node *ASTExpressionNode) {
	// Visit the expression
	node.Expr.Accept(v)
}

func (v *SemanticVisitor) VisitFloatNode(node *ASTFloatNode) {
	// Do nothing
}
func (v *SemanticVisitor) VisitFormalParamNode(node *ASTFormalParamNode) {
	// Check if the parameter is already declared in the current scope
	if _, ok := v.SymbolTable.Lookup(node.Name); ok {
		panic("Parameter already declared: " + node.Name)
	}
	v.SymbolTable.Insert(node.Name, node)
}

func (v *SemanticVisitor) VisitFuncDeclNode(node *ASTFuncDeclNode) {
	if _, ok := v.SymbolTable.Lookup(node.Name); ok {
		panic("Function already declared: " + node.Name)
	}

	v.SymbolTable.Insert(node.Name, node)
	v.SymbolTable.Push()
	defer v.SymbolTable.Pop()
	node.Params.Accept(v)
	node.Block.Accept(v)
	// Check if the return type is valid
	funcBlock, _ := node.Block.(*ASTBlockNode)
	hasReturn := hasReturnStatement(funcBlock, v, node.ReturnType)
	if !hasReturn {
		panic("Function must have a return statement")
	}
}

func hasReturnStatement(node ASTNode, v *SemanticVisitor, expectedType string) bool {
	switch n := node.(type) {
	case *ASTReturnNode:
		returnNode := node.(*ASTReturnNode)
		if getExpressionType(returnNode.Expr, *v.SymbolTable) != expectedType {
			panic(fmt.Sprintf("Return type mismatch: expected %v, got %v", expectedType, getExpressionType(returnNode.Expr, *v.SymbolTable)))
		}
		return true
	case *ASTBlockNode:
		for _, stmt := range n.Stmts {
			if hasReturnStatement(stmt, v, expectedType) {
				return true
			}
		}
	case *ASTIfNode:
		if hasReturnStatement(n.ThenBlock, v, expectedType) {
			if n.ElseBlock != nil {
				return hasReturnStatement(n.ElseBlock, v, expectedType)
			}
			return true
		}
	}
	return false
}

func (v *SemanticVisitor) VisitSimpleExpressionNode(node *ASTSimpleExpression) {
}

func (v *SemanticVisitor) VisitArrayNode(node *ASTArrayNode) {
	// Check if the array is already declared in the current scope
	if node.Size < 0 {
		panic("Array size must be greater than 0: " + fmt.Sprint(node.Size))
	}
	if node.Size < len(node.Items) {
		panic("Array size must be greater than the number of items: " + fmt.Sprint(node.Size) + " < " + fmt.Sprint(len(node.Items)))
	}

	// Check if the types of the items match the array type
	for _, item := range node.Items {
		item.Accept(v)
		itemType := getExpressionType(item, *v.SymbolTable)

		if itemType != getArrayType(node) {
			panic(fmt.Sprintf("Array item type mismatch: expected %v, got %v", getArrayType(node), itemType))
		}
	}
}
func getArrayType(node *ASTArrayNode) string {
	return strings.Split(node.Type, "[")[0]
}
