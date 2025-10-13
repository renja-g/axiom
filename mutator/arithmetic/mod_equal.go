package arithmetic

import (
	"go/ast"
	"go/token"
)

// ModEqual mutates %= to *=
// REM_ASSIGN -> MUL_ASSIGN
type ModEqual struct{}

func (m ModEqual) Name() string {
	return "Arithmetic_REM_ASSIGN"
}

func (m ModEqual) CanMutate(node ast.Node) bool {
	assign, ok := node.(*ast.AssignStmt)
	return ok && assign.Tok == token.REM_ASSIGN
}

func (m ModEqual) Mutate(node ast.Node) ast.Node {
	assign := node.(*ast.AssignStmt)
	return &ast.AssignStmt{
		Tok: token.MUL_ASSIGN,
		Lhs: assign.Lhs,
		Rhs: assign.Rhs,
	}
}
