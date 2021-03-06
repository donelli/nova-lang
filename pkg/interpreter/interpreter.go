package interpreter

import (
	"fmt"
	"nova-lang/pkg/lexer"
	"nova-lang/pkg/parser"
	"nova-lang/pkg/screen"
	"nova-lang/pkg/shared"
)

type Interpreter struct {
	context *Context
	screen  screen.Screen
}

func (interpreter *Interpreter) Start(node parser.Node, testMode bool) *RuntimeResult {

	interpreter.context.CurrentInterpreter = interpreter

	if len(BuiltInFunctions) == 0 {
		InitBuiltInFunctions()
	}

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

	if testMode {
		interpreter.screen = screen.NewConsoleScreen(screen.OutputType_Test)
	} else {
		interpreter.screen = screen.NewConsoleScreen(screen.OutputType_Console)
	}

	err := interpreter.screen.Init()
	if err != nil {
		panic(err.Error())
	}

	defer func() {
		if err := recover(); err != nil {
			interpreter.screen.Close()
			fmt.Println("-------- Internal error --------")
			fmt.Println(err)
		}
	}()

	result := interpreter.visit(node)

	interpreter.screen.Close()

	return result
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
		return res.Failure(shared.NewRuntimeErrorRange(forLoopNode.Range(), fmt.Sprintf("Invalid end value for `for` loop. Expected number got `%v`", endValue.Type())))
	}

	var stepValue Value = NewNumber(1)
	if forLoopNode.StepNode != nil {

		stepValue = res.Register(interpreter.visit(forLoopNode.StepNode))
		if res.ShouldReturn() {
			return res
		}

		if stepValue.Type() != ValueType_Number {
			return res.Failure(shared.NewRuntimeErrorRange(forLoopNode.Range(), fmt.Sprintf("Invalid step value for `for` loop. Expected number got `%v`", stepValue.Type())))
		}

		// TODO step cannot be zero

	}

	stepNumber := stepValue.(*Number).Value
	endValueNumber := endValue.(*Number).Value

	// TODO when static types are implemented, we can declare this variable as a number, and than get rid of the non number value check
	interpreter.context.AssignValueToVariable(forLoopNode.VarName, startValue, forLoopNode.StartNode.Range())

	interpreter.context.RegisterLoopEnter()

	for {

		variable, _ := interpreter.context.GetVariableAtCurrentLevel(forLoopNode.VarName)

		// TODO in the cases of the two errors below, report the last location of the last places that the variables has changed

		if variable == nil {
			return res.Failure(shared.NewRuntimeErrorRange(forLoopNode.Range(), fmt.Sprintf("Variable `%s` cannot be released during the loop", forLoopNode.VarName)))
		}

		if variable.Value.Type() != ValueType_Number {
			return res.Failure(shared.NewRuntimeErrorRange(forLoopNode.Range(), fmt.Sprintf("Variable `%s` cannot be changed to non number value during the loop", forLoopNode.VarName)))
		}

		if stepNumber > 0 {
			if variable.Value.(*Number).Value > endValueNumber {
				break
			}
		} else {
			if variable.Value.(*Number).Value < endValueNumber {
				break
			}
		}

		res.Register(interpreter.visit(forLoopNode.BodyNode))
		if res.Error != nil || res.FunctionReturnValue != nil || res.LoopShouldExit {
			break
		}

		newValue, _ := variable.Value.Add(stepValue)

		interpreter.context.AssignValueToVariable(forLoopNode.VarName, newValue, forLoopNode.StartNode.Range())

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
	} else if commandNode.CommandType == parser.CommandType_Assert {
		return interpreter.visitAssertNode(commandNode)
	} else if commandNode.CommandType == parser.CommandType_Store {
		return interpreter.visitStoreNode(commandNode)
	} else if commandNode.CommandType == parser.CommandType_Say {
		return interpreter.visitSayNode(commandNode)
	} else if commandNode.CommandType == parser.CommandType_Get {
		return interpreter.visitGetNode(commandNode)
	} else if commandNode.CommandType == parser.CommandType_Read {
		return interpreter.visitReadNode(commandNode)
	}

	panic("command interpretation not implemented yet for " + fmt.Sprint(commandNode.CommandType))
}

func (interpreter *Interpreter) visitReadNode(commandNode *parser.CommandNode) *RuntimeResult {

	res := NewRuntimeResult()

	if len(interpreter.context.ActiveGets) == 0 {
		return res.Failure(shared.NewRuntimeErrorRange(commandNode.Range(), "No gets"))
	}

	getIndex := 0

	for {

		getArgs := interpreter.context.ActiveGets[getIndex]

		variable, _ := interpreter.context.GetVariable(getArgs.VarName)

		if variable == nil {
			return res.Failure(shared.NewRuntimeErrorRange(commandNode.Range(), "Variable '"+getArgs.VarName+"' is not defined anymore"))
		}

		value := []rune{}

		if variable.Value.Type() == ValueType_String {
			value = variable.Value.(*String).Value
		} else {
			return res.Failure(shared.NewRuntimeErrorRange(commandNode.Range(), fmt.Sprintf("Variable of type %s cannot be used in a Read command", variable.Value.Type())))
		}

		interpreter.screen.SayWithModifier(getArgs.Row, getArgs.Column, value, screen.Modif_Reverse)

		interpreter.screen.ReadStr(getArgs.Row, getArgs.Column, value)

		panic("stopped here")

	}

	return res
}

