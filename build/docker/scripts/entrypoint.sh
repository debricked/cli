#!/usr/bin/env bash

# shellcheck disable=SC1091
source /setup_env.sh

debricked scan "${DEBRICKED_SCAN_PATH}" \
  --access-token "${DEBRICKED_TOKEN}" \
  --repository "${DEBRICKED_SCAN_REPOSITORY}" \
  --commit "${DEBRICKED_SCAN_COMMIT}" \
  --branch "${DEBRICKED_SCAN_BRANCH}" \
  --author "${DEBRICKED_SCAN_AUTHOR}" \
  --repository-url "${DEBRICKED_SCAN_REPOSITORY_URL}" \
  --integration "${DEBRICKED_SCAN_INTEGRATION}"
