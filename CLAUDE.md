# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

This repository provides Ansible-based automation for configuring macOS developer environments. It's a clean, modern laptop setup tool that installs essential development tools, configures SSH access to GitHub and GitLab, and sets up a productive shell environment with Oh My Zsh and curated dotfiles.

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
  - **ohmyzsh-setup.yml**: Oh My Zsh installation and theme selection
  - **dotfiles.yml**: Symlink dotfiles from repo to home directory
  - **git-config.yml**: Git user.name and user.email
  - **github-setup.yml**: GitHub SSH authentication via API
  - **gitlab-setup.yml**: GitLab SSH authentication via API
  - **dev-tools.yml**: Common development tools via Homebrew
  - **golang.yml**: Go installation and GOPATH configuration
  - **nvm.yml**: NVM installation with latest LTS Node.js
  - **pnpm.yml**: pnpm package manager installation
  - **python.yml**: Python with pip, pipx, and virtualenv
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
5. **tasks/ohmyzsh-setup.yml**: Install Oh My Zsh with theme selection
6. **tasks/dotfiles.yml**: Symlink dotfiles to home directory
7. **tasks/git-config.yml**: Configure git user.name and user.email
8. **tasks/github-setup.yml**: GitHub SSH authentication and key upload
9. **tasks/gitlab-setup.yml**: GitLab SSH authentication and key upload
10. **tasks/dev-tools.yml**: Install development tools
11. **tasks/golang.yml**: Install Go and configure GOPATH
12. **tasks/nvm.yml**: Install NVM and latest LTS Node.js
13. **tasks/pnpm.yml**: Install pnpm package manager
14. **tasks/python.yml**: Install Python with pip, pipx, and virtualenv

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

## Important Notes

- The repository is located at `~/dev/new-laptop-setup` by default
- SSH keys are generated as 4096-bit RSA for maximum compatibility
- The automation prompts for sudo once and keeps it alive with a background refresh loop
- Directory permissions are fixed for macOS Sonoma compatibility (`chmod 0755 /opt /private/etc`)
- Git HTTPS URLs are automatically converted to SSH for GitHub and GitLab
- Tokens are stored in git config (acceptable for development machines, consider keychain for production)
