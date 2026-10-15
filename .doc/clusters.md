# [☰](../README.md) Clusters

The Clusters page is the entry point for all cluster management. It shows every cluster Ivory knows
about and lets you add new ones in three different ways depending on whether the cluster already
exists or needs to be deployed from scratch.

## Cluster List

Each row in the list shows:

- **Cluster name** — click to open its [overview](overview.md)
- **Nodes** — colour-coded by role: green = leader, blue = replica, grey = unknown
- **Health indicators** — warning icons when nodes have problems, error icons when Ivory cannot
  reach a node
- **Tag chips** — cluster tags for filtering (manage tags in [Cluster Options](overview.md#cluster-options))

The filter bar above the list lets you narrow clusters by tag. A warning counter on the right shows
how many clusters have at least one unhealthy node — useful when the list spans multiple pages.

---

## Adding a Cluster

There are three ways to add a cluster depending on your situation.

### Manual (existing cluster, known node addresses)

Click the **+** button to open the manual form. Provide:

- Cluster **name**
- One or more **nodes** — for each node a **name** (unique within the cluster), its host, keeper port
  (for example `8008` for Patroni), database port (for example `5432` for Postgres), and optionally an SSH
  port
- **Keeper plugin** — Patroni, Postgres, Etcd, Redis, ClickHouse, ZooKeeper or MongoDB
- **Database plugin** — the engine behind that keeper; selecting a keeper preselects its usual pair
- Optional: TLS certificates for keeper and database connections, vault references for credentials,
  and tags

Use this when you already know all node addresses or when auto-detection is not available.

### Auto-Detect (existing cluster, only one node address known)

Click the **auto** button (the wand/detect icon next to +). Provide just one node's host and keeper
port — Ivory asks the keeper on that node for the full cluster membership, then creates every node
entry with the correct hosts, keeper ports and database ports. Where the keeper names its members
itself (Patroni, etcd), those names become the node names; otherwise a node is named after its host.

All detected nodes share the same keeper plugin, database plugin, TLS, and vault settings that you
configure in the wizard.

### Deploy (new cluster, nothing running yet)

Click the **deploy** button to create a new cluster on your own hosts over SSH. You pick a
[deployment template](deployment.md) — the command that runs on each node — fill in the hosts and
credentials, and Ivory registers the cluster and starts every node. [Deployment](deployment.md)
covers the whole flow.

---

## Cluster Options

Credentials, TLS certificates, tags and the main instance are changed later from the cluster's
[Overview → Cluster Options](overview.md#cluster-options).
