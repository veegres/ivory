# [☰](../../README.md) ZooKeeper

A ZooKeeper ensemble — the coordination service ClickHouse replicates through, and a store other HA tools can
use too. ZooKeeper elects its leader by itself, so Ivory watches it rather than operating it.

**Stage:** BETA

## At a glance

- **Cluster side** — ZooKeeper's four-letter-word commands (`mntr`, `conf`) on the client port, used as both **Keeper** and **Database port** (`2181`)
- **Keeper credentials** — none; the four-letter words have no authentication
- **Replication** — one leader, followers and observers
- **Lag** — not measured; ZooKeeper does not report it per follower

## How Ivory reads the cluster

- **Members** — the `server.N` lines of the node's configuration, when they carry a client port
  (`…;2181`). The single-host template's lines do not, so there each node shows only itself.
- **Cluster identity** — each node reports the ensemble it is configured into, so a node from another
  ensemble is flagged.
- **Roles** — `leader` and `standalone` are the leader, `follower` and `observer` are replicas.
- **Also shown** — the ZooKeeper version, the number of znodes and connections, and the server state when
  it is neither leader nor follower.
- **Warnings** — followers not in sync with the leader, a server not serving requests, and a configuration
  Ivory could not read.

ZooKeeper only answers commands on its whitelist, so add at least `mntr` and `conf` to
`4lw.commands.whitelist`.

## Operations

- **View** the configuration (`conf`)

Leader election is automatic and cannot be triggered by hand, so nothing else is offered.

## Query console

A znode console: `ls`, `get`, `create`, `set`, `delete` and `exists`, plus the `mntr`, `srvr`, `stat`, `ruok`
and `conf` commands. Database credentials, if given, are used as digest authentication. Saved queries include
the root children, the ensemble configuration and the znodes Patroni and ClickHouse keep.

## Deploying

Both templates run the official `zookeeper:3.9` image with the four-letter words Ivory needs whitelisted.

- **ZooKeeper (Multi Host)** — three nodes `zookeeper1`–`zookeeper3`, one per machine, client port `2181`
  (quorum `2888`, election `3888`). Server ids are 1–3 in order and the server list names
  `10.0.0.1`–`10.0.0.3`; replace the addresses and keep the order.
- **ZooKeeper (Single Host)** — all three on one machine on client ports `2181`–`2183`. It is what the
  ClickHouse single-host template points at.

See [Deployment](../deployment.md).
