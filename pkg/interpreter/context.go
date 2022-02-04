package interpreter

import (
	"nova-lang/pkg/parser"
	"nova-lang/pkg/shared"
)

type Context struct {
	CurrentInterpreter *Interpreter
	Stack              []parser.Node
	CurrentLevel       int
	VariablesPerLevel  []map[string]*Variable
	LoopCountPerLevel  []int
	Functions          map[string]*parser.FunctionNode
}

func (context *Context) IncreaseLevel(node parser.Node) {

	context.CurrentLevel++

	context.LoopCountPerLevel = append(context.LoopCountPerLevel, 0)
	context.VariablesPerLevel = append(context.VariablesPerLevel, make(map[string]*Variable))

	context.Stack = append(context.Stack, node)

}

func (context *Context) DecreaseLevel() {
	context.CurrentLevel--
	context.Stack = context.Stack[:len(context.Stack)-1]
	context.LoopCountPerLevel = context.LoopCountPerLevel[:len(context.LoopCountPerLevel)-1]
	context.VariablesPerLevel = context.VariablesPerLevel[:len(context.VariablesPerLevel)-1]
}

func (context *Context) GetFunction(funcName string) *parser.FunctionNode {

	if function, ok := context.Functions[funcName]; ok {
		return function
	}

	return nil
}

func (context *Context) DeclareFunction(function *parser.FunctionNode) *shared.Error {

	if _, ok := context.Functions[function.FuncName]; ok {
		return shared.NewRuntimeErrorRange(function.Range(), "Function already declared")
	}

	context.Functions[function.FuncName] = function

	return nil
}

func (context *Context) GetVariableAtCurrentLevel(name string) (*Variable, int) {

	if value, ok := context.VariablesPerLevel[context.CurrentLevel][name]; ok {
		return value, context.CurrentLevel
	}

	return nil, -1
}

func (context *Context) GetVariable(name string) (*Variable, int) {
	for i := context.CurrentLevel; i >= 0; i-- {
		if value, ok := context.VariablesPerLevel[i][name]; ok {
			return value, i
		}
	}

	return nil, -1
}

func (context *Context) AssignValueToVariablePointer(variable *Variable, value Value) {
	variable.Value = value

	if variable.ReferenceVar != nil {
		context.AssignValueToVariablePointer(variable.ReferenceVar, value)
	}
}

func (context *Context) AssignValueToVariable(variableName string, value Value, nodeRange *shared.Range) {

	variable, _ := context.GetVariable(variableName)

	if variable == nil {
		context.DeclareVariable(variableName, Visibility_Private, nodeRange)
		variable, _ = context.GetVariable(variableName)
	}

	variable.Value = value

	if variable.ReferenceVar != nil {
		context.AssignValueToVariablePointer(variable.ReferenceVar, value)
	}

}

func (context *Context) DeclareVariable(variableName string, visibility Visibility, nodeRange *shared.Range) {

	levelToAdd := context.CurrentLevel
	if visibility == Visibility_Public {
		levelToAdd = 0
	}

	oldVar, level := context.GetVariable(variableName)

	if level == levelToAdd {
		// Already declared at the current level
	} else {

		if oldVar != nil && visibility == Visibility_Public {

			delete(context.VariablesPerLevel[level], variableName)

			variable := NewVariable(variableName, oldVar.Value, oldVar.Visibility)
			context.VariablesPerLevel[levelToAdd][variableName] = variable

		} else {

			variable := NewVariable(variableName, NewBoolean(false).UpdateRange(nodeRange), visibility)
			context.VariablesPerLevel[levelToAdd][variableName] = variable

		}

	}

}

func (context *Context) IsInsideLoop() bool {
	return context.LoopCountPerLevel[context.CurrentLevel] > 0
}

func (context *Context) RegisterLoopEnter() {
	context.LoopCountPerLevel[context.CurrentLevel]++
}

func (context *Context) RegisterLoopExit() {
	context.LoopCountPerLevel[context.CurrentLevel]--
}

func NewContext() *Context {

	context := &Context{
		CurrentLevel:      1,
		VariablesPerLevel: []map[string]*Variable{},
		LoopCountPerLevel: []int{0, 0},
		Functions:         map[string]*parser.FunctionNode{},
	}

	// context.VariablesPerLevel[0] = Public variables
	context.VariablesPerLevel = append(context.VariablesPerLevel, make(map[string]*Variable))

	// context.VariablesPerLevel[...] = Variables per level, starting from level 1
	context.VariablesPerLevel = append(context.VariablesPerLevel, make(map[string]*Variable))

	return context
}
