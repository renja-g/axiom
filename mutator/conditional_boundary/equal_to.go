package conditional_boundary

import (
	"go/ast"
	"go/token"
)

// EqualTo mutates == to !=
// EQL -> NEQ
type EqualTo struct{}

func (m EqualTo) Name() string {
	return "ConditionalBoundary_EQL_NEQ"
}

func (m EqualTo) CanMutate(node ast.Node) bool {
	bin, ok := node.(*ast.BinaryExpr)
	return ok && bin.Op == token.EQL
}

func (m EqualTo) Mutate(node ast.Node) ast.Node {
	bin := node.(*ast.BinaryExpr)
	return &ast.BinaryExpr{
		Op: token.NEQ,
		X:  bin.X,
		Y:  bin.Y,
	}
}
