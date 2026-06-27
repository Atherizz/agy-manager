package profile

import (
	"os"
	"path/filepath"
)

func Exists(p *PathResolver, name string) bool {
	info, err := os.Stat(p.ProfileDir(name))
	return err == nil && info.IsDir()
}

func List(p *PathResolver) ([]string, error) {
	entries, err := os.ReadDir(p.ProfilesDir())
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	var profiles []string
	for _, e := range entries {
		if e.IsDir() {
			profiles = append(profiles, e.Name())
		}
	}
	return profiles, nil
}

func Create(p *PathResolver, name string) error {
	profileDir := p.ProfileDir(name)
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		return err
	}
	for _, d := range p.IsolatedDirs() {
		if err := os.MkdirAll(filepath.Join(profileDir, d), 0755); err != nil {
			return err
		}
	}
	return nil
}