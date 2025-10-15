package arithmetic

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestMinusEqualName(t *testing.T) {
	mut := MinusEqual{}

	if got, want := mut.Name(), "Arithmetic_SUB_ASSIGN"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestMinusEqualCanMutate(t *testing.T) {
	mut := MinusEqual{}

	tests := []struct {
		name string
		node ast.Node
		want bool
	}{
		{
			name: "subtraction assignment is mutable",
			node: &ast.AssignStmt{Tok: token.SUB_ASSIGN},
			want: true,
		},
		{
			name: "addition assignment is not mutable",
			node: &ast.AssignStmt{Tok: token.ADD_ASSIGN},
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

func TestMinusEqualMutate(t *testing.T) {
	mut := MinusEqual{}

	lhs := []ast.Expr{&ast.Ident{Name: "x"}}
	rhs := []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "5"}}
	original := &ast.AssignStmt{
		Tok: token.SUB_ASSIGN,
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

	if mutatedAssign.Tok != token.ADD_ASSIGN {
		t.Fatalf("Mutate() Tok = %v, want token.ADD_ASSIGN", mutatedAssign.Tok)
	}

	if !sameExprSlices(mutatedAssign.Lhs, lhs) {
		t.Fatalf("Mutate() Lhs = %#v, want %#v", mutatedAssign.Lhs, lhs)
	}

	if !sameExprSlices(mutatedAssign.Rhs, rhs) {
		t.Fatalf("Mutate() Rhs = %#v, want %#v", mutatedAssign.Rhs, rhs)
	}
}
