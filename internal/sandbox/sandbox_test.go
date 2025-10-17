package sandbox

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewCopiesTree(t *testing.T) {
	root := t.TempDir()

	writeFile(t, filepath.Join(root, "file.txt"), "hello")
	writeFile(t, filepath.Join(root, "dir", "nested.txt"), "nested")

	sb, err := New(root)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	t.Cleanup(func() {
		if cerr := sb.Cleanup(); cerr != nil {
			t.Fatalf("cleanup failed: %v", cerr)
		}
	})

	if sb.Root() == root {
		t.Fatalf("expected sandbox root to differ from original root")
	}

	sandboxFile := filepath.Join(sb.Root(), "file.txt")
	if data := readFile(t, sandboxFile); string(data) != "hello" {
		t.Fatalf("unexpected sandbox file contents: %q", string(data))
	}

	sandboxNested := filepath.Join(sb.Root(), "dir", "nested.txt")
	if data := readFile(t, sandboxNested); string(data) != "nested" {
		t.Fatalf("unexpected sandbox nested file contents: %q", string(data))
	}
}

func TestNewRequiresDirectory(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "notadir")
	writeFile(t, filePath, "content")

	if _, err := New(filePath); err == nil {
		t.Fatal("expected error when creating sandbox from file")
	}
}

func TestMirrorAndOriginalPath(t *testing.T) {
	root := t.TempDir()
	original := filepath.Join(root, "sub", "file.txt")
	writeFile(t, original, "value")

	sb, err := New(root)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	t.Cleanup(func() {
		if cerr := sb.Cleanup(); cerr != nil {
			t.Fatalf("cleanup failed: %v", cerr)
		}
	})

	mirror := sb.MirrorPath(original)
	expectedMirror := filepath.Join(sb.Root(), "sub", "file.txt")
	if mirror != expectedMirror {
		t.Fatalf("MirrorPath = %q, want %q", mirror, expectedMirror)
	}

	roundTrip := sb.OriginalPath(mirror)
	if roundTrip != original {
		t.Fatalf("OriginalPath = %q, want %q", roundTrip, original)
	}
}

func TestCleanupRemovesSandbox(t *testing.T) {
	root := t.TempDir()
	sb, err := New(root)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	sandboxRoot := sb.Root()

	if err := sb.Cleanup(); err != nil {
		t.Fatalf("Cleanup returned error: %v", err)
	}

	if _, err := os.Stat(sandboxRoot); !os.IsNotExist(err) {
		t.Fatalf("expected sandbox directory to be removed, stat err = %v", err)
	}
}

func TestCleanupNilSandbox(t *testing.T) {
	var sb *Sandbox
	if err := sb.Cleanup(); err != nil {
		t.Fatalf("Cleanup on nil sandbox returned error: %v", err)
	}
}

func writeFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed creating directories for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatalf("failed writing file %s: %v", path, err)
	}
}

func readFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed reading file %s: %v", path, err)
	}
	return data
}
