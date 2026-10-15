# [☰](../../README.md) etcd

An etcd cluster — the key-value store that coordinates Raft consensus among its members, and the store
Patroni elects its leader through. Ivory manages it directly because so many HA setups depend on it.

**Stage:** BETA

## At a glance

- **Cluster side** — etcd's v3 API on the client port, used as both **Keeper** and **Database port** (`2379`)
- **Keeper credentials** — optional; only if etcd authentication is enabled
- **Replication** — Raft: one leader, the rest follow
- **Lag** — Raft entries a member is behind the leader; warns past 1000

## How Ivory reads the cluster

- **Members** — etcd's member list comes from whichever node is asked, so the whole cluster is visible from
  one. Each member keeps its etcd name; members other than the node itself are shown as reported until their
  own status arrives.
- **Cluster identity** — each node reports its etcd cluster id, so a node from a different cluster is flagged.
- **Roles** — the member etcd names as leader is the leader, the rest are replicas.
- **Lag** — the gap between a member's applied raft index and the leader's.
- **Also shown** — the etcd version, the raft term, and whether a member is a learner.
- **Warnings** — anything the member itself reports (alarms such as `NOSPACE`, a lost leader), a backlog of
  committed but unapplied entries large enough for etcd to refuse requests, and a member or leader Ivory
  cannot reach, naming the address it tried.

## Operations

- **Switchover** — moves leadership to the member you choose; a candidate is required.

etcd elects its leader by itself, so there is no failover, restart or configuration to manage from Ivory.

## Query console

An `etcdctl`-style console: `get` (with `--prefix`, `--limit`, `--keys-only`), `put`, `del`, `member list`,
`endpoint status`, `endpoint health` and `alarm list`. Saved queries include the member list, endpoint status
and health, alarms, and the keys Patroni stores for its clusters.

## Deploying

Both templates run `quay.io/coreos/etcd` and leave authentication off — on purpose, so the
[Patroni](patroni.md) templates can use the cluster as their coordination store.

- **Etcd (Multi Host)** — three nodes `etcd1`–`etcd3`, one per machine, client port `2379` and peer port
  `2380`. The member list in each command names `etcd1`–`etcd3` at `10.0.0.1`–`10.0.0.3`; replace the
  addresses, and keep the node names matching the list.
- **Etcd (Single Host)** — all three on one machine on client ports `2479`, `2481`, `2483` (peers one above
  each). It avoids the default ports deliberately, and is what the Patroni single-host template points at.

See [Deployment](../deployment.md).
