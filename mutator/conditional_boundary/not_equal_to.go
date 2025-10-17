package conditional_boundary

import (
	"go/ast"
	"go/token"
)

// NotEqualTo mutates != to ==
// NEQ -> EQL
type NotEqualTo struct{}

func (m NotEqualTo) Name() string {
	return "ConditionalBoundary_NEQ_EQL"
}

func (m NotEqualTo) CanMutate(node ast.Node) bool {
	bin, ok := node.(*ast.BinaryExpr)
	return ok && bin.Op == token.NEQ
}

func (m NotEqualTo) Mutate(node ast.Node) ast.Node {
	bin := node.(*ast.BinaryExpr)
	return &ast.BinaryExpr{
		Op: token.EQL,
		X:  bin.X,
		Y:  bin.Y,
	}
}
