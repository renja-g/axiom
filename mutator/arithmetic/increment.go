package arithmetic

import (
	"go/ast"
	"go/token"
)

// Increment mutates i++ to i--
// INC -> DEC
type Increment struct{}

func (m Increment) Name() string {
	return "Arithmetic_INC"
}

func (m Increment) CanMutate(node ast.Node) bool {
	stmt, ok := node.(*ast.IncDecStmt)
	return ok && stmt.Tok == token.INC
}

func (m Increment) Mutate(node ast.Node) ast.Node {
	stmt := node.(*ast.IncDecStmt)
	return &ast.IncDecStmt{
		X:      stmt.X,
		TokPos: stmt.TokPos,
		Tok:    token.DEC,
	}
}
