package arithmetic

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestDecrementName(t *testing.T) {
	mut := Decrement{}

	if got, want := mut.Name(), "Arithmetic_DEC"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestDecrementCanMutate(t *testing.T) {
	mut := Decrement{}

	tests := []struct {
		name string
		node ast.Node
		want bool
	}{
		{
			name: "DEC is mutable",
			node: &ast.IncDecStmt{Tok: token.DEC},
			want: true,
		},
		{
			name: "INC is not mutable",
			node: &ast.IncDecStmt{Tok: token.INC},
			want: false,
		},
		{
			name: "non-incdec node",
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

func TestDecrementMutate(t *testing.T) {
	mut := Decrement{}

	x := &ast.Ident{Name: "i"}
	original := &ast.IncDecStmt{X: x, Tok: token.DEC, TokPos: token.Pos(11)}

	mutated := mut.Mutate(original)

	stmt, ok := mutated.(*ast.IncDecStmt)
	if !ok {
		t.Fatalf("Mutate() returned %T, want *ast.IncDecStmt", mutated)
	}

	if stmt == original {
		t.Fatalf("Mutate() returned the original node, want a new node")
	}

	if stmt.Tok != token.INC {
		t.Fatalf("Mutate() Tok = %v, want token.INC", stmt.Tok)
	}

	if stmt.X != x {
		t.Fatalf("Mutate() X = %p, want %p", stmt.X, x)
	}

	if stmt.TokPos != original.TokPos {
		t.Fatalf("Mutate() TokPos = %v, want %v", stmt.TokPos, original.TokPos)
	}
}
