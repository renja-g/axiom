package runner

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"os/exec"

	"github.com/renja-g/go-mutation-testing/internal/model"
	"github.com/renja-g/go-mutation-testing/internal/sandbox"
)

// Runner applies mutations and runs tests inside a sandbox copy.
type Runner struct {
	sandbox *sandbox.Sandbox
}

func New(sb *sandbox.Sandbox) *Runner { return &Runner{sandbox: sb} }

// TestMutation applies a single mutation, runs `go test` on the given package, restores the file, and returns the result.
func (r *Runner) TestMutation(m model.Mutation, pkg string) (result model.Result, err error) {
	result = model.Result{Mutation: m}

	// determine sandbox path equivalent
	path := m.FilePath
	if r.sandbox != nil {
		path = r.sandbox.MirrorPath(m.FilePath)
	}

	// read original
	original, rerr := os.ReadFile(path)
	if rerr != nil {
		return result, rerr
	}

	// parse
	fset := token.NewFileSet()
	astFile, perr := parser.ParseFile(fset, path, original, 0)
	if perr != nil {
		err = perr
		return
	}

	// apply mutation
	ast.Inspect(astFile, func(n ast.Node) bool {
		bin, ok := n.(*ast.BinaryExpr)
		if !ok {
			return true
		}
		pos := fset.Position(bin.OpPos)
		if pos.Line == m.Line && pos.Column == m.Column {
			mut := m.Mutator.Mutate(n)
			if mutatedBin, ok := mut.(*ast.BinaryExpr); ok {
				bin.Op = mutatedBin.Op
			}
			return false
		}
		return true
	})

	// write mutated
	var buf bytes.Buffer
	printer.Fprint(&buf, fset, astFile)
	if werr := os.WriteFile(path, buf.Bytes(), 0644); werr != nil {
		err = werr
		return
	}
	defer func() {
		if restoreErr := os.WriteFile(path, original, 0644); err == nil && restoreErr != nil {
			err = restoreErr
		}
	}()

	// run tests
	cmd := exec.Command("go", "test", pkg)
	if r.sandbox != nil {
		cmd.Dir = r.sandbox.Root()
	}
	output, cmdErr := cmd.CombinedOutput()
	result.Output = string(output)
	if cmdErr != nil {
		if exitErr, ok := cmdErr.(*exec.ExitError); ok {
			// Non-zero exit means the mutation was killed.
			result.Killed = exitErr.ExitCode() != 0
			return
		}
		err = cmdErr
		return
	}
	result.Killed = false
	return
}
