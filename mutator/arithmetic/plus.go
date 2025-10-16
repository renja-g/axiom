package arithmetic

import (
	"go/ast"
	"go/token"
	"go/types"
)

// Plus mutates + to -
// ADD -> SUB
type Plus struct{}

func (m Plus) Name() string {
	return "Arithmetic_ADD"
}

func (m Plus) CanMutate(node ast.Node) bool {
	bin, ok := node.(*ast.BinaryExpr)
	return ok && bin.Op == token.ADD
}

// CanMutateWithType checks if the node can be mutated, excluding string concatenation.
// This method is called by the generator when type information is available.
func (m Plus) CanMutateWithType(node ast.Node, typeInfo *types.Info) bool {
	bin, ok := node.(*ast.BinaryExpr)
	if !ok || bin.Op != token.ADD {
		return false
	}

	// Check if either operand is a string type
	xType := typeInfo.TypeOf(bin.X)
	yType := typeInfo.TypeOf(bin.Y)

	// If we can't determine the type, fall back to allowing the mutation
	if xType == nil || yType == nil {
		return true
	}

	// Check if either operand is a string
	xBasic, xOk := xType.Underlying().(*types.Basic)
	yBasic, yOk := yType.Underlying().(*types.Basic)

	if (xOk && xBasic.Kind() == types.String) || (yOk && yBasic.Kind() == types.String) {
		// This is string concatenation, don't mutate
		return false
	}

	return true
}

func (m Plus) Mutate(node ast.Node) ast.Node {
	bin := node.(*ast.BinaryExpr)
	return &ast.BinaryExpr{
		Op: token.SUB,
		X:  bin.X,
		Y:  bin.Y,
	}
}
