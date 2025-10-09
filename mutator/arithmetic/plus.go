package arithmetic

import (
	"go/ast"
	"go/token"
)

// Plus mutates + to -
// ADD -> SUB
type Plus struct{}

func (m Plus) Name() string {
	return "Arithmetic_ADD"
}

func (m Plus) CanMutate(node ast.Node) bool {
	bin, ok := node.(*ast.BinaryExpr)
	return ok && bin.Op == token.ADD
}

func (m Plus) Mutate(node ast.Node) ast.Node {
	bin := node.(*ast.BinaryExpr)
	return &ast.BinaryExpr{
		Op: token.SUB,
		X:  bin.X,
		Y:  bin.Y,
	}
}