func (interpreter *Interpreter) visitGetNode(commandNode *parser.CommandNode) *RuntimeResult {

	res := NewRuntimeResult()

	row := commandNode.Args["row"].(parser.Node)
	column := commandNode.Args["column"].(parser.Node)
	varName := commandNode.Args["varName"].(string)

	rowValue := res.Register(interpreter.visit(row))
	if res.ShouldReturn() {
		return res
	}

	if rowValue.Type() != ValueType_Number {
		return res.Failure(shared.NewRuntimeErrorRange(row.Range(), fmt.Sprintf("Invalid row value for get command. Expected number got `%v`", rowValue.Type())))
	}

	columnValue := res.Register(interpreter.visit(column))
	if res.ShouldReturn() {
		return res
	}

	if columnValue.Type() != ValueType_Number {
		return res.Failure(shared.NewRuntimeErrorRange(column.Range(), fmt.Sprintf("Invalid column value for get command. Expected number got `%v`", columnValue.Type())))
	}

	variable, _ := interpreter.context.GetVariable(varName)

	if variable == nil {
		return res.Failure(shared.NewRuntimeErrorRange(commandNode.Range(), "Variable '"+varName+"' is not defined"))
	}

	rowNumber := int(rowValue.(*Number).Value)
	columnNumber := int(columnValue.(*Number).Value)

	if variable.Value.Type() == ValueType_String {
		interpreter.screen.SayWithModifier(rowNumber, columnNumber, variable.Value.(*String).Value, screen.Modif_Reverse)
		interpreter.context.ActiveGets = append(interpreter.context.ActiveGets, Get{
			Row:     rowNumber,
			Column:  columnNumber,
			VarName: varName,
		})
	} else {
		return res.Failure(shared.NewRuntimeErrorRange(commandNode.Range(), fmt.Sprintf("Variable of type %s cannot be used in a Get command", variable.Value.Type())))
	}

	return res
}

func (interpreter *Interpreter) visitSayNode(commandNode *parser.CommandNode) *RuntimeResult {

	res := NewRuntimeResult()

	row := commandNode.Args["row"].(parser.Node)
	column := commandNode.Args["column"].(parser.Node)
	value := commandNode.Args["value"].(parser.Node)

	rowValue := res.Register(interpreter.visit(row))
	if res.ShouldReturn() {
		return res
	}

	if rowValue.Type() != ValueType_Number {
		return res.Failure(shared.NewRuntimeErrorRange(row.Range(), fmt.Sprintf("Invalid row value for say command. Expected number got `%v`", rowValue.Type())))
	}

	columnValue := res.Register(interpreter.visit(column))
	if res.ShouldReturn() {
		return res
	}

	if columnValue.Type() != ValueType_Number {
		return res.Failure(shared.NewRuntimeErrorRange(column.Range(), fmt.Sprintf("Invalid column value for say command. Expected number got `%v`", columnValue.Type())))
	}

	valueValue := res.Register(interpreter.visit(value))
	if res.ShouldReturn() {
		return res
	}

	if valueValue.Type() != ValueType_String {
		return res.Failure(shared.NewRuntimeErrorRange(value.Range(), fmt.Sprintf("Invalid value for say command. Expected string got `%v`", valueValue.Type())))
	}

	rowNumber := int(rowValue.(*Number).Value)
	columnNumber := int(columnValue.(*Number).Value)
	valueString := valueValue.(*String).Value

	interpreter.screen.Say(rowNumber, columnNumber, valueString)

	return res
}

func (interpreter *Interpreter) visitStoreNode(commandNode *parser.CommandNode) *RuntimeResult {

	res := NewRuntimeResult()

	valueNode := commandNode.Args["value"].(parser.Node)

	value := res.Register(interpreter.visit(valueNode))
	if res.ShouldReturn() {
		return res
	}

	for _, varName := range commandNode.Args["varNames"].([]string) {
		interpreter.context.AssignValueToVariable(varName, value, commandNode.Range())
	}

	return res
}

func (interpreter *Interpreter) visitAssertNode(commandNode *parser.CommandNode) *RuntimeResult {

	res := NewRuntimeResult()

	expr := commandNode.Args["expr"].(parser.Node)
	exprResult := res.Register(interpreter.visit(expr))

	if res.Error != nil {
		return res
	}

	if exprResult.Type() != ValueType_Boolean {
		return res.Failure(shared.NewRuntimeErrorRange(expr.Range(), fmt.Sprintf("Expected expression to return a boolean. Got `%v`", exprResult.Type())))
	}

	boolVal := exprResult.(*Boolean)

	if !boolVal.Value {

		if commandNode.Args["message"] != nil {
			return res.Failure(shared.NewAssertError(commandNode.Range(), fmt.Sprintf("Assertion failed: `%s`", commandNode.Args["message"])))
		}

		return res.Failure(shared.NewAssertError(commandNode.Range(), "Assertion failed"))
	}

	return res.Success(nil)
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

	// TODO parse macros inside strings

	return res.Success(NewString([]rune(strNode.Value)).UpdateRange(node.Range()))
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

	if leftValue.Type() != rightValue.Type() {
		return res.Failure(shared.NewRuntimeErrorRange(binOpNode.Range(), fmt.Sprintf("Operands must be of the same type (got %v and %v)", leftValue.Type(), rightValue.Type())))
	}

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

	// fmt.Println(value.PrintRepresentation())

	str := value.PrintRepresentation()

	interpreter.screen.Print([]rune(str))

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
