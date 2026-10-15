# [☰](../README.md) Supported databases

Every database Ivory supports is a plugin with two halves:

- **The cluster side** — how Ivory finds the members of a cluster, tells the leader from the replicas, measures
  replication lag and runs operations like a switchover or failover. In Ivory's terms this is the cluster's
  [Keeper](concepts.md#keeper): an HA tool such as Patroni, or the engine's own built-in replication.
  It is reached on the node's **Keeper port**.
- **The database side** — the query console, saved queries and dashboards, reached on the node's
  **Database port**. For engines that manage their own replication, both ports are usually the same.

Because of that split, every cluster looks the same in Ivory — one [overview](overview.md), one set of
operation buttons, one console — and each page below describes what that engine fills in.

| Database                                     | Cluster                                    | Stage  |
|----------------------------------------------|--------------------------------------------|--------|
| [PostgreSQL + Patroni](databases/patroni.md) | High-availability cluster with Patroni     | STABLE |
| [PostgreSQL](databases/postgres.md)          | Primary and standby streaming replication  | STABLE |
| [etcd](databases/etcd.md)                    | Raft cluster                               | BETA   |
| [Redis](databases/redis.md)                  | Master and replica replication             | BETA   |
| [ClickHouse](databases/clickhouse.md)        | Replicated cluster                         | BETA   |
| [MongoDB](databases/mongodb.md)              | Replica set                                | BETA   |
| [ZooKeeper](databases/zookeeper.md)          | Ensemble                                   | BETA   |

**STABLE** is what Ivory has run against for a long time; **BETA** works end to end but has had less time in
the field; new plugins start at **ALPHA**.

## Operations

An operation a cluster cannot perform is simply not offered in the UI. ⏱ marks one that can also be
scheduled for a later time.

| Cluster    | Operations                                                                                     |
|------------|------------------------------------------------------------------------------------------------|
| Patroni    | Switchover ⏱, failover, restart ⏱, reinitialise, reload, pause / resume, view and edit config |
| PostgreSQL | Failover (promote a standby), reload, view config                                              |
| etcd       | Switchover (move the leader)                                                                   |
| Redis      | Failover (promote a replica), view config                                                      |
| ClickHouse | Reload, view config                                                                            |
| MongoDB    | Failover (step the primary down), view and edit config                                         |
| ZooKeeper  | View config                                                                                    |

## Health and replication lag

Ivory warns about a replica whose lag passes its threshold. The threshold is Ivory's own default for that
engine, not a setting of the engine itself.

| Cluster    | Replication  | Lag is measured in                          | Warning threshold |
|------------|--------------|---------------------------------------------|-------------------|
| Patroni    | One leader   | Bytes behind the leader                     | 50 MiB            |
| PostgreSQL | One leader   | Bytes of WAL not yet replayed               | 50 MiB            |
| etcd       | One leader   | Raft entries behind the leader              | 1000 entries      |
| Redis      | One leader   | Seconds since the master was last heard     | 30 s              |
| ClickHouse | Every node   | Seconds of replication delay                | 30 s              |
| MongoDB    | One leader   | Seconds behind the primary                  | 30 s              |
| ZooKeeper  | One leader   | Not measured                                | —                 |

## Query console

| Cluster    | Console                          | Browse                         | Running queries      | Dashboard |
|------------|----------------------------------|--------------------------------|----------------------|-----------|
| Patroni    | SQL                              | Databases, schemas, tables     | List, cancel, kill   | ✅        |
| PostgreSQL | SQL                              | Databases, schemas, tables     | List, cancel, kill   | ✅        |
| etcd       | etcdctl-style commands           | —                              | —                    |           |
| Redis      | Redis commands                   | Numbered databases             | List, kill clients   |           |
| ClickHouse | SQL                              | Databases, tables              | —                    | ✅        |
| MongoDB    | Shell-style `collection.verb()`  | Databases, collections         | List, kill           |           |
| ZooKeeper  | znode commands                   | —                              | —                    |           |

Every engine ships a set of saved queries for monitoring and diagnostics; see [Node → Database](node.md#database).

## Adding your database

Every engine is a self-contained plugin, so supporting a new one needs no change to Ivory's core. The
repository ships an AI agent for exactly this — ask Claude Code to use `keeper-database-plugin` for your engine
and it writes both halves, their tests and the UI wiring by the project's checklist. Then open a pull request;
see [Contributing](../CONTRIBUTING.md).
