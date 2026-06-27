package profile

import (
	"io"
	"os"
	"path/filepath"
)

func StashCurrentProfile(p *PathResolver, name string) error {
	for _, f := range p.IsolatedFiles() {
		src := filepath.Join(p.GeminiRoot(), f)
		if err := copyIfExists(src, filepath.Join(p.ProfileDir(name), f)); err != nil {
			return err
		}
		os.Remove(src)
	}
	for _, d := range p.IsolatedDirs() {
		if err := moveIfExists(filepath.Join(p.GeminiRoot(), d), filepath.Join(p.ProfileDir(name), d)); err != nil {
			return err
		}
	}
	return stashWindowsCred(name)
}

func SwapToProfile(p *PathResolver, name string) error {
	for _, f := range p.IsolatedFiles() {
		os.Remove(filepath.Join(p.GeminiRoot(), f))
	}
	for _, f := range p.IsolatedFiles() {
		if err := copyIfExists(filepath.Join(p.ProfileDir(name), f), filepath.Join(p.GeminiRoot(), f)); err != nil {
			return err
		}
	}
	for _, d := range p.IsolatedDirs() {
		if err := moveIfExists(filepath.Join(p.ProfileDir(name), d), filepath.Join(p.GeminiRoot(), d)); err != nil {
			return err
		}
	}
	return swapWindowsCred(name)
}

func copyIfExists(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer srcFile.Close()
	os.MkdirAll(filepath.Dir(dst), 0755)
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	return err
}

func moveIfExists(src, dst string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return nil
	}
	os.RemoveAll(dst)
	os.MkdirAll(filepath.Dir(dst), 0755)
	return os.Rename(src, dst)
}
