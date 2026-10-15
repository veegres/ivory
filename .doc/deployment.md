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
| **Templates are reusable** | A template suggests node names, ports and usernames, but the hosts and passwords come from each deploy, so the same template works for different clusters. |
| **New engines need no new machinery** | A template is text plus a closed variable list, so a keeper plugin only has to ship its commands to appear in the list. |
| **Secrets stay on the server** | `{{keeperPass}}` and `{{dbPass}}` are never sent to the browser. The preview shows `*****`, and the server substitutes the real password from the vault at execution time. |
| **Templates can be reviewed** | A template is plain text, so it can be read and compared before anyone runs it. |

---

## Anatomy of a template

| Field | Meaning |
|-------|---------|
| **Name** | Unique within one keeper + platform pair — that pair is the list you picked the name in. |
| **Description** | What the template assumes: node naming, which ports differ, what to edit before running it. |
| **Keeper** | The HA engine the commands deploy (Patroni, Etcd, Redis, …). |
| **Platform** | The deployment target the commands are written for (Docker over SSH today). |
| **Cluster Defaults** | What the whole cluster shares: the keeper and database usernames its commands create, and the DCS address its nodes coordinate through. Never a password. |
| **Commands** | An ordered list, one command per node, so a template's command count is the number of nodes it deploys. |

Each command carries three fields:

- **Command** — the whole deployment command for that node, for example a complete `docker run`. It
  is stored as one literal command; Ivory does not assemble it from parts.
- **Post Script** *(optional)* — commands executed **inside** the container once every node is up,
  one per line, in order. Each line runs on its own, with no shell around it, so `&&` and pipes do not
  work. A step that needs the whole cluster running, such as MongoDB's `rs.initiate`, goes on the
  **last** node.
- **Node Defaults** — the node name, host and the keeper and database ports this command answers on.
  They live on the command because only the command knows them: a single-host template gives every
  node its own port, and a member list naming `etcd2` is only right if that node is called `etcd2`.

---

## Variables

A command may reference a **closed, fixed set** of variables. Ivory substitutes each node's own
values into that node's own command, so one node's values can never leak into another's.

**Node variables** have a different value in every node's command:

| Variable | Value | Example |
|----------|-------|---------|
| `{{name}}` | Node name — unique in the cluster, and the deployment's own name | `etcd1` |
| `{{host}}` | Node host | `10.0.0.1` |
| `{{sshPort}}` | Node SSH port | `22` |
| `{{keeperPort}}` | Keeper endpoint port | `8008` |
| `{{dbPort}}` | Database endpoint port | `5432` |

**Cluster variables** have one value shared by every node:

| Variable | Value | Example |
|----------|-------|---------|
| `{{cluster}}` | Cluster name | `my-cluster` |
| `{{dcs}}` | The coordination store the cluster registers with, in whatever format the engine expects | `10.0.0.1:2379,10.0.0.2:2379` |
| `{{keeperUser}}` | Keeper username | `default` |
| `{{keeperPass}}` | Keeper password, substituted on the server | `*****` |
| `{{dbUser}}` | Database username | `postgres` |
| `{{dbPass}}` | Database password, substituted on the server | `*****` |

**Everything else is written literally in the command.** A peer port, a member list or which node
is the leader are values only you know, so they are written as plain text rather than added as
variables:

```
-e ETCD_INITIAL_CLUSTER="etcd1=http://10.0.0.1:2380,etcd2=http://10.0.0.2:2380,etcd3=http://10.0.0.3:2380"
```

The list is fixed, and this is checked in both directions:

- A `{{placeholder}}` outside the tables above is **rejected when the template is saved**, so a
  misspelled variable is reported as an error instead of being treated as a new one.
- A variable with no value at deploy time is **rejected before anything runs** (`missing values for
  placeholders: …`), so a partially interpolated command never reaches a host.

Values are filled in after the command is split into arguments, so a password containing quotes,
spaces or `$` is passed through exactly as typed, without escaping.

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
| **Cluster** | The cluster name, which becomes `{{cluster}}`, and the **DCS** address when the template uses `{{dcs}}`. |
| **SSH** | An existing vault entry whose key is already on the hosts, or a username and password Ivory uses once to install a key it generates. |
| **Keeper** | Credentials for the keeper endpoint: an existing vault entry, a new username and password, or none. |
| **Database** | Credentials for the database, with the same three choices. |
| **Node cards** | One card per command in the template: node name, host, SSH port, keeper port, database port — plus the command itself, read-only and already interpolated. |

Each field starts from the template's defaults. A username the template names is the account its
commands create, so it is filled in and locked, and the vault list only offers entries with that
username — what is left to choose is the password. A template that names no username opens that
credential on **none**. A port the template does not state starts empty and must be typed.

**Parallel deployment** is a checkbox on the cluster section. Leave it off for keepers such as
Patroni that need their nodes brought up one after another.

Deploy stays clickable while fields are empty; clicking it marks each invalid field.

### 3. Deploy

The server then, in order:

1. **Validates** the cluster name, the node names (mandatory and unique), every node's three ports,
   and the credentials. A vault entry together with a typed username and password for the same
   credential is rejected, since Ivory would otherwise have to choose between them silently.
2. **Installs the SSH key** when it generated one: it signs in to each host with the SSH password
   you typed and adds the public key to `authorized_keys`. Every later connection uses the key. If a
   host refuses, the deploy stops here, before anything has started.
3. **Creates vaults** for the credentials you typed.
4. **Registers the cluster** with every node, before any of them starts.
5. **Deploys every node**: interpolates that node's command, connects to its host and runs the
   command verbatim. Nodes run one after another, or together when parallel is checked.
