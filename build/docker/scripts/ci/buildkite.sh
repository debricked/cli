#!/usr/bin/env bash

# shellcheck disable=SC2034

# the repository is determined according to the following rules:
# 1. If BUILDKITE_REPO starts with "http(s)://" and ends with ".git", use capture group to set REPOSITORY.
# 2. If BUILDKITE_REPO starts with "git@" and ends with ".git", use capture group to set REPOSITORY.
# 3. Set REPOSITORY to BUILDKITE_REPO.
http_name_regex="^https?:\/\/.+\.[a-z0-9]+\/(.+)\.git$"
ssh_name_regex="^.*:[0-9]*\/*(.+)\.git$"
if [[ $BUILDKITE_REPO =~ $http_name_regex ]]; then
    DEBRICKED_SCAN_REPOSITORY="${BASH_REMATCH[1]}"
elif [[ $BUILDKITE_REPO =~ $ssh_name_regex ]]; then
    DEBRICKED_SCAN_REPOSITORY="${BASH_REMATCH[1]}"
else
    DEBRICKED_SCAN_REPOSITORY="${BUILDKITE_REPO}"
fi

DEBRICKED_SCAN_COMMIT="${BUILDKITE_COMMIT}"
DEBRICKED_SCAN_BRANCH="${BUILDKITE_BRANCH}"
DEBRICKED_SCAN_INTEGRATION=buildkite

if command -v git &> /dev/null
then
    DEBRICKED_SCAN_AUTHOR="$(git log -1 --pretty=%ae)"
fi

# the repository url is determined according to the following rules:
# 1. If DEBRICKED_SCAN_REPOSITORY_URL is set, always use it as the repo url.
# 2. If BUILDKITE_REPO starts with "http(s)://" and ends with ".git", use capture group to set REPOSITORY_URL.
# 3. If BUILDKITE_REPO is of the form "git@github.com:organisation/reponame.git",
#    rewrite and use "https://github.com/organisation/reponame" as REPOSITORY_URL.
# 4. Otherwise, show warning and set repository url to ""
if [[ -z "$DEBRICKED_SCAN_REPOSITORY_URL" ]]; then
    http_url_regex="^(https?:\/\/.+)\.git$"
    ssh_url_regex="git@(.+):[0-9]*\/?(.+)\.git$"
    if [[ $BUILDKITE_REPO =~ $http_url_regex ]]; then
        DEBRICKED_SCAN_REPOSITORY_URL="${BASH_REMATCH[1]}"
    elif [[ $BUILDKITE_REPO =~ $ssh_url_regex ]]; then
        domain="${BASH_REMATCH[1]}"
        uri="${BASH_REMATCH[2]}"
        DEBRICKED_SCAN_REPOSITORY_URL="https://${domain}/${uri}"
    else
        echo "INFO: Your repository URL could not be found. Set it manually with DEBRICKED_REPOSITORY_URL"
        DEBRICKED_SCAN_REPOSITORY_URL=""
    fi
fi
