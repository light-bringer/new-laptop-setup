# Dotfiles Guide

This setup includes a comprehensive collection of dotfiles to supercharge your development workflow.

## 📂 Included Dotfiles

### 1. `.aliases` - 200+ Shell Aliases

#### Navigation
```bash
..          # cd ..
...         # cd ../..
....        # cd ../../..
ll          # ls -lah (detailed list)
la          # ls -A (all files)
lsd         # List directories only
```

#### Git Shortcuts (40+)
```bash
g           # git
gs          # git status
ga          # git add
gaa         # git add --all
gc          # git commit -v
gcm         # git commit -m
gp          # git push
gpl         # git pull
gco         # git checkout
gcb         # git checkout -b
gb          # git branch
gd          # git diff
gl          # git log --oneline --graph
gst         # git stash
gm          # git merge
```

#### Docker Shortcuts (15+)
```bash
d           # docker
dc          # docker-compose
dps         # docker ps
dpsa        # docker ps -a
di          # docker images
dex         # docker exec -it
dlog        # docker logs -f
dstop       # Stop all running containers
dclean      # Clean up docker system
```

#### Kubernetes Shortcuts (15+)
```bash
k           # kubectl
kgp         # kubectl get pods
kgs         # kubectl get services
kgd         # kubectl get deployments
kdp         # kubectl describe pod
kl          # kubectl logs -f
kex         # kubectl exec -it
```

#### Terraform Shortcuts
```bash
tf          # terraform
tfi         # terraform init
tfp         # terraform plan
tfa         # terraform apply
tfd         # terraform destroy
```

#### NPM/Node/pnpm
```bash
ni          # npm install
nr          # npm run
ns          # npm start
pi          # pnpm install
pr          # pnpm run
```

#### System Utilities
```bash
myip        # Show public IP
localip     # Show local IP
ports       # Show open ports
cleanup     # Remove .DS_Store files
update      # Update Homebrew packages
```

### 2. `.functions` - 30+ Powerful Shell Functions

#### Navigation
```bash
mkcd <dir>          # Create directory and cd into it
up <n>              # Go up n directories
gcl <url>           # Git clone and cd into repo
```

#### File Operations
```bash
extract <file>      # Extract any archive (zip, tar, gz, etc.)
backup <file>       # Create timestamped backup
```

#### Git Helpers
```bash
ginit               # Initialize repo and make first commit
gbrecent            # Show branches by last commit date
gcleanup            # Delete merged branches
```

#### Docker Helpers
```bash
dbash <container>   # Bash into container
dsh <container>     # Sh into container
drmall              # Remove all stopped containers
```

#### Network Utilities
```bash
portcheck <port>    # Check if port is open
killport <port>     # Kill process on port
listening [pattern] # Show what's listening
```

#### Development Tools
```bash
findreplace <old> <new>  # Find and replace in files
loc <extension>          # Count lines of code
genpass [length]         # Generate random password
json                     # Pretty print JSON
server [port]            # Start HTTP server (default: 8000)
```

#### Kubernetes with FZF
```bash
klog                # Fuzzy search pod and tail logs
kexec               # Fuzzy search pod and exec into it
```

#### Python
```bash
mkvenv [name]       # Create and activate venv
```

#### Fun Utilities
```bash
weather [city]      # Get weather forecast
cheat <command>     # Get cheat sheet for command
```

### 3. `.gitconfig_global` - Enhanced Git Configuration

#### Git Aliases (40+)
```bash
git s               # Short status
git st              # Verbose status
git co              # Checkout
git cob             # Checkout new branch
git cm              # Commit with message
git ca              # Amend last commit
git l               # Pretty log graph
git recent          # Branches by last modified
git cleanup         # Delete merged branches
git fm <text>       # Find commits by message
git fc <text>       # Find commits by code
```

#### Configuration
- Color-coded output for branches, diffs, status
- Better diff and merge tools (vimdiff)
- Auto-stash during rebase
- Auto-prune on fetch
- Better conflict resolution (diff3)

### 4. `.inputrc` - Better Readline Experience

Features:
- ✅ Case-insensitive tab completion
- ✅ Up/Down arrows search history by prefix
- ✅ Show all completions without asking
- ✅ Smart autocomplete for symlinks
- ✅ UTF-8 support
- ✅ Don't autocomplete hidden files unless explicit

### 5. `.curlrc` - Sensible Curl Defaults

Features:
- ✅ Follow redirects automatically
- ✅ Show progress bar
- ✅ 60-second timeout
- ✅ Show errors
- ✅ Save cookies automatically

### 6. `.wgetrc` - Sensible Wget Defaults

Features:
- ✅ Use server timestamps
- ✅ Don't ascend in directory tree
- ✅ Retry 3 times on failure
- ✅ Follow FTP links
- ✅ Add .html extension when needed
- ✅ Ignore robots.txt