6. **Runs the post scripts** in node order, once every node is up. If a node failed, the post scripts
   are skipped and the logs say so.

**A failed node is kept, on purpose.** The cluster stays registered with every node, and a container
that failed to start is left on its host. Its logs, exit code and a shell into it are what explain the
failure, and Ivory reaches them through the cluster's own node settings. Open the failed node, find
the cause, then fix and redeploy it.

Every step is written to the response screen as a timestamped log line per node, along with how each
node's deploy and post-script steps ended: done, failed, or skipped. Going back from the logs returns
to the form with its values still filled in.

A successful deploy means `docker run` exited cleanly, not that the database inside is healthy. The
cluster's [overview](overview.md) is what tells you that.

---

## Deploying a single cluster node

**Node → Container → Deploy** runs one node of a template on the host you already have open. Use it
to add a member to a running cluster, or to rebuild one you removed.

Everything on this screen is filled in and disabled: the node is the one the dialog was opened on,
its ports are the ones the chosen command states, and the credentials are the cluster's. A node's
name, host and ports are chosen when the cluster is created, not here. If the chosen command states
no ports, Deploy stays disabled; fix that in the template.

The one choice on this screen is **which of the template's nodes runs on this host** — a toggle above
the node card. A template describes a whole cluster, so Ivory cannot tell which of its members this
host should be. The toggle is hidden for a one-node template.

---

## Default templates

Every keeper plugin ships two default templates for Docker, each deploying a three-node cluster:

| Template | What it deploys |
|----------|-----------------|
| **Patroni** | Three spilo nodes coordinating through an etcd you run, given as the DCS address. |
| **Postgres** | One leader and two streaming replicas, no HA agent. |
| **Etcd** | A three-member static etcd cluster, without authentication. |
| **Redis** | One leader and two replicas. |
| **ClickHouse** | Three replicas in one shard, coordinated through a ZooKeeper or ClickHouse Keeper ensemble you run. |
| **ZooKeeper** | A three-node ensemble. |
| **MongoDB** | A three-member replica set; the last node runs `rs.initiate` as its post script. |

Each comes in two variants:

- **Multi Host** — one node per machine, with ports published normally. This is the usual layout for
  a real cluster, where a lost machine takes only one node with it. Replace the example addresses in
  the commands with your own before deploying.
- **Single Host** — all three nodes on one machine, using `--network host` with a distinct port per
  node, and `localhost` as every node's host. Use it to try a cluster out on one test VM; it gives
  you the same topology to operate, but no protection against losing the machine. The single-host
  Patroni and ClickHouse templates point at the single-host etcd and ZooKeeper templates on the same
  machine, so deploy those first.

Both are ordinary templates rather than two modes of one template, so there is no flag to set —
which one you deploy is which row you pick.

Read the description on each row before running it: it says what the commands assume, such as
"name the nodes `etcd1..3` or edit the member list to match".

---

## Writing your own template

**New template** on the list screen, or **copy** on any row. The same editor is used in all three
cases (a new template, a copy, or editing your own) and returns to the list on save, where the new
template is selected like any other.

The editor gives you a name, a description, a **Cluster Defaults** box, and one collapsible block per
node holding that node's command, post script and **Node Defaults**. A **Variables** box lists every
token next to an example of what it turns into, grouped into node and cluster variables.

A few conventions to follow:

- **Write a complete command** — real member lists, real peer ports, real image tags. A copied
  template should deploy as-is once the hosts are filled in.
- **Describe what the template assumes**, for example "name the nodes `etcd1..3` or edit the member
  list to match". A reader has no other way to know.
- **Write each node's command out in full**, even when two commands differ by only a few flags. A
  template is data rather than code, and the duplication keeps every command readable on its own.
- **A leader/replica difference is just a different command at position 1** — a bootstrap command
  first and a rebase command for the rest, for example.
- **Entry scripts go at the end of the command**, after the image. A multi-line `sh -c '…'` keeps its
  newlines; whitespace between flags is collapsed. Inside such a script, read passwords from
  environment variables the command sets rather than from `{{dbPass}}`, since the script is parsed
  a second time by the container's shell.
- **Write options only.** Ivory adds `docker run` itself, and drops a leading `docker run` if you
  paste one in.

Your templates are stored in Ivory and included in **v2 backups**, so they migrate along with the
rest of your configuration. Default templates are not exported, since they are recomputed on every
install.

---

## Requirements

On each target host:

- An **SSH server** reachable from Ivory, and **Docker** installed.
- For an SSH vault entry you pick yourself, its public key in the host's `authorized_keys`; the vault
  shows it next to the entry. A key Ivory generates during the deploy is installed for you.

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
| `unknown variables: {{…}}` | The template references a placeholder outside the closed set. Fix the spelling, or write the value literally. |
| `missing values for placeholders: {{…}}` | The command references a variable this deploy has no value for — usually a credential that was set to none, or an empty DCS address. |
| `cluster node keeper, database and ssh ports are all required` | A node card is missing a port. Every node needs all three. |
| `failed to authorize the ssh key on …` | Ivory could not sign in with the SSH username and password to install its generated key. Nothing was started. |
| `N of M node(s) deployed` | Some nodes failed. The cluster is registered; open the failed nodes to see why. |
| `skipping post-deploy initialization` | Not every node came up, so the post scripts did not run. Fix the failed node, then run its deploy again. |
| Node deployed but shows as unreachable | Ivory matches cluster members by `host:keeperPort`. Check that the ports on the node card match the ones inside the command. |

---

## See also

- [Clusters](clusters.md) — the three ways to add a cluster
- [Node → Container](node.md#container) — the container lifecycle after deployment
- [Node → System](node.md#system) — VM metrics, logs and processes on the host
