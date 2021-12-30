package interpreter

import "nova-lang/pkg/shared"

type Context struct {
	CurrentLevel      int
	VariablesPerLevel []map[string]*Variable
	LoopCountPerLevel []int
}

func (context *Context) GetVariable(name string) (*Variable, int) {
	for i := context.CurrentLevel; i >= 0; i-- {
		if value, ok := context.VariablesPerLevel[i][name]; ok {
			return value, i
		}
	}

	return nil, -1
}

func (context *Context) AssignValueToVariable(variableName string, value Value, nodeRange *shared.Range) {

	variable, _ := context.GetVariable(variableName)

	if variable == nil {
		context.DeclareVariable(variableName, Visibility_Private, nodeRange)
		variable, _ = context.GetVariable(variableName)
	}

	variable.Value = value

}

func (context *Context) DeclareVariable(variableName string, visibility Visibility, nodeRange *shared.Range) {

	levelToAdd := context.CurrentLevel
	if visibility == Visibility_Public {
		levelToAdd = 0
	}

	_, level := context.GetVariable(variableName)

	if level == levelToAdd {
		// Already declared at the current level
	} else {

		variable := NewVariable(variableName, NewBoolean(false).UpdateRange(nodeRange), visibility)
		context.VariablesPerLevel[levelToAdd][variableName] = variable

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
	}

	// context.VariablesPerLevel[0] = Public variables
	context.VariablesPerLevel = append(context.VariablesPerLevel, make(map[string]*Variable))

	// context.VariablesPerLevel[...] = Variables per level, starting from level 1
	context.VariablesPerLevel = append(context.VariablesPerLevel, make(map[string]*Variable))

	return context
}
