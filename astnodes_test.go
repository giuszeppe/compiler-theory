package main

// func TestPrintingVisitor(t *testing.T) {
// 	printVisitor := NewPrintNodesVisitor()

// 	// === Test Case 1: Assignment x = 23 ===
// 	fmt.Println("Building AST for assignment statement: x = 23")
// 	assignmentLHS := &ASTVariableNode{Name: "xVar", Lexeme: "x"}
// 	assignmentRHS := &ASTIntegerNode{Name: "int23", Value: 23}
// 	root := &ASTAssignmentNode{
// 		Name: "assign_x_23",
// 		Id:   assignmentLHS,
// 		Expr: assignmentRHS,
// 	}
// 	root.Accept(printVisitor)
// 	fmt.Println("Node Count =>", printVisitor.NodeCount)
// 	fmt.Println("----")

// 	// === Test Case 2: Variable node x123 ===
// 	fmt.Println("Building AST for variable x123")
// 	rootVar := &ASTVariableNode{Name: "x123Var", Lexeme: "x123"}
// 	rootVar.Accept(printVisitor)
// 	fmt.Println("Node Count =>", printVisitor.NodeCount)
// }
