package sandbox

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Sandbox copies a source tree into a temporary directory where mutations can run safely.
// It keeps track of both the original root and the mirrored sandbox root to translate paths.
type Sandbox struct {
	originalRoot string
	root         string
}

// New creates a sandbox by recursively copying the given root directory into a temporary directory.
func New(originalRoot string) (*Sandbox, error) {
	info, err := os.Stat(originalRoot)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, &fs.PathError{Op: "sandbox", Path: originalRoot, Err: fs.ErrInvalid}
	}

	tempDir, err := os.MkdirTemp("", "gomut-sandbox-*")
	if err != nil {
		return nil, err
	}

	// Copy the tree into the sandbox.
	if err := copyTree(originalRoot, tempDir); err != nil {
		os.RemoveAll(tempDir)
		return nil, err
	}

	return &Sandbox{originalRoot: originalRoot, root: tempDir}, nil
}

// Root returns the filesystem path to the sandbox copy.
func (s *Sandbox) Root() string {
	return s.root
}

// MirrorPath converts an original file path into the corresponding sandbox path.
func (s *Sandbox) MirrorPath(original string) string {
	rel, err := filepath.Rel(s.originalRoot, original)
	if err != nil {
		return original
	}
	return filepath.Join(s.root, rel)
}

// OriginalPath converts a sandbox file path back to the original path.
func (s *Sandbox) OriginalPath(sandboxPath string) string {
	rel, err := filepath.Rel(s.root, sandboxPath)
	if err != nil {
		return sandboxPath
	}
	return filepath.Join(s.originalRoot, rel)
}

// Cleanup removes the sandbox directory.
func (s *Sandbox) Cleanup() error {
	if s == nil {
		return nil
	}
	return os.RemoveAll(s.root)
}

func copyTree(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return nil
}
