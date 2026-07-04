# [☰](../README.md) Node

A **node** in Ivory represents a single server (physical or virtual) that is a member of a database
cluster. Each node runs two core services side by side: a **Keeper** agent that manages high
availability, and a **Database** engine that stores data. Ivory lets you interact with both from
the same place.

Clicking a node in the [Overview](overview.md) opens the Node view, which is split into five tabs:
**Keeper**, **Container**, **Database**, **Tools**, and **Platform**.

## Platform

The Platform tab provides visibility into the virtual machine (VM) hosting this node over SSH. It does not require any
agent on the remote host — only SSH access.

### VM Metrics

Real-time resource usage is polled from the remote host and displayed on an auto-refreshing dashboard:

| Metric | Details |
|--------|---------|
| **CPU** | Total CPU ticks and idle ticks. Use the difference over an interval to derive utilisation. |
| **Memory** | Total bytes and available bytes on the host. |
| **Network** | Cumulative bytes received and transmitted since boot. |

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

Before deploying, copy Ivory's public SSH key to the remote host using the **Copy ID** button on the Platform tab. This
sets up passwordless SSH access that all container operations depend on.

| Operation | Description |
|-----------|-------------|
| **List** | Show all deployments (Docker containers) currently on the remote host. |
| **Up** | Pull the image and start a new container. Accepts a raw option string with `{{placeholder}}` interpolation (see below). |
| **Start** | Start an existing container that was stopped. |
| **Stop** | Gracefully stop a running container (the container record remains). |
| **Restart** | Stop and immediately start the container. |
| **Down** | Stop and **remove** the container entirely. |
| **Logs** | Stream live output from the container. |

### Option Placeholders

The deployment form accepts a free-form option string (Docker Compose arguments or environment variables) with
Ivory-managed placeholders that are interpolated at deploy time:

| Placeholder | Value injected |
|-------------|---------------|
| `{{host}}` | The node's hostname or IP address |
| `{{cluster}}` | The cluster name configured in Ivory |
| `{{dcs}}` | The DCS connection string used by the Keeper |
| `{{keeperPort}}` | The Keeper API port |
| `{{dbPort}}` | The database port |
| `{{dbUser}}` | The database superuser username |
| `{{dbPass}}` | The database superuser password |

This lets a single image option template work across all nodes of a cluster — each node fills in its own values
automatically.

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

The Tools tab groups maintenance operations that run against the node's database.

### pgcompacttable

[pgcompacttable](https://github.com/dataegret/pgcompacttable) reduces table and index bloat
without taking heavy locks, making it safe to run on production systems. See
[pgcompacttable](pg_compacttable.md) for full details.

Jobs can only be submitted against the **leader** node (Ivory enforces this automatically). Once
submitted, a job runs in the background — you can monitor its status and stream its output from the
job list below the form.
