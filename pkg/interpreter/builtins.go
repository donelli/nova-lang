package interpreter

import (
	"fmt"
	"nova-lang/pkg/lexer"
	"nova-lang/pkg/parser"
	"nova-lang/pkg/shared"
	"strconv"
	"strings"
	"time"
)

type BuiltInFunction interface {
	Call(context *Context, args []Value) *RuntimeResult
}

// TODO create a map with the function contract (parameters and types), and check it on the function call on the interpreter

func checkParameters(funcCallRange *shared.Range, expectedArgTypes []ValueType, args []Value, funcName string) *shared.Error {

	if len(args) != len(expectedArgTypes) {
		return shared.NewRuntimeErrorRange(funcCallRange, fmt.Sprintf("Expected %d arguments in function `%s`, got %d", len(expectedArgTypes), funcName, len(args)))
	}

	for argIndex, arg := range args {
		if arg.Type() != expectedArgTypes[argIndex] {
			return shared.NewRuntimeErrorRange(funcCallRange, fmt.Sprintf("Expected `%v` for argument %d in function `%s`, got `%v`", expectedArgTypes[argIndex], argIndex, funcName, arg.Type()))
		}
	}

	return nil
}

var BuiltInFunctions map[string]func(context *Context, funcCallRange *shared.Range, args []Value) *RuntimeResult

func InitBuiltInFunctions() {
	BuiltInFunctions = map[string]func(context *Context, funcCallRange *shared.Range, args []Value) *RuntimeResult{
		"alltrim": BuiltIn_Alltrim,
		"str":     BuiltIn_Str,
		"sleep":   BuiltIn_Sleep,
		"type":    BuiltIn_Type,
		"val":     BuiltIn_Val,
		"empty":   BuiltIn_Empty,
	}
}

// --------------------------------
// Start functions
// --------------------------------

func BuiltIn_Alltrim(context *Context, funcCallRange *shared.Range, args []Value) *RuntimeResult {

	res := NewRuntimeResult()

	if err := checkParameters(funcCallRange, []ValueType{ValueType_String}, args, "alltrim"); err != nil {
		return res.Failure(err)
	}

	str := args[0].(*String)

	return res.SuccessReturn(NewString(strings.Trim(str.Value, " ")))
}

func BuiltIn_Str(context *Context, funcCallRange *shared.Range, args []Value) *RuntimeResult {

	res := NewRuntimeResult()
	var err *shared.Error

	if len(args) == 1 {
		err = checkParameters(funcCallRange, []ValueType{ValueType_Number}, args, "str")
	} else if len(args) == 2 {
		err = checkParameters(funcCallRange, []ValueType{ValueType_Number, ValueType_Number}, args, "str")
	} else if len(args) == 3 {
		err = checkParameters(funcCallRange, []ValueType{ValueType_Number, ValueType_Number, ValueType_Number}, args, "str")
	} else {
		err = shared.NewRuntimeErrorRange(funcCallRange, "Expected 1-3 arguments for `str` function")
	}

	if err != nil {
		return res.Failure(err)
	}

	number := args[0].(*Number).Value
	width := 10
	decimalPlaces := 0

	if len(args) > 1 {
		width = int(args[1].(*Number).Value)
	}

	if len(args) == 3 {
		decimalPlaces = int(args[2].(*Number).Value)
	}

	// Format decimal places
	str := fmt.Sprintf("%."+fmt.Sprintf("%d", decimalPlaces)+"f", number)

	// Format width
	str = strings.Repeat(" ", width-len(str)) + str

	return res.SuccessReturn(NewString(str))

}

func BuiltIn_Sleep(context *Context, funcCallRange *shared.Range, args []Value) *RuntimeResult {

	res := NewRuntimeResult()

	if err := checkParameters(funcCallRange, []ValueType{ValueType_Number}, args, "sleep"); err != nil {
		return res.Failure(err)
	}

	seconds := int(args[0].(*Number).Value)

	time.Sleep(time.Second * time.Duration(seconds))

	return res.SuccessReturn(NewBoolean(false))
}

func BuiltIn_Type(context *Context, funcCallRange *shared.Range, args []Value) *RuntimeResult {

	res := NewRuntimeResult()

	if err := checkParameters(funcCallRange, []ValueType{ValueType_String}, args, "type"); err != nil {
		return res.Failure(err)
	}

	val := execEmbeddedProgram(context, args[0].(*String).Value)

	return res.SuccessReturn(val)
}

func execEmbeddedProgram(context *Context, program string) Value {

	// TODO dont allow macros inside macros

	// TODO This function takes some times to execute, needs to be optimized

	lex := lexer.NewLexer("embedded", program)
	lexRes := lex.Parse()

	if len(lexRes.Errors) > 0 {
		return NewString("U")
	}

	parser := parser.NewParser(lexRes, true)
	parseRes := parser.Parse()

	if parseRes.Err != nil {
		return NewString("U")
	}

	rtRes := context.CurrentInterpreter.visit(parseRes.Node)
	if rtRes.Error != nil {
		return NewString("U")
	}

	if rtRes.Value.Type() == ValueType_String {
		return NewString("C")
	} else if rtRes.Value.Type() == ValueType_Number {
		return NewString("N")
	} else if rtRes.Value.Type() == ValueType_Boolean {
		return NewString("L")
	}

	return NewString("U")

}

func BuiltIn_Val(context *Context, funcCallRange *shared.Range, args []Value) *RuntimeResult {

	res := NewRuntimeResult()

	if err := checkParameters(funcCallRange, []ValueType{ValueType_String}, args, "val"); err != nil {
		return res.Failure(err)
	}

	str := args[0].(*String).Value

	strToConvert := ""

	for _, char := range str {

		if !strings.Contains(shared.DigitsAndDot, string(char)) {
			break
		}

		strToConvert += string(char)

	}

	convertedVal, err := strconv.ParseFloat(strToConvert, 64)

	if err != nil {
		return res.SuccessReturn(NewNumber(0))
	} else {
		return res.SuccessReturn(NewNumber(convertedVal))
	}

}

func BuiltIn_Empty(context *Context, funcCallRange *shared.Range, args []Value) *RuntimeResult {

	res := NewRuntimeResult()

	if len(args) != 1 {
		return res.Failure(shared.NewRuntimeErrorRange(funcCallRange, fmt.Sprintf("Expected 1 argument in function `empty`, got %d", len(args))))
	}

	return res.SuccessReturn(NewBoolean(args[0].IsEmpty()))

}
