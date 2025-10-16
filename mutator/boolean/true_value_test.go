package boolean

import (
	"go/ast"
	"testing"
)

func TestTrueValueName(t *testing.T) {
	mut := TrueValue{}

	if got, want := mut.Name(), "Boolean_TRUE"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestTrueValueCanMutate(t *testing.T) {
	mut := TrueValue{}

	tests := []struct {
		name string
		node ast.Node
		want bool
	}{
		{
			name: "true identifier is mutable",
			node: &ast.Ident{Name: "true"},
			want: true,
		},
		{
			name: "false identifier is not mutable",
			node: &ast.Ident{Name: "false"},
			want: false,
		},
		{
			name: "other identifier is not mutable",
			node: &ast.Ident{Name: "x"},
			want: false,
		},
		{
			name: "non-identifier node is not mutable",
			node: &ast.BasicLit{Value: "true"},
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

func TestTrueValueMutate(t *testing.T) {
	mut := TrueValue{}

	original := &ast.Ident{
		NamePos: 42,
		Name:    "true",
	}

	mutated := mut.Mutate(original)

	mutatedIdent, ok := mutated.(*ast.Ident)
	if !ok {
		t.Fatalf("Mutate() returned %T, want *ast.Ident", mutated)
	}

	if mutatedIdent == original {
		t.Fatalf("Mutate() returned the original node, want a new node")
	}

	if mutatedIdent.Name != "false" {
		t.Fatalf("Mutate() Name = %q, want %q", mutatedIdent.Name, "false")
	}

	if mutatedIdent.NamePos != original.NamePos {
		t.Fatalf("Mutate() NamePos = %d, want %d", mutatedIdent.NamePos, original.NamePos)
	}
}
