package interpreter

type Visibility int8

const (
	Visibility_Public Visibility = iota + 1
	Visibility_Private
)

type Variable struct {
	Name         string
	Value        Value
	Visibility   Visibility
	ReferenceVar *Variable
}

func NewVariable(name string, value Value, visibility Visibility) *Variable {
	return &Variable{
		Name:         name,
		Value:        value,
		Visibility:   visibility,
		ReferenceVar: nil,
	}
}
