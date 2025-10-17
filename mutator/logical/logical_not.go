package logical

import (
	"go/ast"
	"go/token"
)

// LogicalNot mutates logical NOT operations by removing the '!'.
// !x -> x
type LogicalNot struct{}

func (m LogicalNot) Name() string {
	return "Logical_NOT"
}

func (m LogicalNot) CanMutate(node ast.Node) bool {
	unary, ok := node.(*ast.UnaryExpr)
	return ok && unary.Op == token.NOT
}

func (m LogicalNot) Mutate(node ast.Node) ast.Node {
	unary := node.(*ast.UnaryExpr)
	return unary.X
}
