package parser

import (
	"github.com/ryan-berger/jpl/internal/lexer"
	"github.com/ryan-berger/jpl/internal/types"
)

// parseTypeExpression parses a type expression of the form:
// type : int
// | bool
// | float
// | float3
// | float4
// | <type> [ , ... ]
// | { <type> , ... }
func (p *parser) parseTypeExpression() types.Type {
	t := p.parseType()

	// handle an array type
	for p.expectPeek(lexer.LBrace) {
		rank := 1
		// as long as we see a comma, increment rank
		for p.expectPeek(lexer.Comma) { // TODO: make sure we throw errors here (once errors are ready)
			rank++
		}

		// there should be an ']' as we are in a array type
		if !p.expectPeek(lexer.RBrace) {
			return nil
		}
		t = &types.Array{
			Inner: t,
			Rank:  rank,
		}
	}

	return t
}

func (p *parser) parseType() types.Type {
	// if we have an '{' we could be at the start of a normal type
	if p.curTokenIs(lexer.LCurly) {
		return p.parseTupleType()
	}

	// if we aren't in a tuple type, we need to establish a "base type"
	switch p.cur.Type {
	// if it is a base type, use a map to types.Type
	case lexer.Float, lexer.Int, lexer.Bool:
		return tokenToType[p.cur.Type]
	case lexer.Float3: // expand syntactic sugar for {float, float float}
		return  &types.Tuple{
			Types: []types.Type{types.Float, types.Float, types.Float},
		}
	case lexer.Float4: // expand syntactic sugar for {float, float float, float}
		return &types.Tuple{
			Types: []types.Type{types.Float, types.Float, types.Float, types.Float},
		}
	default:
		p.errorf(p.cur,"expected type received %s", p.cur.Val)
		return nil
	}

}

// parseTupleType is called whenever a tuple type definition needs parsing
func (p *parser) parseTupleType() types.Type {
	tupleType := &types.Tuple{}

	// use parseList to parse types until the next RCurly
	ok := p.parseList(lexer.RCurly, func() bool {
		typ := p.parseTypeExpression()
		if typ == nil {
			return false
		}
		tupleType.Types = append(tupleType.Types, typ)
		return true
	})

	if !ok {
		return nil
	}

	return tupleType
}
