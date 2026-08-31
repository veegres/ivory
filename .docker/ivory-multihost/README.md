## Multi-host lab for deployment templates

The shipped multi-host deployment templates put one node per VM and address their peers by literal
address — `10.0.0.1-3` in every one of them (etcd's member list, ZooKeeper's `ZOO_SERVERS`, Patroni's
`ETCD3_HOSTS`, ClickHouse's replica and keeper lists, Redis' `--replicaof`, Postgres' `pg_basebackup`
host, Mongo's `rs.initiate` member list). Their descriptions tell an operator to replace those with
their own VM addresses.

This lab is three VMs that already have those addresses, so every multi-host template deploys here
**exactly as shipped** — no member list edited, nothing substituted. Each host is a
docker-in-docker container with its **own docker daemon** and an sshd holding an Ivory vault's public
key, which is what makes it a separate deployment target rather than three containers on one daemon:
node names repeat per host, `--hostname` is legal, and a published port belongs to that host alone.

## Run

```
./up.sh          # build, network, three hosts on 10.0.0.1-3
./down.sh        # remove hosts and network, keep the image caches
./down.sh --purge  # also drop each host's image cache and ssh identity
```

`up.sh` reads the SSH_KEY vault out of the running Ivory (`IVORY_URL`, default `http://localhost:8080`)
and installs its public key for its own username. Create that vault in Ivory first — Ivory keeps the
private half, so the key can only travel in this direction. Set `IVORY_VAULT_ID` to pick a specific
one. `NETWORK`, `SUBNET`, `GATEWAY`, `PREFIX` and `COUNT` override the rest.

## Deploying against it

Fill a deploy form with `10.0.0.1`, `10.0.0.2`, `10.0.0.3`, ssh port `22`, and the ssh vault
`up.sh` printed. Node names and both ports come from the template's own per-command defaults. Two
templates coordinate through another cluster and need it deployed first: Patroni's `ETCD3_HOSTS`
names the etcd multi-host cluster, and ClickHouse's `<zookeeper>` names the ZooKeeper one.

## Known lab artifacts

These are the lab, not Ivory or a template. Do not chase them.

- **Host metrics are the physical machine's.** `/node/platform/system/metrics` and `system/info` read
  `/proc` and DMI, which every host shares with the laptop — same CPU, same RAM, same uptime.
- **`Local IP` is empty.** The info script uses `hostname -I`, which busybox and Alpine's `net-tools`
  do not support.
- **The docker socket is world-writable** (`start.sh`), so a non-root ssh account can reach the daemon
  without a docker group. Real hosts put the user in `docker` instead.
- **Nothing is partitioned.** All three hosts sit on one bridge, so this lab exercises deployment, not
  failover under a split.

## Ivory trusts a host key on first use

`clients/console/ssh` remembers a host key in memory for the life of the server process, so a lab host
that comes back with a **new** identity fails every later connection with `host key mismatch` until
Ivory is restarted. `up.sh` keeps `/etc/ssh` in a named volume for exactly this reason, so a
teardown/re-`up.sh` cycle keeps the same identity. `./down.sh --purge` drops it deliberately — after
that, restart Ivory too.
