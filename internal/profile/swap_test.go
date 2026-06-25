package profile

import (

    "os"

    "path/filepath"

    "testing"

)

func TestStashAndSwap(t *testing.T) {

    p, root := setupTestRoot(t)

    Create(p, "alpha")

    Create(p, "beta")

    os.WriteFile(filepath.Join(root, "oauth_creds.json"), []byte("alpha-token"), 0644)

    os.WriteFile(filepath.Join(p.ProfileDir("beta"), "oauth_creds.json"), []byte("beta-token"), 0644)

    if err := StashCurrentProfile(p, "alpha"); err != nil { t.Fatal(err) }

    if err := SwapToProfile(p, "beta"); err != nil { t.Fatal(err) }

    data, _ := os.ReadFile(filepath.Join(root, "oauth_creds.json"))

    if string(data) != "beta-token" { t.Errorf("root = %q, want beta-token", string(data)) }

}
