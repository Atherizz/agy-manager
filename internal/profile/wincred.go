//go:build windows

package profile

import "github.com/danieljoos/wincred"

const credTarget = "gemini:antigravity"

// stashWindowsCred saves the active credential from Windows Credential Manager
// into a profile-specific slot, then deletes the main entry.
func stashWindowsCred(profileName string) error {
	cred, err := wincred.GetGenericCredential(credTarget)
	if err != nil {
		return nil // no credential present, nothing to stash
	}
	profileCred := wincred.NewGenericCredential(credTarget + ":" + profileName)
	profileCred.CredentialBlob = cred.CredentialBlob
	profileCred.UserName = cred.UserName
	if err := profileCred.Write(); err != nil {
		return err
	}
	return cred.Delete()
}

// swapWindowsCred loads the credential for profileName from its profile-specific slot
// into the main Windows Credential Manager entry.
func swapWindowsCred(profileName string) error {
	if existing, err := wincred.GetGenericCredential(credTarget); err == nil {
		existing.Delete()
	}
	profileCred, err := wincred.GetGenericCredential(credTarget + ":" + profileName)
	if err != nil {
		return nil // new profile, no credential yet (will prompt login)
	}
	newCred := wincred.NewGenericCredential(credTarget)
	newCred.CredentialBlob = profileCred.CredentialBlob
	newCred.UserName = profileCred.UserName
	return newCred.Write()
}
