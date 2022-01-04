package interpreter

import (
	"fmt"
	"nova-lang/pkg/lexer"
	"nova-lang/pkg/parser"
	"nova-lang/pkg/shared"
)

type Interpreter struct {
	context *Context
}

func (interpreter *Interpreter) Start(node parser.Node) *RuntimeResult {

	if node.Type() == parser.Node_List {

		listNode := node.(*parser.ListNode)

		for _, node := range listNode.Nodes {

			if node.Type() == parser.Node_Function {
				err := interpreter.context.DeclareFunction(node.(*parser.FunctionNode))
				if err != nil {
					return NewRuntimeResult().Failure(err)
				}
			}

		}

	}

	return interpreter.visit(node)
}

func (interpreter *Interpreter) visit(node parser.Node) *RuntimeResult {

	if node.Type() == parser.Node_List {
		return interpreter.visitListNode(node)
	} else if node.Type() == parser.Node_Number {
		return interpreter.visitNumberNode(node)
	} else if node.Type() == parser.Node_Comment {
		return NewRuntimeResult().Success(nil)
	} else if node.Type() == parser.Node_PrintStdout {
		return interpreter.visitPrintStdoutNode(node)
	} else if node.Type() == parser.Node_BinOp {
		return interpreter.visitBinaryOperationNode(node)
	} else if node.Type() == parser.Node_Return {
		return interpreter.visitReturnNode(node)
	} else if node.Type() == parser.Node_Bool {
		return interpreter.visitBoolNode(node)
	} else if node.Type() == parser.Node_VarDeclar {
		return interpreter.visitVarDeclarationNode(node)
	} else if node.Type() == parser.Node_VarAssign {
		return interpreter.visitVarAssignNode(node)
	} else if node.Type() == parser.Node_VarAccess {
		return interpreter.visitVarAcessNode(node)
	} else if node.Type() == parser.Node_UnaryOp {
		return interpreter.visitUnaryOperationNode(node)
	} else if node.Type() == parser.Node_If {
		return interpreter.visitIfNode(node)
	} else if node.Type() == parser.Node_String {
		return interpreter.visitStringNode(node)
	} else if node.Type() == parser.Node_DoCase {
		return interpreter.visitDoCaseNode(node)
	} else if node.Type() == parser.Node_DoWhile {
		return interpreter.visitDoWhileNode(node)
	} else if node.Type() == parser.Node_Command {
		return interpreter.visitCommandNode(node)
	} else if node.Type() == parser.Node_ForLoop {
		return interpreter.visitForLoopNode(node)
	} else if node.Type() == parser.Node_Function {
		return NewRuntimeResult()
	} else if node.Type() == parser.Node_FuncCall {
		return interpreter.visitFunctionCallNode(node)
	}

	panic("not implemented yet for " + fmt.Sprint(node.Type()))

}

func (interpreter *Interpreter) visitFunctionCallNode(node parser.Node) *RuntimeResult {

	funcCallNode := node.(*parser.FunctionCallNode)
	res := NewRuntimeResult()

	if funcCallNode.FunctionName.Type() != parser.Node_VarAccess {
		return res.Failure(shared.NewRuntimeErrorRange(funcCallNode.FunctionName.Range(), "Expected function name"))
	}

	funcName := funcCallNode.FunctionName.(*parser.VarAccessNode).VarName

	builtInFunction, found := BuiltInFunctions[funcName]

	if found {

		parametersValues := []Value{}

		for i := range funcCallNode.Args {

			value := res.Register(interpreter.visit(funcCallNode.Args[i]))
			if res.ShouldReturn() {
				return res
			}

			parametersValues = append(parametersValues, value)

		}

		res.Register(builtInFunction(interpreter.context, funcCallNode.Range(), parametersValues))
		if res.Error != nil {
			return res
		}

		if res.FunctionReturnValue == nil {
			return res.Failure(shared.NewRuntimeErrorRange(funcCallNode.Range(), fmt.Sprintf("Builtin function `%s` should return a value", funcName)))
		}

		return res.Success(res.FunctionReturnValue)
	}

	function := interpreter.context.GetFunction(funcName)

	if function != nil {

		parametersValues := []Value{}

		for i := range function.Parameters {

			if i >= len(funcCallNode.Args) {
				parametersValues = append(parametersValues, NewBoolean(false))
			} else {

				value := res.Register(interpreter.visit(funcCallNode.Args[i]))
				if res.ShouldReturn() {
					return res
				}

				parametersValues = append(parametersValues, value)

			}

		}

		interpreter.context.IncreaseLevel(function)

		// TODO create parameters with the correct range
		for i, paramName := range function.Parameters {
			interpreter.context.AssignValueToVariable(paramName, parametersValues[i], function.Range())
		}

		res.Register(interpreter.visit(function.Body))

		if res.Error != nil {
			interpreter.context.DecreaseLevel()
			return res
		}

		if res.FunctionReturnValue == nil {
			return res.Failure(shared.NewRuntimeErrorRange(funcCallNode.Range(), fmt.Sprintf("Function `%s` should return a value", funcName)))
		}

		interpreter.context.DecreaseLevel()

		return res.Success(res.FunctionReturnValue)

	}

	return res.Failure(shared.NewRuntimeErrorRange(funcCallNode.FunctionName.Range(), fmt.Sprintf("Function `%s` not declared", funcName)))

}

