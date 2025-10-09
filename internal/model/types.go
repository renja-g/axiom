package model

import (
	"go/ast"
	"go/token"

	"github.com/renja-g/go-mutation-testing/mutator"
)

// Mutation describes a specific change to be applied at a location.
type Mutation struct {
	FilePath   string
	Line       int
	Column     int
	Mutator    mutator.Mutator
	OriginalOp token.Token // for BinaryExpr cases
}

// Target ties a parsed file to its AST and fset for reuse.
type Target struct {
	FilePath string
	Fset     *token.FileSet
	AST      *ast.File
}

// Result captures the outcome of a mutation test run.
type Result struct {
	Mutation Mutation
	Killed   bool
	Output   string
}
