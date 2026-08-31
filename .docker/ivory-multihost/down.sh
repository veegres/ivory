#!/usr/bin/env bash
# Removes the lab hosts and their network. The per-host image caches are kept
# unless --purge is passed, since re-pulling spilo three times is the slowest
# part of a re-run.
set -euo pipefail

NETWORK="${NETWORK:-ivory-mh}"
PREFIX="${PREFIX:-mh-host}"
COUNT="${COUNT:-3}"
PURGE=false
[ "${1:-}" = "--purge" ] && PURGE=true

for i in $(seq 1 "$COUNT"); do
  name="$PREFIX$i"
  docker rm -f "$name" >/dev/null 2>&1 && echo "$name removed" || true
  if $PURGE; then
    for volume in "$name-docker" "$name-ssh"; do
      docker volume rm "$volume" >/dev/null 2>&1 && echo "$volume volume removed" || true
    done
  fi
done

docker network rm "$NETWORK" >/dev/null 2>&1 && echo "network $NETWORK removed" || true
