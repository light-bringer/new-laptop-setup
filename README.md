# Engineering Laptop Setup

Ansible-based automation for configuring macOS developer environments. This tool provides an idempotent setup that installs development tools, configures SSH access to GitHub and GitLab, and sets up a productive shell environment.

## Features

- 🍺 **Homebrew** installation and configuration
- 🐚 **Oh My Zsh** with curated plugins
- ⭐ **Starship** - Fast, customizable cross-shell prompt (default theme)
- 📝 **Dotfiles** management (.vimrc, .tmux.conf, .gitignore_global, .editorconfig)
- 🔑 **SSH key** generation and management
- 🐙 **GitHub** SSH authentication
- 🦊 **GitLab** SSH authentication
- ⚙️ **Git** user configuration
- 🚀 **goto** - Custom shell function for quick Git repository navigation (GitLab default, GitHub optional)
- 🛠️ **CLI Development tools** (50+ tools including git, jq, curl, kubectl, terraform, awscli)
- 🐹 **Go** programming language with GOPATH configuration
- 📦 **NVM** (Node Version Manager) with latest LTS Node.js
- ⚡ **pnpm** - Fast, disk space efficient package manager
- 🐍 **Python 3.12** with pip, pipx, and virtualenv
- 🖥️ **GUI Applications** - Docker Desktop (default), VSCode, iTerm2, Chrome, Slack, Postman (optional)
- 🤖 **AI CLI Tools** - Claude Code CLI, GitHub Copilot CLI (optional)
- 🔄 **Update scripts** for easy maintenance (laptop.update, laptop.upgrade)

## Prerequisites

Before running the setup, ensure you have the following:

### System Requirements

- **macOS**: macOS 12 (Monterey) or later
- **Architecture**: Apple Silicon (M1/M2/M3) or Intel
- **Admin Access**: Sudo privileges required for installation
- **Internet Connection**: Required for downloading packages and dependencies
- **Disk Space**: At least 5GB free for development tools and applications

### Required Credentials

You'll need to create these before starting the setup:

#### 1. GitHub Personal Access Token (Required)
- **Create at**: https://github.com/settings/tokens
- **Required scopes**:
  - `admin:public_key` - Upload SSH keys
  - `read:user` - Read user profile
  - `repo` - Repository access
- **Copy and save** the token - you'll need it during setup

#### 2. GitLab Personal Access Token (Required)
- **Create at**: https://gitlab.com/-/profile/personal_access_tokens
- **Required scopes**:
  - `api` - Full API access
  - `read_user` - Read user profile
  - `write_repository` - Write to repositories
- **Copy and save** the token - you'll need it during setup

### Optional Credentials (for optional features)

#### GitHub Copilot (Optional)
- **Subscription**: Active GitHub Copilot subscription
- **Authentication**: Handled via `gh auth login` after installation

#### Claude Code CLI (Optional)
- **API Key**: Anthropic API key
- **Create at**: https://console.anthropic.com/

### Information to Prepare

During setup, you'll be prompted for:

1. **Git Configuration**:
   - Your full name (for git commits)
   - Your email address (for git commits)

2. **GitHub Information**:
   - GitHub username
   - GitHub email
   - GitHub Personal Access Token (created above)

3. **GitLab Information**:
   - GitLab username
   - GitLab email
   - GitLab Personal Access Token (created above)

4. **goto Function Setup**:
   - Git platform choice: `gitlab` (default) or `github`
   - Default organization/user: Defaults to `vercara`
   - Root directory: Defaults to `~/dev/src`

5. **Oh My Zsh Theme** (choose one):
   - `robbyrussell` (default, minimal)
   - `agnoster` (powerline-style, requires Nerd Font)
   - `powerlevel10k` (highly customizable, requires Nerd Font)
   - `pure` (minimal, async)

### Pre-Setup Checklist

✅ macOS 12+ installed
✅ Admin/sudo access available
✅ GitHub Personal Access Token created and saved
✅ GitLab Personal Access Token created and saved
✅ GitHub username and email ready
✅ GitLab username and email ready
✅ Git name and email decided
✅ goto platform choice decided (gitlab/github)
✅ goto default organization decided (defaults to vercara)
✅ Oh My Zsh theme preference decided

