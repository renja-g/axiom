package function_signature

import (
	"go/ast"
	"unicode"
)

// Exported mutates exported function names to unexported ones
// Public -> private
type Exported struct{}

func (m Exported) Name() string {
	return "FunctionSignature_Exported"
}

func (m Exported) CanMutate(node ast.Node) bool {
	fn, ok := node.(*ast.FuncDecl)
	if !ok || fn.Name == nil {
		return false
	}
	return fn.Name.IsExported() && startsWithLetter(fn.Name.Name)
}

func (m Exported) Mutate(node ast.Node) ast.Node {
	fn, ok := node.(*ast.FuncDecl)
	if !ok {
		return node
	}
	name := fn.Name.Name
	if name == "" {
		return node
	}

	newFn := *fn
	newName := changeFirstRune(name, unicode.ToLower)
	newFn.Name = cloneIdentWithName(fn.Name, newName)
	return &newFn
}
