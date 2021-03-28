


# 1. Expanded struct/interface definitions:

## location.go:
```go
type Location struct {
	Line, Pos int
}

func (l Location) Loc() (int, int) {
	return l.Line, l.Pos
}
```

## arguments.go:
```go
type Argument interface {
	LValue
	argument()
}

type VariableArgument struct {
	Variable string
	Type     types.Type
	Location
}

func (v *VariableArgument) SExpr() string {
	if v.Type != nil {
		return fmt.Sprintf("(VarArgument %s %s)", v.Type, v.Variable)
	}
	return fmt.Sprintf("(VarArgument %s)", v.Variable)
}

func (v *VariableArgument) String() string { return v.Variable }
func (v *VariableArgument) argument()      {}
func (v *VariableArgument) lValue()        {}

type VariableArr struct {
	Variable  string
	Variables []string
	Type      types.Type
	Location
}
```

## binding.go
```go
type Binding interface {
	Node
	binding()
}

type TypeBind struct {
	Argument Argument
	Type     types.Type
	Location
}

func (b *TypeBind) SExpr() string {
	panic("implement me")
}
```

## command.go
```go
type Command interface {
	Node
	command()
}

type Function struct {
	Var        string
	Bindings   []Binding
	ReturnType types.Type
	Statements []Statement
	Location
}

func (f *Function) SExpr() string {
	panic("implement me")
}

func (f *Function) command() {}
func (f *Function) String() string {
	buf := bytes.NewBuffer([]byte{})
	buf.WriteString(fmt.Sprintf("fn %s (", f.Var))
	for i, b := range f.Bindings {
		buf.WriteString(b.String())
		if i != len(f.Bindings)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString(fmt.Sprintf(") : %s {\n", f.ReturnType))
	for _, s := range f.Statements {
		buf.WriteString(fmt.Sprintf("\t%s\n", s))
	}

	buf.WriteString("}\n")
	return buf.String()
}
... // with the normal expressions
```

