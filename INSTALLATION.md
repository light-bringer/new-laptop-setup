# Installation Guide

## Quick Start

### For End Users

If you just want to set up your laptop, it's simple:

```bash
cd ~/dev/new-laptop-setup
./bootstrap.sh
```

That's it! The script will:
1. Install Homebrew (if needed)
2. Install Ansible
3. Add bin/ directory to PATH (for laptop.update, laptop.upgrade)
4. Run the full automation

After the initial setup, you can update your configuration anytime:

```bash
laptop.upgrade
```

## For Developers

### Initial Setup

1. **Clone the repository**:
   ```bash
   git clone <repository-url> ~/dev/new-laptop-setup
   cd ~/dev/new-laptop-setup
   ```

2. **Vendor Ansible dependencies** (optional, for offline availability):
   ```bash
   # First, ensure Ansible is installed
   brew install ansible

   # Then vendor the dependencies
   ./bin/refresh-ansible-collections.sh
   ```

   This downloads and stores Ansible collections and roles locally in:
   - `ansible_collections/`
   - `roles/`

3. **Run the setup**:
   ```bash
   ./bootstrap.sh
   ```

### What Gets Installed

The automation installs and configures:

#### System Tools
- Homebrew package manager
- macOS Command Line Tools
- Ansible (via Homebrew)

#### Shell Environment
- Zsh configuration (.zprofile with Homebrew PATH)
- Oh My Zsh with theme selection
- Curated plugins: git, brew, docker, kubectl
- Update scripts: laptop.update, laptop.upgrade (added to PATH)

#### Dotfiles
Symlinked from `dotfiles/` directory:
- `.gitignore_global` - Global gitignore patterns
- `.vimrc` - Vim configuration
- `.tmux.conf` - Tmux configuration
- `.editorconfig` - EditorConfig settings

#### SSH & Git
- SSH key generation (4096-bit RSA)
- Git user.name and user.email configuration
- GitHub SSH authentication (via Personal Access Token)
- GitLab SSH authentication (via Personal Access Token)
- Git URL rewriting (HTTPS → SSH for GitHub/GitLab)

#### CLI Development Tools (50+ tools)
**Version Control**: git, gh (GitHub CLI), git-lfs
**Text Processing**: jq, yq, ripgrep, fzf, bat, ack
**Network Tools**: curl, wget, httpie
**System Utilities**: tree, htop, ncdu, watch
**Editors**: vim, neovim, tmux
**Build Tools**: make, cmake, pkg-config
**Container/Cloud**: kubectl, helm, terraform, awscli
**Databases**: postgresql, mysql-client, redis
**Security**: gnupg, openssl
**Media**: imagemagick, ffmpeg
**And more**: telnet, nmap, etc.

#### Programming Languages & Tools
- Go programming language with GOPATH configured
- NVM (Node Version Manager) with latest LTS Node.js
- npm (comes with Node.js)
- pnpm (fast, disk space efficient package manager)
- Python 3.12 with pip, pipx, and virtualenv

#### GUI Applications
**Development**: Docker Desktop, Visual Studio Code, iTerm2
**Browsers**: Google Chrome, Firefox
**API Tools**: Postman, Insomnia
**Database**: TablePlus
**Communication**: Slack
**Productivity**: Notion, Rectangle (window manager)
**Design**: Figma

### Required User Input

During the setup, you'll be prompted for:

1. **Git Configuration**:
   - Full name (for commits)
   - Email address (for commits)

