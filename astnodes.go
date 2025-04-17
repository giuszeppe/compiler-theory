package main

// ==== Visitor Interface ====

type ASTVisitor interface {
	VisitIntegerNode(node *ASTIntegerNode)
	VisitAssignmentNode(node *ASTAssignmentNode)
	VisitVariableNode(node *ASTVariableNode)
	VisitBlockNode(node *ASTBlockNode)
	VisitVarDeclNode(node *ASTVarDeclNode)
    VisitBinaryOpNode(node *ASTBinaryOpNode)
    VisitExpressionNode(node *ASTExpressionNode)
    VisitSimpleExpressionNode(node *ASTSimpleExpression)
	IncTabCount()
	DecTabCount()
}

// ==== AST Node Interface ====

type ASTNode interface {
	Accept(visitor ASTVisitor)
}

// ==== AST Node Structs ====

type ASTIntegerNode struct {
	Name  string
	Value int
}

func (n *ASTIntegerNode) Accept(visitor ASTVisitor) {
	visitor.VisitIntegerNode(n)
}

type ASTVariableNode struct {
	Token Token
}

func (n *ASTVariableNode) Accept(visitor ASTVisitor) {
	visitor.VisitVariableNode(n)
}

type ASTAssignmentNode struct {
	Id   ASTVariableNode // usually a VariableNode
	Expr ASTExpressionNode// usually an Expression Node
}

func (n *ASTAssignmentNode) Accept(visitor ASTVisitor) {
	visitor.VisitAssignmentNode(n)
}

type ASTBlockNode struct {
	Name  string
	Stmts []ASTNode
}

func (n *ASTBlockNode) Accept(visitor ASTVisitor) {
	visitor.VisitBlockNode(n)
}

type ASTTypeNode struct {
	Name string
}

func (n *ASTTypeNode) Accept(visitor ASTVisitor) {
}

type ASTVarDeclNode struct {
	Name       string
	Type       string
	Expression ASTExpressionNode
}

func (n *ASTVarDeclNode) Accept(visitor ASTVisitor) {
	visitor.VisitVarDeclNode(n)
}

type ASTExpressionNode struct {
	Expr ASTNode
	Type     string
}

func (n *ASTExpressionNode) Accept(visitor ASTVisitor) {
    visitor.VisitExpressionNode(n)
}

type ASTLiteralNode struct {
	Token Token
}

func (n *ASTLiteralNode) Accept(visitor ASTVisitor) {
}

type ASTSimpleExpression struct {
	Token Token
}

func (n *ASTSimpleExpression) Accept(visitor ASTVisitor) {
    visitor.VisitSimpleExpressionNode(n)
}

type ASTBinaryOpNode struct {
	Operator string
	Left     ASTNode
	Right    ASTNode
}

// Implementing ASTNode interface's Accept method
func (b *ASTBinaryOpNode) Accept(visitor ASTVisitor) {
	visitor.VisitBinaryOpNode(b)
}

