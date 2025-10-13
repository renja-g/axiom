package mutator

import (
	"go/ast"
	"go/types"
)

// Mutator defines the interface that all mutators must implement
type Mutator interface {
	// Name returns the unique identifier for this mutator
	Name() string

	// CanMutate checks if this mutator can be applied to the given node
	CanMutate(node ast.Node) bool

	// Mutate applies the mutation to the node and returns the mutated node
	Mutate(node ast.Node) ast.Node
}

// TypeAwareMutator is an optional interface that mutators can implement
// when they need type information to make mutation decisions.
// The generator will only perform type checking when at least one registered
// mutator implements this interface, ensuring no performance degradation
// for mutators that don't need type information.
type TypeAwareMutator interface {
	Mutator
	// CanMutateWithType checks if this mutator can be applied to the given node
	// using type information. This is called instead of CanMutate when type info is available.
	// Note that the CanMutate should still be implemented to handle cases where type info is not available.
	CanMutateWithType(node ast.Node, typeInfo *types.Info) bool
}
