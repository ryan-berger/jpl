package types

type Type interface {
	Size() int
	Equal(other Type) bool
}

type boolean struct{}

func (b *boolean) Size() int {
	return 4
}

func (b *boolean) Equal(other Type) bool {
	_, ok := other.(*boolean)
	return ok
}

type integer struct{}

func (i *integer) Size() int {
	return 8
}
func (i *integer) Equal(other Type) bool {
	_, ok := other.(*integer)
	return ok
}

type float struct{}

func (f *float) Size() int {
	return 8
}

func (f *float) Equal(other Type) bool {
	_, ok := other.(*float)
	return ok
}

type Array struct {
	Inner Type
	Rank  int
}

func (a *Array) Size() int {
	return 8
}

func (a *Array) Equal(other Type) bool {
	arr, ok := other.(*Array)
	return ok && a.Rank == arr.Rank && a.Inner.Equal(arr.Inner)
}

type Tuple struct {
	Types []Type
}

func (t *Tuple) Size() int {
	size := 0
	for _, typ := range t.Types {
		size += typ.Size()
	}
	return size
}

func (t *Tuple) Equal(other Type) bool {
	tup, ok := other.(*Tuple)
	if !ok {
		return false
	}

	if len(t.Types) != len(tup.Types) {
		return false
	}

	for i, typ := range t.Types {
		if !typ.Equal(tup.Types[i]) {
			return false
		}
	}

	return true
}
