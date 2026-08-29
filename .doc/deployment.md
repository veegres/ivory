# [☰](../README.md) Deployment

Ivory can deploy a new cluster onto your own machines. You provide the hosts and the credentials, it
runs the deployment on each host over SSH and registers the result as a cluster you can operate. The
hosts need an SSH server and Docker; nothing else is installed on them.

A deployment is described by a **template**: the command that will run on each node, saved and
reusable, with a small set of `{{variables}}` Ivory fills in per node.

---

## Why templates

Ivory does not generate deployment commands. A deployment is a command you write, or one of the
commands a keeper plugin ships that you copy and adjust.

| | |
|---|---|
| **The command is visible** | The deploy screen shows the command already interpolated with this node's values, and that is what runs. Nothing rewrites it afterwards. |
| **Default templates are included** | Every keeper plugin ships two working templates, a multi-host and a single-host variant. Deploy one as-is, or copy it and edit what you need. |
| **Templates are reusable** | A template describes how a cluster is deployed. Hosts, ports, node names and credentials belong to a deploy, not to a template, so the same template works for different clusters. |
| **New engines need no new machinery** | A template is text plus a closed variable list, so a keeper plugin only has to ship its commands to appear in the list. |
| **Secrets stay on the server** | `{{dbPass}}` is never sent to the browser. The preview shows `*****`, and the server substitutes the real password from the vault at execution time. |
| **Templates can be reviewed** | A template is plain text, so it can be read and compared before anyone runs it. |

---

## Anatomy of a template

| Field | Meaning |
|-------|---------|
| **Name** | Unique within one keeper + platform pair — that pair is the list you picked the name in. |
| **Description** | What the template assumes: node naming, which ports differ, what to edit before running it. |
| **Keeper** | The HA engine the commands deploy (Patroni, Etcd, Redis, …). |
| **Platform** | The deployment target the commands are written for (Docker over SSH today). |
| **Commands** | An ordered list, one command per node, so a template's command count is the number of nodes it deploys. |

Each command carries two fields:

- **Command** — the whole deployment command for that node, for example a complete `docker run`. It
  is stored as one literal command; Ivory does not assemble it from parts.
