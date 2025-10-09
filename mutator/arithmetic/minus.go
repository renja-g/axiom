package arithmetic

import (
	"go/ast"
	"go/token"
)

// Minus mutates - to +
// SUB -> ADD
type Minus struct{}

func (m Minus) Name() string {
	return "Arithmetic_SUB"
}

func (m Minus) CanMutate(node ast.Node) bool {
	bin, ok := node.(*ast.BinaryExpr)
	return ok && bin.Op == token.SUB
}

func (m Minus) Mutate(node ast.Node) ast.Node {
	bin := node.(*ast.BinaryExpr)
	return &ast.BinaryExpr{
		Op: token.ADD,
		X:  bin.X,
		Y:  bin.Y,
	}
}
