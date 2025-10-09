package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/renja-g/go-mutation-testing/internal/model"
	"github.com/renja-g/go-mutation-testing/mutator"
)

// Generator scans Go files and discovers applicable mutations using the registry.
type Generator struct {
	registry   *mutator.Registry
	pathMapper func(string) string
}

func New(registry *mutator.Registry) *Generator {
	return &Generator{registry: registry, pathMapper: func(p string) string { return p }}
}

// WithPathMapper sets a mapping function to convert discovered sandbox file paths into output paths.
func (g *Generator) WithPathMapper(mapper func(string) string) {
	if mapper != nil {
		g.pathMapper = mapper
	}
}

// Discover walks a directory recursively and returns all discovered mutations.
func (g *Generator) Discover(rootDir string) ([]model.Mutation, error) {
	var mutations []model.Mutation

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// skip vendor and hidden dirs
			base := filepath.Base(path)
			if base == "vendor" || len(base) > 0 && base[0] == '.' {
				return nil
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}

		fset := token.NewFileSet()
		astFile, perr := parser.ParseFile(fset, path, nil, 0)
		if perr != nil {
			return perr
		}

		ast.Inspect(astFile, func(n ast.Node) bool {
			if n == nil {
				return true
			}
			for _, m := range g.registry.GetMutators() {
				if m.CanMutate(n) {
					if bin, ok := n.(*ast.BinaryExpr); ok {
						pos := fset.Position(bin.OpPos)
						mutations = append(mutations, model.Mutation{
							FilePath:   g.pathMapper(path),
							Line:       pos.Line,
							Column:     pos.Column,
							Mutator:    m,
							OriginalOp: bin.Op,
						})
					}
				}
			}
			return true
		})
		return nil
	}

	if err := filepath.Walk(rootDir, walkFn); err != nil {
		return nil, err
	}
	return mutations, nil
}
