package mutator

import (
	"go/ast"
	"reflect"
	"testing"

	"github.com/renja-g/go-mutation-testing/mutator/arithmetic"
	"github.com/renja-g/go-mutation-testing/mutator/boolean"
	"github.com/renja-g/go-mutation-testing/mutator/conditional_boundary"
	"github.com/renja-g/go-mutation-testing/mutator/logical"
)

func TestNewRegistryIncludesExpectedMutators(t *testing.T) {
	registry := NewRegistry()

	found := make(map[string]struct{})
	for _, m := range registry.GetMutators() {
		found[reflect.TypeOf(m).String()] = struct{}{}
	}

	expected := []Mutator{
		conditional_boundary.GreaterThanOrEqualTo{},
		conditional_boundary.LessThanOrEqualTo{},
		conditional_boundary.GreaterThan{},
		conditional_boundary.LessThan{},
		arithmetic.BitwiseAnd{},
		arithmetic.BitwiseNot{},
		arithmetic.BitwiseOr{},
		arithmetic.BitwiseXor{},
		arithmetic.DivEqual{},
		arithmetic.Division{},
		arithmetic.MinusEqual{},
		arithmetic.Minus{},
		arithmetic.ModEqual{},
		arithmetic.Modulus{},
		arithmetic.MulEqual{},
		arithmetic.Multiplication{},
		arithmetic.PlusEqual{},
		arithmetic.Plus{},
		arithmetic.ShiftLeft{},
		arithmetic.ShiftRight{},
		boolean.TrueValue{},
		boolean.FalseValue{},
		logical.LogicalAnd{},
		logical.LogicalOr{},
	}

	for _, mut := range expected {
		typ := reflect.TypeOf(mut).String()
		if _, ok := found[typ]; !ok {
			t.Fatalf("expected mutator %s to be registered", typ)
		}
	}
}

type mockMutator struct {
	name      string
	canMutate bool
}

func (m mockMutator) Name() string {
	return m.name
}

func (m mockMutator) CanMutate(node ast.Node) bool {
	return m.canMutate
}

func (m mockMutator) Mutate(node ast.Node) ast.Node {
	return node
}

func TestRegistryGetApplicableMutators(t *testing.T) {
	applicable := mockMutator{name: "applicable", canMutate: true}
	notApplicable := mockMutator{name: "not-applicable", canMutate: false}

	registry := &Registry{mutators: []Mutator{applicable, notApplicable}}

	result := registry.GetApplicableMutators(nil)
	if len(result) != 1 {
		t.Fatalf("expected 1 applicable mutator, got %d", len(result))
	}
	if result[0].Name() != applicable.Name() {
		t.Fatalf("expected applicable mutator %q, got %q", applicable.Name(), result[0].Name())
	}
}
