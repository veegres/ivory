# MCP server

Ivory exposes its features to [Model Context Protocol](https://modelcontextprotocol.io) clients, so an
AI agent can read a cluster's health and operate it with the same permissions the signed-in user has.
The endpoint speaks streamable HTTP and lives at `/api/mcp` on the same address and port as the web UI.

It is the troubleshooting surface that benefits most: everything Ivory already knows about a cluster —
node roles, replication lag, pending restarts, its own drift warnings, VM metrics, container logs and
the saved diagnostic queries — becomes readable in one conversation instead of one screen at a time.

---

## Connecting a client

Authentication is a bearer token: Ivory's ordinary session JWT, sent in the `Authorization` header.

1. Sign in and take the token out of the login response:

   ```bash
   TOKEN=$(curl -sX POST http://localhost/api/basic/login \
     -H 'Content-Type: application/json' \
     -d '{"username":"alice","password":"..."}' | jq -r .response.token)
   ```

2. Point a client at the endpoint with that token. For Claude Code:

   ```bash
   claude mcp add --transport http ivory http://localhost/api/mcp \
     --header "Authorization: Bearer $TOKEN"
   ```

   Or, in `.mcp.json`:

   ```json
   {
     "mcpServers": {
       "ivory": {
         "type": "http",
         "url": "http://localhost/api/mcp",
         "headers": { "Authorization": "Bearer ${IVORY_TOKEN}" }
       }
     }
   }
   ```

> **The token lasts one hour.** Ivory's session tokens expire after an hour and there is no refresh
> mechanism, so a long agent session will start failing with `401` and the token has to be replaced.
> Issue a fresh one the same way.

When Ivory runs with authentication switched off, no token is needed and every tool is available —
the same as the REST API, which is already unauthenticated in that mode.

---

## Permissions

An MCP client sees **only the tools the user holds permission for**. Each tool reuses the permission of
the REST route it mirrors, so nothing new has to be granted: a user who may view a cluster overview in
the UI gets `ivory_cluster_overview`, and a user who may not run a switchover never sees
`ivory_keeper_switchover` at all.

The tool list is recomputed on every request, so revoking a permission takes effect on the agent's next
call without reconnecting.

Tools are also annotated for the client's approval prompt: read-only tools are marked as such, and the
ones that can break something are marked destructive.

---

## Tools

**Each tool takes the same arguments as the REST route it mirrors.** A tool lives in the feature it
belongs to (`cluster_mcp.go`, `node_mcp.go`, …) and speaks that feature's own request shape, so there is
one contract behind two transports. Only `cluster` tools are addressed by cluster name; `node` and
`query` tools are addressed the way the web UI addresses them — by host, port and vault id.

So the usual flow is two steps, exactly as the UI does it:

```
1. ivory_cluster_get { "name": "prod" }
   -> nodes[{name, host, sshPort, keeperPort, dbPort}],
      plugins{keeper, database}, tls, certs,
      vaults{keeperId, databaseId, sshKeyId}

2. ivory_keeper_node_list {
     "host": "10.0.0.1", "port": 8008,
     "cluster": "prod", "plugin": "patroni",
     "vaultId": "<vaults.keeperId from step 1>"
   }
```

A tool requires an argument only where the REST request does (`binding:"required"`); everything else
defaults the way an omitted JSON field does, and the feature's own validation names whatever it still
needs. `ivory_cluster_overview` is the exception worth knowing: it resolves every node of a cluster
itself, which is why it is one call and the fastest route to a diagnosis.

### Clusters

| Tool | Permission |
|------|------------|
| `ivory_cluster_list` — registered clusters, with nodes, plugins and tags | `view.cluster.list` |
| `ivory_cluster_get` — one cluster's stored configuration | `view.cluster.item` |
| `ivory_cluster_overview` — live health: roles, state, lag, warnings | `view.cluster.overview` |
| `ivory_tag_list` — every tag and the clusters carrying it | `view.tag.list` |
| `ivory_cluster_create` / `ivory_cluster_update` | `manage.cluster.create` / `.update` |
| `ivory_cluster_delete` — removes the registration only | `manage.cluster.delete` |
| `ivory_cluster_detect` — register from one running keeper endpoint | `manage.cluster.create` |
| `ivory_cluster_fix` — replace the node list with what the keeper reports | `manage.cluster.update` |
| `ivory_cluster_deploy` — deploy a whole new cluster | `manage.cluster.create` |

### Keeper

| Tool | Permission |
|------|------------|
| `ivory_keeper_node_list` — one node's view of the members | `view.node.keeper.overview` |
| `ivory_keeper_config_get` | `view.node.keeper.config` |
| `ivory_keeper_config_update` | `manage.node.keeper.config.update` |
| `ivory_keeper_switchover` / `ivory_keeper_switchover_cancel` | `manage.node.keeper.switchover` |
| `ivory_keeper_failover` | `manage.node.keeper.failover` |
| `ivory_keeper_reinitialize` | `manage.node.keeper.reinitialize` |
| `ivory_keeper_restart` / `ivory_keeper_restart_cancel` | `manage.node.keeper.restart` |
| `ivory_keeper_reload` | `manage.node.keeper.reload` |
| `ivory_keeper_activate` / `ivory_keeper_pause` | `manage.node.keeper.activation` |

Keeper tools take `KeeperOneRequest` — the connection (`host`, `port`), the cluster name, the plugin, the
vault id, and a free-form `body` passed straight to the keeper, the same body the REST route forwards:

| operation | body |
|---|---|
| switchover | `{"leader": "<key>", "candidate": "<key>", "scheduled_at": "<iso-8601>"}` |
| failover | `{"candidate": "<key>"}` |
| restart | `{"schedule": "<iso-8601>", "restart_pending": true}` |
| reinitialize | `{"force": true}` |
| config update | the patch, in the shape that keeper's own API expects |

`leader` and `candidate` are the **keeper's own member keys**, which `ivory_cluster_overview` reports as
each node's `keeper.key`. Patroni does not accept a host there. Send a switchover to the current leader.

### Node system and containers

| Tool | Permission |
|------|------------|
| `ivory_system_metrics` / `ivory_system_processes` / `ivory_system_info` / `ivory_system_logs` | `view.node.system` |
| `ivory_system_copy_id` — install a public key in `authorized_keys` | `manage.node.system` |
| `ivory_container_list` / `ivory_container_logs` / `ivory_container_metrics` | `view.node.platform.container` |
| `ivory_container_start` / `_stop` / `_restart` / `_up` / `_down` / `_exec` | `manage.node.platform.container` |
| `ivory_container_keeper_deploy` — deploy one node of a cluster | `manage.node.platform.container` |

System and container tools take a `PlatformVaultConnection` (`host`, `port`, `vaultId`, `platform`) — the
ssh port and the cluster's `vaults.sshKeyId` from `ivory_cluster_get`. Metrics are counters, not
percentages: call twice and take the difference. Log tools return a bounded tail (100 lines unless `tail`
says otherwise) rather than a live stream, and report a failure to read the log in `notes` beside the
lines, because the underlying stream pushes its problems as messages rather than returning an error.

### Database

| Tool | Permission |
|------|------------|
| `ivory_query_list` — the shipped and user-saved queries | `view.query.crud.list` |
| `ivory_query_log_get` | `view.query.log.list` |
| `ivory_query_activity` / `ivory_query_databases` / `ivory_query_schemas` / `ivory_query_tables` | `view.query.db.info` |
| `ivory_query_chart` | `view.query.db.chart` |
| `ivory_query_template` — run a saved query | `manage.query.db.template` |
| `ivory_query_console` — run arbitrary SQL | `manage.query.db.console` |
| `ivory_query_cancel` / `ivory_query_terminate` | `manage.query.db.cancel` / `.terminate` |
| `ivory_query_create` / `_update` / `_delete` | `manage.query.crud.*` |
| `ivory_query_log_delete` | `manage.query.log.delete` |

### Deployment templates

| Tool | Permission |
|------|------------|
| `ivory_deployment_template_list` — shipped and user templates | `view.deployment.template.list` |
| `ivory_deployment_template_create` / `_update` / `_delete` | `manage.deployment.template.*` |

### Ivory itself

`ivory_app_info` reports the version, configuration state and which authentication methods are accepted.
It mirrors the unauthenticated `/info` route and needs no permission.

---

## What this lets an agent do

Granting `manage.*` permissions to a token an agent holds means the agent can take a production cluster
down. That is the same power the user already has in the UI, and Ivory gates it the same way — but a UI
click and a tool call are not equally deliberate, so it is worth being explicit:

- `ivory_keeper_failover` forces a leader change without the checks a switchover makes, and **can lose
  committed data** if the candidate is behind.
- `ivory_keeper_reinitialize` wipes a member's data directory and re-clones it from the leader.
- `ivory_keeper_pause` switches off automatic failover and nothing reminds you later.
- `ivory_container_down` stops **and removes** a container, including the logs that would have explained
  a failure.
- `ivory_container_exec` and `ivory_container_up` run commands you supply on the host or inside a
  container. There is no allow-list — that is what a deployment is.
- `ivory_query_console` runs SQL verbatim. Ivory does not classify statements anywhere, so `DROP` is as
  reachable as `SELECT`.

If an agent should only diagnose, give its user `view.*` permissions and nothing else: the mutating tools
then do not appear in its tool list at all. Most MCP clients also prompt before each call and honour the
destructive annotation Ivory sets.

### Not exposed

Some of Ivory's REST surface has deliberately no tool:

- **Erase, free and secret rotation** — irreversible, whole-install operations with no diagnostic value.
- **Backup export and import** — binary file transport, and the export carries secrets.
- **Vault, certificate, user and permission administration** — these move credentials through tool
  arguments, and letting an agent edit permissions would let it widen its own access.

A deploy names its credentials as **vault ids**, never as inline passwords, for the same reason: no
secret passes through the agent's context. `ivory_cluster_deploy` therefore omits the inline
`sshUser`/`sshPass`/`dbUser`/`dbPass` fields the REST route accepts — create the vault entries in the UI
first and name them in `options.vaults`.

---

## How it is built

MCP is a second engine beside the http one, not a feature of its own:

- `core/engine/engine_mcp.go` is the MCP engine, the sibling of `engine_http.go`. It declares every tool
  next to the permission that gates it and the feature router that serves it — the same thing the route
  table does — and it builds the transport. `main.go` creates it and hands it to the http engine, which
  mounts it at `/api/mcp`.
- Each feature keeps its own MCP implementation in its own package, in `<feature>_mcp.go`, beside the
  router it already has: `cluster_mcp.go`, `node_mcp.go` (plus `node_mcp_keeper.go` /
  `node_mcp_platform.go`), `query_mcp.go`, `deployment_mcp.go`, `tag_mcp.go`, `management_mcp.go`.
- Only `cluster` resolves a cluster into connections, because only `cluster` imports `node` and `query`.
  That is why node and query tools are addressed directly.

## Notes

- The endpoint is stateless: every request carries its own session, and `Mcp-Session-Id` is not used.
- Requests are capped at 4 MB.
- Cross-origin browser requests are rejected. Non-browser MCP clients send no `Origin` and are unaffected.
- Behind a reverse-proxy sub-path, the endpoint moves with it: `IVORY_URL_PATH=/ivory` puts it at
  `/ivory/api/mcp`.
