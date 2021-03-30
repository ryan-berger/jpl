package optimizer

import "github.com/ryan-berger/jpl/internal/ast"

func doesTerminate(fn *ast.Function) bool {

	return true
}

func inline(call ast.CallExpression, fn *ast.Function) {

}

func Inliner(p ast.Program) {
	fns := make(map[string]bool)

	for _, cmd := range p {
		if fn, ok := cmd.(*ast.Function); ok {
			fns[fn.Var] = doesTerminate(fn)
		}
	}


}
