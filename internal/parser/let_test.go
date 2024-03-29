package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ryan-berger/jpl/internal/lexer"
)

func TestParse_Let(t *testing.T) {
	test := `let {x[W,H], y, z} = {(10*3)+10, 3, 4}`

	tokens, ok := lexer.Lex(test)
	assert.True(t, ok)


	fmt.Println(Parse(tokens))
}

func TestParse_Fn(t *testing.T) {
	test := `// OK

/* This generates a mandlebrot set image */

fn mandlestep({ zx : float, zy : float}, {cx : float, cy : float}): {float,float} {
   return { \
     (zx * zx) - (zy * zy) + cx, \
     (2.0 * zx * zy) + cy \
   }
}

fn mandlecount({ zx : float, zy : float}, c : {float,float}, max : int): int {
   return if max == 0 \
   then 0 \
   else if (zx * zx + zy * zy) > 2.0 \
   then 0 \
   else 1 + mandlecount(mandlestep({zx, zy}, c), c, max - 1)
}

fn colormap(c : float) : float4 {
   return if c < 0.5 \
   then { 0.0, 0.0, 2.0 * c, 1.0 } \
   else { (2.0 * c) - 1.0, (2.0 * c) - 1.0, 1.0, 1.0}
}

fn mandlepixel(c : {float,float}, max : int) : float4 {
   let count = mandlecount({0.0, 0.0}, c, max)
   let c2 = if count == max then 0.0 else (float(count + 1) / float(max))
   return colormap(c2)
}

let iter = if argnum > 1 then float(args[1]) / 60.0 else 0.0
let {xlo, xhi} = {-1.5, 1.0 - 2.3 * iter}
let {ylo, yhi} = {-1.0 + 0.9 * iter, 1.0 - 0.9 * iter}
let {W, H} = {500, 400}
time let mandlebrot = array[i : W, j : H] mandlepixel( \
    { xlo + float(i) / float(W) * (xhi - xlo) \
    , ylo + float(j) / float(H) * (yhi - ylo)}, \
    25)
write image mandlebrot to "mandlebrot.png"
`
	tokens, ok := lexer.Lex(test)
	assert.True(t, ok)

	cmds, err := Parse(tokens)
	fmt.Println(cmds)
	assert.Nil(t, err)
}

func TestParse_Cmd(t *testing.T) {
	test := `show x + 10`

	tokens, ok := lexer.Lex(test)
	assert.True(t, ok)

	Parse(tokens)
}
