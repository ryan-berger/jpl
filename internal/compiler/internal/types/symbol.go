package types

type SymbolTable map[string]Symbol

type Symbol interface {
	symbol()
}

type Identifier struct {
	Type Type
}

func (i *Identifier) symbol() {}

type Function struct {
	Args   []Type
	Return Type
}

func (f *Function) symbol() {}

func (s SymbolTable) Copy() SymbolTable {
	newTable := make(SymbolTable, len(s))
	for k, v := range s {
		newTable[k] = v
	}
	return newTable
}

var Integer = &integer{}
var Float = &float{}
var Boolean = &boolean{}
var pict = &Pict{}

func NewSymbolTable() SymbolTable {
	return map[string]Symbol{
		"add_ints": &Function{
			Args:   []Type{Integer, Integer},
			Return: Integer,
		},
		"add_floats": &Function{
			Args:   []Type{Float, Float},
			Return: Float,
		},
		"has_size": &Function{
			Args:   []Type{pict, Integer, Integer},
			Return: Boolean,
		},
		"sepia": &Function{
			Args:   []Type{pict},
			Return: pict,
		},
		"blur": &Function{
			Args:   []Type{pict, Float},
			Return: pict,
		},
		"crop": &Function{
			Args:   []Type{pict, Integer, Integer, Integer, Integer},
			Return: pict,
		},
	}
}
