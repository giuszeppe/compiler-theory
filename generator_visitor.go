package main

import (
	"fmt"
)

type GenStack[T any] struct {
	items []T
}

// Push adds an item to the top of the stack, last inserted item is [0]
func (s *GenStack[T]) Push(item T) {
	s.items = append([]T{item}, s.items...)
}

// Pop removes and returns the top item
func (s *GenStack[T]) Pop() (T, error) {
	if len(s.items) == 0 {
		var zero T
		return zero, fmt.Errorf("stack is empty")
	}
	last := s.items[0]
	s.items = s.items[1:]
	return last, nil
}

// Peek returns the top item without removing it
func (s *GenStack[T]) Peek() (T, error) {
	if len(s.items) == 0 {
		var zero T

		return zero, fmt.Errorf("stack is empty")
	}
	return s.items[0], nil
}

// IsEmpty checks if the stack is empty
func (s *GenStack[T]) IsEmpty() bool {
	return len(s.items) == 0
}

// Size returns the number of items in the stack
func (s *GenStack[T]) Size() int {
	return len(s.items)
}

type GeneratorVisitor struct {
	SymbolTable  *FrameStack
	Instructions []string
}

type SymbolGen struct {
	Name       string
	FrameIndex int // index inside its own frame
}

type Frame struct {
	Symbols map[string]SymbolGen
}

type FrameStack struct {
	Frames GenStack[Frame] // Top of stack is Frames[0]
}

// NewFrameStack creates a new stack
func NewFrameStack() *FrameStack {
	return &FrameStack{
		Frames: GenStack[Frame]{},
	}
}

// PushFrame adds a new empty frame at the top
func (fs *FrameStack) PushFrame() {
	frame := Frame{
		Symbols: make(map[string]SymbolGen),
	}
	fs.Frames.Push(frame)
}

// PopFrame removes the top frame
func (fs *FrameStack) PopFrame() {
	if fs.Frames.Size() > 0 {
		fs.Frames.Pop()
	}
}

// Define adds a symbol to the top frame
func (fs *FrameStack) Define(name string) SymbolGen {
	if fs.Frames.Size() == 0 {
		panic("no frame to define symbol in")
	}
	frame, _ := fs.Frames.Peek()
	sym := SymbolGen{
		Name:       name,
		FrameIndex: len(frame.Symbols), // number of symbols already in the frame
	}
	frame.Symbols[name] = sym
	return sym
}

// Resolve looks for a symbol starting from top frame
func (fs *FrameStack) Resolve(name string) (SymbolGen, int, bool) {

	for frameIndex, frame := range fs.Frames.items {
		if sym, ok := frame.Symbols[name]; ok {
			return sym, frameIndex, true
		}
	}
	return SymbolGen{}, -1, false
}

func (v *GeneratorVisitor) emit(instr string) int {
	v.Instructions = append(v.Instructions, instr)
	return len(v.Instructions) - 1
}

func NewGeneratorVisitor() *GeneratorVisitor {
	return &GeneratorVisitor{
		SymbolTable: &FrameStack{},
	}
}

// ====================================================== Entry Points========================================== //
func (v *GeneratorVisitor) VisitProgramNode(node *ASTProgramNode) {
	v.emit(".main")
	v.emit("push #PC+3")
	v.emit("jmp")
	v.emit("halt")

	v.SymbolTable.PushFrame()
	openFrameAndPopIfBlock(v, &node.Block)
	v.SymbolTable.PopFrame()
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
		for i := len(node.Args) - 1; 0 <= i; i-- {
			node.Args[i].Accept(v)
		}
		v.emit("writebox")
	case "__print":
		for i := len(node.Args) - 1; 0 <= i; i-- {
			node.Args[i].Accept(v)
		}
		v.emit("print")
	case "__random_int":
		node.Args[0].Accept(v)
		v.emit("irnd")
	case "__clear":
		node.Args[0].Accept(v)
		v.emit("clear")
	}
}

// Functions node
func (v *GeneratorVisitor) VisitFuncDeclNode(node *ASTFuncDeclNode) {
	// push frame
	skipFunctionBodyIdx := v.emit("push TBD")
	v.emit("jmp")
	v.SymbolTable.PushFrame()
	v.emit("." + node.Name)
	v.emit(fmt.Sprintf("push %d", CountVarDecls(node.Block)))
	v.emit("alloc")

	// visit params
	node.Params.Accept(v)

	// visit block
	node.Block.Accept(v)

	// pop frame, not needed since return node places it
	// endIdx := v.emit("cframe")
	v.Instructions[skipFunctionBodyIdx] = fmt.Sprint("push #PC+", len(v.Instructions)-skipFunctionBodyIdx)
	v.SymbolTable.PopFrame()
}

func (v *GeneratorVisitor) VisitFormalParamsNode(node *ASTFormalParamsNode) {
	for _, param := range node.Params {
		v.SymbolTable.Define(param.(*ASTVarDeclNode).Name)
	}
}

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
func (v *GeneratorVisitor) VisitFuncCallNode(node *ASTFuncCallNode) {

	params := node.Params.(*ASTActualParamsNode)
	for i := len(params.Params) - 1; i >= 0; i-- {
		params.Params[i].Accept(v)
	}
	v.emit("push " + fmt.Sprint(len(params.Params)))
	v.emit("push ." + node.Name)
	v.emit("call")
}

