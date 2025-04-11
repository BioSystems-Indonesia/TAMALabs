#!/usr/bin/env bash

# Script to run atlas migrate diff with a required argument

# Check if exactly one argument (the migration name) was provided.
# If you want to allow names with spaces like "add new feature",
# you might want to use "$@" to capture all arguments instead of just $1.
# Let's assume for now the name is a single word or quoted properly.
if [ $# -eq 0 ]; then
  # Print error message to standard error (>&2)
  echo "Error: Missing migration name." >&2
  echo "Usage: ./migrate-diff.sh <migration_name>" >&2
  exit 1 # Exit with a non-zero status code to indicate failure
fi

# Get the migration name from the first argument
MIGRATION_NAME="$1"
# If you wanted to capture ALL arguments (e.g., "add new column"):
# MIGRATION_NAME="$@"

# Echo the command that will be run (optional, but helpful for debugging)
echo "Running: atlas migrate diff --env gorm ${MIGRATION_NAME}"

# Execute the actual command
atlas migrate diff --env gorm ${MIGRATION_NAME}

# Optional: Check if the atlas command succeeded
ATLAS_EXIT_CODE=$?
if [ ${ATLAS_EXIT_CODE} -ne 0 ]; then
    echo "Error: Atlas command failed with exit code ${ATLAS_EXIT_CODE}" >&2
    exit ${ATLAS_EXIT_CODE} # Exit with the same error code as atlas
fi

echo "Atlas migrate diff completed."
exit 0 # Explicitly exit with success
