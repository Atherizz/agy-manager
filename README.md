# agym - Antigravity Identity Manager

`agym` is a CLI tool for managing multiple Antigravity CLI accounts on a single machine. It works by swapping per-profile credentials in and out of the location that `agy` reads from, so switching accounts is one command, with no re-login required after the initial setup.

---

## Requirements

- Go 1.22+
- Antigravity CLI (`agy`) installed
- Windows (primary support)

---

## Installation

Clone the repository and build the binary:

```bash
git clone https://github.com/Atherizz/agy-manager
cd agy-manager
go build -o agym.exe .
```

To run `agym` from anywhere without navigating to the project folder, copy the binary to a directory in your PATH:

```powershell
# Create a personal bin folder (if it doesn't exist)
New-Item -ItemType Directory -Path "$env:USERPROFILE\bin" -Force

# Copy the binary
Copy-Item ".\agym.exe" "$env:USERPROFILE\bin\agym.exe" -Force

# Add to User PATH (run once)
$currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
[Environment]::SetEnvironmentVariable("PATH", "$currentPath;$env:USERPROFILE\bin", "User")
```

Open a new terminal and verify:

```bash
agym --help
```

When you rebuild after a code change, re-run the copy step:

```bash
go build -o agym.exe .
Copy-Item ".\agym.exe" "$env:USERPROFILE\bin\agym.exe" -Force
```

---

## Commands

```
agym create <name>    Create a new profile
agym use <name>       Switch to a profile
agym list             List all profiles (marks the active one)
agym status           Show the currently active profile
agym delete <name>    Delete a profile
agym run <name> -- <cmd>   Run a command under a specific profile without switching
agym version          Print version info
```

### Examples

```bash
agym create user1
agym create user2

agym use user1
agym list

agym status

agym delete user2

agym run user1 -- agy
```

---

## How It Works

`agy` reads credentials from a fixed location: `~/.gemini/`. It has no built-in concept of multiple accounts.

`agym` works around this by keeping a vault per profile under `~/.gemini/profiles/<name>/`. When you switch profiles, it:

1. **Stashes** the current credentials (files + Windows Credential Manager entry) into the active profile's vault
2. **Loads** the target profile's credentials into `~/.gemini/`
3. **Updates** `state.json` to record which profile is now active

From `agy`'s perspective, it just sees valid credentials at the expected location, it never knows a swap happened.

### What gets isolated per profile

| Item | Location |
|---|---|
| OAuth token | `oauth_creds.json` |
| Linked accounts | `google_accounts.json` |
| Runtime state | `state.json` |
| Session cache | `antigravity-cli/implicit/` |
| Windows session | Credential Manager (`gemini:antigravity`) |

### What is shared across all profiles

- Plugins, skills, and built-in tools (`config/`, `builtin/`)
- Conversation history (`antigravity-cli/brain/`)
- Installation ID
- Settings and trusted folders

---

## Setting Up Multiple Accounts

### First-time setup

**1. Create your profiles**

```bash
agym create user1
agym create user2
```

**2. Clear existing credentials from Windows Credential Manager**

```bash
cmdkey /delete:gemini:antigravity
```

> If you see "CREDENTIAL_NOT_FOUND", that's fine, nothing was stored yet.

**3. Set the starting profile**

```powershell
'{"active_profile":"user1"}' | Set-Content "$env:USERPROFILE\.gemini\profiles\state.json"
```

**4. Log in with the first account**

```bash
agy
```

Log in with your first account (e.g. `user1@gmail.com`), then exit with `/exit`.

**5. Switch to the second profile and log in**

```bash
agym use user2
agy
```

Log in with your second account (e.g. `user2@gmail.com`), then exit with `/exit`.

That's it. The initial setup is done.

---

### Switching accounts (day-to-day)

```bash
agym use user1   # switch to first account
agym use user2   # switch to second account
```

No logout or re-login needed.

---

### Adding more accounts later

```bash
agym create user3
agym use user3
agy
# Log in with the new account, then /exit
```

---

## Troubleshooting

**Switching profile but `agy` still shows the same account**

Check what's stored in the Credential Manager:

```bash
cmdkey /list | findstr gemini
```

You should see `gemini:antigravity:<profile-name>` entries for each profile that has been set up. If not, the initial setup (clearing with `cmdkey /delete` and re-logging in) may need to be repeated.

**Resetting everything from scratch**

```powershell
cmdkey /delete:gemini:antigravity
Remove-Item "$env:USERPROFILE\.gemini\profiles" -Recurse -Force
```

Then start the first-time setup from the beginning.

---

## Acknowledgements

Inspired by [AntigravityManager](https://github.com/Draculabo/AntigravityManager) by Draculabo — a GUI-based account manager for Antigravity IDE. `agym` is a CLI-focused alternative built specifically for Antigravity CLI users.
