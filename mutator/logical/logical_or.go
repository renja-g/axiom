package logical

import (
	"go/ast"
	"go/token"
)

// LogicalOr mutates logical OR operations to AND.
// || -> &&
type LogicalOr struct{}

func (m LogicalOr) Name() string {
	return "Logical_OR"
}

func (m LogicalOr) CanMutate(node ast.Node) bool {
	bin, ok := node.(*ast.BinaryExpr)
	return ok && bin.Op == token.LOR
}

func (m LogicalOr) Mutate(node ast.Node) ast.Node {
	bin := node.(*ast.BinaryExpr)
	return &ast.BinaryExpr{
		X:     bin.X,
		Op:    token.LAND,
		Y:     bin.Y,
		OpPos: bin.OpPos,
	}
}
