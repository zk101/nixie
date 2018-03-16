#!/bin/bash

# Vars
FOLDER_BASE="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FOLDER_CICD="${FOLDER_BASE}/scripts/linux/cicd"

BUILD_SCRIPTS=("deps.sh" "test_pre.sh" "build.sh" "test_post.sh")

# Main
for SCRIPT in ${BUILD_SCRIPTS[@]}
do
  echo "Running ${SCRIPT}..."
  ${FOLDER_CICD}/${SCRIPT}
  if [[ $? != 0 ]]
  then
    echo "Build Failed..."
    exit 1
  fi
done

# EOF
