#!/usr/bin/env bash
# -*- mode: bash -*-
#
# brew-helpers.sh - Standard (non-admin) macOS account support for Homebrew setup.
#
# Rationale:
#   - We keep Homebrew at its DEFAULT prefix (/opt/homebrew on Apple Silicon,
#     /usr/local on Intel) so users benefit from prebuilt bottles. Anything
#     off the default prefix forces building from source and loses CI cache
#     hits upstream.
#   - The official Homebrew installer (install.sh) only requires `sudo`
#     access to create/chown the prefix directory. It does NOT require the
#     current user to be a member of the macOS `admin` group.
#   - Therefore this library never invokes any command that would add or
#     modify group membership. We only ever probe read-only
#     account/group state and rely on `sudo` for the one-time prefix setup,
#     matching how a standard (non-admin) account with sudo rights granted
#     via MDM/Jamf can still bootstrap Homebrew.
#
# Intended usage: `source` this file from bootstrap.sh and bin/laptop.run.

# Double-source guard.
[[ -n "${_BREW_HELPERS_LOADED:-}" ]] && return
_BREW_HELPERS_LOADED=1

# detect_account_type
#   Read-only check of admin-group membership. Never mutates groups.
#   Echoes exactly one word: "admin" or "standard".
detect_account_type() {
  local user_groups
  user_groups="$(id -Gn "$USER" 2>/dev/null || groups 2>/dev/null)"

  if [[ " ${user_groups} " == *" admin "* ]]; then
    # Cross-verify against the authoritative admin group membership list.
    if dscl . -read /Groups/admin GroupMembership 2>/dev/null | grep -qw "$USER"; then
      echo "admin"
      return 0
    fi
    # id/groups says admin but dscl disagrees (e.g. stubbed/edge case) -
    # trust the read-only id/groups signal as the primary source.
    echo "admin"
    return 0
  fi

  echo "standard"
}

# detect_homebrew_prefix
#   Echoes the resolved (or inferred) Homebrew prefix.
detect_homebrew_prefix() {
  if command -v brew >/dev/null 2>&1; then
    brew --prefix
    return 0
  fi

  if [[ "$(uname -m)" == "arm64" ]]; then
    echo "/opt/homebrew"
  else
    echo "/usr/local"
  fi
}

# check_command_line_tools
#   Verifies Xcode Command Line Tools are installed. Does NOT auto-invoke
#   `xcode-select --install` (that opens a GUI prompt which would hang in
#   unattended/MDM contexts).
check_command_line_tools() {
  if ! xcode-select -p >/dev/null 2>&1; then
    echo >&2 "Command Line Tools not found. On an MDM-managed machine install via Self Service or run 'xcode-select --install', then re-run this script."
    return 1
  fi
  return 0
}

