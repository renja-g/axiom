package arithmetic

import (
	"go/ast"
	"go/token"
)

// MulEqual mutates *= to /=
// MUL_ASSIGN -> QUO_ASSIGN
type MulEqual struct{}

func (m MulEqual) Name() string {
	return "Arithmetic_MUL_ASSIGN"
}

func (m MulEqual) CanMutate(node ast.Node) bool {
	assign, ok := node.(*ast.AssignStmt)
	return ok && assign.Tok == token.MUL_ASSIGN
}

func (m MulEqual) Mutate(node ast.Node) ast.Node {
	assign := node.(*ast.AssignStmt)
	return &ast.AssignStmt{
		Tok: token.QUO_ASSIGN,
		Lhs: assign.Lhs,
		Rhs: assign.Rhs,
	}
}
