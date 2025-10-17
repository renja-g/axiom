package mutator

import (
	"go/ast"

	"github.com/renja-g/axiom/mutator/arithmetic"
	"github.com/renja-g/axiom/mutator/boolean"
	"github.com/renja-g/axiom/mutator/conditional_boundary"
	"github.com/renja-g/axiom/mutator/logical"
)

// Registry holds all available mutators
type Registry struct {
	mutators []Mutator
}

// NewRegistry creates a new registry with all available mutators
func NewRegistry() *Registry {
	return &Registry{
		mutators: []Mutator{
			// Conditional Boundary Mutators
			conditional_boundary.EqualTo{},
			conditional_boundary.NotEqualTo{},
			conditional_boundary.GreaterThanOrEqualTo{},
			conditional_boundary.LessThanOrEqualTo{},
			conditional_boundary.GreaterThan{},
			conditional_boundary.LessThan{},

			// Arithmetic Mutators
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

			// Boolean Mutators
			boolean.TrueValue{},
			boolean.FalseValue{},

			// Logical Mutators
			logical.LogicalAnd{},
			logical.LogicalOr{},
		},
	}
}

// GetMutators returns all registered mutators
func (r *Registry) GetMutators() []Mutator {
	return r.mutators
}

// GetApplicableMutators returns mutators that can mutate the given node
func (r *Registry) GetApplicableMutators(node ast.Node) []Mutator {
	var applicable []Mutator
	for _, m := range r.mutators {
		if m.CanMutate(node) {
			applicable = append(applicable, m)
		}
	}
	return applicable
}
