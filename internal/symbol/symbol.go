package symbol

import "github.com/ryan-berger/jpl/internal/types"

type Table map[string]Symbol

type Symbol interface {
	symbol()
}

type Identifier struct {
	Type types.Type
}

func (i *Identifier) symbol() {}

type Function struct {
	Args   []types.Type
	Return types.Type
}

func (f *Function) symbol() {}

func (s Table) Copy() Table {
	newTable := make(Table, len(s))
	for k, v := range s {
		newTable[k] = v
	}
	return newTable
}

func NewSymbolTable() Table {
	return map[string]Symbol{
		"sub_ints": &Function{
			Args:   []types.Type{types.Integer, types.Integer},
			Return: types.Integer,
		},
		"sub_floats": &Function{
			Args:   []types.Type{types.Float, types.Float},
			Return: types.Float,
		},
		"has_size": &Function{
			Args:   []types.Type{types.Pict, types.Integer, types.Integer},
			Return: types.Boolean,
		},
		"sepia": &Function{
			Args:   []types.Type{types.Pict},
			Return: types.Pict,
		},
		"resize": &Function{
			Args:   []types.Type{types.Pict, types.Integer, types.Integer},
			Return: types.Pict,
		},
		"blur": &Function{
			Args:   []types.Type{types.Pict, types.Float},
			Return: types.Pict,
		},
		"crop": &Function{
			Args:   []types.Type{types.Pict, types.Integer, types.Integer, types.Integer, types.Integer},
			Return: types.Pict,
		},
		"float": &Function{
			Args:   []types.Type{types.Integer},
			Return: types.Float,
		},
		"int": &Function{
			Args:   []types.Type{types.Float},
			Return: types.Integer,
		},
		"read.image": &Function{
			Args: []types.Type{
				&types.Array{
					Inner: types.Integer,
					Rank:  1,
				},
			},
			Return: types.Pict,
		},
	}
}
