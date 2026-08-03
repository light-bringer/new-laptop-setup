#!/usr/bin/env bash
# -*- mode: bash -*-

set -eo pipefail
[[ -n $DEBUG ]] && set -x

: "${LAPTOP_SETUP_DIR:="$HOME/dev/new-laptop-setup"}"
export LAPTOP_SETUP_DIR

# shellcheck source=bin/lib/brew-helpers.sh disable=SC1091
source "$(dirname "${BASH_SOURCE[0]:-$0}")/bin/lib/brew-helpers.sh"

main() {
  (
    ensure_sudo_available
    ensure_directory_permissions
    ensure_homebrew
    add_bin_to_path
  )
  ./bin/laptop.run "${@}"
}

ensure_directory_permissions() {
  ## workaround for recurring directory permission changes post-Sonoma upgrade
  guarded_system_chmod 0755 /opt /private/etc
}

ensure_homebrew() {
  if command -v brew &>/dev/null; then
    return
  fi

  echo 'Ensuring Homebrew is installed...'

  if [[ "$(detect_account_type)" == "admin" ]]; then
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
    return
  fi

  check_command_line_tools || exit 1
  preflight_capability_check "$(detect_homebrew_prefix)" || exit 1
  install_homebrew_official || exit 1
}

ensure_sudo_available() {
  ## ensure sudo is fresh for unattended Homebrew installation
  if ! sudo -n true &>/dev/null; then
    echo >&2 'Prompting for sudo permissions...'
    local rc=0
    sudo -v || rc=$?
    if [[ ${rc} -ne 0 ]]; then
      echo >&2 "sudo -v failed (exit ${rc}). Possible causes: (a) this account does not have sudo access at all - contact your MDM administrator to grant sudo rights; (b) no interactive terminal is available for the password prompt - try running this directly in Terminal.app or iTerm rather than through another wrapper; (c) an incorrect password was entered."

      local detected_pam_tool
      detected_pam_tool="$(detect_sudo_pam_interception)"
      if [[ -n "${detected_pam_tool}" ]]; then
        echo >&2 "Detected: ${detected_pam_tool} is intercepting sudo on this machine (found in /etc/pam.d/sudo)."
        echo >&2 "This is a third-party security policy layer separate from standard sudoers. Your"
        echo >&2 "organization's security/IT team needs to update the ${detected_pam_tool} policy to authorize"
        echo >&2 "sudo for this account (scoped to mkdir/chown/chgrp/installer at minimum for Homebrew,"
        echo >&2 "or broader temporary elevation for the install). This is NOT something this script"
        echo >&2 "can resolve - it requires a policy change on the security tool's management console."
      fi

      exit "${rc}"
    fi
  fi

  ## sudo refresh trick via https://github.com/geerlingguy/dotfiles/blob/8489a049/.osx#L32
  while true; do
    sudo -n true
    sleep 60
    kill -0 "$$" || exit
  done 2>/dev/null &
  SUDO_KEEPALIVE_PID=$!
  trap 'kill "${SUDO_KEEPALIVE_PID}" 2>/dev/null' EXIT
}

add_bin_to_path() {
  local setup_dir="${LAPTOP_SETUP_DIR}"
  local zshenv="${HOME}/.zshenv"

  # Create .zshenv if it doesn't exist
  touch "${zshenv}"

  # Check if PATH is already configured
  if grep -q "# Engineering Laptop Setup PATH" "${zshenv}"; then
    echo "PATH already configured"
    return
  fi

  # Add bin directory to PATH
  cat >> "${zshenv}" <<EOF

# Engineering Laptop Setup PATH
export PATH="${setup_dir}/bin:\${PATH}"
EOF

  echo "Laptop setup commands installed:"
  echo "  laptop.update  - Pull latest changes from repository"
  echo "  laptop.upgrade - Update and run full setup"
  echo ""
  echo "These commands will be available in new shell sessions."
  echo "For this session, run: export PATH=\"${setup_dir}/bin:\${PATH}\""
}

main "${@}"
