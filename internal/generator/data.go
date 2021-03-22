package generator

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

func isFloat(expression ast.Expression) bool {
	_, ok := expression.(*ast.FloatExpression)
	return ok
}

func isInteger(expression ast.Expression) bool {
	_, ok := expression.(*ast.IntExpression)
	return ok
}

type namer func() string
type constantMapper map[ast.Node]string

func dataWalk(n ast.Node, namer namer, names constantMapper) []data {
	switch node := n.(type) {
	case ast.Program:
		var datum []data
		for _, cmd := range node {
			datum = append(datum, dataWalk(cmd, namer, names)...)
		}
		return datum
	case *ast.Show:
		name := namer()
		names[node] = name
		typeStr := node.Expr.
			Typ().
			String()

		return []data{{Name: name, Value: typeStr, Type: str}}
	case *ast.Read:
		name := namer()
		names[node] = name

		return []data{{Name: name, Value: node.Src[1 : len(node.Src)-1], Type: str}}
	case *ast.Write:
		name := namer()
		names[node] = name

		return []data{{Name: name, Value: node.Dest[1 : len(node.Dest)-1], Type: str}}
	case *ast.LetStatement:
		if isFloat(node.Expr) || isInteger(node.Expr) {
			name := namer()
			names[node] = name

			return []data{{Name: name, Value: node.Expr.String()}}
		}
		return nil
	case *ast.AssertStatement:
		name := namer()
		names[node] = name
		return []data{
			{
				Name:  name,
				Value: fmt.Sprintf(`%s\n`, node.Message[1:len(node.Message)-1]),
				Type:  str,
			},
		}
	case *ast.Print:
		name := namer()
		names[node] = name
		return []data{
			{
				Name:  name,
				Value: fmt.Sprintf(`%s\n`, node.Str[1:len(node.Str)-1]),
				Type:  str,
			},
		}
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

	datum := dataWalk(program, namer, mapper)
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
