#!/usr/bin/env bash

# shellcheck disable=SC2034

if [[ ! -f /ci_runned ]] ; then
    DEBRICKED_SCAN_REPOSITORY="${CI_PROJECT_PATH}"
    DEBRICKED_SCAN_COMMIT="${CI_COMMIT_SHA}"
    DEBRICKED_SCAN_BRANCH="${CI_COMMIT_REF_NAME}"
#    if [[ ! -z "${CI_DEFAULT_BRANCH}" ]]; then
#      DEBRICKED_SCAN_DEFAULT_BRANCH="${CI_DEFAULT_BRANCH}"
#    else
#      echo -e "You are probably using a version of gitlab before 12.4. This means we can not know what your default branch is. This might impact your experience using Debricked's tools"
#    fi
    DEBRICKED_SCAN_PATH="${BASE_DIRECTORY:=$CI_PROJECT_DIR}"
    DEBRICKED_SCAN_REPOSITORY_URL="${CI_PROJECT_URL}"
    DEBRICKED_SCAN_INTEGRATION=gitlab
    if [[ -n "${CI_COMMIT_AUTHOR}" ]]; then
      DEBRICKED_SCAN_AUTHOR="${CI_COMMIT_AUTHOR}"
    elif command -v git &> /dev/null
    then
      DEBRICKED_SCAN_AUTHOR="$(git log -1 --pretty=%ae)"
    fi

    #[[ $CI ]] && touch /ci_runned
fi