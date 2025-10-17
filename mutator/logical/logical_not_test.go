package logical

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestLogicalNotName(t *testing.T) {
	mut := LogicalNot{}

	if got, want := mut.Name(), "Logical_NOT"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestLogicalNotCanMutate(t *testing.T) {
	mut := LogicalNot{}

	tests := []struct {
		name string
		node ast.Node
		want bool
	}{
		{
			name: "unary NOT is mutable",
			node: &ast.UnaryExpr{Op: token.NOT},
			want: true,
		},
		{
			name: "unary ADD is not mutable",
			node: &ast.UnaryExpr{Op: token.ADD},
			want: false,
		},
		{
			name: "non-unary node",
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

func TestLogicalNotMutate(t *testing.T) {
	mut := LogicalNot{}

	operand := &ast.Ident{Name: "x"}
	original := &ast.UnaryExpr{Op: token.NOT, X: operand}

	mutated := mut.Mutate(original)

	if mutated != operand {
		t.Fatalf("Mutate() = %p, want operand %p", mutated, operand)
	}
}
