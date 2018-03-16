#!/bin/bash

# Vars
FOLDER_BASE="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FOLDER_BIN="${FOLDER_BASE}/bin"
FOLDER_ETC="${FOLDER_BASE}/etc"

# Base Functions
function exit_help {
  echo "${0}: Er00r: ${1}"
  echo
  echo "-a <application>"
  echo "-c <relative/path/to/config.conf>"
  echo
  echo "e.g: ${0} -a auth -c auth/sandbox.toml"
  echo
  exit 1
}

# Commandline Options
ARGV=`getopt -o a:c: -n "${0}" -- "$@"`

eval set -- "${ARGV}"
while true
do
  case "${1}" in
    -a) OPT_APP=${2} ; shift 2 ;;
    -c) OPT_CONFIG=${2} ; shift 2 ;;
    --) shift ; break ;;
    *) exit_help "GetOpt Broke!" ; exit 1 ;;
  esac
done

# SANITY
if [[ ! ${OPT_APP} || ! ${OPT_CONFIG} ]]
then
  exit_help "Required Vars not found!"
fi

if [[ ! -f "${FOLDER_BIN}/${OPT_APP}" ]]
then
  exit_help "App not found!"
fi

if [[ ! -f "${FOLDER_ETC}/${OPT_CONFIG}" ]]
then
  exit_help "Config not found!"
fi

if [[ ! -x "${FOLDER_BIN}/${OPT_APP}" ]]
then
  chmod +x ${FOLDER_BIN}/${OPT_APP}
fi

# Main
${FOLDER_BIN}/${OPT_APP} -configfile ${FOLDER_ETC}/${OPT_CONFIG}

# EOF
