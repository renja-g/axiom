package conditional_boundary

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestGreaterThanName(t *testing.T) {
	mut := GreaterThan{}

	if got, want := mut.Name(), "ConditionalBoundary_GTR_GEQ"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestGreaterThanCanMutate(t *testing.T) {
	mut := GreaterThan{}

	tests := []struct {
		name string
		node ast.Node
		want bool
	}{
		{
			name: "binary GTR is mutable",
			node: &ast.BinaryExpr{Op: token.GTR},
			want: true,
		},
		{
			name: "binary GEQ is not mutable",
			node: &ast.BinaryExpr{Op: token.GEQ},
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

func TestGreaterThanMutate(t *testing.T) {
	mut := GreaterThan{}

	original := &ast.BinaryExpr{
		Op: token.GTR,
		X:  &ast.BasicLit{Kind: token.INT, Value: "2"},
		Y:  &ast.BasicLit{Kind: token.INT, Value: "1"},
	}

	mutated := mut.Mutate(original)

	mutatedExpr, ok := mutated.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("Mutate() returned %T, want *ast.BinaryExpr", mutated)
	}

	if mutatedExpr == original {
		t.Fatalf("Mutate() returned the original node, want a new node")
	}

	if mutatedExpr.Op != token.GEQ {
		t.Fatalf("Mutate() Op = %v, want token.GEQ", mutatedExpr.Op)
	}

	if mutatedExpr.X != original.X {
		t.Fatalf("Mutate() X = %p, want %p", mutatedExpr.X, original.X)
	}

	if mutatedExpr.Y != original.Y {
		t.Fatalf("Mutate() Y = %p, want %p", mutatedExpr.Y, original.Y)
	}
}
