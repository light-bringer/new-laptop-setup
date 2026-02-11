# Starship Prompt Configuration

Starship is a fast, customizable cross-shell prompt written in Rust. It's the default prompt for this setup.

## What You'll See

The prompt displays:

```
┌───────────────────>
│  ~/dev/my-project  main ✓ via  v1.21.0 via  v20.10.0
└─>
```

### Prompt Elements

- **OS Symbol**: `` (macOS)
- **Directory**: Current directory path (truncated if too long)
- **Git Branch**:  `main` - Current git branch
- **Git Status**:
  - `✓` - Clean working tree
  - `📝` - Modified files
  - `++N` - Staged files
  - `🤷` - Untracked files
  - `⇡N` - Ahead of remote
  - `⇣N` - Behind remote
- **Language Versions**:
  - ` v1.21.0` - Go version (when in Go project)
  - ` v20.10.0` - Node.js version (when in Node project)
  - ` 3.12.0` - Python version (when in Python project)
  - ` 1.75.0` - Rust version (when in Rust project)
- **Docker Context**:  `context` - Active Docker context
- **Command Duration**: `took 2s` - For slow commands (>500ms)

## Configuration File

Location: `~/.config/starship.toml`

The default configuration is automatically created during installation.

## Customization

### Change Colors

Edit `~/.config/starship.toml`:

```toml
[directory]
style = "bold cyan"  # Change to any color
```

### Add/Remove Modules

```toml
# Disable a module
[package]
disabled = true

# Enable a module
[aws]
disabled = false
format = 'on [$symbol($profile )(\($region\) )]($style)'
```

### Change Symbols

```toml
[git_branch]
symbol = "🌱 "  # Change branch symbol

[nodejs]
symbol = "⬢ "  # Change Node symbol
```

## Available Modules

Starship supports 100+ modules. Commonly used:

- **Version Control**: git_branch, git_status, git_commit, git_state
- **Languages**: nodejs, python, golang, rust, ruby, php, java, etc.
- **Cloud**: aws, gcloud, azure
- **Tools**: docker_context, kubernetes, terraform
- **System**: battery, memory_usage, time, cmd_duration
- **And many more**: See https://starship.rs/config/

## Preset Configurations

Starship provides several presets:

```bash
# Nerd Font Symbols preset
starship preset nerd-font-symbols -o ~/.config/starship.toml

# Bracketed Segments preset
starship preset bracketed-segments -o ~/.config/starship.toml

# Plain Text Symbols (no icons)
starship preset plain-text-symbols -o ~/.config/starship.toml

# No Runtime Versions
starship preset no-runtime-versions -o ~/.config/starship.toml
```

**Note**: Backup your current config before applying a preset!

## Performance

Starship is extremely fast:
- Written in Rust for maximum performance
- Parallel execution of commands
- Smart caching
- Typical render time: <10ms

## Nerd Fonts

For the best experience, install a Nerd Font:

```bash
# Install a Nerd Font via Homebrew
brew tap homebrew/cask-fonts
brew install font-fira-code-nerd-font
brew install font-meslo-lg-nerd-font
brew install font-hack-nerd-font
```

Then set your terminal to use the font:
- **iTerm2**: Preferences → Profiles → Text → Font
- **Terminal.app**: Preferences → Profiles → Font

## Troubleshooting

### Prompt not showing
```bash
# Verify Starship is installed
starship --version

# Check if init is in .zshrc
grep "starship init" ~/.zshrc

# Reload shell
exec zsh
```

### Symbols not displaying
- Install a Nerd Font (see above)
- Set terminal to use the Nerd Font
- Ensure terminal supports Unicode

### Slow prompt
```bash
# Enable debug mode
export STARSHIP_LOG=trace
starship timings

# Disable slow modules in starship.toml
[nodejs]
disabled = true
```

## Examples

### Minimal Configuration

```toml
format = """
$directory\
$git_branch\
$git_status\
$line_break\
$character"""

[directory]
truncation_length = 3

[character]
success_symbol = "[➜](bold green)"
error_symbol = "[➜](bold red)"
```

### Full-Featured Configuration

```toml
format = """
[┌───────────────────>](bold green)
[│](bold green)$os\
$username\
$hostname\
$directory\
$git_branch\
$git_status\
$golang\
$nodejs\
$python\
$rust\
$docker_context\
$kubernetes\
$aws
[└─>](bold green) """

# Configure each module...
```

### Single-Line Compact

```toml
format = "$directory$git_branch$git_status$character"

[directory]
truncation_length = 1

[character]
success_symbol = "❯"
error_symbol = "❯"
```

## Resources

- **Official Website**: https://starship.rs
- **Configuration Guide**: https://starship.rs/config/
- **Presets**: https://starship.rs/presets/
- **GitHub**: https://github.com/starship/starship
- **Discord**: https://discord.gg/starship

## Quick Tips

1. **Test changes without restarting**:
   ```bash
   starship init zsh | source
   ```

2. **Print current config**:
   ```bash
   starship print-config
   ```

3. **Explain current prompt**:
   ```bash
   starship explain
   ```

4. **Toggle modules on/off**:
   ```bash
   starship toggle python  # Toggle Python module
   ```

5. **Benchmark performance**:
   ```bash
   starship timings
   ```

---

**Installed by**: `./bin/laptop.run --tags starship`
**Configuration**: `~/.config/starship.toml`
**Reload**: `exec zsh` or restart terminal
