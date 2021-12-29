package interpreter

type SymbolTable struct {
	symbols map[string]Value
}

func (symbolTable *SymbolTable) Get(name string) (Value, bool) {
	value, found := symbolTable.symbols[name]
	return value, found
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		symbols: make(map[string]Value),
	}
}
