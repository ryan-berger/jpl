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
var Pict = &Array{
	Inner: &Tuple{
		Types: []Type{Integer, Integer, Integer, Float},
	},
	Rank: 2,
}

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
			Args:   []Type{Pict, Integer, Integer},
			Return: Boolean,
		},
		"sepia": &Function{
			Args:   []Type{Pict},
			Return: Pict,
		},
		"blur": &Function{
			Args:   []Type{Pict, Float},
			Return: Pict,
		},
		"crop": &Function{
			Args:   []Type{Pict, Integer, Integer, Integer, Integer},
			Return: Pict,
		},
		"float": &Function{
			Args:   []Type{Integer},
			Return: Float,
		},
		"int": &Function{
			Args:   []Type{Float},
			Return: Integer,
		},
		"read.image": &Function{
			Args: []Type{
				&Array{
					Inner: Integer,
					Rank:  1,
				},
			},
			Return: Pict,
		},
	}
}