func (v *GeneratorVisitor) VisitPrintNode(node *ASTPrintNode) {}

func (v *GeneratorVisitor) VisitReturnNode(node *ASTReturnNode) {
	node.Expr.Accept(v)
	v.emit("cframe")
	v.emit("ret")
}

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
	case "and":
		v.emit("and")
	case "or":
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
	case "not":
		v.emit("not")
	}
}
func (v *GeneratorVisitor) VisitTypeNode(node *ASTTypeNode) {}

// ========================================== Variables and assignments ========================================== //
// ===== Declarations & Assignments =====
func (v *GeneratorVisitor) VisitVarDeclNode(node *ASTVarDeclNode) {
	// store value
	// val, _ := v.SymbolTable.Frames.Peek()
	// item := FrameItem{Type: node.Type, Node: node, IndexInFrame: len(val), LevelInSoF: 0}
	item := v.SymbolTable.Define(node.Name)
	_, a, _ := v.SymbolTable.Resolve(node.Name)

	// evaluate expression
	node.Expression.Accept(v)
	v.emit(fmt.Sprintf("push %d", item.FrameIndex))
	v.emit(fmt.Sprintf("push %d", a))
	v.emit("st")
}
func (v *GeneratorVisitor) VisitAssignmentNode(node *ASTAssignmentNode) {
	// evaluate RHS
	node.Expr.Accept(v)
	// lookup var
	item, level, _ := v.SymbolTable.Resolve(node.Id.Token.Lexeme)
	v.emit(fmt.Sprintf("push %d", item.FrameIndex))
	v.emit(fmt.Sprintf("push %d", level))
	v.emit("st")
}

// ===== Variables =====
func (v *GeneratorVisitor) VisitVariableNode(node *ASTVariableNode) {
	item, level, _ := v.SymbolTable.Resolve(node.Token.Lexeme)
	v.emit(fmt.Sprintf("push [%d:%d]", item.FrameIndex, level))
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
func (v *GeneratorVisitor) VisitWhileNode(node *ASTWhileNode) {
	v.SymbolTable.PushFrame()
	v.emit("push " + fmt.Sprint(CountVarDecls(node.Block)))
	idxCondition := v.emit("oframe")
	idxCondition++

	node.Condition.Accept(v)
	idxEnd := v.emit("")

	v.Instructions[idxEnd] = fmt.Sprintf("push #PC+%d", 4)
	idx := v.emit("cjmp")

	exitLoopInstructionIdx := v.emit("push #TBD")
	v.emit("jmp")
	// previousInstructionCount := len(v.Instructions)

	node.Block.Accept(v)

	// blockCount := len(v.Instructions) - previousInstructionCount
	v.emit("push " + fmt.Sprint(idxCondition)) // change to #PC+n where n is the number of instructions in the block to go back to the condition
	v.emit("jmp")
	endIdx := v.emit("cframe")

	v.Instructions[exitLoopInstructionIdx] = fmt.Sprintf("push #PC+%d", endIdx-idx-1)

	v.SymbolTable.PopFrame()
}

func (v *GeneratorVisitor) VisitForNode(node *ASTForNode) {
	v.SymbolTable.PushFrame()
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

	v.Instructions[exitLoopInstructionIdx] = fmt.Sprintf("push #PC+%d", endIdx-idx-1)

	v.SymbolTable.PopFrame()
}

func (v *GeneratorVisitor) VisitIfNode(node *ASTIfNode) {
	v.SymbolTable.PushFrame()
	v.emit("push " + fmt.Sprint(CountVarDecls(node.ThenBlock)))
	v.emit("oframe")

	node.Condition.Accept(v)
	v.emit("push #PC+4")
	idx := v.emit("cjmp")

	elseBlockLocation := v.emit("push #TBD")
	v.emit("jmp")

	node.ThenBlock.Accept(v)

	endIdx := v.emit("cframe")
	skipElseIdx := v.emit("push #PC+TBD")
	v.emit("jmp")

	if node.ElseBlock != nil {
		v.Instructions[elseBlockLocation] = fmt.Sprintf("push #PC+%d", skipElseIdx-idx+1)
	} else {
		v.Instructions[elseBlockLocation] = fmt.Sprintf("push #PC+%d", endIdx-idx-2)
	}

	elseSize := len(v.Instructions)
	if node.ElseBlock != nil {
		node.ElseBlock.Accept(v)
	}
	v.emit("cframe")
	elseSize = len(v.Instructions) - elseSize
	v.Instructions[skipElseIdx] = fmt.Sprintf("push #PC+%d", elseSize+2)

	v.SymbolTable.PopFrame()

}

func (v *GeneratorVisitor) VisitTypeCastNode(node *ASTTypeCastNode) {
	node.Expr.Accept(v)
}

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
