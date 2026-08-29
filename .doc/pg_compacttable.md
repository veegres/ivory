# [☰](../README.md) pgcompacttable

**This tool is Postgres only.** It runs against Postgres nodes, whichever keeper manages them, and
has no equivalent for the other database engines Ivory supports.

PostgreSQL does not reclaim disk space immediately when rows are updated or deleted — old row
versions (dead tuples) accumulate in the heap as **bloat**. Left unchecked, bloat inflates table
and index sizes, slows sequential scans, and increases I/O. Ivory integrates
[pgcompacttable](https://github.com/dataegret/pgcompacttable) to let you measure and clean bloat
without taking your database offline.

## How it works

pgcompacttable compacts tables and rebuilds indexes incrementally in small batches, acquiring only
`ROW EXCLUSIVE` locks. Unlike `VACUUM FULL` or `CLUSTER`, it is safe to run on active production
databases because it never blocks normal reads and writes for more than a moment at a time.

> **Requirement** — The `pgstattuple` extension must be installed in the target database:
> `CREATE EXTENSION pgstattuple;`

## Jobs

Jobs can only be submitted against the **leader** node. Submit a job by configuring:

- **Target database** — which database to compact
- **Table filter** — optionally limit to specific tables
- **Delay ratio** — pause between batches to throttle I/O
- **Size thresholds** — skip tables below a minimum size or above a maximum
- **Vacuum options** — control the `VACUUM` pass that runs before compaction
- **Reindex options** — control how indexes are rebuilt

Once submitted, the job runs in the background. The job list below the form shows all past and
active jobs:

| Status | Meaning |
|--------|---------|
| **Pending** | Queued, waiting to start |
| **Running** | Currently executing |
| **Completed** | Finished successfully |
| **Failed** | Terminated with an error |

Click a job entry to view its full console output and the exact command that was executed.

## Bloat Queries

The **Queries** tab provides SQL-based bloat estimation queries you can run without the external
tool. Switch the target database above the form — the selector applies to both the compaction form
and the bloat queries.

The bloat tool is also available per-node under the [Node → Tools](node.md#tools) tab when you
need to target a specific cluster member directly.
