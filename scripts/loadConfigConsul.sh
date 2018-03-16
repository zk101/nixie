#!/bin/bash

# Vars
FOLDER_BASE="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FOLDER_BIN="${FOLDER_BASE}/bin"
FOLDER_ETC="${FOLDER_BASE}/etc"

PROG_CONFIGLOAD="${FOLDER_BIN}/configload"

CONSUL_HOST="localhost:8500"

APP_LIST=("auth" "chat" "loadtest" "telemetry" "ws")

# Expect one argument
if [[ -z ${1} ]]; then
  echo "Expect an Environment name"
  exit 1
fi

# Main
for APP in ${APP_LIST[@]}
do
  echo "Loading ${APP}..."

  if [[ ! -f "${FOLDER_ETC}/${APP}/${1}.toml" ]]; then
    echo "Unsupported Environment for ${APP}, skipping..."
    continue
  fi

  ${PROG_CONFIGLOAD} -configtype ${APP} -consuladdr ${CONSUL_HOST} \
    -consulprefix "nixie/${1}/${APP}/v1" -configfile "${FOLDER_ETC}/${APP}/${1}.toml"
done

# EOF
