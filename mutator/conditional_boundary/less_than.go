package conditional_boundary

import (
	"go/ast"
	"go/token"
)

// LessThan mutates < to <=
// LSS -> LEQ
type LessThan struct{}

func (m LessThan) Name() string {
	return "ConditionalBoundary_LSS_LEQ"
}

func (m LessThan) CanMutate(node ast.Node) bool {
	bin, ok := node.(*ast.BinaryExpr)
	return ok && bin.Op == token.LSS
}

func (m LessThan) Mutate(node ast.Node) ast.Node {
	bin := node.(*ast.BinaryExpr)
	return &ast.BinaryExpr{
		Op: token.LEQ,
		X:  bin.X,
		Y:  bin.Y,
	}
}
