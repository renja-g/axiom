package arithmetic

import (
	"go/ast"
	"go/token"
)

// ShiftRight mutates >> to <<
// SHR -> SHL
type ShiftRight struct{}

func (m ShiftRight) Name() string {
	return "Arithmetic_SHR"
}

func (m ShiftRight) CanMutate(node ast.Node) bool {
	bin, ok := node.(*ast.BinaryExpr)
	return ok && bin.Op == token.SHR
}

func (m ShiftRight) Mutate(node ast.Node) ast.Node {
	bin := node.(*ast.BinaryExpr)
	return &ast.BinaryExpr{
		Op: token.SHL,
		X:  bin.X,
		Y:  bin.Y,
	}
}
