package arithmetic

import (
	"go/ast"
	"go/token"
)

// ShiftLeft mutates << to >>
// SHL -> SHR
type ShiftLeft struct{}

func (m ShiftLeft) Name() string {
	return "Arithmetic_SHL"
}

func (m ShiftLeft) CanMutate(node ast.Node) bool {
	bin, ok := node.(*ast.BinaryExpr)
	return ok && bin.Op == token.SHL
}

func (m ShiftLeft) Mutate(node ast.Node) ast.Node {
	bin := node.(*ast.BinaryExpr)
	return &ast.BinaryExpr{
		Op: token.SHR,
		X:  bin.X,
		Y:  bin.Y,
	}
}
