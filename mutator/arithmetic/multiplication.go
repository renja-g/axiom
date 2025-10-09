package arithmetic

import (
	"go/ast"
	"go/token"
)

// Multiplication mutates * to /
// MUL -> QUO
type Multiplication struct{}

func (m Multiplication) Name() string {
	return "Arithmetic_QUO"
}

func (m Multiplication) CanMutate(node ast.Node) bool {
	bin, ok := node.(*ast.BinaryExpr)
	return ok && bin.Op == token.MUL
}

func (m Multiplication) Mutate(node ast.Node) ast.Node {
	bin := node.(*ast.BinaryExpr)
	return &ast.BinaryExpr{
		Op: token.QUO,
		X:  bin.X,
		Y:  bin.Y,
	}
}
