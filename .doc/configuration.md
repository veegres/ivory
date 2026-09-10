# [☰](../README.md) Configuration

How to run Ivory in production: where it keeps its data, how to upgrade it, and every environment
variable it reads.

## Data and Persistence

All Ivory data — clusters, queries, users, permissions, vaults, certificates, deployment templates —
lives in `/opt/ivory/data` inside the container. The image declares a Docker volume for that path, so
a container restart never loses it.

To keep the data across *different* containers you have two options:

- **Bind mount a local directory**

  ```bash
  docker run -p 80:80 \
    --mount type=bind,source=YOUR_LOCAL_PATH,target=/opt/ivory/data \
    --restart always veegres/ivory
  ```

- **Reuse the old container's volume**

  ```bash
  docker run -p 80:80 --volumes-from OLD_CONTAINER_NAME --restart always veegres/ivory
  ```

## Upgrading to a New Version

The recommended path is **Backup / Restore** from the Settings page. It exports your configuration —
users, clusters, queries, permissions, deployment templates — into a single file that imports into a
newer Ivory. The format is versioned and backward compatible, so an older backup still restores into
a newer release.

A backup says **who** your users are and how they sign in, and never carries a password or an
outstanding registration link. After a restore, hand everybody who signs in with a password a fresh
registration link from the _User Manager_ (see [Authentication](authentication.md)).

Mounting the data directory between containers also works, but reliably only for **patch** releases.
For minor and major upgrades, check the [supported versions page](../SECURITY.md) and prefer
Backup / Restore.

## Environment Variables

| Variable                   | Default                       | What it does                                                                           |
|----------------------------|-------------------------------|----------------------------------------------------------------------------------------|
| `IVORY_URL_ADDRESS`        | `:80`, or `:443` with TLS     | Address and port the server listens on                                                 |
| `IVORY_URL_PATH`           | `/`                           | Sub-path Ivory is served under, behind a reverse proxy. Must start with a `/`           |
| `IVORY_CERT_FILE_PATH`     | —                             | Path to the TLS certificate inside the container                                        |
| `IVORY_CERT_KEY_FILE_PATH` | —                             | Path to the TLS private key inside the container                                        |
| `IVORY_STATIC_FILES_PATH`  | set by the image entrypoint   | Where the frontend build is served from                                                 |
| `IVORY_VERSION_TAG`        | set by the build              | Version shown in the UI                                                                 |
| `IVORY_VERSION_COMMIT`     | set by the build              | Commit shown in the UI                                                                  |

### Running Under a Sub-path (Reverse Proxy)

Set `IVORY_URL_PATH` when Ivory sits behind nginx, Traefik or any other proxy under a prefix. The
path must start with a leading slash:

```bash
docker run -p 80:80 --env IVORY_URL_PATH=/ivory --restart always veegres/ivory
```

### Running Under TLS / HTTPS

Mount your certificate and key into the container — `/opt/certs` is the recommended location — and
point both variables at them. Ivory switches its default port to **443** as soon as both are set:

```bash
docker run -p 443:443 \
  --mount type=bind,source=YOUR_CERTS_PATH,target=/opt/certs \
  --env IVORY_CERT_FILE_PATH=/opt/certs/YOUR_CERT_NAME.crt \
  --env IVORY_CERT_KEY_FILE_PATH=/opt/certs/YOUR_KEY_NAME.key \
  --restart always veegres/ivory
```

## Secrets

Every secret Ivory stores — SSH keys, database and keeper passwords, LDAP and OIDC client secrets —
is encrypted with the **secret word** you choose during the initial setup wizard. Keep it: it is
required to unlock Ivory after a restart, and it is what makes it safe to hand Ivory your
credentials.
