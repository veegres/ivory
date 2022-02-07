## Docker environment for local development

This Dockerfile and docker-compose will run for you patroni cluster with 3 instances.

### Run

1. Add to your hosts (`/etc/hosts`) these lines to the file

``` 
127.0.0.1 patroni1
127.0.0.1 patroni2
127.0.0.1 patroni3
```

2. Run Docker Compose

```
docker-compose up -d
```

### Connection

- **Patroni Rest API:** `patroni[1-3]:8008` (example http://patroni1:8008)
- **Postgres ports:** `are not exposed` (you can do it by yourself in docker-compose file)
- **HAProxy statistics:** `localhost:8408`