## expression.go
```go
type Expression interface {
	Node
	Typ() types.Type
	expression()
}

// IntExpression
type IntExpression struct {
	Val  int64
	Type types.Type
	Location
}

func (i *IntExpression) SExpr() string {
	if i.Type != nil {
		return fmt.Sprintf("(IntExpr %s %d)", i.Type, i.Val)
	}
	return fmt.Sprintf("(IntExpr %d)", i.Val)
}

func (i *IntExpression) String() string {
	return fmt.Sprintf("%d", i.Val)
}

func (i *IntExpression) Typ() types.Type {
	return i.Type
}

func (i *IntExpression) expression() {}

// IdentifierExpression
type IdentifierExpression struct {
	Identifier string
	Type       types.Type
	Location
}

func (i *IdentifierExpression) SExpr() string {
	if i.Type != nil {
		return fmt.Sprintf("(VarExpr %s %s)", i.Type, i.Identifier)
	}
	return fmt.Sprintf("(VarExpr %s)", i.Identifier)
}

func (i *IdentifierExpression) String() string {
	return i.Identifier
}

func (i *IdentifierExpression) Typ() types.Type {
	return i.Type
}

func (i *IdentifierExpression) expression() {}

type CallExpression struct {
	Identifier string
	Arguments  []Expression
	Type       types.Type
	Location
}

func (c *CallExpression) SExpr() string {
	strs := make([]string, len(c.Arguments))
	for i, expr := range c.Arguments {
		strs[i] = expr.SExpr()
	}

	if c.Type != nil {
		return fmt.Sprintf("(CallExpr %s %s %s)", c.Type, c.Identifier, strings.Join(strs, " "))
	}

	return fmt.Sprintf("(CallExpr %s %s)", c.Identifier, strings.Join(strs, " "))
}

func (c *CallExpression) String() string {
	strs := make([]string, len(c.Arguments))
	for i, expr := range c.Arguments {
		strs[i] = expr.String()
	}

	return fmt.Sprintf("%s(%s)", c.Identifier, strings.Join(strs, ", "))
}

func (c *CallExpression) Typ() types.Type {
	return c.Type
}

func (c *CallExpression) expression() {}

// FloatExpression
type FloatExpression struct {
	Val  float64
	Type types.Type
	Location
}

func (f *FloatExpression) SExpr() string {
	if f.Type != nil {
		return fmt.Sprintf("(FloatExpr %s %d)", f.Type, int64(f.Val))
	}
	return fmt.Sprintf("(FloatExpr %d)", int64(f.Val))
}

func (f *FloatExpression) String() string {
	return fmt.Sprintf("%f", f.Val)
}

func (f *FloatExpression) Typ() types.Type {
	return f.Type
}

func (f *FloatExpression) expression() {}

type BooleanExpression struct {
	Val  bool
	Type types.Type
	Location
}

func (b *BooleanExpression) SExpr() string {
	if b.Type != nil {
		return fmt.Sprintf("(VarExpr %s %t)", b.Type, b.Val)
	}
	return fmt.Sprintf("(VarExpr %t)", b.Val)
}

func (b *BooleanExpression) String() string {
	if b.Val {
		return "true"
	}
	return "false"
}

func (b *BooleanExpression) Typ() types.Type {
	return b.Type
}

func (b *BooleanExpression) expression() {}

type TupleExpression struct {
	Expressions []Expression
	Type        types.Type
	Location
}

func (t *TupleExpression) SExpr() string {
	panic("implement me")
}

func (t *TupleExpression) String() string {
	strs := make([]string, len(t.Expressions))
	for i, expr := range t.Expressions {
		strs[i] = expr.String()
	}
	return fmt.Sprintf("{%s}", strings.Join(strs, ", "))
}

func (t *TupleExpression) Typ() types.Type {
	return t.Type
}

func (t *TupleExpression) expression() {}

type ArrayExpression struct {
	Expressions []Expression
	Type        types.Type
	Location
}

func (a *ArrayExpression) SExpr() string {
	return ""
}

func (a *ArrayExpression) String() string {
	strs := make([]string, len(a.Expressions))
	for i, e := range a.Expressions {
		strs[i] = e.String()
	}
	return fmt.Sprintf("[%s]", strings.Join(strs, ", "))
}

func (a *ArrayExpression) Typ() types.Type {
	return a.Type
}

func (a *ArrayExpression) expression() {}

type ArrayRefExpression struct {
	Array   Expression
	Indexes []Expression
	Type    types.Type
	Location
}

func (a *ArrayRefExpression) SExpr() string {
	return ""
}
func (a *ArrayRefExpression) String() string {
	strs := make([]string, len(a.Indexes))
	for i, idx := range a.Indexes {
		strs[i] = fmt.Sprintf("%s", idx)
	}
	return fmt.Sprintf("%s[%s]", a.Array, strings.Join(strs, ", "))
}

func (a *ArrayRefExpression) Typ() types.Type {
	return a.Type
}

func (a *ArrayRefExpression) expression() {}

type TupleRefExpression struct {
	Tuple Expression
	Index Expression
	Type  types.Type
	Location
}

func (t *TupleRefExpression) SExpr() string {
	panic("implement me")
}

func (t *TupleRefExpression) String() string {
	return fmt.Sprintf("%s{%s}", t.Tuple, t.Index)
}

func (t *TupleRefExpression) Typ() types.Type {
	return t.Type
}

func (t *TupleRefExpression) expression() {}

type IfExpression struct {
	Condition   Expression
	Consequence Expression
	Otherwise   Expression
	Type        types.Type
	Location
}

func (i *IfExpression) SExpr() string {
	panic("implement me")
}

func (i *IfExpression) String() string {
	return fmt.Sprintf("if %s then %s else %s",
		i.Condition.String(), i.Consequence.String(), i.Otherwise.String())
}

func (i *IfExpression) Typ() types.Type {
	return i.Type
}

func (i *IfExpression) expression() {}

type OpBinding struct {
	Variable string
	Expr     Expression
	Type     types.Type
	Location
}

func (o *OpBinding) String() string {
	return fmt.Sprintf("%s : %s", o.Variable, o.Expr)
}

type ArrayTransform struct {
	OpBindings []OpBinding
	Expr       Expression
	Type       types.Type
	Location
}

func (a *ArrayTransform) SExpr() string {
	panic("implement me")
}

func (a *ArrayTransform) Typ() types.Type {
	return a.Type
}

func (a *ArrayTransform) String() string {
	bindings := make([]string, len(a.OpBindings))
	for i, b := range a.OpBindings {
		bindings[i] = b.String()
	}

	return fmt.Sprintf("array[%s] %s", strings.Join(bindings, ", "), a.Expr)
}

func (a *ArrayTransform) expression() {}

type SumTransform struct {
	OpBindings []OpBinding
	Expr       Expression
	Type       types.Type
	Location
}

func (s *SumTransform) SExpr() string {
	panic("implement me")
}

func (s *SumTransform) String() string {
	bindings := make([]string, len(s.OpBindings))
	for i, b := range s.OpBindings {
		bindings[i] = b.String()
	}
	return fmt.Sprintf("sum[%s] %s", strings.Join(bindings, ", "), s.Expr)
}

func (s *SumTransform) Typ() types.Type {
	return s.Type
}
func (s *SumTransform) expression() {}

type InfixExpression struct {
	Left  Expression
	Right Expression
	Op    string
	Type  types.Type
	Location
}

func (i *InfixExpression) SExpr() string {
	panic("shouldn't run")
}

func (i *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", i.Left, i.Op, i.Right)
}

func (i *InfixExpression) Typ() types.Type {
	return i.Type
}

func (i *InfixExpression) expression() {}

type PrefixExpression struct {
	Op   string
	Expr Expression
	Type types.Type
	Location
}

func (p *PrefixExpression) Typ() types.Type {
	return p.Type
}

func (p *PrefixExpression) SExpr() string {
	return p.Expr.SExpr()
}

func (p *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", p.Op, p.Expr)
}

func (p *PrefixExpression) expression() {}
```

