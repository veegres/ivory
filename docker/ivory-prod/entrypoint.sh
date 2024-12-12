#!/bin/bash

# print service info
grep PRETTY_NAME /etc/os-release
pgcompacttable --version

# run go service
export GIN_MODE=release
export IVORY_STATIC_FILES_PATH="/opt/web"
/opt/service/ivory &

# wait for any process to exit
wait -n

# exit with status of process that exited first
exit $?
