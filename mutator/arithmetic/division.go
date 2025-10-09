package arithmetic

import (
	"go/ast"
	"go/token"
)

// Division mutates / to *
// QUO -> MUL
type Division struct{}

func (m Division) Name() string {
	return "Arithmetic_QUO"
}

func (m Division) CanMutate(node ast.Node) bool {
	bin, ok := node.(*ast.BinaryExpr)
	return ok && bin.Op == token.QUO
}

func (m Division) Mutate(node ast.Node) ast.Node {
	bin := node.(*ast.BinaryExpr)
	return &ast.BinaryExpr{
		Op: token.MUL,
		X:  bin.X,
		Y:  bin.Y,
	}
}
