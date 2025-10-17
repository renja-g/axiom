package arithmetic

import (
	"go/ast"
	"go/token"
	"strconv"
)

// IntegerLiteralBoundary mutates integer literals to nearby boundary values.
// Rules:
// 0 <-> 1; 1 <-> 0; n>1 => n-1; n<-1 => n+1; -1 => 0 (by rule n<-1 => n+1 gives 0?)
// We special-case -1 to 0 to keep single-step mutation.
type IntegerLiteralBoundary struct{}

func (m IntegerLiteralBoundary) Name() string { return "Arithmetic_INT_LITERAL_BOUNDARY" }

func (m IntegerLiteralBoundary) CanMutate(node ast.Node) bool {
	lit, ok := node.(*ast.BasicLit)
	if !ok || lit.Kind != token.INT {
		return false
	}
	// Ensure it parses as a base-10 integer; skip hex/oct/bin formats for simplicity
	if _, err := strconv.Atoi(lit.Value); err != nil {
		return false
	}
	return true
}

func (m IntegerLiteralBoundary) Mutate(node ast.Node) ast.Node {
	lit, ok := node.(*ast.BasicLit)
	if !ok || lit.Kind != token.INT {
		return node
	}
	v, err := strconv.Atoi(lit.Value)
	if err != nil {
		return node
	}

	var mutated string
	switch v {
	case 0:
		mutated = "1"
	case 1:
		mutated = "0"
	case -1:
		mutated = "0"
	default:
		if v > 1 {
			mutated = strconv.Itoa(v - 1)
		} else { // v < -1
			mutated = strconv.Itoa(v + 1)
		}
	}

	return &ast.BasicLit{Kind: token.INT, Value: mutated}
}