func (interpreter *Interpreter) visitForLoopNode(node parser.Node) *RuntimeResult {

	forLoopNode := node.(*parser.ForNode)
	res := NewRuntimeResult()

	startValue := res.Register(interpreter.visit(forLoopNode.StartNode))
	if res.ShouldReturn() {
		return res
	}

	if startValue.Type() != ValueType_Number {
		return res.Failure(shared.NewRuntimeErrorRange(forLoopNode.Range(), fmt.Sprintf("Invalid start value for for loop. Expected number got `%v`", startValue.Type())))
	}

	endValue := res.Register(interpreter.visit(forLoopNode.EndNode))
	if res.ShouldReturn() {
		return res
	}

	if endValue.Type() != ValueType_Number {
		return res.Failure(shared.NewRuntimeErrorRange(forLoopNode.Range(), fmt.Sprintf("Invalid end value for for loop. Expected number got `%v`", endValue.Type())))
	}

	var stepValue Value = NewNumber(1)
	if forLoopNode.StepNode != nil {

		stepValue = res.Register(interpreter.visit(forLoopNode.StepNode))
		if res.ShouldReturn() {
			return res
		}

		if stepValue.Type() != ValueType_Number {
			return res.Failure(shared.NewRuntimeErrorRange(forLoopNode.Range(), fmt.Sprintf("Invalid step value for for loop. Expected number got `%v`", stepValue.Type())))
		}

	}

	interpreter.context.AssignValueToVariable(forLoopNode.VarName, startValue, forLoopNode.StartNode.Range())

	interpreter.context.RegisterLoopEnter()

	for {

		res.Register(interpreter.visit(forLoopNode.BodyNode))
		if res.Error != nil || res.FunctionReturnValue != nil || res.LoopShouldExit {
			break
		}

		variable, _ := interpreter.context.GetVariable(forLoopNode.VarName)
		newValue, _ := variable.Value.Add(stepValue)

		interpreter.context.AssignValueToVariable(forLoopNode.VarName, newValue, forLoopNode.StartNode.Range())

		if stepValue.(*Number).Value > 0 {
			if newValue.(*Number).Value > endValue.(*Number).Value {
				break
			}
		} else if stepValue.(*Number).Value < 0 {
			if newValue.(*Number).Value < endValue.(*Number).Value {
				break
			}
		}

	}

	res.LoopShouldExit = false
	res.LoopShouldLoop = false

	interpreter.context.RegisterLoopExit()

	return res
}

func (interpreter *Interpreter) visitCommandNode(node parser.Node) *RuntimeResult {

	commandNode := node.(*parser.CommandNode)
	res := NewRuntimeResult()

	if commandNode.CommandType == parser.CommandType_Exit {
		return res.SuccessExit()
	} else if commandNode.CommandType == parser.CommandType_Loop {
		return res.SuccessLoop()
	}

	panic("not implemented yet for " + fmt.Sprint(commandNode.CommandType))
}

func (interpreter *Interpreter) visitDoWhileNode(node parser.Node) *RuntimeResult {

	doWhileNode := node.(*parser.DoWhileNode)
	res := NewRuntimeResult()

	interpreter.context.RegisterLoopEnter()

	for {

		loopConditionValue := res.Register(interpreter.visit(doWhileNode.Condition))
		if res.ShouldReturn() {
			break
		}

		isTrue, err := loopConditionValue.IsTrue()
		if err != nil {
			interpreter.context.RegisterLoopExit()
			return res.Failure(err)
		}

		if !isTrue {
			break
		}

		res.Register(interpreter.visit(doWhileNode.Body))
		if res.Error != nil || res.FunctionReturnValue != nil || res.LoopShouldExit {
			break
		}

	}

	interpreter.context.RegisterLoopExit()

	res.LoopShouldExit = false
	res.LoopShouldLoop = false

	return res
}

