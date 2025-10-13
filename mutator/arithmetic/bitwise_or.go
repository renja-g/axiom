package arithmetic

import (
	"go/ast"
	"go/token"
)

// BitwiseOr mutates | to &
// OR -> AND
type BitwiseOr struct{}

func (m BitwiseOr) Name() string {
	return "Arithmetic_OR"
}

func (m BitwiseOr) CanMutate(node ast.Node) bool {
	bin, ok := node.(*ast.BinaryExpr)
	return ok && bin.Op == token.OR
}

func (m BitwiseOr) Mutate(node ast.Node) ast.Node {
	bin := node.(*ast.BinaryExpr)
	return &ast.BinaryExpr{
		Op: token.AND,
		X:  bin.X,
		Y:  bin.Y,
	}
}
