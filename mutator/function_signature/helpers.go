package function_signature

import (
	"go/ast"
	"unicode"
)

func startsWithLetter(name string) bool {
	for _, r := range name {
		return unicode.IsLetter(r)
	}
	return false
}

func changeFirstRune(name string, transform func(rune) rune) string {
	runes := []rune(name)
	if len(runes) == 0 {
		return name
	}
	runes[0] = transform(runes[0])
	return string(runes)
}

func cloneIdentWithName(ident *ast.Ident, newName string) *ast.Ident {
	if ident == nil {
		return ast.NewIdent(newName)
	}
	newIdent := ast.NewIdent(newName)
	newIdent.NamePos = ident.NamePos
	newIdent.Obj = ident.Obj
	return newIdent
}
