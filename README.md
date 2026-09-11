<div style="text-align: center;" align="center">
   <img src="app/shared/assets/ivory.png" alt="Ivory - open source database cluster management UI" height="200" width="200" />

# Ivory

### Deploy it. Watch it. Fix it. Anywhere.

**Open-source web UI for high-availability database clusters — PostgreSQL & Patroni, etcd, Redis, ClickHouse, MongoDB and ZooKeeper.**

   <img src="https://img.shields.io/github/deployments/veegres/ivory/production?style=flat-square&link=https%3A%2F%2Fgithub.com%2Fveegres%2Fivory%2Fdeployments%2Fproduction" alt="deployment" />
   <img src="https://img.shields.io/docker/v/veegres/ivory/latest?label=stable&style=flat-square&link=https%3A%2F%2Fhub.docker.com%2Fr%2Fveegres%2Fivory" alt="stable version" />
   <img src="https://img.shields.io/docker/v/veegres/ivory?label=latest&style=flat-square&link=https%3A%2F%2Fhub.docker.com%2Fr%2Fveegres%2Fivory" alt="latest version" />
   <img src="https://img.shields.io/docker/pulls/veegres/ivory?style=flat-square&link=https%3A%2F%2Fhub.docker.com%2Fr%2Fveegres%2Fivory" alt="docker pulls" />
</div>

<br>

**Ivory is a self-hosted, open-source database cluster management tool that puts control in your pocket.**
One web UI to **deploy a cluster**, watch its **health and replication lag**, run a **switchover or failover**,
open a **SQL query console**, control **containers over SSH**, and read **VM metrics and logs** — from your
browser or your phone, without dropping into the CLI for every task.

