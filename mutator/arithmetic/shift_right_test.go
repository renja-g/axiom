package arithmetic

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestShiftRightName(t *testing.T) {
	mut := ShiftRight{}

	if got, want := mut.Name(), "Arithmetic_SHR"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestShiftRightCanMutate(t *testing.T) {
	mut := ShiftRight{}

	tests := []struct {
		name string
		node ast.Node
		want bool
	}{
		{
			name: "binary SHR is mutable",
			node: &ast.BinaryExpr{Op: token.SHR},
			want: true,
		},
		{
			name: "binary SHL is not mutable",
			node: &ast.BinaryExpr{Op: token.SHL},
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

func TestShiftRightMutate(t *testing.T) {
	mut := ShiftRight{}

	original := &ast.BinaryExpr{
		Op: token.SHR,
		X:  &ast.BasicLit{Kind: token.INT, Value: "8"},
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

	if mutatedExpr.Op != token.SHL {
		t.Fatalf("Mutate() Op = %v, want token.SHL", mutatedExpr.Op)
	}

	if mutatedExpr.X != original.X {
		t.Fatalf("Mutate() X = %p, want %p", mutatedExpr.X, original.X)
	}

	if mutatedExpr.Y != original.Y {
		t.Fatalf("Mutate() Y = %p, want %p", mutatedExpr.Y, original.Y)
	}
}
