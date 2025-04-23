package main

import "fmt"

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
	VisitProgramNode(node *ASTProgramNode)
	VisitPrintNode(node *ASTPrintNode)
	VisitIfNode(node *ASTIfNode)
	VisitWhileNode(node *ASTWhileNode)
	VisitEpsilon(node *ASTEpsilon)
	VisitTypeCastNode(node *ASTTypeCastNode)
	VisitForNode(node *ASTForNode)

	VisitFuncDeclNode(node *ASTFuncDeclNode)
	VisitFormalParamsNode(node *ASTFormalParamsNode)
	VisitFormalParamNode(node *ASTFormalParamNode)
	VisitTypeNode(node *ASTTypeNode)

	IncTabCount()
	DecTabCount()
}

// ==== AST Node Interface ====

type ASTNode interface {
	Accept(visitor ASTVisitor)
}

// ==== AST Node Structs ====

type ASTProgramNode struct {
	Block ASTBlockNode
}

func (p *ASTProgramNode) Accept(visitor ASTVisitor) {
	visitor.VisitProgramNode(p)
}

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
	Expr ASTNode         // usually an Expression Node
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
	Expression ASTNode
}

func (n *ASTVarDeclNode) Accept(visitor ASTVisitor) {
	visitor.VisitVarDeclNode(n)
}

type ASTExpressionNode struct {
	Expr ASTNode
	Type string
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

type ASTPrintNode struct {
	Expr ASTExpressionNode
}

// Implementing ASTNode interface's Accept method
func (b *ASTPrintNode) Accept(visitor ASTVisitor) {
	visitor.VisitPrintNode(b)
}

type ASTEpsilon struct{}

func (b *ASTEpsilon) Accept(visitor ASTVisitor) {
	visitor.VisitEpsilon(b)
}

type ASTOpList struct {
	Pairs []struct {
		Op    string
		Right ASTNode
	}
}

func (b *ASTOpList) Accept(visitor ASTVisitor) {
	fmt.Println("ASTOpList")
}

type ASTIfNode struct {
	Condition ASTNode
	ThenBlock ASTNode
	ElseBlock ASTNode
}

func (n *ASTIfNode) Accept(visitor ASTVisitor) {
	visitor.VisitIfNode(n)
}

type ASTWhileNode struct {
	Condition ASTNode
	Block     ASTNode
}

func (n *ASTWhileNode) Accept(visitor ASTVisitor) {
	visitor.VisitWhileNode(n)
}

type ASTTypeCastNode struct {
	Type string
	Expr ASTNode
}

func (n *ASTTypeCastNode) Accept(visitor ASTVisitor) {
	visitor.VisitTypeCastNode(n)
}

type ASTForNode struct {
	VarDecl   ASTNode
	Condition ASTNode
	Increment ASTNode
	Block     ASTNode
}

func (n *ASTForNode) Accept(visitor ASTVisitor) {
	visitor.VisitForNode(n)
}

type ASTFormalParamsNode struct {
	Params []ASTNode
}

func (n *ASTFormalParamsNode) Accept(visitor ASTVisitor) {
	visitor.VisitFormalParamsNode(n)
}

type ASTFuncDeclNode struct {
	Name       string
	ReturnType string
	Params     ASTNode
	Block      ASTNode
}

func (n *ASTFuncDeclNode) Accept(visitor ASTVisitor) {
	visitor.VisitFuncDeclNode(n)
}

type ASTFormalParamNode struct {
	Name string
	Type string
}

func (n *ASTFormalParamNode) Accept(visitor ASTVisitor) {
	visitor.VisitFormalParamNode(n)
}
