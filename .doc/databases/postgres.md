# [☰](../../README.md) PostgreSQL

Plain PostgreSQL with streaming replication — a primary and its standbys, no HA tool in between. Not every
setup needs automatic failover; this lets you watch and operate one without installing Patroni.

**Stage:** STABLE

## At a glance

- **Cluster side** — PostgreSQL itself, on the **Keeper port**, the same port as the database (`5432`)
- **Database side** — PostgreSQL on the **Database port**
- **Keeper credentials** — **required**: a PostgreSQL user; the cluster side logs into the database to read its state
- **Replication** — one primary, any number of standbys
- **Lag** — bytes of WAL a standby has not replayed yet; warns past 50 MiB

## How Ivory reads the cluster

Each node is asked one connection's worth of questions per refresh:

- **Role** — `pg_is_in_recovery()`: a node not in recovery is the primary, one in recovery is a standby.
- **From the primary** — `pg_stat_replication` marks which standbys are synchronous. Standbys are matched by
  their `application_name`, so name each standby's replication connection after its node.
- **From a standby** — `pg_stat_wal_receiver` names the server it streams from. That server is shown with an
  unknown role, because with cascading replication it may itself be a standby.
- **Lag** — measured against the primary's reported WAL end; a standby with no WAL receiver running shows no
  lag at all rather than a misleading zero.
- **Also shown** — the PostgreSQL version and connections in use out of `max_connections`.
- **Warnings** — a standby that is not streaming, a replication slot past `max_slot_wal_keep_size` or one
  whose WAL is already lost, and replication details Ivory could not read.

Reading the replication views in full needs a superuser or a member of `pg_read_all_stats`.

## Operations

- **Failover** — promotes a standby with `pg_promote()`; refused on the primary
- **Reload** the configuration with `pg_reload_conf()`
- **View** the configuration (`pg_settings`)

Switchover, restart, reinitialise and pause need an HA tool to coordinate them, so they are not offered — use
[PostgreSQL + Patroni](patroni.md) for those.

## Query console

The full PostgreSQL console — the same as for [Patroni](patroni.md#query-console): SQL, schema browsing,
dashboard, running queries with cancel and terminate, saved queries, and
[pgcompacttable](../tools/postgres/pg_compacttable.md).

## Deploying

Both templates run the official `postgres:18` image and create the user `postgres`. The first node is the
primary; the other two clone it with `pg_basebackup` on first start and stream from it.

- **Postgres (Multi Host)** — three nodes `postgres1`–`postgres3`, one per machine, all on port `5432`.
  The standbys' commands name the primary's address literally (`10.0.0.1:5432`) — replace it with your
  primary's host before deploying.
- **Postgres (Single Host)** — all three on one machine using host networking, on ports `5432`–`5434`.

See [Deployment](../deployment.md).
