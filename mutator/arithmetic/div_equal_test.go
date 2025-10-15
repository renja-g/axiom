package arithmetic

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestDivEqualName(t *testing.T) {
	mut := DivEqual{}

	if got, want := mut.Name(), "Arithmetic_QUO_ASSIGN"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestDivEqualCanMutate(t *testing.T) {
	mut := DivEqual{}

	tests := []struct {
		name string
		node ast.Node
		want bool
	}{
		{
			name: "division assignment is mutable",
			node: &ast.AssignStmt{Tok: token.QUO_ASSIGN},
			want: true,
		},
		{
			name: "multiplication assignment is not mutable",
			node: &ast.AssignStmt{Tok: token.MUL_ASSIGN},
			want: false,
		},
		{
			name: "non-assignment node",
			node: &ast.Ident{Name: "x"},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := mut.CanMutate(tc.node)

			if got != tc.want {
				t.Fatalf("CanMutate() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDivEqualMutate(t *testing.T) {
	mut := DivEqual{}

	lhs := []ast.Expr{&ast.Ident{Name: "x"}}
	rhs := []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "3"}}
	original := &ast.AssignStmt{
		Tok: token.QUO_ASSIGN,
		Lhs: lhs,
		Rhs: rhs,
	}

	mutated := mut.Mutate(original)

	mutatedAssign, ok := mutated.(*ast.AssignStmt)
	if !ok {
		t.Fatalf("Mutate() returned %T, want *ast.AssignStmt", mutated)
	}

	if mutatedAssign == original {
		t.Fatalf("Mutate() returned the original node, want a new node")
	}

	if mutatedAssign.Tok != token.MUL_ASSIGN {
		t.Fatalf("Mutate() Tok = %v, want token.MUL_ASSIGN", mutatedAssign.Tok)
	}

	if !sameExprSlices(mutatedAssign.Lhs, lhs) {
		t.Fatalf("Mutate() Lhs = %#v, want %#v", mutatedAssign.Lhs, lhs)
	}

	if !sameExprSlices(mutatedAssign.Rhs, rhs) {
		t.Fatalf("Mutate() Rhs = %#v, want %#v", mutatedAssign.Rhs, rhs)
	}
}

func sameExprSlices(a, b []ast.Expr) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
