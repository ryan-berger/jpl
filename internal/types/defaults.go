package types

var Integer = &integer{}
var Float = &float{}
var Boolean = &boolean{}
var Pict = &Array{
	Inner: &Tuple{
		Types: []Type{Integer, Integer, Integer, Float},
	},
	Rank: 2,
}
