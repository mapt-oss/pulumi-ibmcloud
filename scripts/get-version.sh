#!/usr/bin/env bash
# Script to determine the provider version from git tags
# Returns a semantic version without the 'v' prefix

set -euo pipefail

# Function to get version from git
get_git_version() {
    # Try to get the exact tag for current commit (excluding sdk/ tags)
    EXACT_TAG=$(git describe --exact-match --tags HEAD 2>/dev/null | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+' || true)
    if [ -n "$EXACT_TAG" ]; then
        # Current commit has a version tag
        echo "$EXACT_TAG" | sed 's/^v//'
        return 0
    fi

    # Try to get the most recent version tag
    LATEST_TAG_WITH_V=$(git describe --tags --abbrev=0 --match 'v[0-9]*.[0-9]*.[0-9]*' HEAD 2>/dev/null || true)
    if [ -n "$LATEST_TAG_WITH_V" ]; then
        # Get the latest version tag and append commit info
        LATEST_TAG=$(echo "$LATEST_TAG_WITH_V" | sed 's/^v//')
        COMMIT_COUNT=$(git rev-list "${LATEST_TAG_WITH_V}"..HEAD --count)
        SHORT_SHA=$(git rev-parse --short=8 HEAD)

        if [ "$COMMIT_COUNT" -gt 0 ]; then
            echo "${LATEST_TAG}+dev.${COMMIT_COUNT}.g${SHORT_SHA}"
        else
            echo "${LATEST_TAG}"
        fi
        return 0
    fi

    # No version tags found, use default
    echo "0.0.0+dev.$(git rev-parse --short=8 HEAD 2>/dev/null || echo 'unknown')"
}

# Main execution
if [ -d .git ] || git rev-parse --git-dir > /dev/null 2>&1; then
    get_git_version
else
    # Not in a git repository, check for PROVIDER_VERSION env var
    if [ -n "${PROVIDER_VERSION:-}" ]; then
        echo "${PROVIDER_VERSION}" | sed 's/^v//'
    else
        echo "1.0.0-alpha.0+dev"
    fi
fi
