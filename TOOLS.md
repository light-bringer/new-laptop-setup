# Complete Tool List

This document lists all tools and applications installed by this setup.

## 🖥️ GUI Applications (Homebrew Cask)

### Development Tools
- **Docker Desktop** - Container platform for building and running applications
- **Visual Studio Code** - Modern code editor with extensions
- **iTerm2** - Better terminal emulator for macOS

### Web Browsers
- **Google Chrome** - Web browser
- **Firefox** - Open-source web browser

### API & Database Tools
- **Postman** - API development and testing platform
- **Insomnia** - Alternative API client with GraphQL support
- **TablePlus** - Modern database management tool

### Communication
- **Slack** - Team communication and collaboration

### Productivity
- **Notion** - Notes, docs, and project management
- **Rectangle** - Window management for macOS (keyboard shortcuts)

### Design
- **Figma** - Collaborative design and prototyping tool

## 🛠️ CLI Development Tools (50+ tools)

### Version Control
- **git** - Distributed version control system
- **gh** - Official GitHub CLI
- **git-lfs** - Git Large File Storage extension

### Text Processing & Search
- **jq** - Command-line JSON processor
- **yq** - Command-line YAML processor
- **ripgrep (rg)** - Fast recursive search tool
- **fzf** - Fuzzy finder for command-line
- **bat** - Cat clone with syntax highlighting
- **ack** - Code search tool optimized for programmers

### Network Tools
- **curl** - Transfer data with URLs
- **wget** - Network downloader
- **httpie** - User-friendly HTTP client
- **telnet** - Telnet client
- **nmap** - Network scanner and security tool

### System Utilities
- **tree** - Display directory tree structure
- **htop** - Interactive process viewer
- **ncdu** - NCurses disk usage analyzer
- **watch** - Execute commands periodically

### Text Editors & Terminal
- **vim** - Classic text editor
- **neovim** - Modern, extensible vim
- **tmux** - Terminal multiplexer

### Build Tools
- **make** - Build automation tool
- **cmake** - Cross-platform build system
- **pkg-config** - Package configuration helper

### Container & Orchestration
- **kubectl** - Kubernetes command-line tool
- **helm** - Kubernetes package manager
- **docker** - Container CLI (via Docker Desktop)

### Infrastructure as Code
- **terraform** - Infrastructure provisioning and management

### Cloud Platforms
- **awscli** - Amazon Web Services CLI
- **aws-vault** - Secure AWS credential storage (optional)

### Database CLIs
- **postgresql@16** - PostgreSQL client tools
- **mysql-client** - MySQL command-line client
- **redis** - Redis CLI and server

### Security & Encryption
- **gnupg** - GPG encryption tool
- **openssl** - SSL/TLS toolkit and crypto library

### Media Processing
- **imagemagick** - Image manipulation and conversion
- **ffmpeg** - Video and audio processing

## 🐹 Programming Languages

### Go
- **Go** - Programming language
- **GOPATH** - Configured at `~/go`
- **Go bins** - Added to PATH

### Node.js Ecosystem
- **NVM** - Node Version Manager
- **Node.js** - Latest LTS version
- **npm** - Node package manager (bundled with Node)
- **pnpm** - Fast, disk space efficient package manager

### Python
- **Python 3.12** - Latest stable Python
- **pip3** - Python package installer
- **pipx** - Install and run Python applications in isolated environments
- **virtualenv** - Python virtual environment creator

## 🐚 Shell Environment

### Shell Configuration
- **Zsh** - Modern shell (macOS default)
- **Oh My Zsh** - Zsh framework with plugins and themes

### Oh My Zsh Themes
- **robbyrussell** - Simple and clean (default)
- **agnoster** - Powerline-style with git info
- **powerlevel10k** - Feature-rich with customization
- **refined** - Minimal and elegant

### Oh My Zsh Plugins
- **git** - Git aliases and functions
- **brew** - Homebrew completions
- **docker** - Docker completions
- **kubectl** - Kubernetes completions

## 📝 Dotfiles

### Included Configurations
- **.vimrc** - Vim editor configuration
  - Line numbers, syntax highlighting
  - Smart indentation, search settings
  - Mouse support

- **.tmux.conf** - Tmux terminal multiplexer
  - Custom prefix (C-a)
  - Mouse mode enabled
  - Vi keybindings
  - Status bar customization

- **.gitignore_global** - Global git ignore patterns
  - macOS system files
  - IDE configurations
  - Build artifacts
  - Environment files

- **.editorconfig** - Editor configuration
  - Consistent coding styles
  - Language-specific settings

