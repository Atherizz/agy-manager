package profile

import (
	"os"
	"path/filepath"
)

type PathResolver struct {
	geminiRoot string
}

func NewPathResolver(geminiRoot string) *PathResolver {
	return &PathResolver{geminiRoot: geminiRoot}
}

func DefaultGeminiRoot() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".gemini"), nil
}

func (p *PathResolver) GeminiRoot() string { return p.geminiRoot }

func (p *PathResolver) ProfilesDir() string { return filepath.Join(p.geminiRoot, "profiles") }

func (p *PathResolver) ProfileDir(name string) string {
	return filepath.Join(p.geminiRoot, "profiles", name)
}

func (p *PathResolver) StateFile() string {
	return filepath.Join(p.geminiRoot, "profiles", "state.json")
}

// IsolatedFiles returns per-identity authentication files managed by agy.
func (p *PathResolver) IsolatedFiles() []string {
	return []string{"oauth_creds.json", "google_accounts.json", "state.json"}
}

// IsolatedDirs returns session cache directories (.pb files) managed per account by agy.
func (p *PathResolver) IsolatedDirs() []string {
	return []string{filepath.Join("antigravity-cli", "implicit")}
}
