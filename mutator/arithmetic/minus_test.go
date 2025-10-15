package arithmetic

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestMinusName(t *testing.T) {
	mut := Minus{}

	if got, want := mut.Name(), "Arithmetic_SUB"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestMinusCanMutate(t *testing.T) {
	mut := Minus{}

	tests := []struct {
		name string
		node ast.Node
		want bool
	}{
		{
			name: "binary subtraction is mutable",
			node: &ast.BinaryExpr{Op: token.SUB},
			want: true,
		},
		{
			name: "binary addition is not mutable",
			node: &ast.BinaryExpr{Op: token.ADD},
			want: false,
		},
		{
			name: "non-binary node",
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

func TestMinusMutate(t *testing.T) {
	mut := Minus{}

	original := &ast.BinaryExpr{
		Op: token.SUB,
		X:  &ast.BasicLit{Kind: token.INT, Value: "5"},
		Y:  &ast.BasicLit{Kind: token.INT, Value: "3"},
	}

	mutated := mut.Mutate(original)

	mutatedExpr, ok := mutated.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("Mutate() returned %T, want *ast.BinaryExpr", mutated)
	}

	if mutatedExpr == original {
		t.Fatalf("Mutate() returned the original node, want a new node")
	}

	if mutatedExpr.Op != token.ADD {
		t.Fatalf("Mutate() Op = %v, want token.ADD", mutatedExpr.Op)
	}

	if mutatedExpr.X != original.X {
		t.Fatalf("Mutate() X = %p, want %p", mutatedExpr.X, original.X)
	}

	if mutatedExpr.Y != original.Y {
		t.Fatalf("Mutate() Y = %p, want %p", mutatedExpr.Y, original.Y)
	}
}
