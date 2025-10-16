package arithmetic

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestShiftLeftName(t *testing.T) {
	mut := ShiftLeft{}

	if got, want := mut.Name(), "Arithmetic_SHL"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestShiftLeftCanMutate(t *testing.T) {
	mut := ShiftLeft{}

	tests := []struct {
		name string
		node ast.Node
		want bool
	}{
		{
			name: "binary SHL is mutable",
			node: &ast.BinaryExpr{Op: token.SHL},
			want: true,
		},
		{
			name: "binary SHR is not mutable",
			node: &ast.BinaryExpr{Op: token.SHR},
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

func TestShiftLeftMutate(t *testing.T) {
	mut := ShiftLeft{}

	original := &ast.BinaryExpr{
		Op: token.SHL,
		X:  &ast.BasicLit{Kind: token.INT, Value: "2"},
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

	if mutatedExpr.Op != token.SHR {
		t.Fatalf("Mutate() Op = %v, want token.SHR", mutatedExpr.Op)
	}

	if mutatedExpr.X != original.X {
		t.Fatalf("Mutate() X = %p, want %p", mutatedExpr.X, original.X)
	}

	if mutatedExpr.Y != original.Y {
		t.Fatalf("Mutate() Y = %p, want %p", mutatedExpr.Y, original.Y)
	}
}
