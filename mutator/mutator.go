package mutator

import "go/ast"

// Mutator defines the interface that all mutators must implement
type Mutator interface {
	// Name returns the unique identifier for this mutator
	Name() string

	// CanMutate checks if this mutator can be applied to the given node
	CanMutate(node ast.Node) bool

	// Mutate applies the mutation to the node and returns the mutated node
	Mutate(node ast.Node) ast.Node
}
