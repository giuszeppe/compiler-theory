package main

import "fmt"

type Frame map[string]FrameItem

type GeneratorVisitor struct {
	SymbolTable  *GeneratorSymbolTable
	Instructions []string
}

type FrameItem struct {
	Node         ASTNode
	Type         string
	IndexInFrame int
	LevelInSoF   int
}

type GeneratorSymbolTable struct {
	Frames Stack[Frame]
}

func (st *GeneratorSymbolTable) Push() {
	newFrame := make(Frame)
	if !st.Frames.IsEmpty() {
		// Copy the current scope to the new scope

		// update all frame items to add 1 to their level in the stack of frames
		newItems := Stack[Frame]{}
		for _, v := range st.Frames.items {
			for key, item := range v {
				v[key] = FrameItem{
					Node:         item.Node,
					Type:         item.Type,
					IndexInFrame: item.IndexInFrame,
					LevelInSoF:   item.LevelInSoF + 1,
				}
			}
			newItems.Push(v)
		}
		st.Frames = newItems
	}
	st.Frames.Push(newFrame)
}
func (st *GeneratorSymbolTable) Pop() {
	// update all frame items to remove 1 from their level in the stack of frames
	newItems := Stack[Frame]{}
	for _, v := range st.Frames.items {
		for key, item := range v {
			v[key] = FrameItem{
				Node:         item.Node,
				Type:         item.Type,
				IndexInFrame: item.IndexInFrame,
				LevelInSoF:   item.LevelInSoF - 1,
			}
		}
		newItems.Push(v)
	}
	st.Frames = newItems
	st.Frames.Pop()
}
func (st *GeneratorSymbolTable) Lookup(name string) (FrameItem, bool) {
	currentFrame, err := st.Frames.Peek()
	if err != nil {
		return FrameItem{}, false
	}
	if node, ok := currentFrame[name]; ok {
		return node, true
	}
	return FrameItem{}, false
}
func (st *GeneratorSymbolTable) Insert(name string, node FrameItem) {
	currentFrame, err := st.Frames.Peek()
	if err != nil {
		return
	}
	currentFrame[name] = node
}

func (v *GeneratorVisitor) emit(instr string) int {
	v.Instructions = append(v.Instructions, instr)
	return len(v.Instructions) - 1
}

func NewGeneratorVisitor() *GeneratorVisitor {
	return &GeneratorVisitor{
		SymbolTable: &GeneratorSymbolTable{},
	}
}

// ====================================================== Entry Points========================================== //
func (v *GeneratorVisitor) VisitProgramNode(node *ASTProgramNode) {
	v.emit(".main")
	openFrameAndPopIfBlock(v, &node.Block)
	v.emit("halt")
}
func (v *GeneratorVisitor) VisitBlockNode(node *ASTBlockNode) {
	for _, stmt := range node.Stmts {
		openFrameAndPopIfBlock(v, stmt)
	}
}

func openFrameAndPopIfBlock(v *GeneratorVisitor, node ASTNode) {
	// if node is a block, push and pop the frame
	if blockNode, ok := node.(*ASTBlockNode); ok {
		varCount := 0
		for _, stmt := range blockNode.Stmts {
			_, isVarDecl := stmt.(*ASTVarDeclNode)
			if isVarDecl {
				varCount++
			}
		}
		v.emit(fmt.Sprintf("push %d", varCount))
		v.emit("oframe")
		node.Accept(v)
		v.emit("cframe") // pop frame
	} else {
		// if node is not a block, just accept it
		node.Accept(v)
	}
}

// ===== Builtins =====

func (v *GeneratorVisitor) VisitBuiltinFuncNode(node *ASTBuiltinFuncNode) {
	switch node.Name {
	case "__delay":
		node.Args[0].Accept(v)
		v.emit("delay")
	case "__width":
		v.emit("width")
	case "__height":
		v.emit("height")
	case "__write":
		for i := len(node.Args) - 1; 0 <= i; i-- {
			node.Args[i].Accept(v)
		}
		v.emit("write")
	case "__write_box":
		for _, arg := range node.Args {
			arg.Accept(v)
		}
		v.emit("writebox")
	case "__clear":
		node.Args[0].Accept(v)
		v.emit("clear")
	}
}

