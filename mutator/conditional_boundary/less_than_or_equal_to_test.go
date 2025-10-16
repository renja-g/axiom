package conditional_boundary

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestLessThanOrEqualToName(t *testing.T) {
	mut := LessThanOrEqualTo{}

	if got, want := mut.Name(), "ConditionalBoundary_LEQ_LSS"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestLessThanOrEqualToCanMutate(t *testing.T) {
	mut := LessThanOrEqualTo{}

	tests := []struct {
		name string
		node ast.Node
		want bool
	}{
		{
			name: "binary LEQ is mutable",
			node: &ast.BinaryExpr{Op: token.LEQ},
			want: true,
		},
		{
			name: "binary LSS is not mutable",
			node: &ast.BinaryExpr{Op: token.LSS},
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

func TestLessThanOrEqualToMutate(t *testing.T) {
	mut := LessThanOrEqualTo{}

	original := &ast.BinaryExpr{
		Op: token.LEQ,
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

	if mutatedExpr.Op != token.LSS {
		t.Fatalf("Mutate() Op = %v, want token.LSS", mutatedExpr.Op)
	}

	if mutatedExpr.X != original.X {
		t.Fatalf("Mutate() X = %p, want %p", mutatedExpr.X, original.X)
	}

	if mutatedExpr.Y != original.Y {
		t.Fatalf("Mutate() Y = %p, want %p", mutatedExpr.Y, original.Y)
	}
}
