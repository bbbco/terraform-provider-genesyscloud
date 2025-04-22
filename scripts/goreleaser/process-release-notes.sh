#!/bin/bash
VERSION=$1

# Ensure we're in the root of the git repository
if [ ! -d .git ]; then
    echo "Must be run from root of git repository"
    exit 1
fi

# If next.md exists, check for actual content
if [ -f docs/releases/next.md ]; then
    # Create a temporary file with only content after "## Additional Release Information"
    sed -n '/^## Additional Release Information/,$p' docs/releases/next.md >temp.md

    # Check if there's any content beyond the header
    if [ "$(cat temp.md | wc -l)" -gt 1 ]; then
        # Copy the original file (with any formatting) to version-specific file
        cp docs/releases/next.md "docs/releases/${VERSION}.md"

        # Reset next.md with template
        cat >docs/releases/next.md <<'EOF'
<!--
This file is for important supplementary information that needs to be included in the next release.
This is not meant to duplicate information from PR descriptions or commit messages.
Examples might include:
- Breaking changes that need extra explanation
- Migration guides
- Complex upgrade instructions
- Important dependency changes that need special attention
-->

## Additional Release Information

EOF
    fi

    # Clean up temporary file
    rm temp.md
fi
