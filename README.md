
# agym — Antigravity Identity Manager

  

Manage multiple Antigravity CLI profiles on one machine.

  

## Installation

  

```bash

go build -o agym.exe .

# Move agym.exe to a directory in your PATH

```

  

## Usage

  

```bash

# Create profiles

agym create personal

agym create work

  

# Switch profile

agym use work

  

# Check active profile

agym status

  

# List all profiles

agym list

  

# Run a one-off command with a specific profile

agym run personal -- agy

  

# Delete a profile

agym delete personal

```

  

## How It Works

  

Each profile stores its own isolated copy of:

- `oauth_creds.json` — authentication credentials

- `google_accounts.json` — linked accounts

- `projects.json` — GCP project settings

- `settings.json` — user preferences

- `state.json` — runtime state

- `antigravity-cli/brain/` — conversation history

- `history/` — command history

  

Global components are shared across all profiles:

- `installation_id` — machine identifier

- `trustedFolders.json` — trusted workspace folders

- `config/plugins/` — installed plugins

- `config/skills/` — installed skills

- `builtin/` — built-in capabilities

  

## Architecture

  

```

~/.gemini/

├── profiles/

│   ├── state.json          ← tracks active profile

│   ├── personal/

│   │   ├── oauth_creds.json

│   │   ├── antigravity-cli/brain/

│   │   └── ...

│   └── work/

│       ├── oauth_creds.json

│       ├── antigravity-cli/brain/

│       └── ...

├── installation_id         ← shared

├── config/                 ← shared

└── oauth_creds.json        ← active profile's copy

```