# preflight_capability_check <prefix>
#   Verifies the machine can host Homebrew at the given prefix without
#   actually installing anything.
preflight_capability_check() {
  local prefix="$1"
  local err_file="/tmp/brew_preflight.err"

  # _preflight_check_sudo_authorization
  #   Mirrors Homebrew's official installer's own have_sudo_access() check:
  #   `sudo -l mkdir` is an AUTHORIZATION query ("would sudo permit me to run
  #   mkdir"), not just "do I have a cached credential". MDM policies that
  #   scope sudo access to specific whitelisted commands can pass a plain
  #   `sudo mkdir -p` probe while still failing this stricter check, which is
  #   exactly what Homebrew's installer runs before aborting. We use `-n` to
  #   avoid an unexpected password prompt here, since credentials should
  #   already be cached from an earlier interactive `sudo -v`.
  _preflight_check_sudo_authorization() {
    if sudo -n -l mkdir &>/dev/null; then
      return 0
    fi

    echo >&2 "Pre-flight FAILED: your account has some sudo access, but is not authorized to run 'mkdir' via sudo (checked via 'sudo -l mkdir')."
    echo >&2 "This usually means your MDM has granted sudo access scoped to specific commands rather than broad/unrestricted access."
    echo >&2 "Homebrew's official installer requires broader sudo authorization (it runs mkdir, chown, and chgrp as root)."
    echo >&2 "Contact your IT/MDM administrator to request a broader sudo grant - command-scoped sudo access is not sufficient for the official Homebrew installer."
    return 1
  }

  if [[ -e "${prefix}" ]]; then
    local owner
    owner="$(stat -f '%Su' "${prefix}" 2>/dev/null)"
    if [[ "${owner}" == "$(whoami)" ]]; then
      _preflight_check_sudo_authorization || return 1
      echo "Pre-flight OK: '${prefix}' exists and is owned by $(whoami)."
      return 0
    fi

    _preflight_check_sudo_authorization || return 1
    echo "Pre-flight WARNING: '${prefix}' exists but is owned by '${owner}', not $(whoami). The installer will run 'sudo chown' to fix ownership."
    return 0
  fi

  # Prefix does not exist yet - test whether sudo can create it.
  if sudo mkdir -p "${prefix}" 2>"${err_file}"; then
    sudo rmdir "${prefix}" 2>/dev/null
    _preflight_check_sudo_authorization || return 1
    echo "Pre-flight OK: sudo can create '${prefix}'."
    return 0
  fi

  local captured_err=""
  [[ -f "${err_file}" ]] && captured_err="$(cat "${err_file}")"

  echo >&2 "Pre-flight FAILED: could not create '${prefix}'."
  [[ -n "${captured_err}" ]] && echo >&2 "${captured_err}"
  echo >&2 "sudo could not create ${prefix}. This looks like it is blocked by MDM policy. Contact IT or your MDM administrator."
  return 1
}

# detect_sudo_pam_interception [pam_sudo_path]
#   Read-only inspection of the sudo PAM config (default: /etc/pam.d/sudo,
#   world-readable on stock macOS, no sudo required) for known signatures of
#   third-party sudo-gatekeeping PAM modules (e.g. enterprise PAM/MDM security
#   tooling that intercepts sudo independently of standard sudoers). An
#   optional path argument is accepted so tests can point this at a temp file
#   instead of the real system file.
#   Echoes the detected tool name and returns 0 on match; returns 1 (no
#   output) if the file is missing/unreadable or no signature matches.
detect_sudo_pam_interception() {
  local pam_sudo_path="${1:-/etc/pam.d/sudo}"

  [[ -r "${pam_sudo_path}" ]] || return 1

  if grep -q 'CyberArkEPMPAM' "${pam_sudo_path}" 2>/dev/null; then
    echo "CyberArk EPM"
    return 0
  fi

  if grep -qE 'auth[[:space:]]+sufficient[[:space:]]+/private/cyberark/' "${pam_sudo_path}" 2>/dev/null; then
    echo "CyberArk EPM"
    return 0
  fi

  return 1
}

# guarded_system_chmod <mode> <path> [path...]
#   Only runs `sudo chmod` when the current permissions differ from the
#   requested mode. Missing paths are warned about and skipped (non-fatal).
guarded_system_chmod() {
  local mode="$1"
  shift

  local path
  for path in "$@"; do
    if [[ ! -e "${path}" ]]; then
      echo >&2 "Warning: '${path}' does not exist, skipping chmod."
      continue
    fi

    local current_perms
    current_perms="$(stat -f '%Lp' "${path}" 2>/dev/null)"

    # Compare as octal numbers so "755" (stat output) matches "0755" (requested mode).
    if [[ -n "${current_perms}" ]] && (( 8#${current_perms} == 8#${mode} )); then
      continue
    fi

    sudo chmod "${mode}" "${path}"
  done

  return 0
}

# install_homebrew_official
#   Runs the official Homebrew installer non-interactively.
install_homebrew_official() {
  local rc=0
  NONINTERACTIVE=1 /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)" || rc=$?

  if [[ ${rc} -ne 0 ]]; then
    echo >&2 "Homebrew installation failed (exit ${rc}). On an MDM-managed machine, this may be blocked by policy restricting sudo or /opt access. Contact IT or your MDM administrator."
    return 1
  fi

  return 0
}
