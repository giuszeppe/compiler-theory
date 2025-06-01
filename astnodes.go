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
	VisitFloatNode(node *ASTFloatNode)
	VisitBuiltinFuncNode(node *ASTBuiltinFuncNode)
	VisitFuncCallNode(node *ASTFuncCallNode)
	VisitActualParamsNode(node *ASTActualParamsNode)
	VisitActualParamNode(node *ASTActualParamNode)
	VisitUnaryOpNode(node *ASTUnaryOpNode)
	VisitBooleanNode(node *ASTBooleanNode)
	VisitColorNode(node *ASTColorNode)
	VisitReturnNode(node *ASTReturnNode)
	VisitArrayNode(node *ASTArrayNode)
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
	Token  Token
	Offset ASTNode
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
	Token      Token
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
	Token    Token
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
		Op    Token
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
	Token      Token
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

type ASTFloatNode struct {
	Name  string
	Value float64
}

func (n *ASTFloatNode) Accept(visitor ASTVisitor) {
	visitor.VisitFloatNode(n)
}

type ASTBuiltinFuncNode struct {
	Token Token
	Name  string
	Args  []ASTNode
}

func (n *ASTBuiltinFuncNode) Accept(visitor ASTVisitor) {
	visitor.VisitBuiltinFuncNode(n)
}

type ASTFuncCallNode struct {
	Name   Token
	Params ASTNode
}

func (n *ASTFuncCallNode) Accept(visitor ASTVisitor) {
	visitor.VisitFuncCallNode(n)
}

type ASTActualParamsNode struct {
	Params []ASTNode
}

func (n *ASTActualParamsNode) Accept(visitor ASTVisitor) {
	visitor.VisitActualParamsNode(n)
}

type ASTActualParamNode struct {
	Value ASTNode
	Type  string
}

func (n *ASTActualParamNode) Accept(visitor ASTVisitor) {
	visitor.VisitActualParamNode(n)
}

type ASTUnaryOpNode struct {
	Operator string
	Operand  ASTNode
}

func (n *ASTUnaryOpNode) Accept(visitor ASTVisitor) {
	// Implement the Accept method for ASTUnaryOpNode
	visitor.VisitUnaryOpNode(n)
}

type ASTBooleanNode struct {
	Value bool
}

func (n *ASTBooleanNode) Accept(visitor ASTVisitor) {
	// Implement the Accept method for ASTBooleanNode
	visitor.VisitBooleanNode(n)
}

type ASTColorNode struct {
	Token Token
	Value string
}

func (n *ASTColorNode) Accept(visitor ASTVisitor) {
	// Implement the Accept method for ASTColorNode
	visitor.VisitColorNode(n)
}

type ASTReturnNode struct {
	Token Token
	Expr  ASTNode
}

func (n *ASTReturnNode) Accept(visitor ASTVisitor) {
	// Implement the Accept method for ASTReturnNode
	visitor.VisitReturnNode(n)
}

type ASTArrayNode struct {
	Type  string
	Items []ASTNode
	Size  int
	Token Token
}

func (n *ASTArrayNode) Accept(visitor ASTVisitor) {
	// Implement the Accept method for ASTArrayNode
	visitor.VisitArrayNode(n)
}
