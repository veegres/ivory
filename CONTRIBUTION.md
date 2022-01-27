### [General]
- [ ] Choose licence
- [ ] Add README.md
- [ ] Add docker container with configured (storage)
- [ ] Add versioning and push image to docker hub

### [Backend]
- [ ] Improve architecture and package structure
- [ ] Add default value for patroni ports

### [Web]
- [ ] Refactor all nested function components <> to render() or extract to separate file (diff behave could cause a problems)
- [ ] Add remove button for empty elements in the center of nodes array
- [ ] Refactor all strings from `""` to `''` and props in jsx from `=""` to `={''}`
- [ ] Add more strict Eslint rules (imports, spaces near objects `{}`, etc), Hooks for git and Prettier

### [Long Term Plans]
- [ ] Add opportunity to use tools as [pgcompacttable](https://github.com/dataegret/pgcompacttable) and [pg_repack](https://github.com/reorg/pg_repack)
- [ ] Add opportunity to collect and see different [Postgres Checkups](https://gitlab.com/postgres-ai/postgres-checkup)
- [ ] Think about support other tools as Patroni - [PgPool 2](https://www.pgpool.net/), [PAF](http://clusterlabs.github.io/PAF/), [repmgr](https://repmgr.org/), [Stolon](https://github.com/sorintlab/stolon), etc
- [ ] Think about support sharding tools as [Citus](https://www.citusdata.com/), [Postgres-XL](https://www.postgres-xl.org/), etc (can we?)
