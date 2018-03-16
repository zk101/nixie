#!/bin/bash

# Vars
FOLDER_BASE="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
FOLDER_APP="${FOLDER_BASE}/app"
FOLDER_BIN="${FOLDER_BASE}/bin"
FOLDER_CICD="${FOLDER_BASE}/scripts/linux/cicd"
FOLDER_SCRIPTS="${FOLDER_BASE}/scripts"
FOLDER_TOOLS="${FOLDER_BASE}/tools"

PROG_AWK=$(which awk)
PROG_INSTALL_GOVENDOR="${FOLDER_SCRIPTS}/installGovendor.sh"
PROG_GENPROTOBUF="${FOLDER_SCRIPTS}/genProtoBuf.sh"
PROG_GOVENDOR=$(which govendor)
PROG_GREP=$(which grep)

GOVENDOR_MISSING_IGNORE=()
GOVENDOR_MISSING_ERROR=0

# Include defaults
source ${FOLDER_CICD}/defaults

# SANITY
if [[ ! -f ${PROG_GENPROTOBUF} ]]
then
  echo "genProtoBuf script not found!"
  exit 1
fi

if [[ ! -x ${PROG_GENPROTOBUF} ]]
then
  chmod +x ${PROG_GENPROTOBUF}
fi

if [[ ! -f ${PROG_GOVENDOR} ]]
then
  ${PROG_INSTALL_GOVENDOR}
  if [[ $? != 0 ]]
  then
    echo "govendor install failed!"
    exit 1
  fi
  PROG_GOVENDOR=$(which govendor)
fi

if [[ ! -x ${PROG_GOVENDOR} ]]
then
  chmod +x ${PROG_GOVENDOR}
fi

# Main
if [[ ! -e "/usr/local/bin/protoc" ]]; then
  cp -f ${FOLDER_TOOLS}/linux/protoc/bin/protoc /usr/local/bin
  cp -fr ${FOLDER_TOOLS}/linux/protoc/include/* /usr/local/include/
  chmod +x ${FOLDER_TOOLS}/linux/protoc/bin/protoc
fi

${PROG_GENPROTOBUF}
if [[ $? != 0 ]]
then
  exit 1
fi

for PACKAGE in $(cd ${FOLDER_BASE}; ${PROG_GOVENDOR} list | ${PROG_GREP} " m " | ${PROG_AWK} '{print $2}')
do
  for IGNORE in ${GOVENDOR_MISSING_IGNORE[@]}
  do
    if [[ ${IGNORE} == ${PACKAGE} ]]
    then
      continue
    fi
    echo "Missing Package: ${PACKAGE}"
    GOVENDOR_MISSING_ERROR=1
  done
done

if [[ ${GOVENDOR_MISSING_ERROR} != 0 ]]
then
  echo "Failed!"
  exit 1
fi

# EOF
