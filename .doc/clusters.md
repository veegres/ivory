# [☰](../README.md) Clusters

The Clusters page is the entry point for all cluster management. It shows every cluster Ivory knows
about and lets you add new ones in three different ways depending on whether the cluster already
exists or needs to be deployed from scratch.

## Cluster List

Each row in the list shows:

- **Cluster name** — click to open the [Overview](overview.md) for that cluster
- **Nodes** — colour-coded by role: green = leader, blue = replica, grey = unknown
- **Health indicators** — warning icons when nodes have problems, error icons when Ivory cannot
  reach a node
- **Tag chips** — cluster tags for filtering (manage tags in the Overview settings)

The filter bar above the list lets you narrow clusters by tag. A warning counter on the right shows
how many clusters have at least one unhealthy node — useful when the list spans multiple pages.

---

## Adding a Cluster

There are three ways to add a cluster depending on your situation.

### Manual (existing cluster, known node addresses)

Click the **+** button to open the manual form. Provide:

- Cluster **name**
- One or more **nodes** — for each node specify its host, keeper port (default `8008`),
  database port (default `5432`), and optionally an SSH port
- **Keeper plugin** — currently Patroni
- **Database plugin** — currently PostgreSQL
- Optional: TLS certificates for keeper and database connections, vault references for credentials,
  and tags

Use this when you already know all node addresses or when auto-detection is not available.

### Auto-Detect (existing cluster, only one node address known)

Click the **auto** button (the wand/detect icon next to +). Provide just one node's host and keeper
port — Ivory queries Patroni on that node, which returns the full cluster membership, and then
automatically creates all node entries with the correct hosts, keeper ports, and database ports.

All detected nodes share the same keeper plugin, database plugin, TLS, and vault settings that you
configure in the wizard.

### Deploy (new cluster, no PostgreSQL running yet)

Click the **deploy** button to open the deployment wizard. This is for provisioning a brand-new
cluster on remote hosts. The wizard collects:

- Cluster **name** and a list of **target hosts** with their intended ports
- **SSH credentials** — used to connect to each host and run Docker commands; optionally stored in
  vault for reuse
- **Database credentials** — the superuser username and password for the new PostgreSQL instances
- **Image options** — a free-form Docker option string with `{{placeholder}}` interpolation
  (see [Node → Container](node.md#container) for the full placeholder reference)

Ivory then deploys the container on each host sequentially (or in parallel if configured), sets up
Patroni, and registers all nodes in the cluster.

After deployment you can manage each node's container lifecycle from the [Node → Container](node.md#container)
tab at any time.

---

## Cluster Options

In the [Overview](overview.md), the settings button (top-right corner of the cluster header) opens
the cluster options panel. Here you can:

- Set or rotate **database credentials** (postgres user password) and **Patroni REST API credentials**
- Attach or change **TLS certificates** for keeper and database connections
- Link a **vault** entry to avoid storing credentials in plain text
- Add or remove **tags** for filtering
- Change the **main instance** — the specific node that Ivory sends keeper operations to. By default
  Ivory always picks the current leader automatically; override this here if needed
