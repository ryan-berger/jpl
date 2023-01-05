package types

// Integer is an exported constant of type integer
var Integer = &integer{}

// Str is an exported constant of type str
var Str = &str{}

// Float is an exported constant of type float
var Float = &float{}

// Boolean is an exported constant of type boolean
var Boolean = &boolean{}

// Pict is shorthand for [,]float4
var Pict = &Array{
	Inner: &Tuple{
		[]Type{Float, Float, Float, Float},
	},
	Rank: 2,
}
