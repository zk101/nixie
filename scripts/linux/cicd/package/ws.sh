#!/bin/bash

# Packaging Script

# Vars
PROG_ECHO=$(which echo)
PROG_GETOPT=$(which getopt)

# Base Functions
function exit_help {
  ${PROG_ECHO} "${0}: Er00r: ${1}"
  ${PROG_ECHO}
  ${PROG_ECHO} "-d <dst_root>"
  ${PROG_ECHO} "-s <src_root>"
  ${PROG_ECHO}
  ${PROG_ECHO} "e.g: ${0} -d /path/to/dst -s /path/to/dst"
  ${PROG_ECHO}
  exit 1
}

# Commandline Options
ARGV=`${PROG_GETOPT} -o d:s: -n "${0}" -- "$@"`

eval set -- "${ARGV}"
while true
do
  case "${1}" in
    -d) OPT_DST_ROOT=${2} ; shift 2 ;;
    -s) OPT_SRC_ROOT=${2} ; shift 2 ;;
    --) shift ; break ;;
    *) exit_help "GetOpt Broke!" ; exit 1 ;;
  esac
done

# SANITY
if [[ ! ${OPT_DST_ROOT} || ! ${OPT_SRC_ROOT} ]]
then
  exit_help "Required Vars not found!"
fi

if [[ ! -d ${OPT_DST_ROOT} ]]
then
  exit_help "Dst Root does not exist!"
fi

if [[ ! -d ${OPT_SRC_ROOT} ]]
then
  exit_help "Src Root does not exist!"
fi

# Main
mkdir -p ${OPT_DST_ROOT}/bin
mkdir -p ${OPT_DST_ROOT}/etc
rsync -avP ${OPT_SRC_ROOT}/bin/ws ${OPT_DST_ROOT}/bin/
rsync -avP ${OPT_SRC_ROOT}/etc/ws/ ${OPT_DST_ROOT}/etc/

# EOF
