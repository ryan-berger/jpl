package meta

type Locationer interface {
	Loc() (int, int)
}