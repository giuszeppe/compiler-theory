package main

import (
	"fmt"
	"os"
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
	varDecl, ok := v.SymbolTable.Lookup(node.Token.Lexeme)
	if !ok {
		panic(ErrVariableNotDeclared(node.Token))
	}
	if _, isEpsilon := node.Offset.(*ASTEpsilon); !isEpsilon {
		offsetType := getExpressionType(node.Offset, *v.SymbolTable)
		if offsetType != "int" {
			panic(ErrInvalidOffsetType("int", offsetType, node.Token))
		}
		varDeclNode, ok := varDecl.(*ASTVarDeclNode)
		if !ok {
			panic(ErrNotVariableDeclaration(node.Token))
		}
		// check if type in variable declaration is an array
		if !strings.Contains(varDeclNode.Type, "[") {
			panic(ErrNotAnArray(node.Token))
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
		return "colour"
	case *ASTVariableNode:
		variableNode, ok := node.(*ASTVariableNode)
		if !ok {
			panic(ErrUnknownExpressionType(node))
		}
		val, ok := symbolTable.Lookup(variableNode.Token.Lexeme)
		if !ok {
			panic(ErrVariableNotDeclared(variableNode.Token))
		}
		varDeclNode, ok := val.(*ASTVarDeclNode)
		if !ok {
			panic(ErrNotVariableDeclaration(variableNode.Token))
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
			panic(ErrTypeMismatch(leftType, rightType, n.Token))
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
		val, ok := symbolTable.Lookup(n.Name.Lexeme)
		if !ok {
			panic(ErrFunctionNotDeclared(n.Name))
		}
		funcDeclNode, ok := val.(*ASTFuncDeclNode)
		if !ok {
			panic(ErrNotVariableDeclaration(funcDeclNode.Token))
		}
		formalParamsNode, _ := funcDeclNode.Params.(*ASTFormalParamsNode)
		actualParamsNode, _ := n.Params.(*ASTActualParamsNode)
		if len(actualParamsNode.Params) != len(formalParamsNode.Params) {
			panic(ErrArgumentCountMismatch(len(formalParamsNode.Params), len(actualParamsNode.Params), funcDeclNode.Token))
		}
		for i, param := range actualParamsNode.Params {
			paramType := getExpressionType(param, symbolTable)
			formParamNode := formalParamsNode.Params[i].(*ASTVarDeclNode)
			funcParamType := formParamNode.Type
			if paramType != funcParamType {
				panic(ErrTypeMismatch(funcParamType, paramType, formParamNode.Token))
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
		arrNode := node.(*ASTArrayNode)
		return arrNode.Type
	case *ASTEpsilon:
		return ""
	case *ASTBuiltinFuncNode:
		switch n.Token.Lexeme {
		case "__random_int":
			if len(n.Args) != 1 {
				panic(ErrArgumentCountMismatch(1, len(n.Args), n.Token))
			}
			argType := getExpressionType(n.Args[0], symbolTable)
			if argType != "int" {
				panic(ErrTypeMismatch("int", argType, n.Token))
			}
			return "int"
		case "__delay":
			if len(n.Args) != 1 {
				panic(ErrArgumentCountMismatch(1, len(n.Args), n.Token))
			}
			argType := getExpressionType(n.Args[0], symbolTable)
			if argType != "int" {
				panic(ErrTypeMismatch("int", argType, n.Token))
			}
			return ""
		case "__height", "__width":
			return "int"
		case "__write":
			if len(n.Args) != 3 {
				panic(ErrArgumentCountMismatch(3, len(n.Args), n.Token))
			}
			arg1Type := getExpressionType(n.Args[0], symbolTable)
			arg2Type := getExpressionType(n.Args[1], symbolTable)
			arg3Type := getExpressionType(n.Args[2], symbolTable)
			if arg1Type != "int" {
				panic(ErrTypeMismatch("int", arg1Type, n.Token))
			}
			if arg2Type != "int" {
				panic(ErrTypeMismatch("int", arg2Type, n.Token))
			}
			if arg3Type != "colour" {
				panic(ErrTypeMismatch("colour", arg3Type, n.Token))
			}
			return ""
		case "__print":
			if len(n.Args) != 1 {
				panic(ErrArgumentCountMismatch(1, len(n.Args), n.Token))
			}
		case "__write_box":
			if len(n.Args) != 5 {
				panic(ErrArgumentCountMismatch(5, len(n.Args), n.Token))
			}
			arg1Type := getExpressionType(n.Args[0], symbolTable)
			arg2Type := getExpressionType(n.Args[1], symbolTable)
			arg3Type := getExpressionType(n.Args[2], symbolTable)
			arg4Type := getExpressionType(n.Args[3], symbolTable)
			arg5Type := getExpressionType(n.Args[4], symbolTable)
			if arg1Type != "int" {
				panic(ErrTypeMismatch("int", arg1Type, n.Token))
			}
			if arg2Type != "int" {
				panic(ErrTypeMismatch("int", arg2Type, n.Token))
			}
			if arg3Type != "int" {
				panic(ErrTypeMismatch("int", arg3Type, n.Token))
			}
			if arg4Type != "int" {
				panic(ErrTypeMismatch("int", arg4Type, n.Token))
			}
			if arg5Type != "colour" {
				panic(ErrTypeMismatch("colour", arg5Type, n.Token))
			}
			return "int"
		case "__read":
			if len(n.Args) != 2 {
				panic(ErrArgumentCountMismatch(2, len(n.Args), n.Token))
			}
			arg1Type := getExpressionType(n.Args[0], symbolTable)
			arg2Type := getExpressionType(n.Args[1], symbolTable)
			if arg1Type != "int" {
				panic(ErrTypeMismatch("int", arg1Type, n.Token))
			}
			if arg2Type != "int" {
				panic(ErrTypeMismatch("int", arg2Type, n.Token))
			}
			return "colour"
		default:
			panic(ErrUnknownExpressionType(n.Token))
		}
		return ""
	default:
		panic(ErrUnknownExpressionType(node))
	}
}

func (v *SemanticVisitor) VisitAssignmentNode(node *ASTAssignmentNode) {
	val, ok := v.SymbolTable.Lookup(node.Id.Token.Lexeme)
	if !ok {
		panic(ErrVariableNotDeclared(node.Id.Token))
	}
	varDeclNode, ok := val.(*ASTVarDeclNode)
	if !ok {
		panic(ErrNotVariableDeclaration(node.Id.Token))
	}
	if _, isEpsilon := node.Id.Offset.(*ASTEpsilon); !isEpsilon {
		offsetType := getExpressionType(node.Id.Offset, *v.SymbolTable)
		if offsetType != "int" {
			panic(ErrInvalidOffsetType("int", offsetType, node.Id.Token))
		}
		if getExpressionType(node.Expr, *v.SymbolTable) != varDeclNode.Type[:strings.Index(varDeclNode.Type, "[")] {
			panic(ErrTypeMismatch(varDeclNode.Type[:strings.Index(varDeclNode.Type, "[")], getExpressionType(node.Expr, *v.SymbolTable), node.Id.Token))
		}
	} else {
		if getExpressionType(node.Expr, *v.SymbolTable) != varDeclNode.Type {
			panic(ErrTypeMismatch(varDeclNode.Type, getExpressionType(node.Expr, *v.SymbolTable), node.Id.Token))
		}
	}
	node.Expr.Accept(v)
}

func (v *SemanticVisitor) VisitVarDeclNode(node *ASTVarDeclNode) {
	if _, ok := v.SymbolTable.Lookup(node.Token.Lexeme); ok {
		panic(ErrVariableAlreadyDeclared(node.Token))
	}
	nodeType := getExpressionType(node.Expression, *v.SymbolTable)
	if nodeType != "" && nodeType != node.Type {
		panic(ErrTypeMismatch(node.Type, getExpressionType(node.Expression, *v.SymbolTable), node.Token))
	}
	v.SymbolTable.Insert(node.Token.Lexeme, node)
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

	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", r)
			os.Exit(1)
		}
	}()
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
	_, ok := v.SymbolTable.Lookup(node.Name.Lexeme)
	if !ok {
		panic(ErrFunctionNotDeclared(node.Name))
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

	// Check if type is the same
	leftType := getExpressionType(node.Left, *v.SymbolTable)
	rightType := getExpressionType(node.Right, *v.SymbolTable)
	if leftType != rightType {
		panic(ErrTypeMismatch(leftType, rightType, node.Token))
	}
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
		panic(ErrInvalidColorValue(hexValue, node.Token))
	}
	// Check if the color value is valid
	if _, err := fmt.Sscanf(hexValue, "#%x", new(int)); err != nil {
		panic(ErrInvalidColorValue(hexValue, node.Token))
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
		panic(ErrParameterAlreadyDeclared(node.Name))
	}
	v.SymbolTable.Insert(node.Name, node)
}

func (v *SemanticVisitor) VisitFuncDeclNode(node *ASTFuncDeclNode) {
	if _, ok := v.SymbolTable.Lookup(node.Token.Lexeme); ok {
		panic(ErrFunctionAlreadyDeclared(node.Token))
	}

	v.SymbolTable.Insert(node.Token.Lexeme, node)
	v.SymbolTable.Push()
	defer v.SymbolTable.Pop()
	node.Params.Accept(v)
	node.Block.Accept(v)
	// Check if the return type is valid
	funcBlock, _ := node.Block.(*ASTBlockNode)
	hasReturn := hasReturnStatement(funcBlock, v, node.ReturnType)
	if !hasReturn {
		panic(ErrFunctionMustHaveReturn(node.Token))
	}
}

func hasReturnStatement(node ASTNode, v *SemanticVisitor, expectedType string) bool {
	switch n := node.(type) {
	case *ASTReturnNode:
		returnNode := node.(*ASTReturnNode)
		if getExpressionType(returnNode.Expr, *v.SymbolTable) != expectedType {
			panic(ErrReturnTypeMismatch(expectedType, getExpressionType(returnNode.Expr, *v.SymbolTable), returnNode.Token))
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
	case *ASTWhileNode:
		return hasReturnStatement(n.Block, v, expectedType)
	case *ASTForNode:
		return hasReturnStatement(n.Block, v, expectedType)
	}

	return false
}

func (v *SemanticVisitor) VisitSimpleExpressionNode(node *ASTSimpleExpression) {
}

func (v *SemanticVisitor) VisitArrayNode(node *ASTArrayNode) {
	if node.Size < 0 {
		panic(ErrArraySizeNegative(node.Size, node.Token))
	}
	if node.Size < len(node.Items) {
		panic(ErrArraySize(node.Size, len(node.Items), node.Token))
	}
	for _, item := range node.Items {
		item.Accept(v)
		itemType := getExpressionType(item, *v.SymbolTable)
		if itemType != getArrayType(node) {
			panic(ErrTypeMismatch(getArrayType(node), itemType, node.Token))
		}
	}
}
func getArrayType(node *ASTArrayNode) string {
	return strings.Split(node.Type, "[")[0]
}