// Functions node
func (v *GeneratorVisitor) VisitFuncDeclNode(node *ASTFuncDeclNode) {}

func (v *GeneratorVisitor) VisitFormalParamsNode(node *ASTFormalParamsNode) {}

func (v *GeneratorVisitor) VisitFormalParamNode(node *ASTFormalParamNode) {}

func (v *GeneratorVisitor) VisitActualParamsNode(node *ASTActualParamsNode) {}

func (v *GeneratorVisitor) VisitActualParamNode(node *ASTActualParamNode) {}

/*
call
1 Push current program counter + 1 to address stack.
2 Pop a = program counter (address) of function to call.
3 Pop c = parameters count.
4 Pop c parameter values.
5 Create frame f and store all c values.
6 Push frame to memory stack.
7 Set program counter to a.
ret - Pop value a from address stack and set program
counter to a. Pops frame from memory stack i.e. closes the
current scope.
*/
func (v *GeneratorVisitor) VisitFuncCallNode(node *ASTFuncCallNode) {}

func (v *GeneratorVisitor) VisitPrintNode(node *ASTPrintNode) {}

func (v *GeneratorVisitor) VisitReturnNode(node *ASTReturnNode) {}

/*
delay - Pops x from the operand stack and delays (pauses)
execution of the program by x milliseconds.
write - Pops x, y and c from the operand stack and sets the
colour of pixel at location x, y to c.
writebox - Pops x, y, w, h and c from the operand stack
and sets the colour of pixel region at location x, y,w, h to c.
clear - Pops colour c from the operand stack and clears the
display using that colour.
*/

// ============================================== Literals and operators ========================================== //
/*
add, sub, mul, inc, dec, max, min, mod, irnd, lt, le, eq, gt, ge
Relative operators push 0 or 1 to the operand stack
inc, dec and irnd requiire a single operand, the others require two operands
Pops values from operand stack and pushes back result.
*/
func (v *GeneratorVisitor) VisitIntegerNode(node *ASTIntegerNode) {
	v.emit(fmt.Sprintf("push %d", node.Value))
}

func (v *GeneratorVisitor) VisitFloatNode(node *ASTFloatNode) {
	v.emit(fmt.Sprintf("push %f", node.Value))
}

func (v *GeneratorVisitor) VisitBooleanNode(node *ASTBooleanNode) {
	val := 0
	if node.Value {
		val = 1
	}
	v.emit(fmt.Sprintf("push %d", val))
}

func (v *GeneratorVisitor) VisitColorNode(node *ASTColorNode) {
	// assume node.Value is hex string, remove '#' and parse
	v.emit(fmt.Sprintf("push %s", node.Value))
}

// ===== Expressions =====
func (v *GeneratorVisitor) VisitBinaryOpNode(node *ASTBinaryOpNode) {
	node.Right.Accept(v)
	node.Left.Accept(v)
	switch node.Operator {
	case "+":
		v.emit("add")
	case "-":
		v.emit("sub")
	case "*":
		v.emit("mul")
	case "/":
		v.emit("div")
	case "%":
		v.emit("mod")
	case "&&":
		v.emit("and")
	case "||":
		v.emit("or")
	case "==":
		v.emit("eq")
	case "<":
		v.emit("lt")
	case "<=":
		v.emit("le")
	case ">":
		v.emit("gt")
	case ">=":
		v.emit("ge")
	}
}

func (v *GeneratorVisitor) VisitUnaryOpNode(node *ASTUnaryOpNode) {
	node.Operand.Accept(v)
	switch node.Operator {
	case "-":
		v.emit("dec")
	case "+":
		v.emit("inc")
	case "!":
		v.emit("not")
	}
}
func (v *GeneratorVisitor) VisitTypeNode(node *ASTTypeNode) {}

