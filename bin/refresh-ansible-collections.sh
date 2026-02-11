#!/usr/bin/env bash
# -*- mode: bash -*-
#
# Vendor Ansible collections and roles locally for offline availability

set -eo pipefail
[[ -n $DEBUG ]] && set -x

main() {
  repo_dir="$(dirname "${BASH_SOURCE[0]}")/.."
  cd "$repo_dir"

  echo "Installing Ansible collections and roles..."
  ansible-galaxy install -r requirements.yml --force

  echo ""
  echo "Dependencies installed successfully!"
  echo ""
  echo "Collections installed to: ./ansible_collections/"
  echo "Roles installed to: ./roles/"
  echo ""
  echo "These are vendored locally and committed to the repository."
}

main "${@}"
