## Docker environment for local development

###### This Dockerfile and docker-compose will run for you patroni cluster with 3 instances.

1. Edit hosts (`/etc/hosts`), add these lines to the file

``` 
127.0.0.1 patroni1
127.0.0.1 patroni2
127.0.0.1 patroni3
```

2. Run Docker

```
docker-compose up -d
```

**Patroni Rest API:** `patroni[1-3]:8008` (http://patroni1:8008)

**Postgres ports:** `are not exported` (you can do it by yourself in docker-compose file)

**HAProxy statistics:** `localhost:8408`