# 2. Expanded grammar
```
expr : <integer>
     | <float>
     | true
     | false

postfix : <expr> { <integer> }
        | <expr> [expr, ...]
        | <expr>
         

prefix : - <postfix>
       | ! <postfix>
       | <postfix>

multiplicative : <prefix> * <prefix>
               | <prefx> / <prefix>
               | <prefix> % <prefix>
               | <prefix>

additive : <multiplicative> + <multiplicative>
         | <multiplicative> - <multiplicative>
         | <multiplicative>

ordered-comp : <additve> < <additive>
             | <additive> > <additive>
             | <additive> >= <additive>
             | <additive> <= <additive>
             | <additive>

unordered-comp : <ordered-comp> == <ordered-comp> 
               | <ordered-comp> != <unordered-comp>
               | <ordered-comp>

boolean-ops : <unordered-comp> && <unordered-comp>
            | <unordered-comp> || <unordered-comp>
            | <unordered-comp>

rest : array[<variable> : <boolean-ops>, ...] <boolean-ops>
     | sum [<variable> : <boolean-ops>, ...] <boolean-ops>
     | if <boolean-ops> then <boolean-ops> else < booleanops
 
```


# 3. Typechecking signature
```go
// Check takes in an `ast.Program` and returns a type-annotated ast.Program, along with a symbol table.
// Both the program and table will return nil if there is an error (so check for an error first!)
func Check(program ast.Program) (ast.Program, *symbol.Table, error)
```

# 4. Array expression typecheck
```go
func checkArrayTransform(expr *ast.ArrayTransform, table *symbol.Table) (types.Type, error) {
	// copy symbol table to use as a local copy
	cpy := table.Copy()

	// loop over bindings and
	for _, binding := range expr.OpBindings {
		// make sure no variable shadowing is going on
		if _, ok := cpy.Get(binding.Variable); ok {
			return nil, NewError(binding,
				"illegal shadowing in sum expr, var: %s", binding.Variable)
		}

		// typecheck the binding type <var> : <expression> <-- this here
		bindType, err := expressionType(binding.Expr, cpy)
		if err != nil {
			return nil, err
		}

		// the binding type must be an integer
		if !bindType.Equal(types.Integer) {
			return nil, NewError(binding, "bindArg expr initializer for %s returns non-integer",
				binding.Variable)
		}

		// set the variable up in the local symbol table as an integer
		cpy.Set(binding.Variable, &symbol.Identifier{Type: types.Integer})
	}

	// get the expression type of the right hand side of the expression
	exprType, err := expressionType(expr.Expr, cpy)
	if err != nil {
		return nil, err
	}

	// make sure that it is some sort of array expression. If it isn't, then maybe we should be
	// using a sum[]?
	if arr, ok := exprType.(*types.Array); ok {
		if arr.Rank != len(expr.OpBindings) {
			return nil, NewError(expr.Expr, "return type of array expression must be of equal rank of number of bindings")
		}
		return exprType, nil
	}
	
	expr.Type = exprType
	return nil, NewError(expr.Expr, "return type of array expression must be array")
}
```
