package types

// Integer is an exported constant of type integer
var Integer = &integer{}

// Float is an exported constant of type float
var Float = &float{}

// Boolean is an exported constant of type boolean
var Boolean = &boolean{}

// Pict is shorthand for [,]{float, float, float, float}
var Pict = &Array{
	Inner: &Tuple{
		Types: []Type{Float, Float, Float, Float},
	},
	Rank: 2,
}
