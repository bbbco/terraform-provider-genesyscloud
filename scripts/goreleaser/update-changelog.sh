#!/bin/bash
VERSION=$1
CHANGELOG_FILE="CHANGELOG.md"
TEMP_NEW_ENTRY="CHANGELOG.md.tmp"

# Ensure we're in the root of the git repository
if [ ! -d .git ]; then
    echo "Must be run from root of git repository"
    exit 1
fi

if [ ! -f $TEMP_NEW_ENTRY ]; then
    echo "Error: $TEMP_NEW_ENTRY does not exist. This script should not be executed directly. The GoReleaser tool should only call this script."
    exit 1
fi

if [ -f $CHANGELOG_FILE ]; then
    # If CHANGELOG.md exists, insert the new entry after the header
    sed -i.bak '2i\' $CHANGELOG_FILE # Add a blank line after header
    sed -i.bak "2r $TEMP_NEW_ENTRY" $CHANGELOG_FILE
    # Clean up backup files
    rm $CHANGELOG_FILE.bak
else
    # If CHANGELOG.md doesn't exist, create it with a header
    echo "# Changelog" >$CHANGELOG_FILE
    echo "" >>$CHANGELOG_FILE
    cat $TEMP_NEW_ENTRY >>$CHANGELOG_FILE
fi

# Clean up temp file
rm $TEMP_NEW_ENTRY