func (interpreter *Interpreter) visitDoCaseNode(node parser.Node) *RuntimeResult {

	doCaseNode := node.(*parser.CaseNode)
	res := NewRuntimeResult()

	for _, caseCase := range doCaseNode.Cases {

		conditionValue := res.Register(interpreter.visit(caseCase.CaseExpr))
		if res.ShouldReturn() {
			return res
		}

		isTrue, err := conditionValue.IsTrue()
		if err != nil {
			return res.Failure(err)
		}

		if isTrue {
			res.Register(interpreter.visit(caseCase.Body))
			if res.ShouldReturn() {
				return res
			}
			return res.Success(nil)
		}

	}

	if doCaseNode.OtherwiseCase != nil {
		res.Register(interpreter.visit(doCaseNode.OtherwiseCase))
		if res.ShouldReturn() {
			return res
		}
		return res.Success(nil)
	}

	return res.Success(nil)
}

func (interpreter *Interpreter) visitStringNode(node parser.Node) *RuntimeResult {

	strNode := node.(*parser.StringNode)
	res := NewRuntimeResult()

	return res.Success(NewString(strNode.Value).UpdateRange(node.Range()))
}

func (interpreter *Interpreter) visitIfNode(node parser.Node) *RuntimeResult {

	ifNode := node.(*parser.IfNode)
	res := NewRuntimeResult()

	for _, ifCase := range ifNode.IfCases {

		conditionValue := res.Register(interpreter.visit(ifCase.CaseExpr))
		if res.ShouldReturn() {
			return res
		}

		isTrue, err := conditionValue.IsTrue()
		if err != nil {
			return res.Failure(err)
		}

		if isTrue {
			res.Register(interpreter.visit(ifCase.Body))
			if res.ShouldReturn() {
				return res
			}
			return res.Success(nil)
		}

	}

	if ifNode.ElseCase != nil {
		res.Register(interpreter.visit(ifNode.ElseCase))
		if res.ShouldReturn() {
			return res
		}
		return res.Success(nil)
	}

	return res.Success(nil)
}

func (interpreter *Interpreter) visitUnaryOperationNode(node parser.Node) *RuntimeResult {

	unaryOperNode := node.(*parser.UnaryOperationNode)
	res := NewRuntimeResult()

	value := res.Register(interpreter.visit(unaryOperNode.Node))
	if res.ShouldReturn() {
		return res
	}

	if unaryOperNode.OperationToken.Type == lexer.TokenType_Minus {

		if value.Type() != ValueType_Number {
			return res.Failure(shared.NewRuntimeErrorRange(unaryOperNode.Range(), "Operand must be a number"))
		}

		number := value.(*Number)
		number.Value = -number.Value
		number.UpdateRange(unaryOperNode.Range())

		return res.Success(number)

	} else if unaryOperNode.OperationToken.Type == lexer.TokenType_Not {

		if value.Type() != ValueType_Boolean {
			return res.Failure(shared.NewRuntimeErrorRange(unaryOperNode.Range(), "Operand must be a boolean/logic"))
		}

		boolean := value.(*Boolean)
		boolean.Value = !boolean.Value
		boolean.UpdateRange(unaryOperNode.Range())

		return res.Success(boolean)

	}

	panic("unreachable")
}

func (interpreter *Interpreter) visitVarAcessNode(node parser.Node) *RuntimeResult {

	varAcessNode := node.(*parser.VarAccessNode)
	res := NewRuntimeResult()

	variable, _ := interpreter.context.GetVariable(varAcessNode.VarName)

	if variable == nil {
		return res.Failure(shared.NewRuntimeErrorRange(varAcessNode.Range(), "Variable '"+varAcessNode.VarName+"' is not defined"))
	}

	return res.Success(variable.Value)
}

func (interpreter *Interpreter) visitVarAssignNode(node parser.Node) *RuntimeResult {

	varAssignNode := node.(*parser.VarAssignmentNode)
	res := NewRuntimeResult()

	variableName := varAssignNode.VarName
	value := res.Register(interpreter.visit(varAssignNode.Expr))
	if res.ShouldReturn() {
		return res
	}

	interpreter.context.AssignValueToVariable(variableName, value, varAssignNode.Range())

	return res
}

func (interpreter *Interpreter) visitVarDeclarationNode(node parser.Node) *RuntimeResult {

	varDeclNode := node.(*parser.VarDeclarationNode)
	res := NewRuntimeResult()

	for _, variableName := range varDeclNode.VarNames {

		if varDeclNode.Modifier == "public" {
			interpreter.context.DeclareVariable(variableName, Visibility_Public, node.Range())
		} else {
			interpreter.context.DeclareVariable(variableName, Visibility_Private, node.Range())
		}

	}

	return res.Success(nil)
}

func (interpreter *Interpreter) visitBoolNode(node parser.Node) *RuntimeResult {
	boolNode := node.(*parser.BooleanNode)
	return NewRuntimeResult().Success(NewBoolean(boolNode.Value).UpdateRange(node.Range()))
}

