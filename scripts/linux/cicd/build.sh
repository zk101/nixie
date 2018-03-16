#!/bin/bash

# Vars
FOLDER_BASE="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
FOLDER_APP="${FOLDER_BASE}/app"
FOLDER_BIN="${FOLDER_BASE}/bin"
FOLDER_CICD="${FOLDER_BASE}/scripts/linux/cicd"

ERROR=0

# Include defaults
source ${FOLDER_CICD}/defaults

# SANITY
if [[ ! -d ${FOLDER_BIN} ]]
then
  mkdir -p ${FOLDER_BIN}
fi

if [[ -z "${BUILD_REVISION}" ]]; then
  BUILD_REVISION="Unavailable"
fi

if [[ -z "${BUILD_NUMBER}" ]]; then
  BUILD_NUMBER="Unavailable"
fi

# Main
for i in $(find ${FOLDER_APP} -type f -iname "main.go")
do
  DIRNAME=$(cd "$(dirname $i)" && pwd)
  FILENAME=${DIRNAME##*/}
  echo "Building ${FILENAME}..."
  cd ${DIRNAME}
  BUILD_DATE=$(date -u '+%Y-%m-%d %H:%M:%S UTC')
  BUILD_REV="Mercurial ${BUILD_REVISION}"
  BUILD_TYPE="Jenkins (${BUILD_NUMBER})"
  CGO_ENABLED=0 go build -o ${FOLDER_BIN}/${FILENAME} -ldflags "-w -s -extldflags '-static' -X 'github.com/zk101/nixie/lib/buildinfo.BuildStamp=${BUILD_DATE}' -X 'github.com/zk101/nixie/lib/buildinfo.BuildRevision=${BUILD_REV}' -X 'github.com/zk101/nixie/lib/buildinfo.BuildType=${BUILD_TYPE}'"
  if [[ $? != 0 ]]
  then
    ERROR=1
  fi
done

exit ${ERROR}

# EOF
