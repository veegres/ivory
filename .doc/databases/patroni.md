# [☰](../../README.md) PostgreSQL + Patroni

A high-availability PostgreSQL cluster managed by [Patroni](https://patroni.readthedocs.io/): Patroni elects
the leader through a coordination store (etcd), fails over automatically and exposes a REST API that Ivory
drives. This is the cluster type Ivory was originally built around, and the one with the most operations.

**Stage:** STABLE

## At a glance

- **Cluster side** — Patroni REST API over HTTP(S) on the **Keeper port** (Patroni's default is `8008`)
- **Database side** — PostgreSQL on the **Database port** (`5432` by default)
- **Keeper credentials** — optional; sent as HTTP basic auth when Patroni's REST API requires it
- **Replication** — one leader, any number of replicas; synchronous replicas are marked
- **Lag** — bytes behind the leader, as Patroni reports it; warns past 50 MiB

## How Ivory reads the cluster

- **Members** — one call to Patroni's `/cluster` endpoint returns every member, so reaching any one node is
  enough to see the whole cluster. Each member keeps the name Patroni gives it.
- **Cluster identity** — each node reports its Patroni `scope` (Patroni 3.0.4+). A node configured into the
  wrong Ivory cluster shows a warning naming the scope it belongs to.
- **Roles** — `leader`, `master` and `standby_leader` show as the leader; `replica`, `sync_standby` and
  `quorum_standby` as replicas, the last two marked as synchronous.
- **Also shown** — pending restarts, a scheduled restart or switchover, the timeline, whether the cluster is
  paused, and every tag set on the member in Patroni's own configuration.
- **Warnings** — a member whose `api_url` Ivory cannot parse is still shown, with a warning naming it.

## Operations

Everything in the [HA operations](../overview.md#ha-operations) list is supported:

- **Switchover**, to a chosen replica or one Patroni picks, now or **scheduled** for later
- **Failover** to a replica
- **Restart**, now or **scheduled**
- **Reinitialise** a replica from the leader
- **Reload** the configuration
- **Pause / resume** Patroni's automatic failover
- **View and edit** the dynamic configuration (`PATCH /config`)

## Query console

The full PostgreSQL console: SQL, browsing of databases, schemas and tables, a dashboard (databases,
connections and sizes, uptime), the list of running queries with cancel and terminate, and saved queries for
activity, bloat, statistics and replication. PostgreSQL clusters also get the
[pgcompacttable](../tools/postgres/pg_compacttable.md) tool.

## Deploying

Both templates run Zalando's [Spilo](https://github.com/zalando/spilo) image (PostgreSQL 18 with Patroni) and
create the database user `postgres`. Patroni needs an **etcd** to coordinate through — the deploy screen asks
for its address once, for the whole cluster.

- **Patroni (Multi Host)** — three nodes `patroni1`–`patroni3`, one per machine, Keeper port `8008` and
  Database port `5432`. Point it at an etcd you already run (or deploy one with the [etcd](etcd.md) template).
- **Patroni (Single Host)** — all three nodes on one machine using host networking, Keeper ports `8008`–`8010`
  and Database ports `5442`–`5444`. Its etcd address defaults to the etcd **single-host** template, so deploy
  that first on the same machine.

See [Deployment](../deployment.md) for how templates and the deploy screen work.

## Good to know

- Patroni's REST API is unauthenticated in the shipped templates; if yours requires auth, give the cluster
  Keeper credentials.
- The cluster name in Ivory becomes Patroni's `scope` when you deploy, so it is also how Patroni names the
  cluster in etcd.
