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
)

// Runner applies mutations and runs tests.
type Runner struct {
	workDir string
}

func New(workDir string) *Runner { return &Runner{workDir: workDir} }

// TestMutation applies a single mutation, runs `go test` on the given package, restores the file, and returns the result.
func (r *Runner) TestMutation(m model.Mutation, pkg string) (model.Result, error) {
	result := model.Result{Mutation: m}

	// read original
	original, rerr := os.ReadFile(m.FilePath)
	if rerr != nil {
		return result, rerr
	}

	// parse
	fset := token.NewFileSet()
	astFile, perr := parser.ParseFile(fset, m.FilePath, original, 0)
	if perr != nil {
		return result, perr
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
	os.WriteFile(m.FilePath, buf.Bytes(), 0644)
	defer os.WriteFile(m.FilePath, original, 0644)

	// run tests
	cmd := exec.Command("go", "test", pkg)
	if r.workDir != "" {
		cmd.Dir = r.workDir
	}
	output, _ := cmd.CombinedOutput()
	result.Output = string(output)
	if bytes.Contains(output, []byte("FAIL")) {
		result.Killed = true
	}
	return result, nil
}
