package arithmetic

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestIntegerLiteralBoundaryName(t *testing.T) {
	mut := IntegerLiteralBoundary{}

	if got, want := mut.Name(), "Arithmetic_INT_LITERAL_BOUNDARY"; got != want {
		t.Fatalf("Name() = %q, want %q", got, want)
	}
}

func TestIntegerLiteralBoundaryCanMutate_PositiveCases(t *testing.T) {
	mut := IntegerLiteralBoundary{}

	cases := []ast.Node{
		&ast.BasicLit{Kind: token.INT, Value: "0"},
		&ast.BasicLit{Kind: token.INT, Value: "1"},
		&ast.BasicLit{Kind: token.INT, Value: "42"},
		&ast.BasicLit{Kind: token.INT, Value: "-7"},
	}

	for _, n := range cases {
		if !mut.CanMutate(n) {
			t.Fatalf("expected CanMutate to be true for node: %#v", n)
		}
	}
}

func TestIntegerLiteralBoundaryCanMutate_NegativeCases(t *testing.T) {
	mut := IntegerLiteralBoundary{}

	negatives := []ast.Node{
		&ast.BasicLit{Kind: token.FLOAT, Value: "1.0"},
		&ast.BasicLit{Kind: token.STRING, Value: "\"str\""},
		&ast.Ident{Name: "x"},
		&ast.UnaryExpr{Op: token.SUB, X: &ast.BasicLit{Kind: token.INT, Value: "1"}},
	}

	for _, n := range negatives {
		if mut.CanMutate(n) {
			t.Fatalf("expected CanMutate to be false for node: %#v", n)
		}
	}
}

func TestIntegerLiteralBoundaryMutate_ZeroAndOneSwap(t *testing.T) {
	mut := IntegerLiteralBoundary{}

	zero := &ast.BasicLit{Kind: token.INT, Value: "0"}
	one := &ast.BasicLit{Kind: token.INT, Value: "1"}

	mutatedZero := mut.Mutate(zero)
	mz, ok := mutatedZero.(*ast.BasicLit)
	if !ok {
		t.Fatalf("Mutate(0) returned %T, want *ast.BasicLit", mutatedZero)
	}
	if mz == zero {
		t.Fatalf("Mutate(0) returned original node")
	}
	if mz.Kind != token.INT || mz.Value != "1" {
		t.Fatalf("Mutate(0) = %v %q, want INT \"1\"", mz.Kind, mz.Value)
	}

	mutatedOne := mut.Mutate(one)
	mo, ok := mutatedOne.(*ast.BasicLit)
	if !ok {
		t.Fatalf("Mutate(1) returned %T, want *ast.BasicLit", mutatedOne)
	}
	if mo.Kind != token.INT || mo.Value != "0" {
		t.Fatalf("Mutate(1) = %v %q, want INT \"0\"", mo.Kind, mo.Value)
	}
}

func TestIntegerLiteralBoundaryMutate_IncrementAndDecrement(t *testing.T) {
	mut := IntegerLiteralBoundary{}

	// positive number decremented by 1
	fortyTwo := &ast.BasicLit{Kind: token.INT, Value: "42"}
	m1 := mut.Mutate(fortyTwo)
	b1, ok := m1.(*ast.BasicLit)
	if !ok {
		t.Fatalf("Mutate(42) returned %T, want *ast.BasicLit", m1)
	}
	if b1.Value != "41" {
		t.Fatalf("Mutate(42) Value = %q, want \"41\"", b1.Value)
	}

	// negative number incremented by 1
	negSeven := &ast.BasicLit{Kind: token.INT, Value: "-7"}
	m2 := mut.Mutate(negSeven)
	b2, ok := m2.(*ast.BasicLit)
	if !ok {
		t.Fatalf("Mutate(-7) returned %T, want *ast.BasicLit", m2)
	}
	if b2.Value != "-6" {
		t.Fatalf("Mutate(-7) Value = %q, want \"-6\"", b2.Value)
	}
}

func TestIntegerLiteralBoundaryMutate_NonIntegerIsNoop(t *testing.T) {
	mut := IntegerLiteralBoundary{}
	node := &ast.BasicLit{Kind: token.FLOAT, Value: "3.14"}
	if got := mut.Mutate(node); got != node {
		t.Fatalf("expected non-integer mutate to be noop; got %T", got)
	}
}