2. **GitHub Setup**:
   - Personal Access Token (create at https://github.com/settings/tokens)
     - Required scopes: `admin:public_key`, `read:user`, `repo`
   - GitHub username
   - GitHub email

3. **GitLab Setup**:
   - Personal Access Token (create at https://gitlab.com/-/profile/personal_access_tokens)
     - Required scopes: `api`, `read_user`, `write_repository`
   - GitLab username
   - GitLab email

4. **Oh My Zsh Theme** (optional):
   - Choose from: robbyrussell (default), agnoster, powerlevel10k, pure

### Selective Installation

Run only specific components:

```bash
# Just SSH keys
./bin/laptop.run --tags ssh

# Just Oh My Zsh
./bin/laptop.run --tags ohmyzsh

# Just GitHub setup
./bin/laptop.run --tags github

# Multiple components
./bin/laptop.run --tags ssh,git,github,dotfiles

# Everything except Oh My Zsh
./bin/laptop.run --skip-tags ohmyzsh
```

### Available Tags

- `cli-tools` - macOS Command Line Tools
- `homebrew` - Homebrew installation
- `ssh` - SSH key generation
- `zsh` - Zsh configuration
- `ohmyzsh` - Oh My Zsh installation
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
- Individual apps: `docker`, `vscode`, `iterm2`, `chrome`, `firefox`, `slack`, `postman`, `notion`, `tableplus`, `figma`, `insomnia`, `rectangle`

### Dry Run

Preview what would change without making any changes:

```bash
./bootstrap.sh --check
# or after initial setup
laptop.upgrade --check
```

## Verification

After installation, verify everything works:

### 1. Check Homebrew
```bash
brew --version
```

### 2. Check Oh My Zsh
```bash
echo $ZSH
# Should output: /Users/yourusername/.oh-my-zsh
```

### 3. Check SSH Keys
```bash
ls -la ~/.ssh/
# Should see id_rsa and id_rsa.pub (or id_ed25519, id_ecdsa)
```

### 4. Check Dotfiles
```bash
ls -la ~/ | grep -E "(\.vimrc|\.tmux\.conf|\.editorconfig|\.gitignore_global)"
# All should be symlinks pointing to ~/dev/new-laptop-setup/dotfiles/
```

### 5. Check Git Configuration
```bash
git config --global user.name
git config --global user.email
git config --global core.excludesfile
```

### 6. Test GitHub SSH
```bash
ssh -T git@github.com
# Should output: Hi username! You've successfully authenticated...
```

### 7. Test GitLab SSH
```bash
ssh -T git@gitlab.com
# Should output: Welcome to GitLab, @username!
```

### 8. Test Git URL Rewriting
```bash
git config --global url."ssh://git@github.com/".insteadOf
# Should output: https://github.com/

git config --global url."ssh://git@gitlab.com/".insteadOf
# Should output: https://gitlab.com/
```

### 9. Check Update Scripts
```bash
which laptop.update laptop.upgrade
# Should show paths in ~/dev/new-laptop-setup/bin/
```

### 10. Check Development Tools
```bash
jq --version
rg --version
fzf --version
```

### 11. Check Go Installation
```bash
go version
echo $GOPATH
# Should show Go version and GOPATH set to ~/go
```

### 12. Check NVM and Node.js
```bash
nvm --version
node --version
npm --version
# Should show installed versions
```

### 13. Check pnpm
```bash
pnpm --version
# Should show installed version
```

### 14. Check Python
```bash
python3 --version
pip3 --version
virtualenv --version
# Should show installed versions
```

### 15. Check GUI Applications
```bash
# Docker Desktop
open -a Docker

# VSCode
code --version

# Check if apps are installed
ls /Applications | grep -E "(Docker|Visual Studio Code|iTerm|Slack|Postman)"
```

## Troubleshooting

### Homebrew Not in PATH

```bash
# ARM Macs (M1/M2/M3)
eval "$(/opt/homebrew/bin/brew shellenv)"

# Intel Macs
eval "$(/usr/local/bin/brew shellenv)"
```

Then add to `~/.zprofile`:
```bash
echo 'eval "$(/opt/homebrew/bin/brew shellenv)"' >> ~/.zprofile
```

### SSH Authentication Fails

1. Check your SSH key exists:
   ```bash
   ls -la ~/.ssh/
   ```

2. Verify the key is uploaded:
   ```bash
   # GitHub
   curl -H "Authorization: token YOUR_TOKEN" https://api.github.com/user/keys

   # GitLab
   curl -H "PRIVATE-TOKEN: YOUR_TOKEN" https://gitlab.com/api/v4/user/keys
   ```

3. Test SSH connection:
   ```bash
   ssh -Tv git@github.com
   ssh -Tv git@gitlab.com
   ```

### Oh My Zsh Not Loading

1. Check installation:
   ```bash
   ls -la ~/.oh-my-zsh
   ```

2. Check .zshrc:
   ```bash
   cat ~/.zshrc | grep "oh-my-zsh"
   ```

3. Ensure .zshrc sources .zprofile:
   ```bash
   cat ~/.zshrc | grep ".zprofile"
   ```

### Ansible Errors

1. Ensure Ansible is installed:
   ```bash
   brew install ansible
   ```

2. Check Python version:
   ```bash
   python3 --version
   ```

3. Reinstall dependencies:
   ```bash
   ./bin/refresh-ansible-collections.sh
   ```

### Permission Denied Errors

The script automatically fixes macOS Sonoma permission issues, but if you encounter problems:

```bash
sudo chmod 0755 /opt /private/etc
```

## Uninstallation

To remove the setup:

1. **Remove bin/ from PATH** in `~/.zshenv`:
   ```bash
   # Edit ~/.zshenv and remove the "Engineering Laptop Setup PATH" section
   ```

2. **Remove dotfiles symlinks**:
   ```bash
   rm ~/.vimrc ~/.tmux.conf ~/.editorconfig ~/.gitignore_global
   ```

3. **Remove Oh My Zsh** (optional):
   ```bash
   rm -rf ~/.oh-my-zsh
   ```

4. **Remove SSH keys** (optional, be careful!):
   ```bash
   # Only if you want to completely remove SSH keys
   rm ~/.ssh/id_rsa ~/.ssh/id_rsa.pub
   ```

5. **Remove Git configuration** (optional):
   ```bash
   # Remove GitHub/GitLab tokens
   git config --global --unset github.token
   git config --global --unset gitlab.token

   # Remove URL rewriting
   git config --global --unset url."ssh://git@github.com/".insteadOf
   git config --global --unset url."ssh://git@gitlab.com/".insteadOf
   ```

6. **Uninstall Homebrew** (optional):
   ```bash
   # Only if you want to completely remove Homebrew
   /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/uninstall.sh)"
   ```

## Security Considerations

### Personal Access Tokens

Tokens are stored in `~/.gitconfig` which has permissions `0644` by default. This means:

- ✅ Readable by owner
- ✅ Protected from modification by other users
- ⚠️ Potentially readable by other users on multi-user systems

**Recommendations**:
- Use a credential helper for enhanced security
- Limit token scopes to minimum required
- Rotate tokens periodically
- Never commit .gitconfig to public repositories

### SSH Keys

- Keys are generated with 4096-bit RSA by default (strong security)
- Private keys are stored with permissions `0600` (owner read/write only)
- Public keys are uploaded to GitHub/GitLab automatically

## Updates

To update your setup with the latest changes:

```bash
laptop.upgrade
```

This will:
1. Pull the latest changes from the repository
2. Re-run the full automation
3. Apply any new configurations or tools

The automation is idempotent, so running it multiple times is safe and won't duplicate work.
