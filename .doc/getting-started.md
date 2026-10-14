# [☰](../README.md) Getting Started

Ivory ships as one Docker image with no external dependencies — no agent on your hosts, no orchestrator,
no cloud account. This page takes you from nothing to your first monitored cluster.

## 1. Run the container

Pick either registry, they publish the same image:

```bash
# Docker Hub
docker run -p 80:80 --restart always veegres/ivory

# GitHub Container Registry
docker run -p 80:80 --restart always ghcr.io/veegres/ivory
```

For anything beyond a quick try, keep the data in a directory you control, so it survives replacing the
container:

```bash
docker run -p 80:80 \
  --mount type=bind,source=YOUR_LOCAL_PATH,target=/opt/ivory/data \
  --restart always veegres/ivory
```

See [Configuration](configuration.md) for persistence, upgrades, TLS, running under a sub-path and every
environment variable.

## 2. Complete the setup wizard

Open http://localhost:80. On first start Ivory asks for:

- **A secret word** — every stored secret (SSH keys, database passwords, LDAP and OIDC client secrets) is
  encrypted with it. See [Configuration → Secrets](configuration.md#secrets).
- **Authentication** — run without it locally, or enable Basic, LDAP or OIDC/SSO and register the first
  superuser. See [Authentication](authentication.md).

## 3. Add your first cluster

There are three ways in, depending on what already exists:

| You have                                   | Use             | Read more                                                                                        |
|--------------------------------------------|-----------------|--------------------------------------------------------------------------------------------------|
| A running cluster and every node's address | **Manual**      | [Clusters → Manual](clusters.md#manual-existing-cluster-known-node-addresses)                    |
| A running cluster and one node's address   | **Auto-detect** | [Clusters → Auto-Detect](clusters.md#auto-detect-existing-cluster-only-one-node-address-known)   |
| Machines, but no cluster yet               | **Deploy**      | [Deployment](deployment.md)                                                                      |

## 4. Start operating

- [Overview](overview.md) — node roles, replication lag, warnings, switchover and failover
- [Node](node.md) — VM metrics, containers, logs, configuration and the query console
- [MCP server](mcp.md) — connect Claude, Cursor, Codex or any MCP client to the same cluster
