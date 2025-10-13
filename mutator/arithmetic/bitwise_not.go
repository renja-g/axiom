package arithmetic

import (
	"go/ast"
	"go/token"
)

// BitwiseNot mutates ^x to x (removes bitwise NOT)
// NOT -> (remove)
type BitwiseNot struct{}

func (m BitwiseNot) Name() string {
	return "Arithmetic_NOT"
}

func (m BitwiseNot) CanMutate(node ast.Node) bool {
	unary, ok := node.(*ast.UnaryExpr)
	return ok && unary.Op == token.XOR
}

func (m BitwiseNot) Mutate(node ast.Node) ast.Node {
	unary := node.(*ast.UnaryExpr)
	return unary.X
}
