#!/bin/bash

# Vars
FOLDER_BASE="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FOLDER_PROTO="${FOLDER_BASE}/proto"
FOLDER_SCRIPTS="${FOLDER_BASE}/scripts"

PROG_INSTALL_PROTOC_GEN_GO="${FOLDER_SCRIPTS}/installProtocGenGo.sh"
PROG_PROTOC=$(which protoc)
PROG_PROTOC_GEN_GO=$(which protoc-gen-go)

GOPATH_SRC="${GOPATH}/src"

PROTO_FOLDERS=()

# SANITY
if [[ -z "${GOPATH}" || ! -d "${GOPATH}" ]]
then
  echo "GOPATH must be set, failing..."
  exit 1
fi

if [[ ! -f ${PROG_PROTOC} || ! -x ${PROG_PROTOC} ]]
then
  echo "Failed to find protoc!"
  exit 1
fi

if [[ ! -f ${PROG_PROTOC_GEN_GO} || ! -x ${PROG_PROTOC_GEN_GO} ]]
then
  ${PROG_INSTALL_PROTOC_GEN_GO}
  if [[ $? != 0 ]]
  then
    echo "protoc-gen-go install failed!"
    exit 1
  fi
fi

# Main
for i in $(find ${FOLDER_PROTO} -type f -regex ".*\.proto$")
do
  PROTO_FOLDERS+=($(dirname ${i}))
done
PROTO_FOLDERS=($(echo "${PROTO_FOLDERS[@]}" | tr ' ' '\n' | sort -u | tr '\n' ' '))

for i in ${PROTO_FOLDERS[@]}
do
  echo "Processing: ${i##*/}"
  protoc -I=${i%/*} -I=${GOPATH_SRC} --go_out=${i%/*} ${i}/*.proto
  if [[ $? != 0 ]]
  then
    ERROR=1
  fi
done

if [[ ${ERROR} == 1 ]]
then
  echo "Something Failed..."
  exit 1
fi

# EOF
