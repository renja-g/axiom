package conditional_boundary

import (
	"go/ast"
	"go/token"
)

// GreaterThan mutates > to >=
// GTR -> GEQ
type GreaterThan struct{}

func (m GreaterThan) Name() string {
	return "ConditionalBoundary_GTR_GEQ"
}

func (m GreaterThan) CanMutate(node ast.Node) bool {
	bin, ok := node.(*ast.BinaryExpr)
	return ok && bin.Op == token.GTR
}

func (m GreaterThan) Mutate(node ast.Node) ast.Node {
	bin := node.(*ast.BinaryExpr)
	return &ast.BinaryExpr{
		Op: token.GEQ,
		X:  bin.X,
		Y:  bin.Y,
	}
}
