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
- One or more **nodes** — for each node a **name** (unique within the cluster), its host, keeper port
  (default `8008` for Patroni), database port (default `5432` for Postgres), and optionally an SSH
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

Click the **deploy** button to deploy a new cluster onto remote hosts over SSH. The deployment is
described by a **[deployment template](deployment.md)**: the command that will run on each node,
either one of the templates a keeper plugin ships or one you wrote yourself.

The dialog has three screens:

1. **Pick a template** for the selected keeper and platform, or copy an existing one and adjust it.
2. **Fill it in** — cluster name, SSH credentials, database credentials, and one node card per
   command in the template (name, host, ports). Each card shows the command it will run, already
   interpolated with that node's values.
3. **Deploy**, and read the per-node logs. Ivory registers the cluster once a node is up, then runs
   any post scripts.

See [Deployment](deployment.md) for the details: what a template contains, the variables a command
may use, the default templates, and how a deploy is executed.

After deployment you can manage each node's container lifecycle from the [Node → Container](node.md#container)
tab at any time.

---

## Cluster Options

In the [Overview](overview.md), the settings button (top-right corner of the cluster header) opens
the cluster options panel. Here you can:

- Set or rotate **database credentials** and **Keeper API credentials**
- Attach or change **TLS certificates** for keeper and database connections
- Link a **vault** entry to avoid storing credentials in plain text
- Add or remove **tags** for filtering
- Change the **main instance** — the specific node that Ivory sends keeper operations to. By default
  Ivory always picks the current leader automatically; override this here if needed
