//go:build !windows && !darwin

package profile

// Stubs for non-Windows, non-Darwin platforms (e.g. Linux).
func stashWindowsCred(profileName string) error  { return nil }
func swapWindowsCred(profileName string) error   { return nil }
func DeleteWindowsCred(profileName string) error  { return nil }
func RenameWindowsCred(oldName, newName string) error { return nil }
func GetAccountEmail(profileName, activeProfile string) string { return "" }
func HasWindowsCred(profileName, activeProfile string) bool { return false }
func GetAccessToken(profileName, activeProfile string) string { return "" }
