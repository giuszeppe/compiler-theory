package main

import "fmt"

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
		return varDeclNode.Type

	case *ASTBinaryOpNode:
		leftType := getExpressionType(n.Left, symbolTable)
		rightType := getExpressionType(n.Right, symbolTable)
		if leftType != rightType {
			panic("Type mismatch: expected " + leftType + ", got " + rightType)
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
	case *ASTEpsilon:
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

	// Check if the type of the expression matches the variable type
	if getExpressionType(node.Expr, *v.SymbolTable) != varDeclNode.Type {
		panic(fmt.Sprintf("Type mismatch: expected %v, got %v", varDeclNode.Type, getExpressionType(node.Expr, *v.SymbolTable)))
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
	hasReturn := false
	for _, stmt := range funcBlock.Stmts {
		if ret, ok := stmt.(*ASTReturnNode); ok {
			hasReturn = true
			// Check if the return type matches the function return type
			if getExpressionType(ret.Expr, *v.SymbolTable) != node.ReturnType {
				panic(fmt.Sprintf("Return type mismatch: expected %v, got %v", node.ReturnType, getExpressionType(ret, *v.SymbolTable)))
			}
		}
	}
	if !hasReturn {
		panic("Function must have a return statement")
	}
}

func (v *SemanticVisitor) VisitSimpleExpressionNode(node *ASTSimpleExpression) {
}
