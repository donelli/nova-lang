package interpreter

type Value interface {
	Copy() Value
}
