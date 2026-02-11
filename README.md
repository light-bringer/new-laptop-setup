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
- 🚀 **goto** - Custom shell function for quick GitHub repository navigation
- 🛠️ **CLI Development tools** (50+ tools including git, jq, curl, kubectl, terraform, awscli)
- 🐹 **Go** programming language with GOPATH configuration
- 📦 **NVM** (Node Version Manager) with latest LTS Node.js
- ⚡ **pnpm** - Fast, disk space efficient package manager
- 🐍 **Python 3.12** with pip, pipx, and virtualenv
- 🖥️ **GUI Applications** - Docker Desktop, VSCode, iTerm2, Chrome, Slack, Postman, and more
- 🔄 **Update scripts** for easy maintenance (laptop.update, laptop.upgrade)

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
# Install Claude Code CLI (not installed by default)
./bin/laptop.run --tags claude-code
```

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
- `claude-code` / `claude` - Claude Code CLI (requires explicit `--tags claude-code`)

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
│   ├── github-setup.yml     # GitHub SSH setup
│   ├── gitlab-setup.yml     # GitLab SSH setup
│   ├── dev-tools.yml        # Development tools
│   ├── golang.yml           # Go installation and configuration
│   ├── nvm.yml              # NVM and Node.js installation
│   ├── pnpm.yml             # pnpm package manager
│   ├── python.yml           # Python with pip, pipx, and virtualenv
│   ├── applications.yml     # GUI applications (Docker, VSCode, etc.)
│   └── claude-code.yml      # Claude Code CLI (optional)
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
