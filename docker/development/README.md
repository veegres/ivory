## Docker environment for local development

This Dockerfile and docker-compose will run for you patroni cluster with 3 instances.

## Run

```
1. docker-compose build
2. docker-compose up -d
```

## Connection

- **Patroni Rest API:** `localhost:[8001-8005]` (example http://patroni1:8008)
- **Postgres Ports:** `localhost:[5001-5005]` (example `psql --host=localhost --port=5003 --username=postgres`, `password=password`)
- **HAProxy Statistics:** `localhost:8408`
- **Additional Patroni API** `patroni[1-3,-cert,-pass]:8008` you can add them to your hosts
  (`/etc/hosts`) these lines to the file
  ``` 
  127.0.0.1 patroni1
  127.0.0.1 patroni2
  127.0.0.1 patroni3
  127.0.0.1 patroni-cert
  127.0.0.1 patroni-pass
  ```
### Requirements
- `patroni-cert:8004` - cluster is needed to check patroni certificates and required client certificates
- `patroni-pass:8005` - cluster is needed to check patroni password and required password for unsafe requests


## pgcompacttable

1. Install Perl DBI library `apt install libdbi-perl libdbd-pg-perl`
2. Download pgcompacttable `curl -o /usr/bin/pgcompacttable https://raw.githubusercontent.com/dataegret/pgcompacttable/master/bin/pgcompacttable`
3. Make file executable `chmod +x /usr/bin/pgcompacttable`

Now you can use pgcompacttable from your console, Ivory will use it as well

## Certificates

These certificates are only needed for testing purposes, please, don't use them
in real VMs, you can follow instruction and generate them. Certificates were generated
for 24855 days (End date 4 May 2091), probably it should be enough not to change them :)

### Generate Certificates

First go to the certs package `cd certs`

- Certificate Authority
  - Private Key `openssl ecparam -name prime256v1 -genkey -noout -out ca/ca.key`
  - Certificate `openssl req -new -x509 -sha256 -key ca/ca.key -subj "/O=Ivory" -out ca/ca.crt -days 24855`
- Server
  - Private Key `openssl ecparam -name prime256v1 -genkey -noout -out server/patroni-cert.key`
  - Signing Request `openssl req -new -sha256 -key server/patroni-cert.key -subj "/CN=patroni-cert/O=Ivory" -addext "subjectAltName=DNS:patroni-cert" -out server/patroni-cert.csr`
  - Certificate `openssl x509 -req -in server/patroni-cert.csr -CA ca/ca.crt -CAkey ca/ca.key -CAcreateserial -extfile <(printf "subjectAltName=DNS:patroni-cert") -out server/patroni-cert.crt -days 24855 -sha256`
- Client
  - Private Key `openssl ecparam -name prime256v1 -genkey -noout -out client/client.key`
  - Signing Request `openssl req -new -sha256 -key client/client.key -subj "/CN=development/O=Ivory" -out client/client.csr`
  - Certificate `openssl x509 -req -in client/client.csr -CA ca/ca.crt -CAkey ca/ca.key -CAcreateserial -out client/client.crt -days 24855 -sha256`

### Check Certificates

- Certificate Authority
  - CRT `openssl x509 -noout -text -in ca/ca.crt`
- Server
  - Verify `openssl verify -CAfile ca/ca.crt server/patroni-cert.crt`
  - CSR `openssl req -noout -text -in server/patroni-cert.csr`
  - CRT `openssl x509 -noout -text -in server/patroni-cert.crt`
- Client
  - Verify `openssl verify -CAfile ca/ca.crt client/client.crt`
  - CSR `openssl req -noout -text -in client/client.csr`
  - CRT `openssl x509 -noout -text -in client/client.crt`

### Package structure

```
certs
├── ca
│   ├── ca.key            -- Certificate Authority (CA) Private Key
│   └── ca.crt            -- Certificate Authority (CA) Certificate
├── server
│   ├── [SERVER_NAME].key -- Server Certificate Private Key
│   ├── [SERVER_NAME].scr -- Server Certificate Signing Request
│   └── [SERVER_NAME].crt -- Server Certificate
└── client
    ├── client.key        -- Client Certificate Private Key
    ├── client.scr        -- Client Certificate Signing Request
    └── client.crt        -- Client Certificate
```
