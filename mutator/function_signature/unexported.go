package function_signature

import (
	"go/ast"
	"unicode"
)

// Unexported mutates unexported function names to exported ones
// private -> Public
type Unexported struct{}

func (m Unexported) Name() string {
	return "FunctionSignature_Unexported"
}

func (m Unexported) CanMutate(node ast.Node) bool {
	fn, ok := node.(*ast.FuncDecl)
	if !ok || fn.Name == nil {
		return false
	}
	return !fn.Name.IsExported() && startsWithLetter(fn.Name.Name)
}

func (m Unexported) Mutate(node ast.Node) ast.Node {
	fn, ok := node.(*ast.FuncDecl)
	if !ok {
		return node
	}
	name := fn.Name.Name
	if name == "" {
		return node
	}

	newFn := *fn
	newName := changeFirstRune(name, unicode.ToUpper)
	newFn.Name = cloneIdentWithName(fn.Name, newName)
	return &newFn
}
