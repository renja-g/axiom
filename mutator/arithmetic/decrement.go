package arithmetic

import (
	"go/ast"
	"go/token"
)

// Decrement mutates i-- to i++
// DEC -> INC
type Decrement struct{}

func (m Decrement) Name() string {
	return "Arithmetic_DEC"
}

func (m Decrement) CanMutate(node ast.Node) bool {
	stmt, ok := node.(*ast.IncDecStmt)
	return ok && stmt.Tok == token.DEC
}

func (m Decrement) Mutate(node ast.Node) ast.Node {
	stmt := node.(*ast.IncDecStmt)
	return &ast.IncDecStmt{
		X:      stmt.X,
		TokPos: stmt.TokPos,
		Tok:    token.INC,
	}
}
