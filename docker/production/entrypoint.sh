#!/bin/sh

# print service info
grep PRETTY_NAME /etc/os-release
nginx -v

# run go service
export GIN_MODE=release
/opt/service/ivory &

# run nginx server
/usr/sbin/nginx &

# wait for any process to exit
# shellcheck disable=SC2039
wait -n

# exit with status of process that exited first
exit $?
