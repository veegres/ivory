# [☰](../README.md) About Ivory

## What Ivory Is

Ivory is a self-hosted web UI for operating high-availability database clusters: deploy a cluster, watch its
health and replication lag, run a switchover or failover, open a query console, control containers over SSH
and read VM metrics and logs — from a browser or a phone, without dropping into the CLI for every task.

It is designed for **DBAs, SREs and backend developers** who run such clusters and want one place to operate,
troubleshoot and deploy them. It runs as a local tool on your machine or as a shared service on a VM for team
use, ships as a single Docker container and stores everything itself: no agent on your hosts, no orchestrator,
no cloud account.

## The Keeper

Ivory is built around the concept of a **Keeper** — a generic management layer responsible for cluster
manipulations. A Keeper can be a standalone agent running beside the database, or a management system
embedded directly in the database engine. [Patroni](https://patroni.readthedocs.io/) is the Keeper
implementation for PostgreSQL, and the one Ivory was originally built around. See
[Overview → What is a Keeper?](overview.md#what-is-a-keeper) for how it shapes the cluster view.

## Vision: Beyond Postgres

Ivory started as a Postgres/Patroni tool. It is moving towards a pluggable architecture, where support for
other databases and HA tools is added as a plugin instead of being baked into the core.

| Version | |
|---------|-|
| **v1** | Hard-wired to Patroni and Postgres, with Postgres-centric node management and no mobile support. |
| **v2** | Pluggable Keeper and database engines — each database is managed through a simple plugin instead of core behaviour. Node management is VM-centric, with nodes modeled as Hardware + Software behind a generic Platform abstraction (on-prem Docker over SSH today, Kubernetes/OpenShift planned). The UI is mobile friendly, so you can check on and operate your clusters from your phone. |

## Why These Databases

| Keeper     | Why |
|------------|-----|
| Patroni    | The Keeper Ivory was originally built around |
| Postgres   | Not everyone needs HA — sometimes you just want to manage a plain Postgres |
| ETCD       | Patroni relies on etcd as a DCS, so managing it directly made sense too |
| Redis      | A widely used cache/store that often sits right next to the databases Ivory already manages |
| ClickHouse | A popular analytics database and a natural next step for Ivory's plugin set |
| ZooKeeper  | ClickHouse needs the real ZooKeeper protocol, not etcd, to coordinate replication — deployable as its own ensemble and usable as the external DCS other keepers point at |
| MongoDB    | A widely used document store with its own native replica-set HA, a natural fit next to Ivory's other native (non-orchestrated) keepers |

These particular databases weren't picked from a grand roadmap — they're mainly what the maintainer runs day
to day, and Ivory exists to simplify that routine first. Broader support grows from there.