- **Post Script** *(optional)* — a script executed **inside** the container once the node is up. It
  belongs on the command that should run it: a step that needs the whole cluster running (etcd's
  `auth enable`, mongo's `rs.initiate`) goes on the **last** node, because post scripts run only
  after every node has started.

Inside a template a node is identified only by its position in the list. Which host the first command
runs on is chosen at deploy time.

---

## Variables

A command may reference a **closed, fixed set** of variables. Ivory substitutes each node's own
values into that node's own command, so one node's values can never leak into another's.

| Variable | Value | Example |
|----------|-------|---------|
| `{{cluster}}` | Cluster name | `my-cluster` |
| `{{name}}` | Node name — unique in the cluster, and the deployment's own name | `node-1` |
| `{{host}}` | Node host | `10.0.0.1` |
| `{{sshPort}}` | Node SSH port | `22` |
| `{{keeperPort}}` | Keeper endpoint port (the database port when the engine has no separate keeper API) | `8008` |
| `{{dbPort}}` | Database endpoint port | `5432` |
| `{{dbUser}}` | Database username, resolved from the vault on the server | `postgres` |
| `{{dbPass}}` | Database password, resolved from the vault on the server | `*****` |

**Everything else is written literally in the command.** A peer port, a member list, a coordinator
address or which node is the leader are values Ivory does not know, so they are written as plain
text rather than added as variables:

```
-e ETCD3_HOSTS="etcd-1:2379,etcd-2:2379,etcd-3:2379"
```

The list is fixed, and this is checked in both directions:

- A `{{placeholder}}` outside the table above is **rejected when the template is saved**, so a
  misspelled variable is reported as an error instead of being treated as a new one.
- A variable with no value at deploy time is **rejected before anything runs** (`missing values for
  placeholders: …`), so a partially interpolated command never reaches a host.

> Docker's own `{{json .}}`-style format strings are left untouched — only the names above are
> substituted.

---

## Deploying a cluster

**Clusters page → Deploy.** The dialog is three screens: pick a template, fill it in, read the logs.

### 1. Pick a template

The list shows the default templates for the selected keeper and platform alongside your own. Picking
a row opens the deploy form for it; default and manual templates are deployed the same way. The row
buttons cover the other two actions: **edit** your own template, or **copy** any template into one
you own. Default templates are computed from the plugins on every request rather than stored, which
is why they cannot be edited or deleted and always match the version of Ivory you are running.

### 2. Fill it in

| Section | What it asks for |
|---------|------------------|
| **Cluster** | The cluster name. It becomes `{{cluster}}` and the name Ivory registers. |
| **Database Credentials** | An existing vault entry, or a new username/password Ivory stores in the vault before any node starts. Skipped for engines that consume no credentials (ZooKeeper, MongoDB). |
| **SSH Credentials** | The vault entry holding the key used to reach every host, or a username for a key Ivory generates. |
| **Node cards** | One card per command in the template: node name, host, SSH port, keeper port, database port — plus the command itself, read-only and already interpolated. |

**Parallel deployment** is a checkbox on the cluster section. Leave it off for keepers such as
Patroni that need their nodes brought up one after another.

Some engines fix the database username themselves (`postgres` for Patroni, `root` for Etcd, `default`
for Redis). The form prefills it and does not allow changing it, and a vault entry with a different
username is rejected rather than overridden.

Supplying both a vault entry and an inline username/password is rejected as well, since Ivory would
otherwise have to choose between them silently.

Deploy stays clickable while fields are empty; clicking it reports what is missing by marking each
invalid field.

### 3. Deploy

The server then, in order:

1. **Validates** the cluster name, the node names (mandatory and unique), the credentials and the
   template.
2. **Creates vaults** for anything entered inline — the SSH key and the database credentials.
3. **Deploys every node**: interpolates that node's command, opens an SSH connection to its host and
   runs the command verbatim. Nodes run one after another, or together when parallel is checked.
4. **Registers the cluster**, but only once at least one node is up. If no node started, no cluster is
   registered and the vaults this deploy created are removed again.
5. **Runs the post scripts** in node order, once every node is up. If a node failed, the post scripts
   are skipped and the logs report it.

Every step is written to the response screen as a timestamped log line per node. Going back from the
logs returns to the form with its values still filled in, so a deploy that failed on one node can be
corrected and repeated without re-entering everything.

---

## Deploying a single cluster node

**Node → Container → Deploy** runs one node of a template on the host you already have open. Use it
to add a member to a running cluster, or to rebuild one you removed.

Everything on this screen is filled in and disabled: the node is the one the dialog was opened on,
its ports are the plugin's defaults, and the credentials are the cluster's. A node's name, host and
ports are chosen when the cluster is created, not here.

The one choice on this screen is **which of the template's nodes runs on this host** — a toggle above
the node card. A template describes a whole cluster, so Ivory cannot tell which of its members this
host should be. The toggle is hidden for a one-node template.

---

## Default templates

Every keeper plugin ships two default templates for Docker, each deploying a three-node cluster:

| Template | What it deploys |
|----------|-----------------|
| **Patroni** | Three spilo nodes coordinating through an external DCS — point `ETCD3_HOSTS` at the one you run. |
| **Postgres** | One leader and two streaming replicas, no HA agent. |
| **Etcd** | A three-member static etcd cluster; the last node enables authentication as its post script. |
| **Redis** | One leader and two replicas. |
| **ClickHouse** | Three replicas in one shard, coordinated through an external ZooKeeper / ClickHouse Keeper. |
| **ZooKeeper** | A three-node ensemble. |
| **MongoDB** | A three-member replica set; the last node runs `rs.initiate` as its post script. |

Each comes in two variants:

- **Multi Host** — one node per machine, with ports published normally. This is the usual layout for
  a real cluster, where a lost machine takes only one node with it.
- **Single Host** — all three nodes on one machine, using `--network host` with a distinct port per
  node. Use it to try a cluster out locally or on a single test VM; it gives you the same topology
  to operate, but no protection against losing the machine. ClickHouse is the exception here: its
  native port is fixed by the image's `config.xml`, so its single-host variant needs a custom image
  to avoid port collisions.

Both are ordinary templates rather than two modes of one template, so there is no flag to set —
which one you deploy is which row you pick.

Read the description on each row before running it: it says what the commands assume, such as
"name the nodes `etcd-1..3` or edit the member list to match".

---

## Writing your own template

**New template** on the list screen, or **copy** on any row. The same editor is used in all three
cases (a new template, a copy, or editing your own) and returns to the list on save, where the new
template is selected like any other.

The editor gives you a name, a description, and one collapsible block per node holding that node's
command and post script, plus a **Variables** box listing every token next to an example of what it
turns into.

A few conventions to follow:

- **Write a complete command** — real member lists, real peer ports, real image tags. A copied
  template should deploy as-is once the hosts are filled in.
- **Describe what the template assumes**, for example "name the nodes `etcd-1..3` or edit the member
  list to match". A reader has no other way to know.
- **Write each node's command out in full**, even when two commands differ by only a few flags. A
  template is data rather than code, and the duplication keeps every command readable on its own.
- **A leader/replica difference is just a different command at position 1** — a bootstrap command
  first and a rebase command for the rest, for example.
- **Entry scripts go at the end of the command**, after the image. A multi-line `sh -c '…'` keeps its
  newlines; whitespace between flags is collapsed.

Your templates are stored in Ivory and included in **v2 backups**, so they migrate along with the
rest of your configuration. Default templates are not exported, since they are recomputed on every
install.

---

## Requirements

On each target host:

- An **SSH server** reachable from Ivory, and **Docker** installed.
- Ivory's public key in the host's `authorized_keys`. Ivory generates the key pair and keeps it in
  its vault, where the public key is shown next to the entry; copy it onto every host before
  deploying.

In Ivory:

| Permission | Covers |
|------------|--------|
| `view.deployment.template.list` | Seeing templates at all |
| `manage.deployment.template.create` | Creating and copying templates |
| `manage.deployment.template.update` | Editing your own templates |
| `manage.deployment.template.delete` | Deleting your own templates |
| `manage.cluster.create` | Deploying a whole cluster |
| `manage.node.platform.container` | Deploying a single node, and the container lifecycle |

---

## Troubleshooting

| Message | What it means |
|---------|---------------|
| `unknown variables in command: {{…}}` | The command references a placeholder outside the closed set. Fix the spelling, or write the value literally. |
| `missing values for placeholders: {{…}}` | The command references a variable this deploy has no value for — usually credentials for an engine that needs them. |
| `database credentials are required` | The keeper plugin consumes `{{dbUser}}`/`{{dbPass}}` but no vault entry or inline pair was given. |
| `no node deployed: the cluster was not registered` | Every node failed. The per-node log lines above say why; nothing was registered and the vaults this attempt created were removed. |
| `skipping post-deploy initialization` | Not every node came up, so cluster-wide initialization was not attempted. Fix the failed node, then run its deploy again. |
| Node deployed but shows as unreachable | Ivory matches cluster members by `host:keeperPort`. Check that the ports on the node card match the ones inside the command. |

---

## See also

- [Clusters](clusters.md) — the three ways to add a cluster
- [Node → Container](node.md#container) — the container lifecycle after deployment
- [Node → System](node.md#system) — VM metrics, logs and processes on the host
