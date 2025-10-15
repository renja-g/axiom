package arithmetic

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestBitwiseNotName(t *testing.T) {
	mut := BitwiseNot{}

	if got, want := mut.Name(), "Arithmetic_NOT"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestBitwiseNotCanMutate(t *testing.T) {
	mut := BitwiseNot{}

	tests := []struct {
		name string
		node ast.Node
		want bool
	}{
		{
			name: "unary XOR is mutable",
			node: &ast.UnaryExpr{Op: token.XOR},
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

func TestBitwiseNotMutate(t *testing.T) {
	mut := BitwiseNot{}

	operand := &ast.BasicLit{Kind: token.INT, Value: "42"}
	original := &ast.UnaryExpr{Op: token.XOR, X: operand}

	mutated := mut.Mutate(original)

	if mutated != operand {
		t.Fatalf("Mutate() = %p, want operand %p", mutated, operand)
	}
}
