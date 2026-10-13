# [☰](../README.md) FAQ

**How do I upgrade to a new version?**
Back up your configuration from the Settings page and restore it in the new version — the format is
backward compatible. Mounting the data directory between containers works too, but reliably only for patch
releases. See [Configuration → Upgrading](configuration.md#upgrading-to-a-new-version).

**Where does Ivory store its data?**
In `/opt/ivory/data`, backed by a Docker volume. Bind-mount it, or use `--volumes-from`, to keep data across
different containers. See [Configuration → Data and Persistence](configuration.md#data-and-persistence).

**Does Ivory need an agent on my database hosts?**
No. It talks to the Keeper's API, to the database directly, and to the host over SSH.

**How does authentication work?**
Ivory runs with or without it. With it, everybody who signs in is an Ivory user first — you register a username
and choose whether they sign in with a password, LDAP or SSO, and a password is always set by the person
themselves through a one-time link. See [Authentication](authentication.md).

**Are my credentials safe?**
Every secret — SSH keys, database passwords, LDAP and OIDC client secrets — is encrypted with the secret word
you choose during setup. See [Configuration → Secrets](configuration.md#secrets).

**Can I run Ivory behind a reverse proxy or under HTTPS?**
Yes — `IVORY_URL_PATH` for a sub-path, `IVORY_CERT_FILE_PATH` and `IVORY_CERT_KEY_FILE_PATH` for TLS. See
[Configuration → Environment Variables](configuration.md#environment-variables).

**Can an AI agent use Ivory?**
Yes, through the built-in MCP server, gated by the same per-user permissions as the UI. See [MCP server](mcp.md).

**Does it work on a phone?**
Yes. The v2 UI is mobile friendly, so you can check cluster health and run HA operations from your phone.
