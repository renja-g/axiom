package boolean

import (
	"go/ast"
	"testing"
)

func TestFalseValueName(t *testing.T) {
	mut := FalseValue{}

	if got, want := mut.Name(), "Boolean_FALSE"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestFalseValueCanMutate(t *testing.T) {
	mut := FalseValue{}

	tests := []struct {
		name string
		node ast.Node
		want bool
	}{
		{
			name: "false identifier is mutable",
			node: &ast.Ident{Name: "false"},
			want: true,
		},
		{
			name: "true identifier is not mutable",
			node: &ast.Ident{Name: "true"},
			want: false,
		},
		{
			name: "other identifier is not mutable",
			node: &ast.Ident{Name: "x"},
			want: false,
		},
		{
			name: "non-identifier node is not mutable",
			node: &ast.BasicLit{Value: "false"},
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

func TestFalseValueMutate(t *testing.T) {
	mut := FalseValue{}

	original := &ast.Ident{
		NamePos: 42,
		Name:    "false",
	}

	mutated := mut.Mutate(original)

	mutatedIdent, ok := mutated.(*ast.Ident)
	if !ok {
		t.Fatalf("Mutate() returned %T, want *ast.Ident", mutated)
	}

	if mutatedIdent == original {
		t.Fatalf("Mutate() returned the original node, want a new node")
	}

	if mutatedIdent.Name != "true" {
		t.Fatalf("Mutate() Name = %q, want %q", mutatedIdent.Name, "true")
	}

	if mutatedIdent.NamePos != original.NamePos {
		t.Fatalf("Mutate() NamePos = %d, want %d", mutatedIdent.NamePos, original.NamePos)
	}
}
