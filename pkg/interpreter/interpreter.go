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

func (interpreter *Interpreter) Visit(node parser.Node) *RuntimeResult {

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
	}

	panic("not implemented yet for " + fmt.Sprint(node.Type()))

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

		conditionValue := res.Register(interpreter.Visit(ifCase.CaseExpr))
		if res.ShouldReturn() {
			return res
		}

		isTrue, err := conditionValue.IsTrue()
		if err != nil {
			return res.Failure(err)
		}

		if isTrue {
			res.Register(interpreter.Visit(ifCase.Body))
			if res.ShouldReturn() {
				return res
			}
			return res.Success(nil)
		}

	}

	if ifNode.ElseCase != nil {
		res.Register(interpreter.Visit(ifNode.ElseCase))
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

	value := res.Register(interpreter.Visit(unaryOperNode.Node))
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
	value := res.Register(interpreter.Visit(varAssignNode.Expr))
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
		returnValue = res.Register(interpreter.Visit(returnNode.Expr))
		if res.ShouldReturn() {
			return res
		}
	}

	return res.SuccessReturn(returnValue)
}

func (interpreter *Interpreter) visitBinaryOperationNode(node parser.Node) *RuntimeResult {

	binOpNode := node.(*parser.BinaryOperationNode)
	res := NewRuntimeResult()

	leftValue := res.Register(interpreter.Visit(binOpNode.LeftNode))
	if res.ShouldReturn() {
		return res
	}

	rightValue := res.Register(interpreter.Visit(binOpNode.RightNode))
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

	value := res.Register(interpreter.Visit(printNode.Expr))
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

		res.Register(interpreter.Visit(node))
		if res.ShouldReturn() {
			return res
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
