package arithmetic

import (
	"go/ast"
	"go/token"
)

// DivEqual mutates /= to *=
// QUO_ASSIGN -> MUL_ASSIGN
type DivEqual struct{}

func (m DivEqual) Name() string {
	return "Arithmetic_QUO_ASSIGN"
}

func (m DivEqual) CanMutate(node ast.Node) bool {
	assign, ok := node.(*ast.AssignStmt)
	return ok && assign.Tok == token.QUO_ASSIGN
}

func (m DivEqual) Mutate(node ast.Node) ast.Node {
	assign := node.(*ast.AssignStmt)
	return &ast.AssignStmt{
		Tok: token.MUL_ASSIGN,
		Lhs: assign.Lhs,
		Rhs: assign.Rhs,
	}
}
