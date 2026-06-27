//go:build !windows

package profile

// Stub untuk non-Windows (Linux/macOS).
func stashWindowsCred(profileName string) error { return nil }
func swapWindowsCred(profileName string) error  { return nil }
