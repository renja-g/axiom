package logical

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestLogicalAndName(t *testing.T) {
	mut := LogicalAnd{}

	if got, want := mut.Name(), "Logical_AND"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestLogicalAndCanMutate(t *testing.T) {
	mut := LogicalAnd{}

	tests := []struct {
		name string
		node ast.Node
		want bool
	}{
		{
			name: "logical AND is mutable",
			node: &ast.BinaryExpr{Op: token.LAND},
			want: true,
		},
		{
			name: "logical OR is not mutable",
			node: &ast.BinaryExpr{Op: token.LOR},
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

func TestLogicalAndMutate(t *testing.T) {
	mut := LogicalAnd{}

	x := &ast.Ident{Name: "x"}
	y := &ast.Ident{Name: "y"}
	original := &ast.BinaryExpr{
		Op:    token.LAND,
		X:     x,
		Y:     y,
		OpPos: token.Pos(5),
	}

	mutated := mut.Mutate(original)

	mutatedExpr, ok := mutated.(*ast.BinaryExpr)
	if !ok {
		t.Fatalf("Mutate() returned %T, want *ast.BinaryExpr", mutated)
	}

	if mutatedExpr == original {
		t.Fatalf("Mutate() returned the original node, want a new node")
	}

	if mutatedExpr.Op != token.LOR {
		t.Fatalf("Mutate() Op = %v, want token.LOR", mutatedExpr.Op)
	}

	if mutatedExpr.X != x {
		t.Fatalf("Mutate() X = %p, want %p", mutatedExpr.X, x)
	}

	if mutatedExpr.Y != y {
		t.Fatalf("Mutate() Y = %p, want %p", mutatedExpr.Y, y)
	}

	if mutatedExpr.OpPos != original.OpPos {
		t.Fatalf("Mutate() OpPos = %v, want %v", mutatedExpr.OpPos, original.OpPos)
	}
}
