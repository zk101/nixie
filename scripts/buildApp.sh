#!/bin/bash

# Vars
FOLDER_BASE="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FOLDER_APP="${FOLDER_BASE}/app"
FOLDER_BIN="${FOLDER_BASE}/bin"

# Base Functions
function exit_help {
  echo "${0}: Er00r: ${1}"
  echo
  echo "-a <application>"
  echo
  echo "e.g: ${0} -a auth"
  echo
  exit 1
}

# Commandline Options
ARGV=`getopt -o a: -n "${0}" -- "$@"`

eval set -- "${ARGV}"
while true
do
  case "${1}" in
    -a) OPT_APP=${2} ; shift 2 ;;
    --) shift ; break ;;
    *) exit_help "GetOpt Broke!" ; exit 1 ;;
  esac
done

# SANITY
if [[ ! ${OPT_APP} ]]
then
  exit_help "Required Vars not found!"
fi

if [[ ! -d ${FOLDER_BIN} ]]
then
  mkdir -p ${FOLDER_BIN}
fi

# Main
APP_FOLDER="${FOLDER_APP}/${OPT_APP}"
APP_MAIN="${APP_FOLDER}/main.go"

if [[ ! -f "${APP_MAIN}" ]]
then
  exit_help "Application not found!"
fi

cd ${APP_FOLDER}
echo "Building ${OPT_APP} in ${APP_FOLDER}..."
BUILD_DATE=$(date -u '+%Y-%m-%d %H:%M:%S UTC')
BUILD_REVISION="Mercurial $(hg identify)"
BUILD_TYPE="Console"
CGO_ENABLED=0 go build -o ${FOLDER_BIN}/${OPT_APP} -ldflags "-w -s -extldflags '-static' -X 'github.com/zk101/nixie/lib/buildinfo.BuildStamp=${BUILD_DATE}' -X 'github.com/zk101/nixie/lib/buildinfo.BuildRevision=${BUILD_REVISION}' -X 'github.com/zk101/nixie/lib/buildinfo.BuildType=${BUILD_TYPE}'"

# EOF
