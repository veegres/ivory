#!/bin/bash

# print service info
grep PRETTY_NAME /etc/os-release
pgcompacttable --version
nginx -v

# run go service
export GIN_MODE=release
/opt/service/ivory &

# insert base url
envsubst '\$BASE_URL' < /etc/nginx/nginx.conf > "/etc/nginx/nginx.conf.tmp" && mv "/etc/nginx/nginx.conf.tmp" "/etc/nginx/nginx.conf"
if [ -n "${BASE_URL}" ]; then
    sed -i "s/location \//location \\${BASE_URL}/g" /etc/nginx/nginx.conf
    sed -i "s/api/\/api/g" /etc/nginx/nginx.conf
fi

# run nginx server
/usr/sbin/nginx &

# wait for any process to exit
wait -n

# exit with status of process that exited first
exit $?
