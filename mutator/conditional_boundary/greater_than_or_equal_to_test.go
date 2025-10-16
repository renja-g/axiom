package conditional_boundary

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestGreaterThanOrEqualToName(t *testing.T) {
	mut := GreaterThanOrEqualTo{}

	if got, want := mut.Name(), "ConditionalBoundary_GEQ_GTR"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestGreaterThanOrEqualToCanMutate(t *testing.T) {
	mut := GreaterThanOrEqualTo{}

	tests := []struct {
		name string
		node ast.Node
		want bool
	}{
		{
			name: "binary GEQ is mutable",
			node: &ast.BinaryExpr{Op: token.GEQ},
			want: true,
		},
		{
			name: "binary GTR is not mutable",
			node: &ast.BinaryExpr{Op: token.GTR},
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

func TestGreaterThanOrEqualToMutate(t *testing.T) {
	mut := GreaterThanOrEqualTo{}

	original := &ast.BinaryExpr{
		Op: token.GEQ,
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

	if mutatedExpr.Op != token.GTR {
		t.Fatalf("Mutate() Op = %v, want token.GTR", mutatedExpr.Op)
	}

	if mutatedExpr.X != original.X {
		t.Fatalf("Mutate() X = %p, want %p", mutatedExpr.X, original.X)
	}

	if mutatedExpr.Y != original.Y {
		t.Fatalf("Mutate() Y = %p, want %p", mutatedExpr.Y, original.Y)
	}
}
