package profile

import (

    "os"

    "path/filepath"

    "testing"

)

func TestSaveAndLoadState(t *testing.T) {

    root := t.TempDir()

    os.MkdirAll(filepath.Join(root, "profiles"), 0755)

    p := NewPathResolver(root)

    state := &State{ActiveProfile: "team-alpha"}

    if err := SaveState(p, state); err != nil {

        t.Fatalf("SaveState() error = %v", err)

    }

    loaded, err := LoadState(p)

    if err != nil {

        t.Fatalf("LoadState() error = %v", err)

    }

    if loaded.ActiveProfile != "team-alpha" {

        t.Errorf("ActiveProfile = %q, want %q", loaded.ActiveProfile, "team-alpha")

    }

}

func TestLoadState_NoFile(t *testing.T) {

    root := t.TempDir()

    os.MkdirAll(filepath.Join(root, "profiles"), 0755)

    p := NewPathResolver(root)

    state, err := LoadState(p)

    if err != nil {

        t.Fatalf("LoadState() should not error on missing file, got %v", err)

    }

    if state.ActiveProfile != "" {

        t.Errorf("ActiveProfile should be empty, got %q", state.ActiveProfile)

    }

}