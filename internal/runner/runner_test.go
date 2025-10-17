package runner

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/renja-g/axiom/internal/model"
	"github.com/renja-g/axiom/internal/sandbox"
)

type binaryOpMutator struct {
	target token.Token
	name   string
}

func (m binaryOpMutator) Name() string { return m.name }

func (m binaryOpMutator) CanMutate(node ast.Node) bool {
	_, ok := node.(*ast.BinaryExpr)
	return ok
}

func (m binaryOpMutator) Mutate(node ast.Node) ast.Node {
	bin, ok := node.(*ast.BinaryExpr)
	if !ok {
		return node
	}
	cloned := *bin
	cloned.Op = m.target
	return &cloned
}

func TestRunnerTestMutationKills(t *testing.T) {
	fx := newRunnerFixture(t)
	mutation := model.Mutation{
		FilePath: fx.filePath,
		Line:     fx.line,
		Column:   fx.column,
		Mutator:  binaryOpMutator{name: "less-than", target: token.LSS},
	}

	result, err := fx.runner.TestMutation(mutation, ".")
	if err != nil {
		t.Fatalf("TestMutation returned error: %v", err)
	}
	if !result.Killed {
		t.Fatalf("expected mutation to be killed, got result: %+v", result)
	}

	assertFileRestored(t, fx.sandboxPath, fx.originalContent)
}

func TestRunnerTestMutationSurvives(t *testing.T) {
	fx := newRunnerFixture(t)
	mutation := model.Mutation{
		FilePath: fx.filePath,
		Line:     fx.line,
		Column:   fx.column,
		Mutator:  binaryOpMutator{name: "greater-equal", target: token.GEQ},
	}

	result, err := fx.runner.TestMutation(mutation, ".")
	if err != nil {
		t.Fatalf("TestMutation returned error: %v", err)
	}
	if result.Killed {
		t.Fatalf("expected mutation to survive, got killed result: %+v", result)
	}

	assertFileRestored(t, fx.sandboxPath, fx.originalContent)
}

func TestRunnerTestMutationMissingFile(t *testing.T) {
	r := New(nil)
	missingPath := filepath.Join(t.TempDir(), "does", "not", "exist.go")
	mutation := model.Mutation{FilePath: missingPath}

	result, err := r.TestMutation(mutation, ".")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
	if result.Mutation.FilePath != mutation.FilePath {
		t.Fatalf("expected mutation in result to match input path, got %q", result.Mutation.FilePath)
	}
}

type runnerFixture struct {
	runner          *Runner
	filePath        string
	sandboxPath     string
	line            int
	column          int
	originalContent []byte
}

func newRunnerFixture(t *testing.T) runnerFixture {
	t.Helper()

	root := t.TempDir()
	writeFile(t, filepath.Join(root, "go.mod"), "module example.com/runnerfixture\n\ngo 1.21\n")
	sampleSource := "package sample\n\nfunc Compare(a, b int) bool {\n\treturn a > b\n}\n"
	writeFile(t, filepath.Join(root, "sample.go"), sampleSource)
	writeFile(t, filepath.Join(root, "sample_test.go"), `package sample

import "testing"

func TestCompare(t *testing.T) {
	if !Compare(2, 1) {
		t.Fatalf("expected true")
	}
}
`)

	sb, err := sandbox.New(root)
	if err != nil {
		t.Fatalf("failed to create sandbox: %v", err)
	}
	t.Cleanup(func() {
		if cerr := sb.Cleanup(); cerr != nil {
			t.Fatalf("failed to cleanup sandbox: %v", cerr)
		}
	})

	samplePath := filepath.Join(root, "sample.go")
	mirrorPath := sb.MirrorPath(samplePath)
	originalContent, err := os.ReadFile(mirrorPath)
	if err != nil {
		t.Fatalf("failed to read sandbox file: %v", err)
	}

	line, column := findBinaryPosition(t, samplePath, token.GTR)

	return runnerFixture{
		runner:          New(sb),
		filePath:        samplePath,
		sandboxPath:     mirrorPath,
		line:            line,
		column:          column,
		originalContent: originalContent,
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create directories for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write file %s: %v", path, err)
	}
}

func findBinaryPosition(t *testing.T, path string, op token.Token) (int, int) {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		t.Fatalf("failed to parse %s: %v", path, err)
	}

	var position token.Position
	ast.Inspect(file, func(n ast.Node) bool {
		bin, ok := n.(*ast.BinaryExpr)
		if !ok {
			return true
		}
		if bin.Op == op {
			position = fset.Position(bin.OpPos)
			return false
		}
		return true
	})
	if position.Line == 0 {
		t.Fatalf("operator %s not found in %s", op.String(), path)
	}
	return position.Line, position.Column
}

func assertFileRestored(t *testing.T, path string, want []byte) {
	t.Helper()
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}
	if string(got) != string(want) {
		t.Fatalf("sandbox file %s not restored to original content", path)
	}
}
