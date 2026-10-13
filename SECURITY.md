# Security Policy

## Supported Versions

Security fixes land in the latest minor release only. Upgrade to it to receive them.

| Version | Supported                      |
|---------|--------------------------------|
| 1.x     | :triangular_flag_on_post: [no] |
| 2.0.x   | :triangular_flag_on_post: [no] |
| 2.1.x   | :checkered_flag: [yes]         |

## Upgrading and Backward Compatibility

There are two ways to carry your data to a new version
(see [Configuration → Upgrading](.doc/configuration.md#upgrading-to-a-new-version)):

- **Backup and restore** — export a backup from Settings and import it into the new version. Every backup
  format ever released is still readable, so this works across minor and major versions alike. Backups never
  contain passwords: users who sign in with a password need a fresh registration link after a restore.
- **Reusing the data directory** — mounting `/opt/ivory/data` into the new container. Always safe between
  patch releases; across a minor or major release only where the table below says so, because the internal
  storage format is allowed to change there.

| From  | To    | Backup restore            | Data directory              |
|-------|-------|---------------------------|-----------------------------|
| 1.0.x | 1.3.x | :cloud_with_lightning: [no] (backups start at 1.4.0) | :cloud_with_lightning: [no] |
| 1.4.x | 2.0.x | :sunny: [yes]             | :cloud_with_lightning: [no] |
| 2.0.x | 2.1.x | :sunny: [yes]             | :sunny: [yes]               |

## Reporting a Vulnerability

**Please do not report vulnerabilities in a public issue.** Use GitHub's private vulnerability reporting
(the repository's **Security → Report a vulnerability**), or write to [@anselvo](https://github.com/anselvo)
directly. Include the version, how to reproduce it and what an attacker gains. You will get an answer as soon
as the maintainer can look at it, and credit in the release notes unless you would rather not.

### What counts

Ivory deliberately runs commands on hosts: a deployment template, a container exec or a query console runs
exactly what an authorised user wrote. That is the product, not a hole in it — a user holding the matching
`manage.*` permission running a harmful command they can read on screen is not a vulnerability. These are:

- anything that lets someone act **without** the permission the action requires, or widen their own access;
- a command that executes differently from the one shown on screen (a substituted value that adds an
  argument, closes a quote or becomes syntax);
- a secret Ivory holds — vault passwords, SSH keys, tokens, the secret word — leaking into a log, an error,
  a backup, a job id or an API response;
- authentication flaws: forged or extended tokens, registration links that can be reused, bypassing a login
  provider.
