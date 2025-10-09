package arithmetic

import (
	"go/ast"
	"go/token"
)

// Modulus mutates % to *
type Modulus struct{}

func (m Modulus) Name() string {
	return "Arithmetic_REM"
}

func (m Modulus) CanMutate(node ast.Node) bool {
	bin, ok := node.(*ast.BinaryExpr)
	return ok && bin.Op == token.REM
}

func (m Modulus) Mutate(node ast.Node) ast.Node {
	bin := node.(*ast.BinaryExpr)
	return &ast.BinaryExpr{
		Op: token.MUL,
		X:  bin.X,
		Y:  bin.Y,
	}
}
