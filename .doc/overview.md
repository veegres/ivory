# [☰](../README.md) Cluster Overview

Clicking a cluster on the [Clusters](clusters.md) page opens its overview: one row per node with its live state,
the HA operations, and the cluster's settings. The terms used here — keeper, leader, replica — are defined in
[Core concepts](concepts.md).

Ivory asks the keeper on every configured node and merges the answers, so the overview shows the cluster as the
database sees it, not only as you configured it.

---

## Cluster Health Grid

Each node row shows:

| Field               | Description                                                                                               |
|---------------------|-----------------------------------------------------------------------------------------------------------|
| **Host and ports**  | The node's host, keeper port, database port and SSH port                                                   |
| **Role**            | Leader or replica, as the keeper reports it                                                               |
| **State**           | Whether the database process is `running`, `starting`, `stopped`, `failed` or `unreachable`              |
| **Lag**             | How far a replica is behind, in the engine's own unit: bytes for PostgreSQL, seconds for MongoDB, Redis and ClickHouse, raft entries for etcd. It turns into a warning past Ivory's threshold for that engine. |
| **Sync**            | Shown on a replica that replicates synchronously                                                          |
| **Pending restart** | The database needs a restart to apply a configuration change                                              |
| **Scheduled**       | A switchover or restart that is scheduled and has not fired yet                                           |
| **Tags**            | What the engine reports about the node beyond these fields: a version, a size, the operator's own labels  |
| **Warnings**        | Problems with the node, labelled by who found them (see below)                                            |

**Warnings come from two places.** The engine's own warnings are what the database reports about itself — a lost
coordination session, a failed save, a replication slot about to be lost. Ivory's warnings are what it found by
comparing answers: a node the database lists but you never configured, a configured node the database has never
heard of, a node that did not answer, a cluster with no leader, or lag past the threshold. Each engine's page under
[Supported databases](databases.md) lists the warnings it can raise.

Clicking a node opens its [Node](node.md) view for per-node work.

---

## HA Operations

Each node row has its own action buttons. Which ones appear depends on the engine: an operation the engine cannot
perform is not offered. [Supported databases → Operations](databases.md#operations) shows which engine supports what.

**Per-node actions:**

| Operation        | When to use                                                                                                                         |
|------------------|-------------------------------------------------------------------------------------------------------------------------------------|
| **Switchover**   | Planned leadership change. The leader steps down cleanly and a replica is promoted. Pick the replica, or let the keeper decide. Patroni can schedule it. |
| **Failover**     | Promote a replica without waiting for the leader. Use it only when the leader is gone.                                              |
| **Reinitialise** | Wipe a replica's data and copy it again from the leader. Use it when a replica is too far behind or its data is corrupted.          |
| **Restart**      | Restart the database on this node. Patroni can schedule it.                                                                         |
| **Reload**       | Re-read configuration without a restart.                                                                                            |

**Cluster-level actions**, in the cluster header:

| Operation  | When to use                                                                                    |
|------------|------------------------------------------------------------------------------------------------|
| **Pause**  | Put the keeper into maintenance mode. Automatic failover is off until you resume.              |
| **Resume** | Return to normal HA behaviour after a pause.                                                   |

The [Node → Keeper](node.md#keeper) tab shows the same actions for one node on a full page, next to the
configuration editor.

---

## Cluster Options

The settings button in the top-right corner of the cluster header opens the cluster options:

- **Credentials** — the vault entries for keeper and database access
- **TLS certificates** — for keeper and database connections
- **Tags** — used to filter the [cluster list](clusters.md#cluster-list)
- **Main instance** — the node Ivory sends keeper operations to. Ivory picks the leader on its own; override it
  only when you need to pin operations to one node.
