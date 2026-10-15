# [☰](../../README.md) Redis

Redis with master–replica replication. Ivory reads its state with `INFO` and can promote a replica by hand.

**Stage:** BETA

## At a glance

- **Cluster side** — Redis itself, on the **Keeper port**, the same port as the database (`6379`)
- **Keeper credentials** — **password required**; the username is optional
- **Replication** — one master, any number of replicas
- **Lag** — seconds since a replica last heard from its master; warns past 30 s

## How Ivory reads the cluster

- **Members** — a replica reports the master it follows. A master does **not** list its replicas: the
  addresses Redis shows for them are wherever the connection came from, not an address Ivory could reach, so
  listing them would invent nodes. The master shows how many replicas are connected instead. Add every node
  to the cluster yourself.
- **Roles** — a node with `role:slave` is a replica, any other is the master; a node still loading its data
  shows as starting.
- **Also shown** — the Redis version, memory in use, and the master's replica count.
- **Warnings** — a replica's link to its master is down (and for how long), the last background save or AOF
  write failed, a full resync is in progress, and `min-replicas-to-write` is not met or replicas are past
  `min-replicas-max-lag`.

## Operations

- **Failover** — promotes a replica with `REPLICAOF NO ONE`; refused on the master
- **View** the configuration (`CONFIG GET *`)

There is no Sentinel or Redis Cluster integration, so switchover and automatic failover are not offered.

## Query console

Any Redis command, sent as typed. Browse the numbered databases, list connected clients and kill one with
`CLIENT KILL`. Saved queries cover server info, keyspace, replication, clients, slow log, latency, memory,
persistence, command stats and configuration.

## Deploying

Both templates run the official `redis:7.4` image, with the `default` user's password set from the deploy
screen.

- **Redis (Multi Host)** — three nodes `redis1`–`redis3`, one per machine, on port `6379`, with append-only
  persistence on a volume. The replicas' commands name the master literally (`10.0.0.1 6379`) — replace it
  with your first node's host.
- **Redis (Single Host)** — all three on one machine on ports `6379`–`6381`, without persistence — for trying
  it out.

See [Deployment](../deployment.md).
