#!/usr/bin/env bash
set -e

cd "$(dirname "$0")/.."

echo "Installing workspace dependencies..."
pnpm install

echo "StudyBuddy workspace ready."
