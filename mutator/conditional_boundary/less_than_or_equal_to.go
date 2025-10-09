package conditional_boundary

import (
	"go/ast"
	"go/token"
)

// LessThanOrEqualTo mutates <= to <
// LEQ -> LSS
type LessThanOrEqualTo struct{}

func (m LessThanOrEqualTo) Name() string {
	return "ConditionalBoundary_LEQ_LSS"
}

func (m LessThanOrEqualTo) CanMutate(node ast.Node) bool {
	bin, ok := node.(*ast.BinaryExpr)
	return ok && bin.Op == token.LEQ
}

func (m LessThanOrEqualTo) Mutate(node ast.Node) ast.Node {
	bin := node.(*ast.BinaryExpr)
	return &ast.BinaryExpr{
		Op: token.LSS,
		X:  bin.X,
		Y:  bin.Y,
	}
}
