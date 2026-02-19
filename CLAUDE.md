# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This repository provides Ansible-based automation for configuring macOS developer environments. It's a clean, modern laptop setup tool that installs essential development tools, configures SSH access to GitHub and GitLab, and sets up a productive shell environment with Oh My Zsh and curated dotfiles.

## Prerequisites

Users need the following before running setup:

### System Requirements
- macOS 12+ (Monterey or later)
- Admin/sudo access
- Internet connection
- 5GB+ free disk space

### Required Credentials
- **GitHub Personal Access Token**: Scopes: `admin:public_key`, `read:user`, `repo` (https://github.com/settings/tokens)
- **GitLab Personal Access Token**: Scopes: `api`, `read_user`, `write_repository` (https://gitlab.com/-/profile/personal_access_tokens)

### Information Needed During Setup
- Git name and email (for commits)
- GitHub username and email
- GitLab username and email
- goto platform choice: `gitlab` (default) or `github`
- goto default org/user: Defaults to `vercara`
- goto root directory: Defaults to `~/dev/src`
- Oh My Zsh theme: `robbyrussell` (default), `agnoster`, `powerlevel10k`, or `pure`

### Optional Credentials
- **GitHub Copilot**: Active subscription required
- **Claude Code CLI**: Anthropic API key (https://console.anthropic.com/)

The README.md contains detailed step-by-step instructions for creating all required credentials.

## Key Commands

### Initial Setup
```bash
./bootstrap.sh
```
Installs Homebrew and Ansible, sets up shell aliases, then runs the full laptop setup automation. This is the entrypoint for new users.

### Update Scripts (Available After Bootstrap)
```bash
laptop.update    # Pull latest changes from main branch
laptop.upgrade   # Update and run full setup
```
These scripts are located in `bin/` and added to PATH via `~/.zshenv` after the first bootstrap run.

### Selective Execution
```bash
# Run specific components
./bin/laptop.run --tags ssh
./bin/laptop.run --tags ohmyzsh,dotfiles
./bin/laptop.run --skip-tags github

# Install optional components (not installed by default)
./bin/laptop.run --tags claude-code      # Install Claude Code CLI
./bin/laptop.run --tags github-copilot   # Install GitHub Copilot CLI

# Preview changes without applying
./bin/laptop.run --check
```

### Development
```bash
# Install Ansible dependencies
ansible-galaxy install -r requirements.yml

# Vendor dependencies locally
./bin/refresh-ansible-collections.sh

# Syntax check
ansible-playbook main.yml --syntax-check
```

## Architecture

### Execution Flow

1. **bootstrap.sh** → Entry point that installs Homebrew, adds bin/ to PATH, then delegates to `bin/laptop.run`
2. **bin/laptop.run** → Installs Ansible via Homebrew, configures Python environment, runs `ansible-playbook main.yml`
3. **main.yml** → Main Ansible playbook that orchestrates all roles and tasks

### Directory Structure

- **bootstrap.sh**: Initial setup script with Homebrew installation and PATH setup
- **bin/**: Utility scripts (`laptop.run`, `laptop.update`, `laptop.upgrade`, `refresh-ansible-collections.sh`)
- **main.yml**: Main Ansible playbook defining execution order
- **tasks/**: Ansible task files for specific components
  - **ssh-key.yml**: SSH key generation (RSA 4096-bit)
  - **zsh-config.yml**: Zsh + brew shellenv configuration
  - **ohmyzsh-setup.yml**: Oh My Zsh installation
  - **starship.yml**: Starship prompt installation and configuration
  - **dotfiles.yml**: Symlink dotfiles from repo to home directory
  - **git-config.yml**: Git user.name and user.email
  - **goto.yml**: Install goto shell function for quick repository navigation (GitLab default)
  - **github-setup.yml**: GitHub SSH authentication via API
  - **gitlab-setup.yml**: GitLab SSH authentication via API
  - **dev-tools.yml**: Common development tools via Homebrew
  - **golang.yml**: Go installation and GOPATH configuration
  - **nvm.yml**: NVM installation with latest LTS Node.js
  - **pnpm.yml**: pnpm package manager installation
  - **python.yml**: Python with pip, pipx, and virtualenv
  - **applications.yml**: GUI applications via Homebrew Cask (Docker, VSCode, etc.)
  - **claude-code.yml**: Claude Code CLI installation (optional)
  - **github-copilot.yml**: GitHub Copilot CLI installation via gh extension (optional)
- **dotfiles/**: Example dotfiles to symlink
  - **.gitignore_global**: Global gitignore patterns
  - **.vimrc**: Vim configuration
  - **.tmux.conf**: Tmux configuration
  - **.editorconfig**: EditorConfig settings
- **ansible_collections/**: Vendored Ansible collections (offline availability)
- **roles/**: Vendored Ansible roles

### Task Execution Order

Defined in `main.yml`:

1. **elliotweiser.osx-command-line-tools**: Installs macOS Command Line Tools
2. **geerlingguy.mac.homebrew**: Installs and configures Homebrew
3. **tasks/ssh-key.yml**: Generate SSH keys if none exist
4. **tasks/zsh-config.yml**: Configure .zprofile with brew shellenv
5. **tasks/ohmyzsh-setup.yml**: Install Oh My Zsh
6. **tasks/starship.yml**: Install and configure Starship prompt
7. **tasks/dotfiles.yml**: Symlink dotfiles to home directory
8. **tasks/git-config.yml**: Configure git user.name and user.email
9. **tasks/goto.yml**: Install goto shell function (GitLab default, defaults to vercara org)
10. **tasks/github-setup.yml**: GitHub SSH authentication and key upload
11. **tasks/gitlab-setup.yml**: GitLab SSH authentication and key upload
12. **tasks/dev-tools.yml**: Install development tools
13. **tasks/golang.yml**: Install Go and configure GOPATH
14. **tasks/nvm.yml**: Install NVM and latest LTS Node.js
15. **tasks/pnpm.yml**: Install pnpm package manager
16. **tasks/python.yml**: Install Python with pip, pipx, and virtualenv
17. **tasks/applications.yml**: Install GUI applications via Homebrew Cask

**Optional Components** (not installed by default, use `--tags` to enable):
- **tasks/claude-code.yml**: Install Claude Code CLI via npm
- **tasks/github-copilot.yml**: Install GitHub Copilot CLI via gh extension

## GitHub/GitLab Integration

### Authentication Method

This setup uses **Personal Access Tokens (PAT)** with the GitHub/GitLab REST API:

- **GitHub**: Requires token with `admin:public_key`, `read:user`, `repo` scopes
- **GitLab**: Requires token with `api`, `read_user`, `write_repository` scopes

The automation prompts for tokens, stores them in `~/.gitconfig`, and uses them to:
- Upload SSH keys automatically
- Configure Git to use SSH URLs instead of HTTPS

### Storage

- Tokens stored in: `~/.gitconfig` (keys: `github.token`, `gitlab.token`)
- Usernames stored in: `~/.gitconfig` (keys: `github.username`, `gitlab.username`)
- Emails stored in: `~/.gitconfig` (keys: `github.email`, `gitlab.email`)

## Dotfiles Management

The `dotfiles/` directory contains example configurations that are **symlinked** to the home directory:

- **Why symlinks?** Changes to files in the repo automatically apply to the system
- **force: yes** in Ansible replaces existing dotfiles (users can remove symlinks if they prefer custom configs)

Included dotfiles:
- `.gitignore_global` - Global gitignore patterns for macOS, IDEs, build artifacts
- `.vimrc` - Sensible Vim defaults (line numbers, syntax highlighting, search settings)
- `.tmux.conf` - Improved tmux key bindings (C-a prefix, mouse support, vi mode)
- `.editorconfig` - Consistent coding styles across editors

## Oh My Zsh Integration

The `tasks/ohmyzsh-setup.yml` task:

1. Downloads and runs official Oh My Zsh installer with `RUNZSH=no CHSH=no`
2. Prompts user for theme selection (robbyrussell, agnoster, powerlevel10k, pure)
3. Configures plugins: `git brew docker kubectl`
4. Ensures `.zshrc` sources `.zprofile` for Homebrew PATH

**Why RUNZSH=no?** Prevents installer from exec'ing zsh, which would break Ansible execution

## goto Shell Function

The `tasks/goto.yml` task installs a custom shell function for quick repository navigation and cloning:

**Configuration Prompts:**
1. **Platform choice**: GitLab (default) or GitHub
2. **Default organization/user**: Defaults to `vercara` (customizable)
3. **Root directory**: Defaults to `~/dev/src` (customizable)

**Functionality:**
- Navigates to existing local repositories
- Automatically clones repositories via SSH if not found locally
- Directory structure: `$GOTO_ROOT/$PLATFORM/$ORG/$REPO`

**Environment Variables:**
- `GOTO_ROOT` - Base directory for repositories (default: `~/dev/src`)
- `GOTO_DEFAULT_USER` - Default organization/user (default: `vercara`)
- `GOTO_PLATFORM` - Platform choice: `gitlab` or `github` (default: `gitlab`)
- `GOTO_DOMAIN` - Resolved domain: `gitlab.com` or `github.com`

**Usage Examples:**
```bash
goto my-repo              # → ~/dev/src/gitlab.com/vercara/my-repo
goto other-org/project    # → ~/dev/src/gitlab.com/other-org/project
goto                      # Shows current configuration
```

**Implementation Details:**
- Function added to `.zshrc` via `blockinfile` (Ansible-managed)
- Uses SSH for cloning: `git@$GOTO_DOMAIN:$user/$repo.git`
- Creates directory structure automatically
- Idempotent: Re-running won't overwrite function

## Update Scripts

The `bootstrap.sh` script adds the `bin/` directory to PATH in `~/.zshenv`, making these scripts available:

- **bin/laptop.update**: Pulls latest changes from main branch (with uncommitted changes check)
- **bin/laptop.upgrade**: Runs laptop.update, then executes bin/laptop.run

**Pattern**: Follows the Twilio Segment repo pattern of dedicated scripts in bin/ directory

## Idempotency

All tasks are designed to be idempotent:

- SSH keys only generated if none exist (checks id_rsa, id_ed25519, id_ecdsa)
- Git config only prompts if values not already set
- GitHub/GitLab keys only uploaded if not already present
- Dotfiles symlinks use `force: yes` to ensure consistent state
- Oh My Zsh only installed if `~/.oh-my-zsh` doesn't exist

Running `laptop.upgrade` multiple times should show no changes after the first run.

## Common Patterns

### Adding New Development Tools

Add to `tasks/dev-tools.yml`:

```yaml
- name: Install common development tools
  community.general.homebrew:
    name: "{{ item }}"
    state: present
  loop:
    - new-tool-name
```

### Adding New Dotfiles

1. Create file in `dotfiles/` directory
2. Add symlink task to `tasks/dotfiles.yml`:

```yaml
- name: Symlink new dotfile
  ansible.builtin.file:
    src: '{{ playbook_dir }}/dotfiles/.newfile'
    dest: '{{ ansible_env.HOME }}/.newfile'
    state: link
    force: yes
```

### Modifying Oh My Zsh Plugins

Edit `tasks/ohmyzsh-setup.yml` and update the plugins list:

```yaml
- name: Configure Oh My Zsh plugins
  ansible.builtin.lineinfile:
    path: '{{ ansible_env.HOME }}/.zshrc'
    regexp: '^plugins='
    line: 'plugins=(git brew docker kubectl new-plugin)'
```

## Optional AI CLI Tools

### GitHub Copilot CLI

GitHub Copilot CLI provides AI-powered command-line assistance. It's installed as a GitHub CLI extension.

**Installation:**
```bash
./bin/laptop.run --tags github-copilot
```

**Requirements:**
- GitHub CLI (`gh`) is already installed via `dev-tools.yml`
- Requires GitHub authentication: `gh auth login`
- Requires active GitHub Copilot subscription

**Usage:**
- `gh copilot explain` - Explain shell commands
- `gh copilot suggest` - Suggest commands to accomplish a task

**Suggested Shell Aliases:**
```bash
alias ghce='gh copilot explain'
alias ghcs='gh copilot suggest'
```

### Claude Code CLI

Claude Code CLI provides access to Claude AI in the terminal.

**Installation:**
```bash
./bin/laptop.run --tags claude-code
```

**Requirements:**
- Installed via npm (requires NVM from `nvm.yml`)
- Requires Anthropic API key for authentication

**First Run:**
```bash
claude auth  # Authenticate with Anthropic
claude       # Start interactive session
```

## Important Notes

- The repository is located at `~/dev/new-laptop-setup` by default
- SSH keys are generated as 4096-bit RSA for maximum compatibility
- The automation prompts for sudo once and keeps it alive with a background refresh loop
- Directory permissions are fixed for macOS Sonoma compatibility (`chmod 0755 /opt /private/etc`)
- Git HTTPS URLs are automatically converted to SSH for GitHub and GitLab
- Tokens are stored in git config (acceptable for development machines, consider keychain for production)
