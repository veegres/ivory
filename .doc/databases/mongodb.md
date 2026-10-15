# [☰](../../README.md) MongoDB

A MongoDB replica set — MongoDB's own built-in replication and automatic failover, with no separate HA tool.

**Stage:** BETA

## At a glance

- **Cluster side** — MongoDB itself, on the **Keeper port**, the same port as the database (`27017`)
- **Keeper credentials** — optional; only if access control is enabled
- **Replication** — one primary, secondaries follow; MongoDB elects a new primary on its own
- **Lag** — seconds a member is behind the primary; warns past 30 s

## How Ivory reads the cluster

- **Members** — `replSetGetStatus` from any member returns the whole replica set.
- **Cluster identity** — each node reports its replica set name, so a node from another set is flagged.
- **Roles** — `PRIMARY` is the leader, `SECONDARY` a replica. `RECOVERING` and `ROLLBACK` are replicas that
  are catching up, `STARTUP` is starting, an arbiter is running with no data role, and a member that is
  `DOWN` is unreachable.
- **Also shown** — the member state when it is anything but primary or secondary, the sync source when a
  secondary replicates from another secondary, and the heartbeat ping.
- **Warnings** — failed heartbeats, naming which member saw them, and MongoDB's own info message for a member.

## Operations

- **Failover** — asks the current primary to step down (`replSetStepDown`) so the set elects another; it is
  run on the primary
- **View and edit** the replica set configuration (`replSetGetConfig` / `replSetReconfig`) — an edit replaces
  the whole configuration document

## Query console

A shell-style console: `collection.find()`, `findOne`, `insertOne`/`insertMany`, `updateOne`/`updateMany`,
`deleteOne`/`deleteMany`, `countDocuments`, `distinct`, `aggregate`, and `db.runCommand()`. Browse databases
and collections, list running operations and kill one. Saved queries cover server status, replica set status
and config, current operations, index usage and the profiler.

## Deploying

Both templates run the official `mongo:7` image as replica set members named after the Ivory cluster. A
**post-script on the last node** runs `rs.initiate` once all three are up. Access control is left off.

- **Mongo (Multi Host)** — three nodes `mongo1`–`mongo3`, one per machine, on port `27017`. The
  `rs.initiate` member list names `10.0.0.1`–`10.0.0.3`; replace the addresses before deploying.
- **Mongo (Single Host)** — all three on one machine on ports `27017`–`27019`.

See [Deployment](../deployment.md).
