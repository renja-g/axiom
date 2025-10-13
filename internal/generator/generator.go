package generator

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"

	"github.com/renja-g/go-mutation-testing/internal/model"
	"github.com/renja-g/go-mutation-testing/mutator"
)

// Generator scans Go files and discovers applicable mutations using the registry.
type Generator struct {
	registry          *mutator.Registry
	pathMapper        func(string) string
	needsTypeCheck    bool
	typeAwareMutators []mutator.TypeAwareMutator
}

func New(registry *mutator.Registry) *Generator {
	g := &Generator{
		registry:   registry,
		pathMapper: func(p string) string { return p },
	}

	// Check if any mutators need type information
	for _, m := range registry.GetMutators() {
		if tm, ok := m.(mutator.TypeAwareMutator); ok {
			g.needsTypeCheck = true
			g.typeAwareMutators = append(g.typeAwareMutators, tm)
		}
	}

	return g
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

	// Group files by package directory for proper type checking
	pkgFiles := make(map[string][]string)

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// skip vendor and hidden dirs
			base := filepath.Base(path)
			if base == "vendor" || len(base) > 0 && base[0] == '.' {
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}

		pkgDir := filepath.Dir(path)
		pkgFiles[pkgDir] = append(pkgFiles[pkgDir], path)
		return nil
	}

	if err := filepath.Walk(rootDir, walkFn); err != nil {
		return nil, err
	}

	// Process each package
	for pkgDir, files := range pkgFiles {
		pkgMutations, err := g.discoverInPackage(pkgDir, files)
		if err != nil {
			return nil, err
		}
		mutations = append(mutations, pkgMutations...)
	}

	return mutations, nil
}

// discoverInPackage processes all files in a package together for proper type checking
func (g *Generator) discoverInPackage(pkgDir string, filePaths []string) ([]model.Mutation, error) {
	var mutations []model.Mutation

	fset := token.NewFileSet()
	var astFiles []*ast.File
	fileMap := make(map[*ast.File]string) // Map AST files back to their paths

	// Parse all files in the package
	for _, path := range filePaths {
		astFile, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil, err
		}
		astFiles = append(astFiles, astFile)
		fileMap[astFile] = path
	}

	// Perform type checking at package level if needed
	var typeInfo *types.Info
	if g.needsTypeCheck {
		typeInfo = g.performTypeCheckPackage(fset, astFiles)
		// If type checking fails, we can still use non-type-aware mutators
	}

	// Inspect each file for mutations
	for _, astFile := range astFiles {
		filePath := fileMap[astFile]

		ast.Inspect(astFile, func(n ast.Node) bool {
			if n == nil {
				return true
			}
			for _, m := range g.registry.GetMutators() {
				canMutate := false

				// Check if this is a type-aware mutator and we have type info
				if tm, ok := m.(mutator.TypeAwareMutator); ok && typeInfo != nil {
					canMutate = tm.CanMutateWithType(n, typeInfo)
				} else {
					// Fall back to regular CanMutate
					canMutate = m.CanMutate(n)
				}

				if canMutate {
					if bin, ok := n.(*ast.BinaryExpr); ok {
						pos := fset.Position(bin.OpPos)
						mutations = append(mutations, model.Mutation{
							FilePath:   g.pathMapper(filePath),
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
	}

	return mutations, nil
}

// performTypeCheckPackage runs the type checker on a package and returns type information.
// Returns nil if type checking fails (allows graceful degradation).
func (g *Generator) performTypeCheckPackage(fset *token.FileSet, astFiles []*ast.File) *types.Info {
	if len(astFiles) == 0 {
		return nil
	}

	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
	}

	conf := types.Config{
		Importer: importer.Default(),
		Error:    func(err error) {}, // Ignore type errors to allow partial type info
	}

	// Use the package name from the first file
	pkgName := astFiles[0].Name.Name

	_, err := conf.Check(pkgName, fset, astFiles, info)
	if err != nil {
		// Type checking had errors, but we might still have partial type info
		// Return the info anyway as it may be useful
		if len(info.Types) > 0 {
			return info
		}
		return nil
	}

	return info
}
