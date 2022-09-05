#!/usr/bin/env bash

# shellcheck disable=SC2034

# Github gives branches as: refs/heads/master, and tags as refs/tags/v1.1.0.
# Remove prefix refs/{tags,heads} from the name before sending to Debricked.
refonly="${GITHUB_REF#refs/heads/}"
refonly="${refonly#refs/tags/}"

DEBRICKED_SCAN_REPOSITORY="${GITHUB_REPOSITORY}"
DEBRICKED_SCAN_COMMIT="${GITHUB_SHA}"
DEBRICKED_SCAN_BRANCH="${refonly}"
DEBRICKED_SCAN_REPOSITORY_URL="https://github.com/${GITHUB_REPOSITORY}"
DEBRICKED_SCAN_INTEGRATION=github
DEBRICKED_SCAN_PATH="."
if command -v git &> /dev/null
then
    DEBRICKED_SCAN_AUTHOR="$(git log -1 --pretty=%ae)"
fi