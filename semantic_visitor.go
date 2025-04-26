package main

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

func (v *SemanticVisitor) VisitAssignmentNode(node *ASTAssignmentNode) {
	// Check if the variable is declared
	_, ok := v.SymbolTable.Lookup(node.Id.Token.Lexeme)
	if !ok {
		panic("Variable not declared: " + node.Id.Token.Lexeme)
	}
	node.Expr.Accept(v)
}

func (v *SemanticVisitor) VisitVarDeclNode(node *ASTVarDeclNode) {
	// Check if the variable is already declared in the current scope
	if _, ok := v.SymbolTable.Lookup(node.Name); ok {
		panic("Variable already declared: " + node.Name)
	}
	v.SymbolTable.Insert(node.Name, node)
	node.Expression.Accept(v)
}
func (v *SemanticVisitor) VisitBlockNode(node *ASTBlockNode) {
	v.SymbolTable.Push()
	for _, stmt := range node.Stmts {
		stmt.Accept(v)
	}
	v.SymbolTable.Pop()
}
func (v *SemanticVisitor) VisitTypeNode(node *ASTTypeNode) {
	// Do nothing
}
func (v *SemanticVisitor) VisitFunctionNode(node *ASTFuncDeclNode) {
	// Check if the function is already declared in the current scope
	if _, ok := v.SymbolTable.Lookup(node.Name); ok {
		panic("Function already declared: " + node.Name)
	}
	v.SymbolTable.Insert(node.Name, node)
	v.SymbolTable.Push()
	node.Params.Accept(v)
	node.Block.Accept(v)
	v.SymbolTable.Pop()
}

func (v *SemanticVisitor) VisitProgramNode(node *ASTProgramNode) {
	// Visit the block node
	node.Block.Accept(v)
}
func (v *SemanticVisitor) VisitIfNode(node *ASTIfNode) {
	// Visit the condition and the block
	node.Condition.Accept(v)
	node.ThenBlock.Accept(v)
	if node.ElseBlock != nil {
		node.ElseBlock.Accept(v)
	}
}

func (v *SemanticVisitor) VisitWhileNode(node *ASTWhileNode) {
	// Visit the condition and the block
	node.Condition.Accept(v)
	node.Block.Accept(v)
}

func (v *SemanticVisitor) VisitForNode(node *ASTForNode) {
	// Visit the initialization, condition, and block
	node.VarDecl.Accept(v)
	node.Condition.Accept(v)
	node.Increment.Accept(v)
	node.Block.Accept(v)
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
	// Check if the function is already declared in the current scope
	if _, ok := v.SymbolTable.Lookup(node.Name); ok {
		panic("Function already declared: " + node.Name)
	}
	v.SymbolTable.Insert(node.Name, node)
	v.SymbolTable.Push()
	node.Params.Accept(v)
	node.Block.Accept(v)
	v.SymbolTable.Pop()
}

func (v *SemanticVisitor) VisitSimpleExpressionNode(node *ASTSimpleExpression) {
}
