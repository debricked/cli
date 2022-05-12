#!/usr/bin/env bash

DEBRICKED_SCAN_PATH="."

if [[ -n "${TF_BUILD}" ]]; then # https://docs.microsoft.com/en-us/azure/devops/pipelines/build/variables?view=azure-devops&tabs=yaml
  echo "Integration: Azure DevOps"
  # shellcheck disable=SC1091
  source ci/azure.sh
elif [[ -n "${BITBUCKET_BUILD_NUMBER}" ]]; then # https://support.atlassian.com/bitbucket-cloud/docs/variables-and-secrets/
  echo "Integration: Bitbucket"
  # shellcheck disable=SC1091
  source ci/bitbucket.sh
elif [[ -n "${BUILDKITE}" ]]; then # https://buildkite.com/docs/pipelines/environment-variables#bk-env-vars-buildkite
  echo "Integration: Buildkite"
  # shellcheck disable=SC1091
  source ci/buildkite.sh
elif [[ -n "${CIRCLECI}" ]]; then # https://circleci.com/docs/2.0/variables/
  echo "Integration: CircleCI"
  # shellcheck disable=SC1091
  source ci/circleci.sh
elif [[ -n "${GITHUB_ACTION}" ]]; then # https://docs.github.com/en/actions/learn-github-actions/environment-variables
  echo "Integration: GitHub Actions"
  # shellcheck disable=SC1091
  source ci/github.sh
elif [[ -n "${GITLAB_CI}" ]]; then # https://docs.gitlab.com/ee/ci/variables/predefined_variables.html
  echo "Integration: GitLab CI/CD"
  # shellcheck disable=SC1091
  source ci/gitlab.sh
else
  echo "Integration: unknown"
fi

export DEBRICKED_SCAN_PATH
export DEBRICKED_SCAN_REPOSITORY
export DEBRICKED_SCAN_COMMIT
export DEBRICKED_SCAN_BRANCH
export DEBRICKED_SCAN_PATH
export DEBRICKED_SCAN_REPOSITORY_URL
export DEBRICKED_SCAN_INTEGRATION
export DEBRICKED_SCAN_AUTHOR

if [[ -n "${DEBRICKED_DEBUG}" ]]; then
  echo
  printenv | grep DEBRICKED_SCAN || true
fi
