#!/bin/bash

# Create a temporary directory with a unique timestamp
TMP_DIR="tmp_repo_$(date +%s)"
TARGET_DIR="pb_data" # The folder where you want to extract the content
REPO_URL="https://github.com/supotsu-no-ochaya/backend.git"
BRANCH="main" # Branch name
FOLDER_TO_COPY="testdata/$1" # The folder to extract

# Clone the repository into the temporary directory
git clone --depth 1 --branch "$BRANCH" "$REPO_URL" "$TMP_DIR"

# Create the target directory if it doesn't exist
mkdir -p "$TARGET_DIR"

# Copy the specific folder to the target directory
cp -r "$TMP_DIR/$FOLDER_TO_COPY/pb_data"/* "$TARGET_DIR"

# Clean up the temporary directory
rm -rf "$TMP_DIR"

echo "Folder '$FOLDER_TO_COPY' extracted to '$TARGET_DIR' and temporary files cleaned up."