### 7. `.terraformrc` - Terraform Configuration

Features:
- ✅ Plugin cache directory (saves bandwidth)
- ✅ Disable automatic crash reporting
- ✅ Provider installation configuration

### 8. `.vimrc` - Vim Configuration

Features:
- Line numbers and relative numbers
- Syntax highlighting
- Smart indentation (2 spaces)
- Better search (highlight, incremental, smart case)
- Mouse support
- File type detection

### 9. `.tmux.conf` - Tmux Configuration

Features:
- Custom prefix (C-a instead of C-b)
- Better split pane bindings (| and -)
- Mouse mode enabled
- Vi keybindings
- Beautiful status bar
- 10,000 line history

### 10. `.gitignore_global` - Global Git Ignore

Ignores:
- macOS system files (.DS_Store, etc.)
- IDE configurations (.idea, .vscode, *.swp)
- Build artifacts (node_modules, dist, build)
- Environment files (.env)
- Logs and caches

### 11. `.editorconfig` - Editor Configuration

Settings for:
- Universal: 2-space indentation, LF line endings, UTF-8
- Python: 4-space indentation
- Go: Tab indentation
- Makefile: Tab indentation

## 🚀 How to Use

### All Dotfiles Are Automatically Loaded

After running the setup:
- `.aliases` sourced in `.zshrc`
- `.functions` sourced in `.zshrc`
- `.gitconfig_global` included in `.gitconfig`
- All other files symlinked to home directory

### Customization

All dotfiles are symlinked from `~/dev/new-laptop-setup/dotfiles/`, so:

```bash
# Edit any dotfile
cd ~/dev/new-laptop-setup/dotfiles
vim .aliases

# Changes apply immediately (symlinked, not copied)
# Restart shell to reload
exec zsh
```

### Adding Your Own Aliases

```bash
# Edit .aliases
echo "alias myalias='my command'" >> ~/dev/new-laptop-setup/dotfiles/.aliases

# Reload
source ~/.zshrc
```

### Adding Your Own Functions

```bash
# Edit .functions
vim ~/dev/new-laptop-setup/dotfiles/.functions

# Add your function at the end
myfunction() {
  echo "Hello from my function"
}

# Reload
source ~/.zshrc
```

## 📊 Statistics

- **Total Aliases**: 200+
- **Total Functions**: 30+
- **Git Aliases**: 40+
- **Total Lines**: ~800
- **Categories**: 11 different config files

## 🎯 Most Useful Shortcuts

### Daily Git Workflow
```bash
gco -b feature/new    # Create new branch
ga .                  # Stage all
gcm "Add feature"     # Commit
gp                    # Push
```

### Docker Development
```bash
dps                   # Check containers
dlog <container>      # Tail logs
dex <container>       # Exec into container
dclean                # Clean up
```

### Kubernetes Operations
```bash
kgp                   # List pods
klog                  # Fuzzy find and tail logs
kexec                 # Fuzzy find and exec
```

### Quick Navigation
```bash
mkcd ~/dev/project    # Create and cd
..                    # Go up one level
...                   # Go up two levels
```

### File Operations
```bash
extract archive.tar.gz # Extract any archive
backup important.txt   # Create timestamped backup
ff "*.js"             # Find JavaScript files
```

## 💡 Pro Tips

1. **Use Tab Completion**: Most commands have excellent tab completion
2. **Learn FZF**: klog and kexec use fuzzy finder for quick pod selection
3. **Customize Git**: Add your own git aliases to .gitconfig_global
4. **Shell Reload**: After editing, run `source ~/.zshrc` or `exec zsh`
5. **Git Cleanup**: Run `git cleanup` regularly to remove merged branches
6. **Port Management**: Use `killport 3000` to free up stuck ports
7. **Weather**: Run `weather` to get forecast
8. **Cheat Sheets**: Use `cheat <command>` for quick reference

## 🔍 Discovery

### List All Aliases
```bash
alias                 # Show all aliases
alias | grep git      # Show git aliases
alias | grep docker   # Show docker aliases
```

### List All Functions
```bash
declare -F            # List all functions
type <function>       # Show function definition
```

### Git Aliases
```bash
git aliases           # Show all git aliases
```

## 📚 Resources

- **Aliases Guide**: Common shell aliases
- **Functions Guide**: Reusable shell functions
- **Git Aliases**: https://git-scm.com/book/en/v2/Git-Basics-Git-Aliases
- **Readline**: https://www.gnu.org/software/bash/manual/html_node/Readline-Init-File.html
- **EditorConfig**: https://editorconfig.org/

---

**All dotfiles are version-controlled** in this repository and symlinked to your home directory for easy updates and portability.