## 🚀 Custom Shell Functions

### goto Function
A custom shell function for quick navigation to GitHub repositories.

**Features**:
- Jump to existing repositories under `$GOTO_ROOT/github.com/$user/$repo`
- Auto-clone repositories if they don't exist locally
- Support for both `goto my-repo` and `goto user/repo` syntax

**Environment Variables**:
- `GOTO_ROOT` - Root directory for repositories (default: `~/dev/src`)
- `GOTO_DEFAULT_USER` - Default GitHub username

**Usage Examples**:
```bash
# Navigate to your repository (uses GOTO_DEFAULT_USER)
goto my-project

# Navigate to someone else's repository
goto octocat/Hello-World

# If repository doesn't exist locally, it will be cloned automatically
```

## 🔐 SSH & Git

### SSH Configuration
- **SSH keys** - 4096-bit RSA key pair
- **known_hosts** - Configured for GitHub and GitLab

### Git Configuration
- **user.name** - Your full name
- **user.email** - Your email address
- **core.excludesfile** - Global gitignore
- **URL rewriting** - HTTPS → SSH for GitHub/GitLab

### Platform Integration
- **GitHub** - SSH authentication via API
- **GitLab** - SSH authentication via API

## 📦 Package Managers

- **Homebrew** - macOS package manager
- **Homebrew Cask** - GUI application manager
- **npm** - Node.js package manager
- **pnpm** - Fast Node.js package manager
- **pip3** - Python package manager
- **pipx** - Isolated Python app installer

## 🔄 Update Scripts

Located in `bin/` directory:

- **laptop.update** - Pull latest changes from repository
- **laptop.upgrade** - Update and run full setup
- **laptop.run** - Core execution script
- **refresh-ansible-collections.sh** - Vendor Ansible dependencies

## 🏷️ Installation Tags

### Component Tags
- `cli-tools` - macOS Command Line Tools
- `homebrew` - Homebrew installation
- `ssh` - SSH key generation
- `zsh` - Zsh configuration
- `ohmyzsh` - Oh My Zsh installation
- `dotfiles` - Dotfiles symlinking
- `git` - Git configuration
- `goto` - goto shell function
- `github` - GitHub setup
- `gitlab` - GitLab setup
- `dev-tools` - CLI development tools
- `golang` / `go` - Go language
- `nvm` / `node` / `nodejs` - NVM and Node.js
- `pnpm` - pnpm package manager
- `python` / `py` - Python setup
- `applications` / `apps` / `gui` - GUI applications

### Individual Application Tags
- `docker` - Docker Desktop
- `vscode` - Visual Studio Code
- `iterm2` - iTerm2 terminal
- `chrome` - Google Chrome
- `firefox` - Firefox browser
- `slack` - Slack
- `postman` - Postman API tool
- `insomnia` - Insomnia API client
- `notion` - Notion
- `tableplus` - TablePlus database client
- `figma` - Figma design tool
- `rectangle` - Rectangle window manager

## 📊 Summary

- **Total CLI Tools**: 50+
- **Total GUI Apps**: 12
- **Programming Languages**: 3 (Go, Node.js, Python)
- **Package Managers**: 6 (brew, cask, npm, pnpm, pip, pipx)
- **Dotfiles**: 4 (vim, tmux, gitignore, editorconfig)
- **Platform Integrations**: 2 (GitHub, GitLab)
- **Total Tags**: 21+

## 🚀 Usage Examples

### Install Everything
```bash
./bootstrap.sh
```

### Install Specific Categories
```bash
# Development tools only
./bin/laptop.run --tags dev-tools

# Programming languages
./bin/laptop.run --tags golang,nvm,python

# GUI applications only
./bin/laptop.run --tags applications
```

### Install Specific Applications
```bash
# Just Docker and VSCode
./bin/laptop.run --tags docker,vscode

# Just browsers
./bin/laptop.run --tags chrome,firefox
```

### Skip Categories
```bash
# Everything except GUI apps
./bin/laptop.run --skip-tags applications

# Everything except Oh My Zsh
./bin/laptop.run --skip-tags ohmyzsh
```

## 🔍 Verification Commands

```bash
# CLI tools
git --version
go version
node --version
python3 --version
kubectl version --client
terraform version
docker --version

# Package managers
brew --version
npm --version
pnpm --version
pip3 --version

# GUI apps
ls /Applications | grep -E "(Docker|Visual Studio Code|iTerm|Slack)"

# Shell
echo $SHELL
echo $ZSH
```

---

**Note**: This setup is idempotent - running it multiple times is safe and won't duplicate installations.
