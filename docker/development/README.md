## Docker environment for local development

This Dockerfile and docker-compose will run for you patroni cluster with 3 instances.

### Run

```
1. docker-compose build
2. docker-compose up -d
```

### Set up pgcompacttable

```
1. Install Perl DBI library (Debian: apt install libdbi-perl libdbd-pg-perl)
2. curl -o /usr/bin/pgcompacttable https://raw.githubusercontent.com/dataegret/pgcompacttable/master/bin/pgcompacttable
3. chmod +x /usr/bin/pgcompacttable
```
Now you can use pgcompacttable from your console, Ivory will use it as well

### Connection

- **Patroni Rest API:** `localhost:[8001-8003]` (example http://patroni1:8008)
- **Postgres Ports:** `localhost:[5001-5003]` (example `psql --host=localhost --port=5003 --username=postgres`, `password=password`)
- **HAProxy Statistics:** `localhost:8408`


- **Additional Patroni API** `patroni[1-3]:8008` you can add by adding to your hosts
  (`/etc/hosts`) these lines to the file
  ``` 
  127.0.0.1 patroni1
  127.0.0.1 patroni2
  127.0.0.1 patroni3
  ```
