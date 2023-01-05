package nasm

import (
	"bytes"
	"fmt"

	"github.com/ryan-berger/jpl/internal/ast"
)

type dataType int

const (
	number dataType = iota
	str
)

type data struct {
	Name  string
	Value string
	Type  dataType
}

type namer func() string
type constantMapper map[ast.Node]string

// TODO: refactor later
func dataWalk(n ast.Node, namer namer, names constantMapper, reverse map[string]string) []data {
	switch node := n.(type) {
	case ast.Program:
		var datum []data
		for _, cmd := range node {
			datum = append(datum, dataWalk(cmd, namer, names, reverse)...)
		}
		return datum
	case *ast.Function:
		var datum []data
		for _, stmt := range node.Statements {
			datum = append(datum, dataWalk(stmt, namer, names, reverse)...)
		}
		return datum
	case *ast.Show:
		typeStr := node.Expr.
			Typ().
			String()

		if k, ok := reverse[typeStr]; ok {
			names[node] = k
			return nil
		}

		name := namer()
		names[node] = name
		reverse[typeStr] = name
		return []data{{Name: name, Value: typeStr, Type: str}}
	case *ast.Read:
		val := node.Src.String()
		if k, ok := reverse[val]; ok {
			names[node] = k
			return nil
		}

		name := namer()
		reverse[val] = name
		names[node] = name

		return []data{{Name: name, Value: val, Type: str}}
	case *ast.Write:
		val := node.Dest.String()
		if k, ok := reverse[val]; ok {
			names[node] = k
			return nil
		}

		name := namer()
		reverse[val] = name
		names[node] = name

		return []data{{Name: name, Value: val, Type: str}}
	case *ast.LetStatement:
		switch exp := node.Expr.(type) {
		case *ast.FloatExpression, *ast.IntExpression:
			val := node.Expr.String()

			if k, ok := reverse[val]; ok {
				names[node] = k
				return nil
			}

			name := namer()
			names[node] = name
			reverse[val] = name
			return []data{{Name: name, Value: val}}
		case *ast.InfixExpression:
			return append(dataWalk(exp.Left, namer, names, reverse), dataWalk(exp.Left, namer, names, reverse)...)
		case *ast.PrefixExpression:
			return dataWalk(exp.Expr, namer, names, reverse)
		case *ast.ArrayRefExpression:
			datum := dataWalk(exp.Array, namer, names, reverse)
			for _, e := range exp.Indexes {
				datum = append(datum, dataWalk(e, namer, names, reverse)...)
			}
			return datum
		case *ast.ArrayExpression:
			var datum []data
			for _, e := range exp.Expressions {
				datum = append(datum, dataWalk(e, namer, names, reverse)...)
			}
			return datum
		case *ast.TupleExpression:
			var datum []data
			for _, e := range exp.Expressions {
				datum = append(datum, dataWalk(e, namer, names, reverse)...)
			}
			return datum
		case *ast.TupleRefExpression:
			val := fmt.Sprintf("%d", exp.Index)

			if k, ok := reverse[val]; ok {
				names[node] = k
				return nil
			}

			name := namer()
			names[node] = name
			reverse[val] = name
			return append(
				dataWalk(exp.Tuple, namer, names, reverse),
				data{Name: name, Value: val})
		case *ast.ArrayTransform:
			var datum []data
			for _, b := range exp.OpBindings {
				datum = append(datum, dataWalk(b.Expr, namer, names, reverse)...)
			}
			datum = append(datum, dataWalk(exp.Expr, namer, names, reverse)...)
			return datum
		case *ast.SumTransform:
			var datum []data
			for _, b := range exp.OpBindings {
				datum = append(datum, dataWalk(b.Expr, namer, names, reverse)...)
			}
			datum = append(datum, dataWalk(exp.Expr, namer, names, reverse)...)
			return datum
		case *ast.IfExpression:
			var datum []data
			datum = append(datum, dataWalk(exp.Condition, namer, names, reverse)...)
			datum = append(datum, dataWalk(exp.Consequence, namer, names, reverse)...)
			datum = append(datum, dataWalk(exp.Otherwise, namer, names, reverse)...)
			return datum
		}
		return nil
	case *ast.AssertStatement:
		val := fmt.Sprintf(`%s\n`, node.Message[1:len(node.Message)-1])
		if k, ok := reverse[val]; ok {
			names[node] = k
			return nil
		}

		name := namer()
		names[node] = name
		reverse[val] = name
		return []data{{Name: name, Value: val, Type: str}}
	case *ast.Print:
		val := fmt.Sprintf(`%s\n`, node.Str[1:len(node.Str)-1])
		if k, ok := reverse[val]; ok {
			names[node] = k
			return nil
		}

		name := namer()
		names[node] = name
		reverse[val] = name
		return []data{{Name: name, Value: val, Type: str}}
	default:
		return nil
	}
}

func dataSection(program ast.Program) (string, constantMapper) {
	counter := -1
	namer := func() string {
		counter++
		return fmt.Sprintf("const%d", counter)
	}

	mapper := make(constantMapper)
	reverse := make(map[string]string)

	datum := dataWalk(program, namer, mapper, reverse)
	buf := bytes.NewBufferString("section .data\n")

	for _, d := range datum {
		if d.Type == number {
			buf.WriteString(fmt.Sprintf("\t%s: dq %s\n", d.Name, d.Value))
		}

		if d.Type == str {
			buf.WriteString(fmt.Sprintf("\t%s: db `%s`, 0\n", d.Name, d.Value))
		}
	}
	buf.WriteByte('\n')
	return buf.String(), mapper
}
