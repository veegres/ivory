# [☰](../README.md) Core Concepts

Every page in these docs uses the same handful of terms. They are defined once here.

## Cluster

A **cluster** is a group of nodes that run one database together: a Patroni cluster, an etcd cluster, a MongoDB
replica set, a ZooKeeper ensemble. In Ivory a cluster has a name, a list of nodes, the plugins that know how to
talk to it, and optional credentials, certificates and tags.

Ivory treats the cluster you configured and the cluster the database reports as two things that should match. When
they do not — a node the database lists but you never added, or a node you added that the database has never heard
of — the [overview](overview.md) warns about it instead of hiding the difference.

## Node

A **node** is one member of a cluster: a machine (physical or virtual) together with the database running on it.
Each node has a name that is unique within its cluster, a host, and three ports:

| Port            | What Ivory reaches on it                                                         |
|-----------------|----------------------------------------------------------------------------------|
| **Keeper port** | The keeper, to read the cluster's state and run operations such as a switchover  |
| **DB port**     | The database itself, for the query console and saved queries                     |
| **SSH port**    | The machine, for metrics, logs and the containers deployed on it                 |

For engines that manage their own replication, the keeper port and the database port are usually the same.

## Keeper

A **keeper** is whatever decides who leads the cluster and what Ivory asks about its members. It can be a separate
HA agent running next to the database — [Patroni](databases/patroni.md) for PostgreSQL — or the replication built
into the engine itself, as with etcd, Redis, ClickHouse, MongoDB, ZooKeeper and plain PostgreSQL.

Whichever it is, Ivory talks to it through one interface, so every cluster gets the same overview, the same
operation buttons and the same warnings. Each engine fills in what it can; an operation it cannot perform is not
offered. [Supported databases](databases.md) compares them.

Most engines follow a **single-leader** model: one member accepts writes and the rest follow it. ClickHouse is
**multi-leader**: every replica accepts writes and they converge, so its overview has no leader to show.

### Node Roles

| Role        | Meaning                                                                                       |
|-------------|-----------------------------------------------------------------------------------------------|
| **Leader**  | Accepts writes. A single-leader cluster has exactly one.                                      |
| **Replica** | Follows the leader, or in a multi-leader cluster is one of several equal members.             |
| **Unknown** | Ivory could not tell. Treat it as a fault: the node is unreachable or answered unexpectedly.  |

## Plugins

Every supported database is a pair of plugins, which you pick when you add a cluster:

- the **keeper plugin** covers the cluster side — members, roles, lag, health and operations;
- the **database plugin** covers the query side — the console, saved queries and charts.

Selecting a keeper preselects the database it usually runs with.

## Platform

A **platform** is where Ivory deploys and controls things. Today that is Docker on your own Linux hosts, reached
over SSH, with Kubernetes and OpenShift planned. A platform has two halves, and the [node view](node.md) follows
them: the **System** tab is the machine itself, the **Container** tab is what is deployed onto it.

## Deployment Template

A **deployment template** is the command that runs on each node of a new cluster, written out in full and saved.
Ivory fills in a small, fixed set of `{{variables}}` — the node's name, host and ports, the credentials — and runs
exactly what the screen shows. Every keeper plugin ships ready templates; you can copy them or write your own. See
[Deployment](deployment.md).

## Vault

The **vault** is where Ivory keeps credentials: SSH keys, keeper passwords and database passwords. Every secret in
it is encrypted with the secret word chosen during setup (see [Configuration → Secrets](configuration.md#secrets)).
A cluster points at vault entries rather than storing passwords itself, so one entry can serve several clusters and
rotating it updates all of them.
