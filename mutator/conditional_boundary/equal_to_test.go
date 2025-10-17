package conditional_boundary

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestEqualToName(t *testing.T) {
	mut := EqualTo{}

	if got, want := mut.Name(), "ConditionalBoundary_EQL_NEQ"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestEqualToCanMutate(t *testing.T) {
	mut := EqualTo{}

	tests := []struct {
		name string
		node ast.Node
		want bool
	}{
		{
			name: "binary EQL is mutable",
			node: &ast.BinaryExpr{Op: token.EQL},
			want: true,
		},
		{
			name: "binary NEQ is not mutable",
			node: &ast.BinaryExpr{Op: token.NEQ},
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

func TestEqualToMutate(t *testing.T) {
	mut := EqualTo{}

	original := &ast.BinaryExpr{
		Op: token.EQL,
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

	if mutatedExpr.Op != token.NEQ {
		t.Fatalf("Mutate() Op = %v, want token.NEQ", mutatedExpr.Op)
	}

	if mutatedExpr.X != original.X {
		t.Fatalf("Mutate() X = %p, want %p", mutatedExpr.X, original.X)
	}

	if mutatedExpr.Y != original.Y {
		t.Fatalf("Mutate() Y = %p, want %p", mutatedExpr.Y, original.Y)
	}
}
