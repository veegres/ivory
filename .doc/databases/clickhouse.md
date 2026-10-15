# [☰](../../README.md) ClickHouse

A ClickHouse cluster with replicated tables, coordinated through ZooKeeper or ClickHouse Keeper. Every node
accepts writes, so there is no leader — Ivory shows every member as a replica and watches replication health.

**Stage:** BETA

## At a glance

- **Cluster side** — ClickHouse's native protocol, on the **Keeper port**, the same port as the database (`9000`)
- **Keeper credentials** — **password required**
- **Replication** — every node writes; replicas converge through the coordination store
- **Lag** — seconds of replication delay; warns past 30 s

## How Ivory reads the cluster

- **Members** — Ivory reads `system.clusters` for the cluster **named exactly like the cluster in Ivory**, so
  name your Ivory cluster after the `<remote_servers>` entry. A node that declares no such cluster gets a
  warning listing the clusters it does declare.
- **Roles** — every member is a replica: ClickHouse has no leader.
- **Lag** — the largest `absolute_delay` across the node's replicated tables.
- **Warnings** — replicated tables gone read-only (usually a lost session with the coordination store),
  replicas without an active coordination session, stuck tasks in the replication queue with the error that
  blocks them, and delay past ClickHouse's own `max_replica_delay_for_distributed_queries`.

## Operations

- **Reload** the configuration (`SYSTEM RELOAD CONFIG`)
- **View** the settings (`system.settings`)

With no leader, there is nothing to switch over or fail over.

## Query console

SQL, with browsing of databases and tables and a dashboard (databases, connections, sizes, uptime). Saved
queries cover active queries, the replication queue, the coordination connection, detached parts, disks and
more. Running queries cannot be cancelled from Ivory: ClickHouse identifies them by a string id, not a
process id.

## Deploying

Both templates run `clickhouse/clickhouse-server:24`, write the cluster into a `config.d` file before starting
and create the database user named on the deploy screen (`default` unless you change the template). The
macros `{cluster}` and `{replica}` are set to the Ivory cluster and node name. ClickHouse needs a **ZooKeeper** (or ClickHouse Keeper) — the deploy screen asks for its `<node>`
entries once, for the whole cluster.

- **ClickHouse (Multi Host)** — three nodes `clickhouse1`–`clickhouse3`, one per machine, on port `9000`.
  The replica list names `10.0.0.1`–`10.0.0.3`; replace the addresses before deploying.
- **ClickHouse (Single Host)** — all three on one machine on ports `9000`–`9002`. Its coordination store
  defaults to the ZooKeeper **single-host** template, so deploy that first on the same machine.

See [Deployment](../deployment.md).

## Good to know

- The Ivory cluster name is written into ClickHouse's XML config, so it must be a valid XML element name.
- Under docker's bridge networking `CREATE … ON CLUSTER` can hang; create replicated tables on each node
  instead.
