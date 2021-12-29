package interpreter

type Number struct {
	Value float64
}

func (n *Number) Copy() Number {
	return Number{Value: n.Value}
}
