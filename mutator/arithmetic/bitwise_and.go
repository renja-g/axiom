package arithmetic

import (
	"go/ast"
	"go/token"
)

// BitwiseAnd mutates & to |
// AND -> OR
type BitwiseAnd struct{}

func (m BitwiseAnd) Name() string {
	return "Arithmetic_AND"
}

func (m BitwiseAnd) CanMutate(node ast.Node) bool {
	bin, ok := node.(*ast.BinaryExpr)
	return ok && bin.Op == token.AND
}

func (m BitwiseAnd) Mutate(node ast.Node) ast.Node {
	bin := node.(*ast.BinaryExpr)
	return &ast.BinaryExpr{
		Op: token.OR,
		X:  bin.X,
		Y:  bin.Y,
	}
}
