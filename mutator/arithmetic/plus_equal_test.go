package arithmetic

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestPlusEqualName(t *testing.T) {
	mut := PlusEqual{}

	if got, want := mut.Name(), "Arithmetic_ADD_ASSIGN"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestPlusEqualCanMutate(t *testing.T) {
	mut := PlusEqual{}

	tests := []struct {
		name string
		node ast.Node
		want bool
	}{
		{
			name: "addition assignment is mutable",
			node: &ast.AssignStmt{Tok: token.ADD_ASSIGN},
			want: true,
		},
		{
			name: "subtraction assignment is not mutable",
			node: &ast.AssignStmt{Tok: token.SUB_ASSIGN},
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

func TestPlusEqualMutate(t *testing.T) {
	mut := PlusEqual{}

	lhs := []ast.Expr{&ast.Ident{Name: "x"}}
	rhs := []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "10"}}
	original := &ast.AssignStmt{
		Tok: token.ADD_ASSIGN,
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

	if mutatedAssign.Tok != token.SUB_ASSIGN {
		t.Fatalf("Mutate() Tok = %v, want token.SUB_ASSIGN", mutatedAssign.Tok)
	}

	if !sameExprSlices(mutatedAssign.Lhs, lhs) {
		t.Fatalf("Mutate() Lhs = %#v, want %#v", mutatedAssign.Lhs, lhs)
	}

	if !sameExprSlices(mutatedAssign.Rhs, rhs) {
		t.Fatalf("Mutate() Rhs = %#v, want %#v", mutatedAssign.Rhs, rhs)
	}
}
