package profile

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestRoot(t *testing.T) (*PathResolver, string) {
	t.Helper()
	root := t.TempDir()
	return NewPathResolver(root), root
}

func TestCreateAndExists(t *testing.T) {
	p, _ := setupTestRoot(t)
	if err := Create(p, "alpha"); err != nil {
		t.Fatalf("Create error = %v", err)
	}
	if !Exists(p, "alpha") {
		t.Error("profile should exist")
	}
	for _, d := range p.IsolatedDirs() {
		if info, err := os.Stat(filepath.Join(p.ProfileDir("alpha"), d)); err != nil || !info.IsDir() {
			t.Errorf("isolated dir %q missing", d)
		}
	}
}

func TestListProfiles(t *testing.T) {
	p, _ := setupTestRoot(t)
	Create(p, "alpha")
	Create(p, "beta")
	profiles, _ := List(p)
	if len(profiles) != 2 {
		t.Errorf("List() returned %d profiles, want 2", len(profiles))
	}
}