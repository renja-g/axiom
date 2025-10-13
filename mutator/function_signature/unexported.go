package functionsignature

import (
	"go/ast"
	"unicode"
)

// Unexported mutates unexported functions to exported
// Changes first letter from lowercase to uppercase
// add() -> Add()
type Unexported struct{}

func (m Unexported) Name() string {
	return "FunctionSignature_Unexported"
}

func (m Unexported) CanMutate(node ast.Node) bool {
	fn, ok := node.(*ast.FuncDecl)
	if !ok || fn.Name == nil {
		return false
	}

	name := fn.Name.Name
	if len(name) == 0 {
		return false
	}

	firstChar := rune(name[0])
	return unicode.IsLower(firstChar)
}

func (m Unexported) Mutate(node ast.Node) ast.Node {
	fn := node.(*ast.FuncDecl)
	name := fn.Name.Name

	runes := []rune(name)
	runes[0] = unicode.ToUpper(runes[0])
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