### How to Create Required Credentials

#### Creating a GitHub Personal Access Token

1. **Sign in** to GitHub: https://github.com
2. Click your **profile photo** (top right) → **Settings**
3. Scroll down and click **Developer settings** (left sidebar, bottom)
4. Click **Personal access tokens** → **Tokens (classic)**
5. Click **Generate new token** → **Generate new token (classic)**
6. Fill out the form:
   - **Note**: `Laptop Setup - SSH Key Upload` (or any descriptive name)
   - **Expiration**: Choose duration (recommend 90 days or No expiration)
   - **Select scopes**:
     - ✅ `repo` (Full control of private repositories)
     - ✅ `admin:public_key` (Full control of user public keys)
     - ✅ `read:user` (Read ALL user profile data)
7. Click **Generate token**
8. **Copy the token** immediately - you won't be able to see it again!
9. Save it securely (password manager or temporary note)

#### Creating a GitLab Personal Access Token

1. **Sign in** to GitLab: https://gitlab.com
2. Click your **avatar** (top right) → **Edit profile**
3. In the left sidebar, click **Access Tokens**
4. Fill out the form:
   - **Token name**: `Laptop Setup - SSH Key Upload` (or any descriptive name)
   - **Expiration date**: Choose a date (optional but recommended for security)
   - **Select scopes**:
     - ✅ `api` (Grants complete read/write access to the API)
     - ✅ `read_user` (Grants read-only access to your profile)
     - ✅ `write_repository` (Allows read-write access to repositories)
5. Click **Create personal access token**
6. **Copy the token** immediately - you won't be able to see it again!
7. Save it securely (password manager or temporary note)

#### Getting Your GitHub Username and Email

Your GitHub username is visible in your profile URL:
- URL format: `https://github.com/YOUR-USERNAME`
- Or go to https://github.com/settings/profile to see your username

Your GitHub email (for commits):
- Same page: https://github.com/settings/emails
- Use your primary email or a GitHub-provided noreply email

#### Getting Your GitLab Username and Email

Your GitLab username is visible in your profile URL:
- URL format: `https://gitlab.com/YOUR-USERNAME`
- Or go to https://gitlab.com/-/profile to see your username

Your GitLab email (for commits):
- Same page: https://gitlab.com/-/profile
- Use your primary email

#### Optional: Getting Anthropic API Key (for Claude Code CLI)

1. Go to: https://console.anthropic.com/
2. Sign in or create an account
3. Go to **API Keys** section
4. Click **Create Key**
5. Give it a name (e.g., "Claude Code CLI")
6. Copy and save the API key securely

#### Optional: GitHub Copilot Subscription

1. Go to: https://github.com/features/copilot
2. Click **Start free trial** or **Buy now**
3. Follow the subscription process
4. Ensure your GitHub account has active Copilot access

## Quick Start

### Initial Setup

```bash
cd ~/dev/new-laptop-setup
./bootstrap.sh
```

This will:
1. Install Homebrew (if not present)
2. Install Ansible
3. Add `bin/` directory to PATH for update scripts
4. Run the full laptop setup automation

### Updating Your Setup

After the initial bootstrap, use the convenient update scripts:

```bash
# Pull latest changes from repository
laptop.update

# Update repository and run full setup
laptop.upgrade
```

## Optional Components

Some components are not installed by default and must be explicitly requested:

```bash
# Install Claude Code CLI
./bin/laptop.run --tags claude-code

# Install GitHub Copilot CLI
./bin/laptop.run --tags github-copilot

# Install optional GUI applications (VSCode, Chrome, Slack, etc.)
./bin/laptop.run --tags gui-optional

# Install specific apps
./bin/laptop.run --tags vscode,chrome,slack
```

**What's installed by default**: Docker Desktop
**What's optional**: VSCode, iTerm2, Chrome, Firefox, Slack, Postman, Notion, TablePlus, Figma, Insomnia, Rectangle

## Selective Execution

Run only specific components using Ansible tags:

