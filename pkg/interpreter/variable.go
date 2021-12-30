package interpreter

type Visibility int8

const (
	Visibility_Public Visibility = iota + 1
	Visibility_Private
)

type Variable struct {
	Name  string
	Value Value
	// Level      int
	Visibility Visibility
}

func NewVariable(name string, value Value, visibility Visibility) *Variable {
	return &Variable{
		Name:  name,
		Value: value,
		// Level:      level,
		Visibility: visibility,
	}
}
