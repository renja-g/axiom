package conditional_boundary

import (
	"go/ast"
	"go/token"
)

// GreaterThanOrEqualTo mutates >= to >
// GEQ -> GTR
type GreaterThanOrEqualTo struct{}

func (m GreaterThanOrEqualTo) Name() string {
	return "ConditionalBoundary_GEQ_GTR"
}

func (m GreaterThanOrEqualTo) CanMutate(node ast.Node) bool {
	bin, ok := node.(*ast.BinaryExpr)
	return ok && bin.Op == token.GEQ
}

func (m GreaterThanOrEqualTo) Mutate(node ast.Node) ast.Node {
	bin := node.(*ast.BinaryExpr)
	return &ast.BinaryExpr{
		Op: token.GTR,
		X:  bin.X,
		Y:  bin.Y,
	}
}
