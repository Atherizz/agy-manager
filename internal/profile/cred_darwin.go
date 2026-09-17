//go:build darwin

package profile

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var kcServices = []string{"gemini", "gemini:antigravity", "Antigravity Safe Storage"}

const kcAccount = "antigravity"

type tokenClaims struct {
	Email string `json:"email"`
}

type oauthTokenPayload struct {
	Token struct {
		AccessToken string `json:"access_token"`
		Email       string `json:"email"`
	} `json:"token"`
	AccessToken string `json:"access_token"`
	Email       string `json:"email"`
	IDToken     string `json:"id_token"`
}

func kcRead() (string, string, error) {
	for _, svc := range kcServices {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		cmd := exec.CommandContext(ctx, "/usr/bin/security", "find-generic-password", "-s", svc, "-a", kcAccount, "-w")
		out, err := cmd.Output()
		cancel()
		if err == nil && len(bytes.TrimSpace(out)) > 0 {
			return string(bytes.TrimSpace(out)), svc, nil
		}
	}
	return "", "", fmt.Errorf("no credential found in keychain")
}

func kcWrite(token, service string) error {
	if service == "" {
		service = "gemini"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "/usr/bin/security", "add-generic-password", "-U", "-s", service, "-a", kcAccount, "-w", token)
	return cmd.Run()
}

func kcDelete() error {
	for _, svc := range kcServices {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		cmd := exec.CommandContext(ctx, "/usr/bin/security", "delete-generic-password", "-s", svc, "-a", kcAccount)
		_ = cmd.Run()
		cancel()
	}
	return nil
}

func decodeTokenData(raw string) ([]byte, error) {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "go-keyring-base64:") {
		b64 := strings.TrimPrefix(raw, "go-keyring-base64:")
		if rem := len(b64) % 4; rem != 0 {
			b64 += strings.Repeat("=", 4-rem)
		}
		return base64.StdEncoding.DecodeString(b64)
	}
	return []byte(raw), nil
}

func parseAccessToken(raw string) string {
	decoded, err := decodeTokenData(raw)
	if err != nil {
		return ""
	}
	var data oauthTokenPayload
	if err := json.Unmarshal(decoded, &data); err != nil {
		return ""
	}
	if data.Token.AccessToken != "" {
		return data.Token.AccessToken
	}
	return data.AccessToken
}

func parseAccountEmail(raw string) string {
	decoded, err := decodeTokenData(raw)
	if err != nil {
		return ""
	}
	var data oauthTokenPayload
	if err := json.Unmarshal(decoded, &data); err != nil {
		return ""
	}
	if data.Token.Email != "" {
		return data.Token.Email
	}
	if data.Email != "" {
		return data.Email
	}
	if data.IDToken != "" {
		parts := strings.Split(data.IDToken, ".")
		if len(parts) >= 2 {
			payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
			if err != nil {
				payloadStr := parts[1]
				if rem := len(payloadStr) % 4; rem != 0 {
					payloadStr += strings.Repeat("=", 4-rem)
				}
				payloadBytes, _ = base64.URLEncoding.DecodeString(payloadStr)
			}
			if len(payloadBytes) > 0 {
				var claims tokenClaims
				if err := json.Unmarshal(payloadBytes, &claims); err == nil && claims.Email != "" {
					return claims.Email
				}
			}
		}
	}
	return ""
}

// readRawProfileToken retrieves the raw token string from the Keychain (if active) or from the profile vault.
func readRawProfileToken(profileName, activeProfile string) (string, error) {
	if profileName == activeProfile {
		if token, _, err := kcRead(); err == nil && token != "" {
			return token, nil
		}
	}

	root, err := DefaultGeminiRoot()
	if err != nil {
		return "", err
	}
	p := NewPathResolver(root)

	// Check vault oauth_creds.json
	vaultFile := filepath.Join(p.ProfileDir(profileName), "oauth_creds.json")
	if data, err := os.ReadFile(vaultFile); err == nil && len(bytes.TrimSpace(data)) > 0 {
		return strings.TrimSpace(string(data)), nil
	}

	// Check vault antigravity-cli/antigravity-oauth-token
	vaultCliFile := filepath.Join(p.ProfileDir(profileName), "antigravity-cli", "antigravity-oauth-token")
	if data, err := os.ReadFile(vaultCliFile); err == nil && len(bytes.TrimSpace(data)) > 0 {
		return strings.TrimSpace(string(data)), nil
	}

	// If active profile, check root antigravity-cli/antigravity-oauth-token and oauth_creds.json
	if profileName == activeProfile {
		rootCliFile := filepath.Join(p.GeminiRoot(), "antigravity-cli", "antigravity-oauth-token")
		if data, err := os.ReadFile(rootCliFile); err == nil && len(bytes.TrimSpace(data)) > 0 {
			return strings.TrimSpace(string(data)), nil
		}
		rootOAuthFile := filepath.Join(p.GeminiRoot(), "oauth_creds.json")
		if data, err := os.ReadFile(rootOAuthFile); err == nil && len(bytes.TrimSpace(data)) > 0 {
			return strings.TrimSpace(string(data)), nil
		}
	}

	return "", fmt.Errorf("no credential token found for profile %q", profileName)
}

func stashWindowsCred(profileName string) error {
	token, _, err := kcRead()
	if err != nil || token == "" {
		return nil
	}
	root, err := DefaultGeminiRoot()
	if err != nil {
		return err
	}
	p := NewPathResolver(root)
	profileDir := p.ProfileDir(profileName)
	if err := os.MkdirAll(profileDir, 0700); err != nil {
		return err
	}
	dest := filepath.Join(profileDir, "oauth_creds.json")
	if err := os.WriteFile(dest, []byte(token), 0600); err != nil {
		return err
	}
	return kcDelete()
}

func swapWindowsCred(profileName string) error {
	_ = kcDelete()
	root, err := DefaultGeminiRoot()
	if err != nil {
		return err
	}
	p := NewPathResolver(root)
	dest := filepath.Join(p.ProfileDir(profileName), "oauth_creds.json")
	data, err := os.ReadFile(dest)
	if err != nil {
		destCli := filepath.Join(p.ProfileDir(profileName), "antigravity-cli", "antigravity-oauth-token")
		data, err = os.ReadFile(destCli)
		if err != nil {
			return nil // New profile with no stored credentials yet
		}
	}
	token := strings.TrimSpace(string(data))
	if token == "" {
		return nil
	}
	return kcWrite(token, "gemini")
}

func DeleteWindowsCred(profileName string) error {
	root, err := DefaultGeminiRoot()
	if err != nil {
		return err
	}
	p := NewPathResolver(root)
	state, _ := LoadState(p)
	if state != nil && state.ActiveProfile == profileName {
		_ = kcDelete()
	}
	return nil
}

func RenameWindowsCred(oldName, newName string) error {
	return nil
}

func HasWindowsCred(profileName, activeProfile string) bool {
	raw, err := readRawProfileToken(profileName, activeProfile)
	return err == nil && raw != ""
}

func GetAccountEmail(profileName, activeProfile string) string {
	raw, err := readRawProfileToken(profileName, activeProfile)
	if err != nil {
		return ""
	}
	return parseAccountEmail(raw)
}

func GetAccessToken(profileName, activeProfile string) string {
	raw, err := readRawProfileToken(profileName, activeProfile)
	if err != nil {
		return ""
	}
	return parseAccessToken(raw)
}
