package arithmetic

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestBitwiseAndName(t *testing.T) {
	mut := BitwiseAnd{}

	if got, want := mut.Name(), "Arithmetic_AND"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestBitwiseAndCanMutate(t *testing.T) {
	mut := BitwiseAnd{}

	tests := []struct {
		name string
		node ast.Node
		want bool
	}{
		{
			name: "binary AND is mutable",
			node: &ast.BinaryExpr{Op: token.AND},
			want: true,
		},
		{
			name: "binary OR is not mutable",
			node: &ast.BinaryExpr{Op: token.OR},
			want: false,
		},
		{
			name: "non-binary expression",
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

func TestBitwiseAndMutate(t *testing.T) {
	mut := BitwiseAnd{}

	original := &ast.BinaryExpr{
		Op: token.AND,
		X:  &ast.BasicLit{Kind: token.INT, Value: "1"},
		Y:  &ast.BasicLit{Kind: token.INT, Value: "2"},
	}

	mutated := mut.Mutate(original)

	mutatedExpr, ok := mutated.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("Mutate() returned %T, want *ast.BinaryExpr", mutated)
	}

	if mutatedExpr == original {
		t.Fatalf("Mutate() returned the original node, want a new node")
	}

	if mutatedExpr.Op != token.OR {
		t.Fatalf("Mutate() Op = %v, want token.OR", mutatedExpr.Op)
	}

	if mutatedExpr.X != original.X {
		t.Fatalf("Mutate() X = %p, want %p", mutatedExpr.X, original.X)
	}

	if mutatedExpr.Y != original.Y {
		t.Fatalf("Mutate() Y = %p, want %p", mutatedExpr.Y, original.Y)
	}
}
