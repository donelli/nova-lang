package interpreter

type Context struct {
	Level       int
	SymbolTable *SymbolTable
}

func NewContext() *Context {
	return &Context{
		Level:       0,
		SymbolTable: NewSymbolTable(),
	}
}