It is built around the concept of a **Keeper** — a generic management layer responsible for cluster
manipulations. A Keeper can be a standalone agent running beside the database, or a management system
embedded directly in the database engine. [Patroni](https://patroni.readthedocs.io/) is the Keeper
implementation for PostgreSQL, and the one Ivory was originally built around.

Ivory is designed for **DBAs, SRE and backend developers** who operate high-availability database clusters and
want a single UI to run, troubleshoot and deploy them — including from a phone, when you are away from your
laptop. It runs as a local tool on your machine or as a shared service on a VM for team use, ships as a single
Docker container, and stores everything itself: **no agent on your hosts, no orchestrator, no cloud account.**

**Contents** — [Features](#features) · [Get started](#get-started) · [Supported databases](#supported-databases-and-keepers) · [Documentation](#documentation) · [FAQ](#faq) · [Contributing](#contribution)

---

<div align="center">
  <h3>🌟 Support This Project! 🌟</h3>
</div>

If you found this project helpful, interesting, or inspiring, please consider giving it a **star** ⭐! Your support
helps:

✅ **Increase visibility** – More people can discover and benefit from this project.  
✅ **Boost motivation** – It encourages us to keep improving and adding new features.  
✅ **Show appreciation** – A small gesture that means a lot to open-source creators!

Thank you for being part of this journey! 🚀

---

### Vision: Beyond Postgres

Ivory started as a Postgres/Patroni tool. We're working towards a more pluggable architecture, where support for other
databases and HA tools could be added as a plugin, instead of being baked into the core. v2 is the first step of that
rework:

| Version |                                                                                                                                                                                                                                                                                                                                                                                                                                                              |
|---------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **v1**  | Hard-wired to Patroni and Postgres, with Postgres-centric node management and no mobile support.                                                                                                                                                                                                                                                                                                                                                             |
| **v2**  | Pluggable Keeper and database engines — the goal is to manage different databases by implementing a simple plugin for each, instead of baking their behaviour into the core. Node management is now VM-centric, with nodes modeled as Hardware + Software behind a generic Platform abstraction (on-prem Docker over SSH today, Kubernetes/OpenShift planned). The UI is now mobile friendly, so you can check on and operate your clusters from your phone. |

---

## Features

**Deployment — spin up a high-availability cluster over plain SSH**

- [Deploying a cluster](.doc/deployment.md#deploying-a-cluster) — over SSH, with no orchestrator or agent
- [Deploying a single cluster node](.doc/deployment.md#deploying-a-single-cluster-node) — add or rebuild one node of an existing cluster
- [Deployment templates](.doc/deployment.md) — the command that runs on each node, saved and reusable
- [Default templates](.doc/deployment.md#default-templates) — ready to deploy for every supported database, in two variants:
    - **Multi host** — one node per machine, the usual layout for a real cluster
    - **Single host** — the whole cluster on one machine, for trying it out locally or on a test VM

**Cluster management — HA operations from the browser**

- [Cluster list](.doc/clusters.md) — register manually, auto-detect from one address, or deploy
- [Cluster health](.doc/overview.md) — node roles, replication lag, pending restarts, warnings
- [HA operations](.doc/overview.md#ha-operations) — switchover, failover, reinitialise, restart, reload, pause
- [Database configuration](.doc/node.md#configuration) — view and patch settings per node

**Node operations — the machine and the container behind each node**

- [Container lifecycle](.doc/node.md#container) — deploy, start, stop, restart, remove, logs
- [Keeper operations per node](.doc/node.md#keeper) — the same HA actions in a full-page view
- [VM metrics](.doc/node.md#system) — CPU, memory, network and processes on the host
- [Log streaming](.doc/node.md#logs) — any file on the host, or container output

**Database troubleshooting — query console and maintenance tools**

- [Query builder](.doc/node.md#database) — saved SQL queries for monitoring and diagnostics
- [Database tools](.doc/node.md#tools) — engine-specific maintenance tools, run as background jobs with live output:
    - **Postgres** — [pgcompacttable](.doc/pg_compacttable.md), reduces table and index bloat without heavy locks

**Access control and operations**

- [Authentication](.doc/authentication.md) — Basic, LDAP and OIDC/SSO, with granular per-user permissions
- [Configuration](.doc/configuration.md) — data persistence, upgrades, TLS, reverse-proxy sub-path

---

## Get started

Ivory ships as one Docker image with no external dependencies.

1. Start the docker container
    - 🐳 **Docker Hub** `docker run -p 80:80 --restart always veegres/ivory`
    - 📦 **GitHub Container registry** `docker run -p 80:80 --restart always ghcr.io/veegres/ivory`
2. Go to http://localhost:80
3. Complete the initial setup wizard (authentication, secret key)
4. Add your first cluster — **manual** (all node addresses), **auto-detect** (one address, Ivory
   finds the rest), or **[deploy](.doc/deployment.md)** a new one from a template
5. Start monitoring

![Ivory web UI showing a PostgreSQL Patroni cluster overview with node roles and replication lag](.doc/images/demo.png)

---

## Supported databases and keepers

| Keeper                                     | Database   | Stage  | Why                                                                                                                                                                  |
|--------------------------------------------|------------|--------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| [Patroni](https://patroni.readthedocs.io/) | Postgres   | STABLE | The Keeper Ivory was originally built around                                                                                                                         |
| Postgres                                   | Postgres   | STABLE   | Not everyone needs HA — sometimes you just want to manage a plain Postgres                                                                                           |
| ETCD                                       | ETCD       | BETA  | Patroni relies on etcd as a DCS, so managing it directly made sense too                                                                                              |
| Redis                                      | Redis      | BETA  | A widely used cache/store that often sits right next to the databases Ivory already manages                                                                          |
| ClickHouse                                 | ClickHouse | BETA  | A popular analytics database and a natural next step for Ivory's plugin set                                                                                          |
| ZooKeeper                                  | ZooKeeper  | BETA  | ClickHouse needs real ZooKeeper protocol, not etcd, to coordinate replication — deployable as its own ensemble and usable as the external DCS other keepers point at |
| MongoDB                                    | MongoDB    | BETA  | A widely used document store with its own native replica-set HA, a natural fit next to Ivory's other native (non-orchestrated) keepers                               |

P.S. these particular databases weren't picked from a grand roadmap — they're mainly what the maintainer runs day to
day, and Ivory exists to simplify that routine first. Broader support grows from there.

---

## Documentation

| Page                                         | What's in it                                                       |
|----------------------------------------------|--------------------------------------------------------------------|
| [Clusters](.doc/clusters.md)                 | Adding clusters manually, auto-detection, tags and filtering        |
| [Overview](.doc/overview.md)                 | What a Keeper is, cluster health, HA operations                     |
| [Node](.doc/node.md)                         | System, Container, Keeper, Database and Tools tabs                  |
| [Deployment](.doc/deployment.md)             | Deploy a cluster or one node, deployment templates, variables       |
| [Authentication](.doc/authentication.md)     | Basic / LDAP / OIDC, users, superusers, permissions                 |
| [Configuration](.doc/configuration.md)       | Data persistence, upgrades, environment variables, TLS, sub-path    |
| [pgcompacttable](.doc/pg_compacttable.md)    | Reducing Postgres table and index bloat from the UI                 |

---

## FAQ

**How do I upgrade to a new version?**
Back up your configuration from the Settings page and restore it in the new version — the format is
backward compatible. Mounting the data directory between containers works too, but reliably only for patch
releases. See [Configuration → Upgrading](.doc/configuration.md#upgrading-to-a-new-version).

**Where does Ivory store its data?**
In `/opt/ivory/data`, backed by a Docker volume. Bind-mount it, or use `--volumes-from`, to keep data across
different containers. See [Configuration → Data and Persistence](.doc/configuration.md#data-and-persistence).

**Does Ivory need an agent on my database hosts?**
No. It talks to the Keeper's API, to the database directly, and to the host over SSH.

**How does authentication work?**
Ivory runs with or without it. With it, everybody who signs in is an Ivory user first — you register a username
and choose whether they sign in with a password, LDAP or SSO, and a password is always set by the person
themselves through a one-time link. See [Authentication](.doc/authentication.md).

**Are my credentials safe?**
Every secret — SSH keys, database passwords, LDAP and OIDC client secrets — is encrypted with the secret word
you choose during setup.

**Can I run Ivory behind a reverse proxy or under HTTPS?**
Yes — `IVORY_URL_PATH` for a sub-path, `IVORY_CERT_FILE_PATH` and `IVORY_CERT_KEY_FILE_PATH` for TLS. See
[Configuration → Environment Variables](.doc/configuration.md#environment-variables).

**Does it work on a phone?**
Yes. The v2 UI is mobile friendly, so you can check cluster health and run HA operations from your phone.

---

## Contribution

If you're interested in contributing to the Ivory project, consider these options:

- [Enhancements](https://github.com/veegres/ivory/issues)
- [Good for newcomers](https://github.com/veegres/ivory/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
- [Supported versions](SECURITY.md)
- [Setup Local Environment](.docker/ivory-dev/README.md)
- [Build Frontend](app/README.md)
- [Build Backend](server/README.md)

---

<sub>Ivory is an open-source, self-hosted GUI and web dashboard for database cluster management and monitoring:
PostgreSQL high availability with Patroni, plain PostgreSQL, etcd, Redis, ClickHouse, MongoDB and ZooKeeper —
covering cluster deployment over SSH, failover and switchover, replication lag monitoring, SQL query console,
Docker container management and VM metrics, in a single Docker container.</sub>
