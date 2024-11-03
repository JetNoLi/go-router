#!/bin/bash
set -e

# STEP 1: Get Template Content

# Function to display usage instructions
usage() {
    echo "Usage: $0 <PROJECT_NAME> <MODULE_NAME>"
    exit 1
}

# Check if the correct number of arguments is provided
if [ "$#" -ne 2 ]; then
    echo "Error: Incorrect number of arguments."
    usage
fi

# Directory to extract to
DIR="repo"
PROJECT_NAME="$1"
MODULE_NAME="$2"

curl -L https://github.com/jetnoli/go-router/zipball/feat/improve-template-app/ -o repo.zip
unzip repo.zip -d $DIR
mkdir $PROJECT_NAME

# Find the extracted folder
# As we don't know the folder name
FIRST_DIR=$(find "$DIR" -mindepth 1 -maxdepth 1 -type d | head -n 1)

# Check if a directory was found
if [[ -z "$FIRST_DIR" ]]; then
    echo "No directories found in '$DIR'."
    exit 1
fi

# Move directory static contents to this folder
NEW_NAME="$DIR/contents/" 

# Rename the file
mv "$FIRST_DIR" "$NEW_NAME"

cp -r "$DIR/contents/grc/static/" "$PROJECT_NAME/"
rm repo.zip
rm -rf $DIR

# Step 2 Initialize Project and Install Templ
cd $PROJECT_NAME
go mod init $MODULE_NAME
go get github.com/a-h/templ