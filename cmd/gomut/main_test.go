package main

import (
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestNormalizePkgArg(t *testing.T) {
	tempDir := t.TempDir()
	root := filepath.Join(tempDir, "project")
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatalf("failed to create root directory: %v", err)
	}
	other := filepath.Join(tempDir, "other")
	if err := os.MkdirAll(other, 0o755); err != nil {
		t.Fatalf("failed to create other directory: %v", err)
	}

	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalCwd)
	})

	tests := []struct {
		name  string
		pkg   string
		root  string
		chdir string
		want  string
	}{
		{
			name: "empty defaults to all",
			pkg:  "",
			root: root,
			want: "./...",
		},
		{
			name: "absolute path matching root",
			pkg:  root,
			root: root,
			want: "./...",
		},
		{
			name: "absolute path different from root",
			pkg:  other,
			root: root,
			want: other,
		},
		{
			name:  "relative path matching root",
			pkg:   ".",
			root:  root,
			chdir: root,
			want:  "./...",
		},
		{
			name:  "relative path different from root",
			pkg:   "./other",
			root:  root,
			chdir: tempDir,
			want:  "./other",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := os.Chdir(originalCwd); err != nil {
				t.Fatalf("failed to reset working directory: %v", err)
			}
			if tt.chdir != "" {
				if err := os.Chdir(tt.chdir); err != nil {
					t.Fatalf("failed to change working directory: %v", err)
				}
			}

			got := normalizePkgArg(tt.pkg, tt.root)
			if got != tt.want {
				t.Fatalf("normalizePkgArg(%q, %q) = %q, want %q", tt.pkg, tt.root, got, tt.want)
			}
		})
	}
}

func TestSamePath(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{
			name: "identical paths",
			a:    filepath.Join("tmp", "project"),
			b:    filepath.Join("tmp", "project"),
			want: true,
		},
		{
			name: "equivalent after cleaning",
			a:    filepath.Join("tmp", "project", "..", "project", "./module"),
			b:    filepath.Join("tmp", "project", "module"),
			want: true,
		},
		{
			name: "different paths",
			a:    filepath.Join("tmp", "project"),
			b:    filepath.Join("tmp", "another"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := samePath(tt.a, tt.b)
			if got != tt.want {
				t.Fatalf("samePath(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestDisplayPath(t *testing.T) {
	root := filepath.Join(t.TempDir(), "root")
	path := filepath.Join(root, "nested", "file.go")

	got := displayPath(root, path)
	want := filepath.Join("nested", "file.go")
	if got != want {
		t.Fatalf("displayPath(%q, %q) = %q, want %q", root, path, got, want)
	}
}

func TestPercent(t *testing.T) {
	tests := []struct {
		name   string
		killed int
		total  int
		want   float64
	}{
		{
			name: "zero total",
			want: 0,
		},
		{
			name:   "partial",
			killed: 2,
			total:  4,
			want:   50,
		},
		{
			name:   "fractional",
			killed: 1,
			total:  3,
			want:   33.3333333333,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := percent(tt.killed, tt.total)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Fatalf("percent(%d, %d) = %f, want %f", tt.killed, tt.total, got, tt.want)
			}
		})
	}
}
