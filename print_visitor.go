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
func (v *PrintNodesVisitor) VisitProgramNode(node *ASTProgramNode) {
	v.NodeCount++
	fmt.Println(strings.Repeat("\t", v.TabCount), "Program node =>")
	node.Block.Accept(v)
	v.IncTabCount()
	v.DecTabCount()
}
func (v *PrintNodesVisitor) VisitPrintNode(node *ASTPrintNode) {
	v.NodeCount++
	fmt.Println(strings.Repeat("\t", v.TabCount), "Print node =>")
	node.Expr.Accept(v)
	v.IncTabCount()
	v.DecTabCount()
}

func (v *PrintNodesVisitor) VisitIfNode(node *ASTIfNode) {
	v.NodeCount++
	fmt.Println(strings.Repeat("\t", v.TabCount), "If node =>")
	v.IncTabCount()
	node.Condition.Accept(v)
	node.ThenBlock.Accept(v)
	v.DecTabCount()
	v.DecTabCount()
	fmt.Println(strings.Repeat("\t", v.TabCount), "Else Block =>")
	v.IncTabCount()
	if node.ElseBlock != nil {
		node.ElseBlock.Accept(v)
	}
	v.DecTabCount()
}

func (v *PrintNodesVisitor) VisitWhileNode(node *ASTWhileNode) {
	v.NodeCount++
	fmt.Println(strings.Repeat("\t", v.TabCount), "While node =>")
	v.IncTabCount()
	node.Condition.Accept(v)
	node.Block.Accept(v)
	v.DecTabCount()
}

func (v *PrintNodesVisitor) VisitTypeCastNode(node *ASTTypeCastNode) {
	v.NodeCount++
	fmt.Println(strings.Repeat("\t", v.TabCount), "Type cast node:: ", node.Type, " =>")
	v.IncTabCount()
	node.Expr.Accept(v)
	v.DecTabCount()
}
func (v *PrintNodesVisitor) VisitEpsilon(node *ASTEpsilon) {
	v.NodeCount++
	fmt.Println(strings.Repeat("\t", v.TabCount), "Epsilon node")
	v.IncTabCount()
	v.DecTabCount()
}

func (v *PrintNodesVisitor) VisitForNode(node *ASTForNode) {
	v.NodeCount++
	fmt.Println(strings.Repeat("\t", v.TabCount), "For node =>")
	v.IncTabCount()

	fmt.Println(strings.Repeat("\t", v.TabCount), "For var decl =>")
	v.IncTabCount()
	node.VarDecl.Accept(v)
	v.DecTabCount()

	fmt.Println(strings.Repeat("\t", v.TabCount), "For condition =>")
	v.IncTabCount()
	node.Condition.Accept(v)
	v.DecTabCount()

	fmt.Println(strings.Repeat("\t", v.TabCount), "For increment =>")
	v.IncTabCount()
	node.Increment.Accept(v)
	v.DecTabCount()

	fmt.Println(strings.Repeat("\t", v.TabCount), "For block =>")
	v.IncTabCount()
	node.Block.Accept(v)
	v.DecTabCount()
}

func (v *PrintNodesVisitor) VisitFuncDeclNode(node *ASTFuncDeclNode) {
	v.NodeCount++
	fmt.Println(strings.Repeat("\t", v.TabCount), "Function decl node =>", node.Name, ":", node.ReturnType)
	v.IncTabCount()

	fmt.Println(strings.Repeat("\t", v.TabCount), "Function params =>")
	v.IncTabCount()
	node.Params.Accept(v)
	v.DecTabCount()

	fmt.Println(strings.Repeat("\t", v.TabCount), "Function block =>")
	v.IncTabCount()
	node.Block.Accept(v)
	v.DecTabCount()
}
func (v *PrintNodesVisitor) VisitFormalParamsNode(node *ASTFormalParamsNode) {
	v.NodeCount++
	fmt.Println(strings.Repeat("\t", v.TabCount), "Formal params node =>")
	v.IncTabCount()
	for _, param := range node.Params {
		param.Accept(v)
	}
	v.DecTabCount()
}
func (v *PrintNodesVisitor) VisitFormalParamNode(node *ASTFormalParamNode) {
	v.NodeCount++
	fmt.Println(strings.Repeat("\t", v.TabCount), "Formal param node:: ", node.Name, ":", node.Type)
}
func (v *PrintNodesVisitor) VisitTypeNode(node *ASTTypeNode) {
	v.NodeCount++
	fmt.Println(strings.Repeat("\t", v.TabCount), "Type node =>")
	v.IncTabCount()
	fmt.Printf("%v\n", node.Name)
	v.DecTabCount()
}

func (v *PrintNodesVisitor) VisitFloatNode(node *ASTFloatNode) {
	v.NodeCount++
	fmt.Println(strings.Repeat("\t", v.TabCount), "Float value::", node.Value)
}
