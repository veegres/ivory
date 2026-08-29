# [☰](../README.md) Node

A **node** in Ivory represents a single server (physical or virtual) that is a member of a database
cluster. Each node runs two core services side by side: a **Keeper** agent that manages high
availability, and a **Database** engine that stores data. Ivory lets you interact with both from
the same place.

Clicking a node in the [Overview](overview.md) opens the Node view, which is split into five tabs:
**System**, **Container**, **Keeper**, **Database**, and **Tools**.

A **platform** is the deployment target Ivory integrates with — on-prem Docker over SSH today, with
Kubernetes/OpenShift planned. The tabs follow its two parts: the **System** tab covers the machine a
deployment runs on, the **Container** tab covers what is deployed onto it. On a platform that cannot
address the machine itself, the System tab is hidden.

## System

The System tab provides visibility into the virtual machine (VM) hosting this node over SSH. It does not require any
agent on the remote host — only SSH access.

### VM Metrics

Real-time resource usage is polled from the remote host and displayed on an auto-refreshing dashboard:

| Metric | Details |
|--------|---------|
| **CPU** | Total CPU ticks and idle ticks. Use the difference over an interval to derive utilisation. |
| **Memory** | Total bytes and available bytes on the host. |
| **Network** | Cumulative bytes received and transmitted since boot. |

### Processes and Info

The running process list on the host, and basic system information about the machine — enough to answer "is anything
else on this box" without opening a shell.

### Logs

Stream any file from the remote host in real time — useful for watching database or Keeper logs without opening a
separate SSH session. Configure the absolute path to the log file, the tail size (number of lines to start from), and
enable **follow** mode to keep the stream live as new lines are appended.

---

## Container

The Container tab manages the **deployment lifecycle** of the node's services on the remote host. Ivory uses SSH to
execute Docker operations on the target machine, so no agent needs to be installed on the node beyond an SSH server and
Docker.

### Prerequisites

All container operations run over SSH with the key held in Ivory's vault. Copy that key's public half into the host's
`authorized_keys` before using this tab; the vault shows it next to the entry.

| Operation | Description |
|-----------|-------------|
| **Overview** | The container for this node: live metrics, streaming logs, and the lifecycle buttons below. |
| **List** | Show all deployments (Docker containers) currently on the remote host. |
| **Deploy** | Run a [deployment template](deployment.md) on this host to create the container (see below). |
| **Start** | Start an existing container that was stopped. |
| **Stop** | Gracefully stop a running container (the container record remains). |
| **Restart** | Stop and immediately start the container. |
| **Remove** | Stop and **remove** the container entirely. |
| **Logs** | Stream live output from the container. |

### Deploy

**Deploy** opens the same [deployment template](deployment.md) flow used for a whole cluster, applied to this one node.
Use it to add a member to a running cluster, or to rebuild one that was removed.

Everything on the screen is filled in and disabled: the node is the one the dialog was opened on, its ports are the
keeper plugin's defaults, and the credentials are the cluster's. The one choice is **which of the template's nodes runs
on this host**, a toggle above the node card, since a template describes a whole cluster.

The command is read-only and shown already interpolated with this node's values, so it matches what will run. The
database password is the exception: it is displayed as `*****` and substituted on the server from the vault, so it
never reaches the browser.

See [Deployment](deployment.md) for what a template is, the full variable reference, and the default templates.

---

## Keeper

The Overview already shows per-node HA actions on every node card. The Keeper tab provides the
same operations in a dedicated full-page view — all actions visible at once rather than collapsed
into menus — alongside the configuration editor.

### Configuration

The Keeper tab also embeds a **configuration editor** for this node's database settings. Configuration
is managed through the Keeper's dynamic configuration API and propagates to all nodes in the cluster.

Configuration is updated as a **partial patch**:

- Only the fields you specify are changed. Everything else is left untouched.
- Setting a value to **`null`** removes that setting, reverting it to the database default.
- You do not send the full configuration — only the parameters you want to change.

The Keeper decides whether a change can be applied immediately via reload, or requires a database
restart. When a restart is needed the **Pending Restart** flag appears on the affected node(s) in
the [Overview](overview.md) health grid. You can then trigger a restart from the Overview action
bar or from this tab.

---

## Database

The Database tab connects directly to the database instance on this node and provides a
full-featured query and monitoring interface.

### Charts

At-a-glance metrics pulled from the database instance: number of databases, active connections,
database sizes, uptime, schema count, table sizes, index sizes, and total size.

### Query Builder

A template-driven SQL runner. Ivory ships with built-in **system queries** organised by category:

| Category | What it covers |
|----------|---------------|
| **Bloat** | Table and index bloat estimation |
| **Activity** | Long-running queries, locks, idle connections |
| **Replication** | Replication lag, slot status, WAL senders/receivers |
| **Statistic** | Table statistics, autovacuum status, cache hit ratios |
| **Other** | Miscellaneous diagnostic queries |

System queries cannot be deleted. If **Manual Queries** is enabled in your Ivory configuration,
you can also create your own template queries and customise the SQL of any system query.

> **Note** — Ivory executes queries as-is, including `UPDATE`, `INSERT`, and `DELETE`. Access to
> the Query Builder can be restricted per user from the Settings permissions panel.

#### Query Varieties

Each query can carry variety labels that describe where and how it should run:

| Label | Meaning |
|-------|---------|
| **Database Sensitive** | Results vary per database; always run with a specific database selected. |
| **Master Only** | Must run on the leader — will fail or return wrong results on a replica. |
| **Replica Recommended** | Expensive query; prefer running on a replica to avoid load on the leader. |

#### Prepared Statements

Add `$1`, `$2`, … placeholders in your SQL, name each parameter, and Ivory injects the values at
run time. Useful for parameterised diagnostic queries you run repeatedly with different inputs.

### Console

Free-form SQL editor for ad-hoc queries. No templates — type any SQL and run it directly against
the selected database.

---

## Tools

The Tools tab groups maintenance operations that run against the node's database. Each tool is
specific to one database engine and only works on nodes running that engine.

### pgcompacttable (Postgres only)

[pgcompacttable](https://github.com/dataegret/pgcompacttable) is a Postgres tool: it reduces table
and index bloat without taking heavy locks, which makes it safe to run on production systems. It
applies to Postgres nodes only, whichever keeper manages them, and has no equivalent for the other
engines Ivory supports. See [pgcompacttable](pg_compacttable.md) for full details.

Jobs can only be submitted against the **leader** node (Ivory enforces this automatically). Once
submitted, a job runs in the background — you can monitor its status and stream its output from the
job list below the form.
