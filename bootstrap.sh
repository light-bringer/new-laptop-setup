#!/usr/bin/env bash
# -*- mode: bash -*-

set -eo pipefail
[[ -n $DEBUG ]] && set -x

: "${LAPTOP_SETUP_DIR:="$HOME/dev/new-laptop-setup"}"
export LAPTOP_SETUP_DIR

main() {
  (
    ensure_sudo_available
    ensure_directory_permissions
    ensure_homebrew
    setup_shell_aliases
  )
  ./bin/laptop.run "${@}"
}

ensure_directory_permissions() {
  ## workaround for recurring directory permission changes post-Sonoma upgrade
  sudo chmod 0755 /opt /private/etc
}

ensure_homebrew() {
  if ! command -v brew &>/dev/null; then
    echo 'Ensuring Homebrew is installed...'
    (
      ssh_config_fpath=~/.ssh/config
      if [ ! -f "${ssh_config_fpath}" ]; then
        mkdir -p "$(dirname "${ssh_config_fpath}")"
        touch "${ssh_config_fpath}"
        chmod 600 "${ssh_config_fpath}"
      fi

      cat <<'EOT' >> "${ssh_config_fpath}"
Host *
  StrictHostKeyChecking accept-new
EOT

      outputdir="$(mktemp -d -t homebrew-installer)"
      installer="${outputdir}/install.sh"
      curl \
        --output "${installer}" \
        --retry 3 \
        --fail \
        --silent \
        --show-error \
        --location \
        'https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh'
      NONINTERACTIVE=1 "${BASH}" "${installer}"
    )
  fi
}

ensure_sudo_available() {
  ## ensure sudo is fresh for unattended Homebrew installation
  if ! sudo -n true &>/dev/null; then
    echo >&2 'Prompting for sudo permissions...'
    sudo -v
  fi

  ## sudo refresh trick via https://github.com/geerlingguy/dotfiles/blob/8489a049/.osx#L32
  while true; do
    sudo -n true
    sleep 60
    kill -0 "$$" || exit
  done 2>/dev/null &
}

setup_shell_aliases() {
  local setup_dir="${LAPTOP_SETUP_DIR}"
  local zshenv="${HOME}/.zshenv"

  # Create .zshenv if it doesn't exist
  touch "${zshenv}"

  # Check if aliases already exist
  if grep -q "# Engineering Laptop Setup Aliases" "${zshenv}"; then
    echo "Shell aliases already installed"
    return
  fi

  # Add aliases to .zshenv
  cat >> "${zshenv}" <<EOF

# Engineering Laptop Setup Aliases
alias laptop.update='cd ${setup_dir} && git fetch origin && git reset --hard origin/main'
alias laptop.upgrade='laptop.update && cd ${setup_dir} && ./bin/laptop.run'
EOF

  echo "Shell aliases installed:"
  echo "  laptop.update  - Pull latest changes from repository"
  echo "  laptop.upgrade - Update and run full setup"
}

main "${@}"
