#!/bin/bash

# Vars
SCRIPT_FOLDER="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

PROG_AWK=$(which awk)
PROG_GOVENDOR=$(which govendor)

PACKAGES_MISSING=0
PACKAGES_EXTERNAL=0
PACKAGES_UNUSED=0
PACKAGES_VENDOR=0

# SANITY
if [[ ! -x ${PROG_GOVENDOR} ]]
then
  echo "Please check govendor is installed and pathed correctly"
  exit 1
fi

# Main
cd ${SCRIPT_FOLDER}/..

for PACKAGE_STATE in $(${PROG_GOVENDOR} list | ${PROG_AWK} '{print $1}')
do
  if [[ ${PACKAGE_STATE} =~ ^e$ ]]
  then
    let PACKAGES_EXTERNAL+=1
  fi

  if [[ ${PACKAGE_STATE} =~ ^m$ ]]
  then
    let PACKAGES_MISSING+=1
  fi

if [[ ${PACKAGE_STATE} =~ ^u$ ]]
  then
    let PACKAGES_UNUSED+=1
  fi

  if [[ ${PACKAGE_STATE} =~ ^v$ ]]
  then
    let PACKAGES_VENDOR+=1
  fi
done

echo "Found ${PACKAGES_EXTERNAL} external packages"
echo "Found ${PACKAGES_MISSING} missing packages"
echo "Found ${PACKAGES_UNUSED} unused packages"
echo "Found ${PACKAGES_VENDOR} vendored packages"

# EOF
