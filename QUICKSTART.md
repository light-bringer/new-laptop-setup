# Quick Start Guide

## TL;DR

```bash
cd ~/dev/new-laptop-setup
./bootstrap.sh
```

After setup, use these commands to update:

```bash
laptop.update    # Pull latest changes
laptop.upgrade   # Update and run setup
```

## What You'll Need

Before running the setup, prepare:

1. **GitHub Personal Access Token**
   - Go to: https://github.com/settings/tokens
   - Click "Generate new token (classic)"
   - Select scopes: `admin:public_key`, `read:user`, `repo`
   - Copy the token (you'll need it during setup)

2. **GitLab Personal Access Token**
   - Go to: https://gitlab.com/-/profile/personal_access_tokens
   - Click "Add new token"
   - Select scopes: `api`, `read_user`, `write_repository`
   - Copy the token (you'll need it during setup)

3. **Your Information**
   - Full name (for Git commits)
   - Email address (for Git commits)
   - GitHub username
   - GitHub email
   - GitLab username
   - GitLab email

## Setup Process

### Step 1: Run Bootstrap

```bash
cd ~/dev/new-laptop-setup
./bootstrap.sh
```

This will:
- ✅ Install Homebrew
- ✅ Install Ansible
- ✅ Add bin/ directory to PATH
- ✅ Start the automated setup

### Step 2: Answer Prompts

During setup, you'll be asked for:

1. **Git Configuration**
   - Your full name
   - Your email address

2. **GitHub Configuration**
   - Personal Access Token (paste the token you created)
   - GitHub username
   - GitHub email

3. **GitLab Configuration**
   - Personal Access Token (paste the token you created)
   - GitLab username
   - GitLab email

4. **Oh My Zsh Theme** (optional)
   - Press Enter for default (robbyrussell)
   - Or choose: agnoster, powerlevel10k, pure

### Step 3: Verify

After setup completes, verify everything works:

```bash
# Test GitHub SSH
ssh -T git@github.com

# Test GitLab SSH
ssh -T git@gitlab.com

# Check update scripts
which laptop.update laptop.upgrade

# Check Oh My Zsh
echo $ZSH
```

## What Gets Installed

### Tools
- Homebrew, Ansible, git, jq
- curl, wget, tree, htop
- ripgrep, fzf, tmux, vim
- Go programming language
- NVM with Node.js LTS and npm
- pnpm package manager
- Python 3.12 with pip, pipx, and virtualenv

### Shell
- Oh My Zsh with plugins: git, brew, docker, kubectl
- Zsh configured with Homebrew PATH
- Update scripts: laptop.update, laptop.upgrade

### Dotfiles
- `.vimrc` - Vim configuration
- `.tmux.conf` - Tmux configuration
- `.gitignore_global` - Global gitignore
- `.editorconfig` - EditorConfig settings

### SSH & Git
- SSH keys (4096-bit RSA)
- Git user configuration
- GitHub SSH authentication
- GitLab SSH authentication
- Git HTTPS→SSH URL rewriting

## Common Commands

### Update Your Setup
```bash
laptop.upgrade
```

### Install Only Specific Components
```bash
# Just SSH keys
./bin/laptop.run --tags ssh

# Just Oh My Zsh
./bin/laptop.run --tags ohmyzsh

# Just GitHub setup
./bin/laptop.run --tags github

# Just dotfiles
./bin/laptop.run --tags dotfiles
```

### Preview Changes (Dry Run)
```bash
laptop.upgrade --check
```

### Run with Debug Output
```bash
DEBUG=1 ./bootstrap.sh
```

## Troubleshooting

### Homebrew Not Found
```bash
# ARM Macs
eval "$(/opt/homebrew/bin/brew shellenv)"

# Intel Macs
eval "$(/usr/local/bin/brew shellenv)"
```

### SSH Authentication Fails
```bash
# Check SSH key exists
ls -la ~/.ssh/

# Test connection with verbose output
ssh -Tv git@github.com
ssh -Tv git@gitlab.com
```

### Oh My Zsh Not Loading
```bash
# Restart your shell
exec zsh

# Or source .zshrc
source ~/.zshrc
```

## Next Steps

After setup:

1. **Restart your terminal** to activate Oh My Zsh
2. **Test GitHub** by cloning a repository:
   ```bash
   git clone git@github.com:username/repo.git
   ```
3. **Test GitLab** by cloning a repository:
   ```bash
   git clone git@gitlab.com:username/repo.git
   ```
4. **Customize dotfiles** in `~/dev/new-laptop-setup/dotfiles/`
   - Changes apply immediately (they're symlinked)

## Getting Help

- **Full documentation**: See `README.md`
- **Installation guide**: See `INSTALLATION.md`
- **Repository structure**: See `CLAUDE.md`

## Skip Components

Don't want Oh My Zsh? Skip it:
```bash
./bin/laptop.run --skip-tags ohmyzsh
```

Don't want GitLab? Skip it:
```bash
./bin/laptop.run --skip-tags gitlab
```

Only want SSH and Git:
```bash
./bin/laptop.run --tags ssh,git
```

## Safety

- ✅ All operations are idempotent (safe to run multiple times)
- ✅ Existing SSH keys are preserved (only generates if none exist)
- ✅ Git configs are only set if not already configured
- ✅ Dry run available with `--check` flag
- ✅ No destructive operations without confirmation

## Updates

To update your setup with the latest changes from the repository:

```bash
laptop.upgrade
```

This pulls the latest changes and re-runs the automation.

---

**Ready?** Run `./bootstrap.sh` to get started!
