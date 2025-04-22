#!/bin/bash

# Use this script to re-generate the changelog from Github
# Requires you to be logged in via the Github CLI

# Ensure we're in the root of the git repository
if [ ! -d .git ]; then
    echo "Must be run from root of git repository"
    exit 1
fi

# Create new changelog file
echo "# Changelog" >CHANGELOG.md.new
echo "" >>CHANGELOG.md.new

git fetch origin --tags

# Get all tags sorted by version (assuming semantic versioning)
for tag in $(git tag -l | sort -rV); do
    echo "Processing tag: $tag"

    # Get the release notes from GitHub
    gh release view $tag --json body -q .body >release_notes.tmp

    if [ $? -eq 0 ]; then
        # Process Headers
        sed -i '' "s/## Changelog/## ${tag}/" release_notes.tmp
        sed -i '' "s/## What's Changed/## ${tag}/" release_notes.tmp
        sed -i '' "s/## Bug Fixes/### Bug Fixes/" release_notes.tmp
        sed -i '' "s/## New Features/### New Features/" release_notes.tmp
        sed -i '' "s/## New Contributors/### New Contributors/" release_notes.tmp
        # Remove more than one blank line
        sed -i '' '/^$/N;/^\n$/D' release_notes.tmp

        # Add the release notes to the changelog
        cat release_notes.tmp >>CHANGELOG.md.new
        echo -e "\n" >>CHANGELOG.md.new
    else
        echo "Warning: Could not fetch release notes for $tag"
    fi
done

# Clean up
rm release_notes.tmp

# Replace the existing CHANGELOG.md
mv CHANGELOG.md.new CHANGELOG.md
