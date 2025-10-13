package arithmetic

import (
	"go/ast"
	"go/token"
)

// MinusEqual mutates -= to +=
// SUB_ASSIGN -> ADD_ASSIGN
type MinusEqual struct{}

func (m MinusEqual) Name() string {
	return "Arithmetic_SUB_ASSIGN"
}

func (m MinusEqual) CanMutate(node ast.Node) bool {
	assign, ok := node.(*ast.AssignStmt)
	return ok && assign.Tok == token.SUB_ASSIGN
}

func (m MinusEqual) Mutate(node ast.Node) ast.Node {
	assign := node.(*ast.AssignStmt)
	return &ast.AssignStmt{
		Tok: token.ADD_ASSIGN,
		Lhs: assign.Lhs,
		Rhs: assign.Rhs,
	}
}
