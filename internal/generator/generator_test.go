package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go/token"

	"github.com/renja-g/go-mutation-testing/internal/model"
	"github.com/renja-g/go-mutation-testing/mutator"
)

func TestNewDetectsTypeAwareMutators(t *testing.T) {
	registry := mutator.NewRegistry()
	gen := New(registry)

	if !gen.needsTypeCheck {
		t.Fatalf("expected generator to require type checking when type-aware mutators are registered")
	}

	if len(gen.typeAwareMutators) == 0 {
		t.Fatalf("expected generator to collect type-aware mutators")
	}
}

func TestWithPathMapperOverridesDefault(t *testing.T) {
	registry := mutator.NewRegistry()
	gen := New(registry)

	if got := gen.pathMapper("input.go"); got != "input.go" {
		t.Fatalf("default path mapper should return the original path, got %q", got)
	}

	gen.WithPathMapper(func(p string) string {
		return "mapped:" + p
	})

	if got := gen.pathMapper("file.go"); got != "mapped:file.go" {
		t.Fatalf("expected custom path mapper to be used, got %q", got)
	}

	gen.WithPathMapper(nil)

	if got := gen.pathMapper("again.go"); got != "mapped:again.go" {
		t.Fatalf("nil mapper should not override the existing mapper, got %q", got)
	}
}

func TestDiscoverAppliesPathMapperAndCollectsMutations(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "sample.go")

	source := "package sample\n\nfunc cmp(a, b int) bool {\n\treturn a > b\n}\n"

	if err := os.WriteFile(filePath, []byte(source), 0o644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	registry := mutator.NewRegistry()
	gen := New(registry)

	mapper := func(p string) string {
		return strings.Replace(p, dir, "/virtual", 1)
	}
	gen.WithPathMapper(mapper)

	mutations, err := gen.Discover(dir)
	if err != nil {
		t.Fatalf("Discover returned error: %v", err)
	}

	if len(mutations) != 1 {
		t.Fatalf("expected 1 mutation, got %d", len(mutations))
	}

	mutation := mutations[0]
	expectedPath := mapper(filePath)
	if mutation.FilePath != expectedPath {
		t.Fatalf("expected mapped path %q, got %q", expectedPath, mutation.FilePath)
	}

	if mutation.Mutator.Name() != "ConditionalBoundary_GTR_GEQ" {
		t.Fatalf("unexpected mutator name: %q", mutation.Mutator.Name())
	}

	if mutation.OriginalOp != token.GTR {
		t.Fatalf("expected original operator to be token.GTR, got %v", mutation.OriginalOp)
	}

	if mutation.Line == 0 || mutation.Column == 0 {
		t.Fatalf("expected non-zero line and column, got line=%d column=%d", mutation.Line, mutation.Column)
	}
}

func TestDiscoverTypeAwareMutatorsUseTypeInfo(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "sample.go")

	source := "package sample\n\nfunc concat(s string) string {\n\treturn s + \"x\"\n}\n\nfunc sum(a, b int) int {\n\treturn a + b\n}\n"

	if err := os.WriteFile(filePath, []byte(source), 0o644); err != nil {
		t.Fatalf("failed to write source file: %v", err)
	}

	registry := mutator.NewRegistry()
	gen := New(registry)

	mutations, err := gen.Discover(dir)
	if err != nil {
		t.Fatalf("Discover returned error: %v", err)
	}

	var addMutation *model.Mutation
	for i := range mutations {
		mutation := &mutations[i]
		if mutation.Mutator.Name() == "Arithmetic_ADD" {
			if addMutation != nil {
				t.Fatalf("expected exactly one arithmetic addition mutation, found multiple")
			}
			addMutation = mutation
		}
	}

	if addMutation == nil {
		t.Fatalf("expected to find an arithmetic addition mutation, but none were discovered")
	}

	if addMutation.Line != 8 {
		t.Fatalf("expected addition mutation on line 8, got line %d", addMutation.Line)
	}
}
