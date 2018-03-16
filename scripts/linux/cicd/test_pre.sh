#!/bin/bash

# Vars
FOLDER_BASE="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
FOLDER_LIB="${FOLDER_BASE}/lib"
FOLDER_CICD="${FOLDER_BASE}/scripts/linux/cicd"

ERROR=0
TEST_FOLDERS=()

# Include defaults
source ${FOLDER_CICD}/defaults

# Main
for i in $(find ${FOLDER_LIB} -type f -iname "*_test.go")
do
  TEST_FOLDERS+=($(cd "$(dirname ${i})" && pwd))
done

TEST_FOLDERS=($(echo "${TEST_FOLDERS[@]}" | tr ' ' '\n' | sort -u | tr '\n' ' '))
for i in ${TEST_FOLDERS[@]}
do
  cd "${i}"
  go test
  if [[ $? != 0 ]]
  then
    ERROR=1
  fi
done

exit ${ERROR}

# EOF
