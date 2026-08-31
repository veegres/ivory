#!/usr/bin/env bash
# Brings up three docker-in-docker hosts on 10.0.0.1-3, each with its own
# docker daemon and an sshd holding an Ivory vault's public key - the lab the
# shipped multi-host deployment templates are written against. Their example
# addresses are 10.0.0.1-3, so every template deploys here verbatim.
set -euo pipefail

IVORY_URL="${IVORY_URL:-http://localhost:8080}"
IVORY_VAULT_ID="${IVORY_VAULT_ID:-}"
NETWORK="${NETWORK:-ivory-mh}"
SUBNET="${SUBNET:-10.0.0.0/24}"
GATEWAY="${GATEWAY:-10.0.0.254}"
PREFIX="${PREFIX:-mh-host}"
COUNT="${COUNT:-3}"
IMAGE="${IMAGE:-ivory-multihost:latest}"
BASE="${SUBNET%.*}"

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

if docker network inspect "$NETWORK" >/dev/null 2>&1; then
  echo "network $NETWORK exists"
else
  # The gateway is moved off .1 so the first host can hold the address the
  # templates name.
  docker network create --subnet "$SUBNET" --gateway "$GATEWAY" "$NETWORK" >/dev/null
  echo "network $NETWORK created ($SUBNET, gateway $GATEWAY)"
fi

read -r VAULT_ID SSH_USER SSH_KEY <<<"$(
  curl -sf -H "Cookie: session=$(cat /proc/sys/kernel/random/uuid)" "$IVORY_URL/api/vault" |
  IVORY_VAULT_ID="$IVORY_VAULT_ID" python3 -c '
import json, os, sys
vaults = json.load(sys.stdin)["response"]
wanted = os.environ.get("IVORY_VAULT_ID") or ""
keys = {k: v for k, v in vaults.items() if v.get("type") == 3 and v.get("metadata")}
if wanted:
    keys = {k: v for k, v in keys.items() if k == wanted}
if not keys:
    sys.exit("no SSH_KEY vault with a public key found - create one in Ivory first")
vid, vault = next(iter(sorted(keys.items())))
print(vid, vault["username"], vault["metadata"])
'
)"
echo "vault $VAULT_ID (user $SSH_USER)"

docker build -q -t "$IMAGE" "$DIR" >/dev/null
echo "image $IMAGE built"

for i in $(seq 1 "$COUNT"); do
  name="$PREFIX$i"
  ip="$BASE.$i"
  if docker inspect "$name" >/dev/null 2>&1; then
    echo "$name exists, skipping"
    continue
  fi
  # The docker volume keeps each host's image cache across a teardown, which is
  # the difference between a re-run and pulling spilo three times again. The
  # ssh one keeps its host keys: Ivory trusts a host key on first use and holds
  # it in memory for the life of the server process, so a host that comes back
  # with a fresh identity fails every later connection with a host key mismatch
  # until Ivory itself is restarted.
  docker run -d --privileged \
    --name "$name" --hostname "$name" \
    --network "$NETWORK" --ip "$ip" \
    -v "$name-docker:/var/lib/docker" \
    -v "$name-ssh:/etc/ssh" \
    -e SSH_USER="$SSH_USER" \
    -e SSH_PUBLIC_KEY="$SSH_KEY" \
    "$IMAGE" >/dev/null
  echo "$name started on $ip"
done

for i in $(seq 1 "$COUNT"); do
  name="$PREFIX$i"
  until docker exec "$name" docker info >/dev/null 2>&1; do sleep 1; done
done

echo
echo "hosts ready - deploy a multi-host template with sshPort 22 and vault $VAULT_ID:"
for i in $(seq 1 "$COUNT"); do
  echo "  $BASE.$i  ($PREFIX$i)"
done
