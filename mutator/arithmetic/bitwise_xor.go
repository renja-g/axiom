package arithmetic

import (
	"go/ast"
	"go/token"
)

// BitwiseXor mutates ^ to &
// XOR -> AND
type BitwiseXor struct{}

func (m BitwiseXor) Name() string {
	return "Arithmetic_XOR"
}

func (m BitwiseXor) CanMutate(node ast.Node) bool {
	bin, ok := node.(*ast.BinaryExpr)
	return ok && bin.Op == token.XOR
}

func (m BitwiseXor) Mutate(node ast.Node) ast.Node {
	bin := node.(*ast.BinaryExpr)
	return &ast.BinaryExpr{
		Op: token.AND,
		X:  bin.X,
		Y:  bin.Y,
	}
}
