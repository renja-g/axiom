package boolean

import "go/ast"

// TrueValue mutates boolean true literals to false.
// true -> false
type TrueValue struct{}

func (m TrueValue) Name() string {
	return "Boolean_TRUE"
}

func (m TrueValue) CanMutate(node ast.Node) bool {
	ident, ok := node.(*ast.Ident)
	return ok && ident.Name == "true"
}

func (m TrueValue) Mutate(node ast.Node) ast.Node {
	ident := node.(*ast.Ident)
	return &ast.Ident{
		NamePos: ident.NamePos,
		Name:    "false",
	}
}
