package functionsignature

import (
	"go/ast"
	"unicode"
)

// Exported mutates exported functions to unexported
// Changes first letter from uppercase to lowercase
// Add() -> add()
type Exported struct{}

func (m Exported) Name() string {
	return "FunctionSignature_Exported"
}

func (m Exported) CanMutate(node ast.Node) bool {
	fn, ok := node.(*ast.FuncDecl)
	if !ok || fn.Name == nil {
		return false
	}

	name := fn.Name.Name
	if len(name) == 0 {
		return false
	}

	firstChar := rune(name[0])
	return unicode.IsUpper(firstChar)
}

func (m Exported) Mutate(node ast.Node) ast.Node {
	fn := node.(*ast.FuncDecl)
	name := fn.Name.Name

	runes := []rune(name)
	runes[0] = unicode.ToLower(runes[0])
	newName := string(runes)

	mutatedFn := &ast.FuncDecl{
		Doc:  fn.Doc,
		Recv: fn.Recv,
		Name: &ast.Ident{
			Name: newName,
		},
		Type: fn.Type,
		Body: fn.Body,
	}

	return mutatedFn
}