```bash
# SSH key generation only
./bin/laptop.run --tags ssh

# Oh My Zsh setup only
./bin/laptop.run --tags ohmyzsh

# Starship prompt only
./bin/laptop.run --tags starship

# GitHub configuration only
./bin/laptop.run --tags github

# GitLab configuration only
./bin/laptop.run --tags gitlab

# Dotfiles symlinking only
./bin/laptop.run --tags dotfiles

# Go installation only
./bin/laptop.run --tags golang

# NVM and Node.js installation only
./bin/laptop.run --tags nvm

# pnpm installation only
./bin/laptop.run --tags pnpm

# Python installation only
./bin/laptop.run --tags python

# GUI applications only
./bin/laptop.run --tags applications

# Specific apps
./bin/laptop.run --tags docker,vscode,iterm2

# Multiple tags
./bin/laptop.run --tags ssh,git,github

# Skip specific components
./bin/laptop.run --skip-tags ohmyzsh
```

### Available Tags

- `cli-tools` - macOS Command Line Tools
- `homebrew` - Homebrew installation
- `ssh` - SSH key generation
- `zsh` - Zsh configuration
- `ohmyzsh` - Oh My Zsh installation
- `starship` / `prompt` - Starship prompt
- `dotfiles` - Dotfiles symlinking
- `git` - Git user configuration
- `goto` - goto shell function for repository navigation
- `github` - GitHub SSH setup
- `gitlab` - GitLab SSH setup
- `dev-tools` - Development tools
- `golang` / `go` - Go programming language
- `nvm` / `node` / `nodejs` - NVM and Node.js
- `pnpm` - pnpm package manager
- `python` / `py` - Python with pip, pipx, and virtualenv
- `applications` / `apps` / `gui` - GUI applications
- Individual app tags: `docker`, `vscode`, `iterm2`, `chrome`, `firefox`, `slack`, `postman`, `notion`, `tableplus`, `figma`, `insomnia`, `rectangle`

**Optional (not installed by default)**:
- `claude-code` / `claude` - Claude Code CLI
- `github-copilot` / `copilot` - GitHub Copilot CLI
- `gui-optional` - All optional GUI apps (VSCode, Chrome, Slack, etc.)
- Individual apps: `vscode`, `iterm2`, `chrome`, `firefox`, `slack`, `postman`, `notion`, `tableplus`, `figma`, `insomnia`, `rectangle`

**Installed by default**:
- `docker` - Docker Desktop only

## Preview Changes

Run in check mode to see what would change without applying:

```bash
./bootstrap.sh --check
laptop.upgrade --check
```

## GitHub Setup

The automation will prompt you for a GitHub Personal Access Token (PAT):

1. Create a token at: https://github.com/settings/tokens
2. Required scopes: `admin:public_key`, `read:user`, `repo`
3. Enter the token when prompted
4. Your SSH key will be automatically uploaded

The setup configures Git to automatically use SSH for GitHub URLs, so you can use HTTPS clone URLs and they'll be converted to SSH.

## GitLab Setup

Similar to GitHub, you'll need a GitLab Personal Access Token:

1. Create a token at: https://gitlab.com/-/profile/personal_access_tokens
2. Required scopes: `api`, `read_user`, `write_repository`
3. Enter the token when prompted
4. Your SSH key will be automatically uploaded

## goto - Repository Navigation

The `goto` shell function provides quick navigation and cloning for Git repositories:

**Configuration during setup:**
- **Platform**: Choose GitLab (default) or GitHub
- **Default Org/User**: Set to `vercara` by default (customizable during setup)
- **Root Directory**: `~/dev/src` (customizable during setup)

**Usage:**

```bash
# Navigate to or clone a repository under the default org
goto my-repo              # → ~/dev/src/gitlab.com/vercara/my-repo

# Navigate to or clone from a specific org/user
goto other-org/project    # → ~/dev/src/gitlab.com/other-org/project

# Show current configuration
goto
```

**How it works:**
- If the directory exists locally, it navigates to it
- If not found, it automatically clones the repository using SSH
- Directory structure: `$GOTO_ROOT/$PLATFORM/$ORG/$REPO`

**Example workflow:**
```bash
goto backend-api          # Clones vercara/backend-api if needed, then cd's into it
goto teamname/frontend    # Clones teamname/frontend if needed, then cd's into it
```

## Dotfiles

The following dotfiles are symlinked from this repository to your home directory:

