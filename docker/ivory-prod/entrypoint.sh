#!/bin/bash

# make ivory directory as main for service
cd /opt/ivory || exit
echo "Working directory: $(pwd)"

# print service info
. /etc/os-release
echo "$PRETTY_NAME"
pgcompacttable --version

# run go service
export GIN_MODE=release
export IVORY_STATIC_FILES_PATH="./web"
./service/ivory &

# wait for any process to exit
wait -n

# exit with status of process that exited first
exit $?
