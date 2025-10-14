package boolean

import "go/ast"

// FalseValue mutates boolean false literals to true.
// false -> true
type FalseValue struct{}

func (m FalseValue) Name() string {
	return "Boolean_FALSE"
}

func (m FalseValue) CanMutate(node ast.Node) bool {
	ident, ok := node.(*ast.Ident)
	return ok && ident.Name == "false"
}

func (m FalseValue) Mutate(node ast.Node) ast.Node {
	ident := node.(*ast.Ident)
	return &ast.Ident{
		NamePos: ident.NamePos,
		Name:    "true",
	}
}