- `.gitignore_global` - Global gitignore patterns for macOS, IDEs, and common build artifacts
- `.vimrc` - Vim configuration with sensible defaults
- `.tmux.conf` - Tmux configuration with improved key bindings
- `.editorconfig` - EditorConfig settings for consistent coding styles

You can customize these files in the `dotfiles/` directory, and changes will apply immediately (since they're symlinked).

## Architecture

### Execution Flow

1. **bootstrap.sh** - Entry point that installs Homebrew, adds bin/ to PATH, then delegates to `bin/laptop.run`
2. **bin/laptop.run** - Installs Ansible via Homebrew, configures Python environment, runs `ansible-playbook main.yml`
3. **main.yml** - Main Ansible playbook that orchestrates all roles and tasks

### Directory Structure

```
new-laptop-setup/
├── bootstrap.sh              # Entry point
├── main.yml                  # Main Ansible playbook
├── ansible.cfg               # Ansible configuration
├── hosts                     # Inventory file
├── requirements.yml          # Ansible dependencies
├── bin/
│   ├── laptop.run           # Core execution script
│   └── refresh-ansible-collections.sh  # Vendor dependencies
├── tasks/
│   ├── ssh-key.yml          # SSH key generation
│   ├── zsh-config.yml       # Zsh configuration
│   ├── ohmyzsh-setup.yml    # Oh My Zsh installation
│   ├── dotfiles.yml         # Dotfiles symlinking
│   ├── git-config.yml       # Git user configuration
│   ├── goto.yml             # goto shell function (GitLab default, vercara org)
│   ├── github-setup.yml     # GitHub SSH setup
│   ├── gitlab-setup.yml     # GitLab SSH setup
│   ├── dev-tools.yml        # Development tools
│   ├── golang.yml           # Go installation and configuration
│   ├── nvm.yml              # NVM and Node.js installation
│   ├── pnpm.yml             # pnpm package manager
│   ├── python.yml           # Python with pip, pipx, and virtualenv
│   ├── applications.yml     # GUI applications (Docker, VSCode, etc.)
│   ├── claude-code.yml      # Claude Code CLI (optional)
│   └── github-copilot.yml   # GitHub Copilot CLI (optional)
├── dotfiles/                 # Example dotfiles to symlink
│   ├── .gitignore_global
│   ├── .vimrc
│   ├── .tmux.conf
│   └── .editorconfig
├── ansible_collections/      # Vendored Ansible collections
└── roles/                    # Vendored Ansible roles
```

## Idempotency

All tasks are designed to be idempotent - safe to run multiple times. The automation:
- Checks existing state before making changes
- Skips tasks if requirements are already met
- Uses Ansible's declarative approach to ensure desired state

Running `laptop.upgrade` a second time should show no changes.

## Development

### Installing Dependencies

```bash
# Install Ansible collections and roles
ansible-galaxy install -r requirements.yml

# Vendor dependencies for offline use
./bin/refresh-ansible-collections.sh
```

### Testing

```bash
# Run with verbose output
DEBUG=1 ./bootstrap.sh

# Syntax check
ansible-playbook main.yml --syntax-check

# Dry run
ansible-playbook main.yml --check
```

## Troubleshooting

### SSH Key Issues

If SSH authentication fails:

```bash
# Verify your SSH key exists
ls -la ~/.ssh/

# Test GitHub connection
ssh -T git@github.com

# Test GitLab connection
ssh -T git@gitlab.com
```

### Homebrew Issues

If Homebrew is not in PATH:

```bash
# ARM Macs (M1/M2/M3)
eval "$(/opt/homebrew/bin/brew shellenv)"

# Intel Macs
eval "$(/usr/local/bin/brew shellenv)"
```

### Oh My Zsh Issues

If Oh My Zsh doesn't load:

```bash
# Verify installation
ls -la ~/.oh-my-zsh

# Check .zshrc
cat ~/.zshrc | grep "oh-my-zsh"
```

## Security Notes

- Personal Access Tokens are stored in `~/.gitconfig`
- SSH keys are generated with 4096-bit RSA by default
- Tokens have `0644` permissions via git config
- Consider using a credential helper for enhanced security

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
