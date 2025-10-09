package mutator

import (
	"go/ast"

	"github.com/renja-g/go-mutation-testing/mutator/conditional_boundary"
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
			conditional_boundary.GreaterThanOrEqualTo{},
			conditional_boundary.LessThanOrEqualTo{},
			conditional_boundary.GreaterThan{},
			conditional_boundary.LessThan{},
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
