#!/usr/bin/env bash

check_dependencies() {
  # Check if 'bun' command exists
  if ! command -v bun &> /dev/null; then
    echo "Error: 'bun' command not found."
    exit 1
  fi

  # Check if 'air' command exists
  if ! command -v air &> /dev/null; then
    echo "Error: 'air' command not found."
    exit 1
  fi
}

check_dependencies
