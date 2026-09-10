# [☰](../README.md) Authentication & Users

Ivory runs with or without authentication — the initial setup wizard asks which. Locally you
usually don't need it; on a shared VM it is recommended.

## Sign-in Methods

| Method       | What it is                                              |
|--------------|---------------------------------------------------------|
| **Basic**    | Username and password held by Ivory itself              |
| **LDAP**     | Integration with an LDAP directory                      |
| **OIDC/SSO** | Single Sign-On through any OpenID Connect provider      |

Ivory needs the _profile_ or _email_ scope from an SSO provider to identify a person. Client secrets
and LDAP credentials are encrypted with your secret word before they are stored, so it is safe to
give them to Ivory.

## Registration Gates Sign-in

Whichever method you switch on, **everybody who signs in is an Ivory user first**. From the _User
Manager_ on the Settings page you register a username and pick which of the three ways that person
may sign in with. A name your directory knows but Ivory does not is refused — so an unexpected LDAP
or SSO account cannot walk in under a name somebody else holds.

**Nobody ever types somebody else's password.** Registering somebody for basic auth issues a
one-time link instead: the person opens it, sets their own password, and is signed in straight away.
The link works once, expires within hours, and can be revoked at any time — hand it only to the
person it names.

The **superuser** who sets Ivory up is the one exception: their password is typed during the initial
setup, since there is nobody yet to send a link to.

## Superusers

A superuser always holds every permission and cannot have any of them taken away, which is what
guarantees Ivory stays administrable. Only a superuser can register, change, reset or delete another
superuser, and there is always at least one — the last one cannot be deleted. Nobody can delete
themselves.

## What Can and Cannot Change

- A **username** is never changed, and the **superuser flag** is never taken back.
- What you can change about an existing user is **which ways they sign in**.
- The single thing a person can change about **themselves** is their own password, from the same
  _User Manager_, proven by the current one.
- Deleting a user removes their permissions too.

## Forgotten Passwords

From the _User Manager_ you **reset** the person's password, which hands them a new one-time link.
The user and everything granted to it are kept, and the password they already have keeps working
until they use the link — a reset is an offer, not a lockout. Only a superuser can reset a
superuser's password.

## Permissions

Ivory has a granular permission system — view or manage clusters, execute queries, manage
configurations, deploy containers, read node metrics, and so on. Permissions are managed per user
from the _Permission Manager_ on the Settings page once authentication is enabled.

## Running Without Authentication

Running Ivory without authentication switches all of this off: there is no account to change and
nobody to grant anything to, so neither the _User Manager_ nor the _Permission Manager_ is offered
at all, and every action is permitted.
