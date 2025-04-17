package main

import (
	"fmt"
	"strings"
)

type PrintNodesVisitor struct {
	Name      string
	NodeCount int
	TabCount  int
}

func NewPrintNodesVisitor() *PrintNodesVisitor {
	return &PrintNodesVisitor{Name: "Print Tree Visitor"}
}

func (v *PrintNodesVisitor) IncTabCount() {
	v.TabCount++
}

func (v *PrintNodesVisitor) DecTabCount() {
	v.TabCount--
}

func (v *PrintNodesVisitor) VisitIntegerNode(node *ASTIntegerNode) {
	v.NodeCount++
	fmt.Println(strings.Repeat("\t", v.TabCount), "Integer value::", node.Value)
}

func (v *PrintNodesVisitor) VisitVariableNode(node *ASTVariableNode) {
	v.NodeCount++
	fmt.Println(strings.Repeat("\t", v.TabCount), "Variable =>", node.Token.Lexeme)
}

func (v *PrintNodesVisitor) VisitAssignmentNode(node *ASTAssignmentNode) {
	v.NodeCount++
	fmt.Println(strings.Repeat("\t", v.TabCount), "Assignment node =>")
	v.IncTabCount()
	node.Id.Accept(v)
	node.Expr.Accept(v)
	v.DecTabCount()
}
func (v *PrintNodesVisitor) VisitVarDeclNode(node *ASTVarDeclNode) {
	v.NodeCount++
	fmt.Print(strings.Repeat("\t", v.TabCount), "Var decl node => ")
	fmt.Printf("%v %v\n", node.Name, node.Type)
	v.IncTabCount()
	v.IncTabCount()
	node.Expression.Accept(v)

	v.DecTabCount()
}

func (v *PrintNodesVisitor) VisitBlockNode(node *ASTBlockNode) {
	v.NodeCount++
	fmt.Println(strings.Repeat("\t", v.TabCount), "New Block =>")
	v.IncTabCount()
	for _, stmt := range node.Stmts {
		stmt.Accept(v)
	}
	v.DecTabCount()
}

func (v *PrintNodesVisitor) VisitBinaryOpNode(node *ASTBinaryOpNode) {
	v.NodeCount++
	fmt.Println(strings.Repeat("\t", v.TabCount), "Binary Op node =>", node.Operator)
	v.IncTabCount()
	node.Left.Accept(v)
	node.Right.Accept(v)
	v.DecTabCount()
}
func (v *PrintNodesVisitor) VisitExpressionNode(node *ASTExpressionNode) {
	v.NodeCount++
	fmt.Print(strings.Repeat("\t", v.TabCount), "Expression node =>")
	fmt.Printf("%v\n", node.Type)
	v.IncTabCount()
	node.Expr.Accept(v)
	v.DecTabCount()
}
func (v *PrintNodesVisitor) VisitSimpleExpressionNode(node *ASTSimpleExpression) {
	v.NodeCount++
	fmt.Println(strings.Repeat("\t", v.TabCount), "Simple node =>")
	v.IncTabCount()
	v.DecTabCount()
}
