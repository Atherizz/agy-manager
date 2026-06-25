package profile

import(
	"testing"
	"path/filepath"
)

func TestProfilePath(t *testing.T) {

	root := t.TempDir()
	p := NewPathResolver(root)

	want := filepath.Join(root, "profiles")
    if got := p.ProfilesDir(); got != want {

        t.Errorf("ProfilesDir() = %q, want %q", got, want)

    }

}

func TestProfileDir(t *testing.T) {

    root := t.TempDir()

    p := NewPathResolver(root)

  

    want := filepath.Join(root, "profiles", "team-alpha")

    got := p.ProfileDir("team-alpha")

  

    if got != want {

        t.Errorf("ProfileDir() = %q, want %q", got, want)

    }

}

  

func TestIsolatedFiles(t *testing.T) {

    p := NewPathResolver(t.TempDir())

    files := p.IsolatedFiles()

    if len(files) == 0 {

        t.Error("IsolatedFiles() returned empty slice")

    }

    // Pastikan oauth_creds.json ada di daftar

    found := false

    for _, f := range files {
        if f == "oauth_creds.json" {
			found = true
            break
        }

    }

    if !found {
        t.Error("IsolatedFiles() should contain oauth_creds.json")

    }

}

  

func TestIsolatedDirs(t *testing.T) {

    p := NewPathResolver(t.TempDir())

    dirs := p.IsolatedDirs()


    if len(dirs) != 0 {

        t.Error("IsolatedDirs() not returned empty slice")

    }

}