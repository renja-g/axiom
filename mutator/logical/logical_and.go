package logical

import (
	"go/ast"
	"go/token"
)

// LogicalAnd mutates logical AND operations to OR.
// && -> ||
type LogicalAnd struct{}

func (m LogicalAnd) Name() string {
	return "Logical_AND"
}

func (m LogicalAnd) CanMutate(node ast.Node) bool {
	bin, ok := node.(*ast.BinaryExpr)
	return ok && bin.Op == token.LAND
}

func (m LogicalAnd) Mutate(node ast.Node) ast.Node {
	bin := node.(*ast.BinaryExpr)
	return &ast.BinaryExpr{
		X:     bin.X,
		Op:    token.LOR,
		Y:     bin.Y,
		OpPos: bin.OpPos,
	}
}
