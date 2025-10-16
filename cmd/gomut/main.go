package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/renja-g/go-mutation-testing/internal/generator"
	"github.com/renja-g/go-mutation-testing/internal/runner"
	"github.com/renja-g/go-mutation-testing/internal/sandbox"
	"github.com/renja-g/go-mutation-testing/mutator"
)

func main() {
	root := flag.String("path", "./src", "Path to source directory to mutate")
	pkg := flag.String("pkg", "./...", "Go package pattern to test (relative to path)")
	listOnly := flag.Bool("list", false, "List mutations without running tests")
	verbose := flag.Bool("v", false, "Verbose: print test output per mutation")
	flag.Parse()

	abspath, err := filepath.Abs(*root)
	if err != nil {
		panic(err)
	}

	pkgArg := normalizePkgArg(*pkg, abspath)

	sb, err := sandbox.New(abspath)
	if err != nil {
		panic(err)
	}
	defer sb.Cleanup()

	reg := mutator.NewRegistry()
	gen := generator.New(reg)
	gen.WithPathMapper(func(path string) string {
		return sb.OriginalPath(path)
	})
	muts, err := gen.Discover(sb.Root())
	if err != nil {
		panic(err)
	}

	fmt.Printf("Discovered %d mutations\n", len(muts))
	for i, m := range muts {
		fmt.Printf("[%d] %s %s:%d:%d\n", i+1, m.Mutator.Name(), displayPath(abspath, m.FilePath), m.Line, m.Column)
	}

	if *listOnly {
		return
	}

	r := runner.New(sb)
	killed, survived := 0, 0
	for i, m := range muts {
		fmt.Printf("\n[%d/%d] Testing %s at %s:%d:%d\n", i+1, len(muts), m.Mutator.Name(), displayPath(abspath, m.FilePath), m.Line, m.Column)
		res, err := r.TestMutation(m, pkgArg)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		if *verbose {
			fmt.Println(res.Output)
		}
		if res.Killed {
			fmt.Println("  ✓ KILLED")
			killed++
		} else {
			fmt.Println("  ✗ SURVIVED")
			survived++
		}
	}

	fmt.Printf("\nKilled: %d  Survived: %d  Score: %.2f%%\n", killed, survived, percent(killed, len(muts)))
}

func normalizePkgArg(pkg, root string) string {
	if pkg == "" {
		pkg = "./..."
	}

	// attempt to detect when pkg refers to the root directory path and translate to ./...
	if filepath.IsAbs(pkg) {
		if samePath(pkg, root) {
			return "./..."
		}
		return pkg
	}

	cwd, err := os.Getwd()
	if err != nil {
		return pkg
	}
	absPkg := filepath.Clean(filepath.Join(cwd, pkg))
	if samePath(absPkg, root) {
		return "./..."
	}
	return pkg
}

func samePath(a, b string) bool {
	resolve := func(p string) string {
		if p == "" {
			return ""
		}
		if !filepath.IsAbs(p) {
			if cwd, err := os.Getwd(); err == nil {
				p = filepath.Join(cwd, p)
			}
		}
		if resolved, err := filepath.EvalSymlinks(p); err == nil {
			p = resolved
		}
		return filepath.Clean(p)
	}

	return resolve(a) == resolve(b)
}

func displayPath(root, path string) string {
	if rel, err := filepath.Rel(root, path); err == nil {
		return rel
	}
	return path
}

func percent(killed, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(killed) / float64(total) * 100
}
