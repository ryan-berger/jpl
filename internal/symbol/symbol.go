package symbol

import (
	 "github.com/ryan-berger/jpl/internal/ast/types"
)

type Table struct {
	symbols map[string]Symbol
	parent  *Table
}

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

func (t *Table) Copy() *Table {
	tbl := &Table{
		symbols: make(map[string]Symbol),
		parent:  t,
	}

	return tbl
}

func (t *Table) Get(ident string) (Symbol, bool) {
	if sym, ok := t.symbols[ident]; ok {
		return sym, true
	}

	if t.parent == nil {
		return nil, false
	}

	return t.parent.Get(ident)
}

func (t *Table) Set(ident string, sym Symbol) {
	t.symbols[ident] = sym
}


func NewSymbolTable() *Table {
	return &Table{
		symbols: map[string]Symbol{
			"argnum": &Identifier{Type: types.Integer},
			"args": &Identifier{
				Type: &types.Array{
					Inner: types.Integer,
					Rank:  1,
				},
			},
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
			"get_time": &Function{
				Return: types.Float,
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
		},
		parent:  nil,
	}
}
