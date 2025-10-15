package arithmetic

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestMulEqualName(t *testing.T) {
	mut := MulEqual{}

	if got, want := mut.Name(), "Arithmetic_MUL_ASSIGN"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestMulEqualCanMutate(t *testing.T) {
	mut := MulEqual{}

	tests := []struct {
		name string
		node ast.Node
		want bool
	}{
		{
			name: "multiplication assignment is mutable",
			node: &ast.AssignStmt{Tok: token.MUL_ASSIGN},
			want: true,
		},
		{
			name: "division assignment is not mutable",
			node: &ast.AssignStmt{Tok: token.QUO_ASSIGN},
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

func TestMulEqualMutate(t *testing.T) {
	mut := MulEqual{}

	lhs := []ast.Expr{&ast.Ident{Name: "x"}}
	rhs := []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "4"}}
	original := &ast.AssignStmt{
		Tok: token.MUL_ASSIGN,
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

	if mutatedAssign.Tok != token.QUO_ASSIGN {
		t.Fatalf("Mutate() Tok = %v, want token.QUO_ASSIGN", mutatedAssign.Tok)
	}

	if !sameExprSlices(mutatedAssign.Lhs, lhs) {
		t.Fatalf("Mutate() Lhs = %#v, want %#v", mutatedAssign.Lhs, lhs)
	}

	if !sameExprSlices(mutatedAssign.Rhs, rhs) {
		t.Fatalf("Mutate() Rhs = %#v, want %#v", mutatedAssign.Rhs, rhs)
	}
}
