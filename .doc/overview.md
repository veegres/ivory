# [☰](../README.md) Overview

## What is a Keeper?

A **Keeper** is the HA management layer responsible for leader election, automatic failover, and
cluster coordination. Ivory uses "Keeper" as its generic term for this role — the same interface
works regardless of which HA system the cluster uses. A Keeper can be a standalone agent running
beside the database engine or a management system embedded directly in it.
[Patroni](https://patroni.readthedocs.io/) is the first supported implementation, with more planned.

The Keeper's job is to answer one question at all times: *who is the leader?* It does this by:

- Continuously monitoring the local database instance and reporting its health
- Competing with peer Keeper agents for the leadership lock in a shared coordination backend (such
  as a DCS — Distributed Configuration Store like etcd, Consul, or ZooKeeper)
- Triggering automatic **failover** when the current leader loses the lock or becomes unreachable
- Exposing a management API (on the **keeper port**, default `8008`) that Ivory uses to read
  cluster state and issue management commands

Because every node exposes a Keeper and every Keeper knows the full cluster topology, Ivory only
needs to reach any one node's keeper port to get a complete picture of the cluster.

### Node Roles

| Role | Description |
|------|-------------|
| **Leader** | The primary node. Holds the leadership lock. Accepts reads and writes. Exactly one node holds this role at a time. |
| **Replica** | A standby node. Continuously replicates from the leader. Read-only. |
| **Unknown** | Role could not be determined — the Keeper is unreachable or the node is initialising. |

---

## Cluster Health Grid

The Overview block is the cluster-level dashboard. It shows the real-time state of every node as
reported by the Keeper.

| Field | Description |
|-------|-------------|
| **Role** | Leader or Replica |
| **State** | Keeper-reported state: `running`, `streaming`, `paused`, `stopped`, `failed`, … |
| **Lag** | Replication lag in bytes behind the leader (`N/A` for the leader itself) |
| **Pending Restart** | Whether the database needs a restart to apply a pending configuration change |
| **Scheduled** | Any scheduled switchover or restart waiting to fire |
| **Tags** | Key-value pairs set on this node by the Keeper |
| **Warnings** | Issues Ivory detected: connection errors, replication problems, etc. |

Clicking any node card opens its [Node](node.md) view for fine-grained per-node operations.

---

## HA Operations

Every node card in the grid has its own action buttons. The cluster header additionally provides
**Pause** and **Resume** which affect the entire cluster at once.

**Per-node actions** (available on each node card):

| Operation | When to use |
|-----------|------------|
| **Switchover** | Planned leadership change. The current leader steps down cleanly, then a replica is promoted. Choose which replica becomes the new leader, or let the Keeper decide. Can be scheduled. |
| **Failover** | Forced promotion of a replica without waiting for the current leader. Use only when the leader is unreachable. |
| **Reinitialise** | Wipes and re-syncs a node's data directory from the current leader. Use when a replica is too far behind or its data directory is corrupted. |
| **Restart** | Restarts the database process on this node. Can be scheduled. |
| **Reload** | Reloads configuration without a restart — applies changes that take effect immediately. |

**Cluster-level actions** (cluster header):

| Operation | When to use |
|-----------|------------|
| **Pause** | Puts the Keeper into maintenance mode across the entire cluster. Automatic failovers are disabled. Use for planned maintenance. |
| **Resume** | Restores normal HA behaviour after a pause. |

The [Node → Keeper](node.md#keeper) tab shows all actions for a single node in a dedicated
full-page view alongside the configuration editor.

---

## Cluster Options

Click the settings icon in the top-right corner to open the cluster options panel:

- Update **database credentials** and **Keeper API credentials**
- Configure **TLS** for keeper and database connections
- Assign or update a **vault** reference for credential management
- Add, remove, or rename **tags**
- Override the **main instance** — pin operations to a specific node instead of the auto-selected leader
