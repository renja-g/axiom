package arithmetic

import (
	"go/ast"
	"go/token"
)

// PlusEqual mutates += to -=
// ADD_ASSIGN -> SUB_ASSIGN
type PlusEqual struct{}

func (m PlusEqual) Name() string {
	return "Arithmetic_ADD_ASSIGN"
}

func (m PlusEqual) CanMutate(node ast.Node) bool {
	assign, ok := node.(*ast.AssignStmt)
	return ok && assign.Tok == token.ADD_ASSIGN
}

func (m PlusEqual) Mutate(node ast.Node) ast.Node {
	assign := node.(*ast.AssignStmt)
	return &ast.AssignStmt{
		Tok: token.SUB_ASSIGN,
		Lhs: assign.Lhs,
		Rhs: assign.Rhs,
	}
}
