## Ideas and things that would be cool to do

### [General]

- [x] Choose licence
- [ ] Add option to change data folder for docker mount
- [x] Add versioning and push image to docker hub

### [Backend]

- [x] Improve architecture and package structure
- [ ] Add default value for patroni ports
- [ ] Bloat streaming
  - [ ] Implement logs sync when stream is started late
  - [ ] Add jobs filter per cluster

### [Web]

- [ ] Change main cluster table
  - [ ] add request for each cluster to show nodes role in that table
  - [ ] make possibility to open bloat, cluster once or only for master
  - [ ] add action buttons to overview tab
  - [ ] when stream is finished we should reload logs on open button
- [ ] Refactor all nested function components <> to render() or extract to separate file (diff behave could cause a problems)
- [ ] Add remove button for empty elements in the center of nodes array
- [ ] Refactor all strings from `""` to `''` and props in jsx from `=""` to `={''}`
- [ ] Add more strict Eslint rules (imports, spaces near objects `{}`, etc), Hooks for git and Prettier
- [ ] Think about mocks for web, do we need them for easy development

### [Long Term Plans]

- [ ] Add password storage for patroni clusters and postgres nodes (for pgcompact table)
- [ ] Add opportunity to add some grafana boards
- [X] Add opportunity to use tool [pgcompacttable](https://github.com/dataegret/pgcompacttable)
- [ ] Add opportunity to use tool [pg_repack](https://github.com/reorg/pg_repack)
- [ ] Add opportunity to collect and see different [Postgres Checkups](https://gitlab.com/postgres-ai/postgres-checkup)
- [ ] Think about support other tools as Patroni - [PgPool 2](https://www.pgpool.net/), [PAF](http://clusterlabs.github.io/PAF/), [repmgr](https://repmgr.org/)
  , [Stolon](https://github.com/sorintlab/stolon), etc
- [ ] Think about support sharding tools as [Citus](https://www.citusdata.com/), [Postgres-XL](https://www.postgres-xl.org/), etc (can we?)
