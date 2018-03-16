#!/bin/bash

# Vars
FOLDER_BASE="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
FOLDER_SCRIPTS="${FOLDER_BASE}/scripts"

PROG_INSTALL_GOVENDOR="${FOLDER_SCRIPTS}/installGovendor.sh"
PROG_GOVENDOR=$(which govendor)

# SANITY
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
cd ${FOLDER_BASE}

${PROG_GOVENDOR} fetch github.com/BurntSushi/toml@v0.3.0
${PROG_GOVENDOR} fetch github.com/gobwas/ws
${PROG_GOVENDOR} fetch github.com/gobwas/ws/wsutil
${PROG_GOVENDOR} fetch github.com/golang/protobuf/proto
${PROG_GOVENDOR} fetch github.com/gorilla/websocket@v1.2.0
${PROG_GOVENDOR} fetch github.com/hashicorp/consul/api@v1.0.2
${PROG_GOVENDOR} fetch github.com/mailru/easygo/netpoll
${PROG_GOVENDOR} fetch github.com/mavricknz/ldap
${PROG_GOVENDOR} fetch github.com/prometheus/client_golang/prometheus@v0.8.0
${PROG_GOVENDOR} fetch github.com/prometheus/client_golang/prometheus/promhttp@v0.8.0
${PROG_GOVENDOR} fetch github.com/satori/go.uuid@v1.1.0
${PROG_GOVENDOR} fetch golang.org/x/crypto/sha3
#${PROG_GOVENDOR} fetch gopkg.in/couchbase/gocb.v1
${PROG_GOVENDOR} fetch github.com/couchbase/gocb@v1.3.3
${PROG_GOVENDOR} fetch go.uber.org/zap
${PROG_GOVENDOR} fetch gopkg.in/natefinch/lumberjack.v2
${PROG_GOVENDOR} fetch github.com/streadway/amqp
${PROG_GOVENDOR} fetch gopkg.in/yaml.v2
${PROG_GOVENDOR} fetch github.com/go-sql-driver/mysql@v1.3

echo "Check govendor list as there may be external packages to add to support the primary packages added via this script"

# EOF
