#!/usr/bin/env bash

# print service info
grep PRETTY_NAME /etc/os-release
nginx -v

# run go service
export GIN_MODE=release
/opt/service/ivory &

# run nginx server
/usr/sbin/nginx
