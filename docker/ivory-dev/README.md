## Docker environment for local development

This Dockerfile and docker-compose will run for you patroni cluster with 3 cluster. One of them is usual with 3
instances, others are needed only for specific test purposes and have specific configuration (like certs or pass).
All postgres instances are configured with ssl by default, and there are different users for test purposes.

## Run

```
1. docker-compose build
2. docker-compose up -d
```

## Connection

- **Patroni Rest API:** `localhost:[8001-8005]` (example http://patroni1:8003)
- **Postgres Ports:** `localhost:[5001-5005]` (example `psql --host=localhost --port=5003 --dbname=postgres --username=admin`, `password=admin`)
- **HAProxy Statistics:** `localhost:8404` (example http://localhost:8404)
- **Additional Patroni API** `patroni[1-3,-cert,-pass]:8008` you can add them to your hosts
  (`/etc/hosts`) these lines to the file
  ``` 
  127.0.0.1 patroni1
  127.0.0.1 patroni2
  127.0.0.1 patroni3
  127.0.0.1 patroni-cert
  127.0.0.1 patroni-pass
  ```
  
### Clusters

- `patroni[1-3]:[8001-8003]` - cluster is needed for general test
- `patroni-cert:8004` - cluster is needed to check patroni certificates and required client certificates
- `patroni-pass:8005` - cluster is needed to check patroni password and required password for unsafe requests

### Users

#### Database
- `superuser` - used by patroni
- `replicator` - used by patroni for replication set up
- `rewind` - used by patroni for rewind set up
- `trust` - custom user, it doesn't have any password for connection
- `admin` - custom user, it has general password `admin:admin`
- `sslca` - custom user, it requires `verify-ca` connection and has password `sslca:sslca`
- `sslfull` - custom user, it requires `verify-full` connection and has password `sslfull:sslfull`

#### Patroni
- `patroni` - it is used to connect to cluster `patroni-pass:8005` with password `patroni:patroni

### Data

In each cluster there is initial data set with tables
- `users` - it has 1000 rows
- `films` - it has 10000 rows

## pgcompacttable

Note: You need root privilege to use commands below, use `sudo`

1. Install Perl DBI library `apt install libdbi-perl libdbd-pg-perl`
2. Download pgcompacttable `curl -o /usr/bin/pgcompacttable https://raw.githubusercontent.com/dataegret/pgcompacttable/master/bin/pgcompacttable`
3. Make file executable `chmod a+xr /usr/bin/pgcompacttable`

Now you can use pgcompacttable from your console, Ivory will use it as well

## Authentication

Ivory supports 3 type of authentications *BASIC*, *LDAP*, *OIDC*.
You can configure them in initial Ivory load, right after secret set.
Configuration examples:

### BASIC

Just use simple username and password, like `admin:admin`

### LDAP

You can either up local one or use some test ldap providers from the internet.
The one used during development and recommend for use is [ldap.forumsys.com](https://www.forumsys.com/2022/05/10/online-ldap-test-server/).
It is free to use and has open configuration.

#### Configuration
```
URL:            ldap://ldap.forumsys.com
Bind DN:        cn=read-only-admin,dc=example,dc=com
Bind password:  password
Base DN:        dc=example,dc=com
Filter:         (uid=%s)
```

#### Credentials
```
tesla:password
newton:password
einstein:password
euclid:password
```

P.S. for more credentials and options visit providers page

### OIDC

You can either up local one or use some test ldap providers from the internet.
During development only closed one was discovered, but it is free to use.
Provider is [auth0.com](https://auth0.com). You can either create your own provider use that was created already.

#### Configuration

```
Issuer URL:     https://dev-xmg2j154iy3amd7t.us.auth0.com/
Client ID:      lMhKMAPfaEa18DjlPmpnCoOkTSLFF6qL
Client Secret:  EebgFuO_DWPwIlOf6w3JonqOBiikjZYFf74kFjshPp4au-buzmp2yqfx34sfdSzx
```

#### Credentials
```
tesla:password
newton:password
einstein:password
euclid:password
```

#### Create your own provider

1. Register and login to [auth0.com](https://auth0.com)
2. Go to the [dashboard](https://manage.auth0.com/dashboard)
3. Go to application page or find button create application
4. Choose any and provide some name like `Ivory [DEV]`
5. Open the application (there you would see configs to use)
6. Find Application URIs section and add to Allowed callbacks, your local callback `http://localhost:5173/api/oidc/callback`
7. Check that you have users in Authentication DB (you can register new or use Google auth)


## Certificates

These certificates are only needed for testing purposes, please, don't use them in real VMs, you can
follow instruction and generate them. Certificates were generated for 24855 days (End day is 4 May 2091),
probably it should be enough not to keep them for some time :)

### Generate Certificates

First go to the certs package `cd certs`

- Certificate Authority
  - Private Key `openssl ecparam -name prime256v1 -genkey -noout -out ca/ca.key`
  - Certificate `openssl req -new -x509 -sha256 -key ca/ca.key -subj "/O=Ivory" -out ca/ca.crt -days 24855`
- Client
  - Private Key `openssl ecparam -name prime256v1 -genkey -noout -out client/client.key`
  - Signing Request `openssl req -new -sha256 -key client/client.key -subj "/CN=development/O=Ivory" -out client/client.csr`
  - Certificate `openssl x509 -req -in client/client.csr -CA ca/ca.crt -CAkey ca/ca.key -CAcreateserial -out client/client.crt -days 24855 -sha256`
- Server
  - Private Key `openssl ecparam -name prime256v1 -genkey -noout -out server/server-cert.key`
  - Signing Request `openssl req -new -sha256 -key server/server-cert.key -subj "/CN=server-cert/O=Ivory" -addext "subjectAltName=DNS:server-cert" -out server/server-cert.csr`
  - Certificate `openssl x509 -req -in server/server-cert.csr -CA ca/ca.crt -CAkey ca/ca.key -CAcreateserial -extfile <(printf "subjectAltName=DNS:server-cert") -out server/server-cert.crt -days 24855 -sha256`

### Check Certificates

- Certificate Authority
  - CRT `openssl x509 -noout -text -in ca/ca.crt`
- Client
  - Verify `openssl verify -CAfile ca/ca.crt client/client.crt`
  - CSR `openssl req -noout -text -in client/client.csr`
  - CRT `openssl x509 -noout -text -in client/client.crt`
- Server
  - Verify `openssl verify -CAfile ca/ca.crt server/server.crt`
  - CSR `openssl req -noout -text -in server/server.csr`
  - CRT `openssl x509 -noout -text -in server/server.crt`


### Package structure

```
certs
├── ca
│   ├── ca.key     -- Certificate Authority (CA) Private Key
│   └── ca.crt     -- Certificate Authority (CA) Certificate
├── client            (used in ivory ui)
│   ├── client.key -- Client Certificate Private Key
│   ├── client.scr -- Client Certificate Signing Request
│   └── client.crt -- Client Certificate
└── server            (used in any server postgres, patroni, etc)
    ├── server.key -- Server Certificate Private Key
    ├── server.scr -- Server Certificate Signing Request
    └── server.crt -- Server Certificate
```
