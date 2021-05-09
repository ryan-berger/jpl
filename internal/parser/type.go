package parser

import (
	types "github.com/ryan-berger/jpl/internal/ast/types"
	"github.com/ryan-berger/jpl/internal/lexer"
)

// parseTypeExpression parses a type expression of the form:
// type : int
// | bool
// | float
// | float3
// | float4
// | <type> [ , ... ]
// | { <type> , ... }
func (p *parser) parseTypeExpression() (types.Type, error) {
	t, err := p.parseType()
	if err != nil {
		return nil, err
	}

	// handle an array type
	for p.expectPeek(lexer.LBrace) {
		rank := 1
		// as long as we see a comma, increment rank
		for p.expectPeek(lexer.Comma) { // TODO: make sure we throw errors here (once errors are ready)
			rank++
		}

		// there should be an ']' as we are in a array type
		if !p.expectPeek(lexer.RBrace) {
			return nil, p.errorf(p.peek, "expected ']' received %s", p.peek.Val)
		}
		t = &types.Array{
			Inner: t,
			Rank:  rank,
		}
	}

	return t, nil
}

func (p *parser) parseType() (types.Type, error) {
	// if we have an '{' we could be at the start of a normal type
	if p.curTokenIs(lexer.LCurly) {
		return p.parseTupleType()
	}

	// if we aren't in a tuple type, we need to establish a "base type"
	switch p.cur.Type {
	// if it is a base type, use a map to types.Type
	case lexer.Float, lexer.Int, lexer.Bool:
		return tokenToType[p.cur.Type], nil
	case lexer.Float3: // flatten syntactic sugar for {float, float float}
		return &types.Tuple{
			Types: []types.Type{types.Float, types.Float, types.Float},
		}, nil
	case lexer.Float4: // flatten syntactic sugar for {float, float float, float}
		return &types.Tuple{
			Types: []types.Type{types.Float, types.Float, types.Float, types.Float},
		}, nil
	default:
		return nil, p.errorf(p.cur, "expected type received %s", p.cur.Val)
	}

}

// parseTupleType is called whenever a tuple type definition needs parsing
func (p *parser) parseTupleType() (types.Type, error) {
	tupleType := &types.Tuple{}

	// use parseList to parse types until the next RCurly
	listErr := p.parseList(lexer.RCurly, func() error {
		typ, err := p.parseTypeExpression()
		if err != nil {
			return err
		}
		tupleType.Types = append(tupleType.Types, typ)
		return nil
	})

	if listErr != nil {
		return nil, listErr
	}

	return tupleType, nil
}
