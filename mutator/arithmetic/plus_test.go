package arithmetic

import (
	"go/ast"
	"go/token"
	"go/types"
	"testing"
)

func TestPlusName(t *testing.T) {
	mut := Plus{}

	if got, want := mut.Name(), "Arithmetic_ADD"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestPlusCanMutate(t *testing.T) {
	mut := Plus{}

	tests := []struct {
		name string
		node ast.Node
		want bool
	}{
		{
			name: "binary ADD is mutable",
			node: &ast.BinaryExpr{Op: token.ADD},
			want: true,
		},
		{
			name: "binary SUB is not mutable",
			node: &ast.BinaryExpr{Op: token.SUB},
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

func TestPlusCanMutateWithType(t *testing.T) {
	mut := Plus{}

	tests := []struct {
		name  string
		build func() (ast.Node, *types.Info)
		want  bool
	}{
		{
			name: "non-binary expression",
			build: func() (ast.Node, *types.Info) {
				node := &ast.Ident{Name: "x"}
				info := &types.Info{Types: make(map[ast.Expr]types.TypeAndValue)}
				return node, info
			},
			want: false,
		},
		{
			name: "binary ADD with non-string types",
			build: func() (ast.Node, *types.Info) {
				x := &ast.BasicLit{Kind: token.INT, Value: "1"}
				y := &ast.BasicLit{Kind: token.INT, Value: "2"}
				bin := &ast.BinaryExpr{Op: token.ADD, X: x, Y: y}
				info := &types.Info{Types: map[ast.Expr]types.TypeAndValue{
					x: {Type: types.Typ[types.Int]},
					y: {Type: types.Typ[types.Int]},
				}}
				return bin, info
			},
			want: true,
		},
		{
			name: "binary SUB is not mutable",
			build: func() (ast.Node, *types.Info) {
				bin := &ast.BinaryExpr{Op: token.SUB}
				info := &types.Info{Types: make(map[ast.Expr]types.TypeAndValue)}
				return bin, info
			},
			want: false,
		},
		{
			name: "binary ADD with nil type info returns true",
			build: func() (ast.Node, *types.Info) {
				x := &ast.BasicLit{Kind: token.INT, Value: "1"}
				y := &ast.BasicLit{Kind: token.INT, Value: "2"}
				bin := &ast.BinaryExpr{Op: token.ADD, X: x, Y: y}
				info := &types.Info{Types: make(map[ast.Expr]types.TypeAndValue)}
				return bin, info
			},
			want: true,
		},
		{
			name: "binary ADD with string + string is not mutable",
			build: func() (ast.Node, *types.Info) {
				x := &ast.BasicLit{Kind: token.STRING, Value: `"hello"`}
				y := &ast.BasicLit{Kind: token.STRING, Value: `"world"`}
				bin := &ast.BinaryExpr{Op: token.ADD, X: x, Y: y}
				info := &types.Info{Types: map[ast.Expr]types.TypeAndValue{
					x: {Type: types.Typ[types.String]},
					y: {Type: types.Typ[types.String]},
				}}
				return bin, info
			},
			want: false,
		},
		{
			name: "binary ADD with int + string is not mutable",
			build: func() (ast.Node, *types.Info) {
				x := &ast.BasicLit{Kind: token.INT, Value: "1"}
				y := &ast.BasicLit{Kind: token.STRING, Value: `"str"`}
				bin := &ast.BinaryExpr{Op: token.ADD, X: x, Y: y}
				info := &types.Info{Types: map[ast.Expr]types.TypeAndValue{
					x: {Type: types.Typ[types.Int]},
					y: {Type: types.Typ[types.String]},
				}}
				return bin, info
			},
			want: false,
		},
		{
			name: "binary ADD with string + int is not mutable",
			build: func() (ast.Node, *types.Info) {
				x := &ast.BasicLit{Kind: token.STRING, Value: `"str"`}
				y := &ast.BasicLit{Kind: token.INT, Value: "1"}
				bin := &ast.BinaryExpr{Op: token.ADD, X: x, Y: y}
				info := &types.Info{Types: map[ast.Expr]types.TypeAndValue{
					x: {Type: types.Typ[types.String]},
					y: {Type: types.Typ[types.Int]},
				}}
				return bin, info
			},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			node, typeInfo := tc.build()
			got := mut.CanMutateWithType(node, typeInfo)

			if got != tc.want {
				t.Fatalf("CanMutateWithType() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestPlusMutate(t *testing.T) {
	mut := Plus{}

	original := &ast.BinaryExpr{
		Op: token.ADD,
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

	if mutatedExpr.Op != token.SUB {
		t.Fatalf("Mutate() Op = %v, want token.SUB", mutatedExpr.Op)
	}

	if mutatedExpr.X != original.X {
		t.Fatalf("Mutate() X = %p, want %p", mutatedExpr.X, original.X)
	}

	if mutatedExpr.Y != original.Y {
		t.Fatalf("Mutate() Y = %p, want %p", mutatedExpr.Y, original.Y)
	}
}
