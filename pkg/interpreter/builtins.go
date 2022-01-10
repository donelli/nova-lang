package interpreter

import (
	"fmt"
	"nova-lang/pkg/shared"
	"strings"
	"time"
)

type BuiltInFunction interface {
	Call(context *Context, args []Value) *RuntimeResult
}

func checkParameters(funcCallRange *shared.Range, expectedArgTypes []ValueType, args []Value) *shared.Error {

	if len(args) != len(expectedArgTypes) {
		return shared.NewRuntimeErrorRange(funcCallRange, fmt.Sprintf("Expected %d arguments, got %d", len(expectedArgTypes), len(args)))
	}

	for argIndex, arg := range args {
		if arg.Type() != expectedArgTypes[argIndex] {
			return shared.NewRuntimeErrorRange(funcCallRange, fmt.Sprintf("Expected `%v` for argument %d, got `%v`", expectedArgTypes[argIndex], argIndex, arg.Type()))
		}
	}

	return nil
}

var BuiltInFunctions = map[string]func(context *Context, funcCallRange *shared.Range, args []Value) *RuntimeResult{
	"alltrim": BuiltIn_Alltrim,
	"str":     BuiltIn_Str,
	"sleep":   BuiltIn_Sleep,
}

// --------------------------------
// Start functions
// --------------------------------

// The ALLTRIM() function returns a string with all leading and trailing blanks removed.
func BuiltIn_Alltrim(context *Context, funcCallRange *shared.Range, args []Value) *RuntimeResult {

	res := NewRuntimeResult()

	if err := checkParameters(funcCallRange, []ValueType{ValueType_String}, args); err != nil {
		return res.Failure(err)
	}

	str := args[0].(*String)

	return res.SuccessReturn(NewString(strings.Trim(str.Value, " ")))
}

// The STR() function converts the numeric expression <expN1> to a character string of width <expN2>
// with <expN3> decimal places. If <expN3> is not specified then <expN1> is treated as an integer.
func BuiltIn_Str(context *Context, funcCallRange *shared.Range, args []Value) *RuntimeResult {

	res := NewRuntimeResult()
	var err *shared.Error

	if len(args) == 2 {
		err = checkParameters(funcCallRange, []ValueType{ValueType_Number, ValueType_Number}, args)
	} else if len(args) == 3 {
		err = checkParameters(funcCallRange, []ValueType{ValueType_Number, ValueType_Number, ValueType_Number}, args)
	} else {
		err = shared.NewRuntimeErrorRange(funcCallRange, "Expected 2 or 3 arguments")
	}

	if err != nil {
		return res.Failure(err)
	}

	number := args[0].(*Number).Value
	width := int(args[1].(*Number).Value)
	decimalPlaces := 0

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

	if err := checkParameters(funcCallRange, []ValueType{ValueType_Number}, args); err != nil {
		return res.Failure(err)
	}

	seconds := int(args[0].(*Number).Value)

	time.Sleep(time.Second * time.Duration(seconds))

	return res.SuccessReturn(NewBoolean(false))
}