func (interpreter *Interpreter) visitReturnNode(node parser.Node) *RuntimeResult {

	returnNode := node.(*parser.ReturnNode)
	res := NewRuntimeResult()

	var returnValue Value = nil

	if returnNode.Expr == nil {
		// TODO Recital returns .f. in this cases, but i think it should return null or undefined
		returnValue = NewBoolean(false).UpdateRange(node.Range())
	} else {
		returnValue = res.Register(interpreter.visit(returnNode.Expr))
		if res.ShouldReturn() {
			return res
		}
	}

	return res.SuccessReturn(returnValue)
}

func (interpreter *Interpreter) visitBinaryOperationNode(node parser.Node) *RuntimeResult {

	binOpNode := node.(*parser.BinaryOperationNode)
	res := NewRuntimeResult()

	leftValue := res.Register(interpreter.visit(binOpNode.LeftNode))
	if res.ShouldReturn() {
		return res
	}

	rightValue := res.Register(interpreter.visit(binOpNode.RightNode))
	if res.ShouldReturn() {
		return res
	}

	var value Value = nil
	var err *shared.Error = nil

	if binOpNode.OperationToken.Type == lexer.TokenType_Plus {
		value, err = leftValue.Add(rightValue)
	} else if binOpNode.OperationToken.Type == lexer.TokenType_Minus {
		value, err = leftValue.Subtract(rightValue)
	} else if binOpNode.OperationToken.Type == lexer.TokenType_Star {
		value, err = leftValue.Multiply(rightValue)
	} else if binOpNode.OperationToken.Type == lexer.TokenType_Slash {
		value, err = leftValue.Divide(rightValue)
	} else if binOpNode.OperationToken.Type == lexer.TokenType_Exponential {
		value, err = leftValue.Exponential(rightValue)
	} else if binOpNode.OperationToken.Type == lexer.TokenType_Percent {
		value, err = leftValue.Remainder(rightValue)
	} else if binOpNode.OperationToken.Type == lexer.TokenType_Equals {
		value, err = leftValue.Equals(rightValue)
	} else if binOpNode.OperationToken.Type == lexer.TokenType_NotEqual {
		value, err = leftValue.NotEquals(rightValue)
	} else if binOpNode.OperationToken.Type == lexer.TokenType_EqualsEquals {
		value, err = leftValue.EqualsEquals(rightValue)
	} else if binOpNode.OperationToken.Type == lexer.TokenType_LessThan {
		value, err = leftValue.IsLess(rightValue)
	} else if binOpNode.OperationToken.Type == lexer.TokenType_LessThanEqual {
		value, err = leftValue.IsLessEquals(rightValue)
	} else if binOpNode.OperationToken.Type == lexer.TokenType_GreaterThan {
		value, err = leftValue.IsGreater(rightValue)
	} else if binOpNode.OperationToken.Type == lexer.TokenType_GreaterThanEqual {
		value, err = leftValue.IsGreaterEquals(rightValue)
	} else if binOpNode.OperationToken.Type == lexer.TokenType_Keyword {

		if binOpNode.OperationToken.Value == "and" {
			value, err = leftValue.And(rightValue)
		} else if binOpNode.OperationToken.Value == "or" {
			value, err = leftValue.Or(rightValue)
		}

	} else {
		panic("binary operation not implemented yet for " + fmt.Sprint(binOpNode.OperationToken.Type))
	}

	if err != nil {
		return res.Failure(err.UpdateRange(node.Range()))
	}

	return res.Success(value.UpdateRange(node.Range()))
}

func (interpreter *Interpreter) visitPrintStdoutNode(node parser.Node) *RuntimeResult {

	printNode := node.(*parser.PrintStdoutNode)
	res := NewRuntimeResult()

	value := res.Register(interpreter.visit(printNode.Expr))
	if res.ShouldReturn() {
		return res
	}

	fmt.Println(value.PrintRepresentation())

	return res.Success(nil)
}

func (interpreter *Interpreter) visitListNode(node parser.Node) *RuntimeResult {

	res := NewRuntimeResult()
	listNode := node.(*parser.ListNode)

	for _, node := range listNode.Nodes {

		res.Register(interpreter.visit(node))
		if res.Error != nil || res.FunctionReturnValue != nil {
			return res
		}

		if res.LoopShouldExit || res.LoopShouldLoop {

			if !interpreter.context.IsInsideLoop() {
				return res.Failure(shared.NewRuntimeErrorRange(node.Range(), "Cannot use break/loop outside of loop"))
			}

			break
		}

	}

	return res
}

func (interpreter *Interpreter) visitNumberNode(node parser.Node) *RuntimeResult {

	numberNode := node.(*parser.NumberNode)

	numberValue := NewNumber(numberNode.Value).UpdateRange(node.Range())

	return NewRuntimeResult().Success(numberValue)

}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		context: NewContext(),
	}
}
