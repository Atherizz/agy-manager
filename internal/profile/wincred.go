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

// DeleteWindowsCred completely removes a profile's stored credential from the Windows Credential Manager.
func DeleteWindowsCred(profileName string) error {
	if cred, err := wincred.GetGenericCredential(credTarget + ":" + profileName); err == nil {
		return cred.Delete()
	}
	return nil
}

// RenameWindowsCred copies a profile's credential to a new slot and deletes the old one.
func RenameWindowsCred(oldName, newName string) error {
	cred, err := wincred.GetGenericCredential(credTarget + ":" + oldName)
	if err != nil {
		return nil // nothing stored for this profile
	}
	newCred := wincred.NewGenericCredential(credTarget + ":" + newName)
	newCred.CredentialBlob = cred.CredentialBlob
	newCred.UserName = cred.UserName
	if err := newCred.Write(); err != nil {
		return err
	}
	return cred.Delete()
}

// HasWindowsCred returns true if the profile has a stored credential in Windows Credential Manager.
func HasWindowsCred(profileName, activeProfile string) bool {
	target := credTarget + ":" + profileName
	if profileName == activeProfile {
		target = credTarget
	}
	_, err := wincred.GetGenericCredential(target)
	return err == nil
}

// GetAccountEmail returns the Google account email associated with a profile.
func GetAccountEmail(profileName, activeProfile string) string {
	target := credTarget + ":" + profileName
	if profileName == activeProfile {
		target = credTarget
	}
	cred, err := wincred.GetGenericCredential(target)
	if err != nil {
		return ""
	}
	return cred.UserName
}