// ========================================== Variables and assignments ========================================== //
// ===== Declarations & Assignments =====
func (v *GeneratorVisitor) VisitVarDeclNode(node *ASTVarDeclNode) {
	// store value
	val, _ := v.SymbolTable.Frames.Peek()
	item := FrameItem{Type: node.Type, Node: node, IndexInFrame: len(val), LevelInSoF: 0}
	v.SymbolTable.Insert(node.Name, item)

	// evaluate expression
	node.Expression.Accept(v)
	v.emit(fmt.Sprintf("push %d", item.IndexInFrame))
	v.emit(fmt.Sprintf("push %d", item.LevelInSoF))
	v.emit("st")
}
func (v *GeneratorVisitor) VisitAssignmentNode(node *ASTAssignmentNode) {
	// evaluate RHS
	node.Expr.Accept(v)
	// lookup var
	item, _ := v.SymbolTable.Lookup(node.Id.Token.Lexeme)
	v.emit(fmt.Sprintf("push %d", item.IndexInFrame))
	v.emit(fmt.Sprintf("push %d", item.LevelInSoF))
	v.emit("st")
}

// ===== Variables =====
func (v *GeneratorVisitor) VisitVariableNode(node *ASTVariableNode) {
	item, _ := v.SymbolTable.Lookup(node.Token.Lexeme)
	v.emit(fmt.Sprintf("push [%d:%d]", item.IndexInFrame, item.LevelInSoF))
}

func (v *GeneratorVisitor) VisitSimpleExpressionNode(node *ASTSimpleExpression) {}

func (v *GeneratorVisitor) VisitExpressionNode(node *ASTExpressionNode) {}

// ============================================ Control Flow ================================================= //
/*
All control flow instructions update the PC so that program
execution jumps to a specific instruction rather than the next
one.
jmp - unconditional jump; sets PC to value at top of stack.
cjmp, cjmp2 - conditional jump; pops twp values a,b from
stack. If a==1 jump, set PC to b.
*/
func (v *GeneratorVisitor) VisitWhileNode(node *ASTWhileNode) {}

func (v *GeneratorVisitor) VisitForNode(node *ASTForNode) {
	v.SymbolTable.Push()
	v.emit("push " + fmt.Sprint(CountVarDecls(node.Block)+CountVarDecls(node.VarDecl)))
	v.emit("oframe")

	node.VarDecl.Accept(v)

	node.Condition.Accept(v)
	v.emit("push #PC+4")
	idx := v.emit("cjmp")

	exitLoopInstructionIdx := v.emit("push #TBD")
	v.emit("jmp")
	previousInstructionCount := len(v.Instructions)

	node.Block.Accept(v)

	blockCount := len(v.Instructions) - previousInstructionCount
	preIncrementCount := len(v.Instructions)
	node.Increment.Accept(v)
	incrementCount := len(v.Instructions) - preIncrementCount
	v.emit("push #PC-" + fmt.Sprint(7+(incrementCount)+(blockCount)))
	v.emit("jmp")
	endIdx := v.emit("cframe")

	v.Instructions[exitLoopInstructionIdx] = fmt.Sprintf("push #PC+%d", endIdx-idx)

	v.SymbolTable.Pop()
}

func (v *GeneratorVisitor) VisitIfNode(node *ASTIfNode) {}

func (v *GeneratorVisitor) VisitTypeCastNode(node *ASTTypeCastNode) {}

func (v *GeneratorVisitor) VisitEpsilon(node *ASTEpsilon) {}

func CountVarDecls(node ASTNode) int {
	switch node := node.(type) {
	case *ASTBlockNode:
		count := 0
		for _, stmt := range node.Stmts {
			count += CountVarDecls(stmt)
		}
		return count
	case *ASTVarDeclNode:
		return 1
	default:
		return 0
	}
}
