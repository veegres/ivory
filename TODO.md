## Ideas and things that would be cool to do

### [General]

- [x] Choose licence
- [ ] Add option to change data folder for docker mount
- [x] Add versioning and push image to docker hub
- [x] Add opportunity to save password and defer them from postgres & patroni
- [ ] Add opportunity to support certs for patroni and maybe postgres (do we need it?)
- [ ] Add opportunity to use/check more commonly used queries for managing DB (bloat of index, bloat of table, active queries, etc)
  - [ ] active queries - place them in instance block to be able to check queries for each instance
  - [ ] bloat queries - place them in bloat tab to be able to switch and easily start bloat cleaning
  - [ ] other queries - make separate tab for them, think where to place it and how we can add customisation feature for them
- [ ] Add check of active queries to Instance block and function to kill them
- [ ] Support more patroni features like (schedule times, do switchover for direct instance, reload, etc)
- [ ] Add opportunity set up new cluster automatically by entering name and one of the instance
- [ ] Add tags to cluster, we should be able to filter by them

### [Backend]

- [ ] Change packages layout for secret & credentials
- [ ] Move Worker from Job to new service
- [ ] Add default value for patroni ports
- [x] Bloat streaming
  - [x] implement logs sync when stream is started late
  - [x] add jobs filter per cluster
- [ ] Refactor deserialization with using Marshal

### [Web]

- [ ] Change main cluster table
  - [x] add request for each cluster to show nodes role in that table
  - [x] make possibility to open bloat, cluster once or only for master
  - [ ] add action buttons to Instance overview tab
  - [ ] when stream is finished we should reload logs on open button
- [ ] Add navigation (router + urls)
- [ ] Add footer with information about the project, link to github and year, etc
- [ ] Refactor all nested function components `<>` to `render()` or extract to separate file (diff behave could cause a problems)
- [ ] Add remove button for empty elements in the center of nodes array
- [ ] Refactor all strings from `''` to `""` and props in jsx from `=""` to `={""}`
- [ ] Add more strict Eslint rules (imports, spaces near objects `{}`, etc), Hooks for git and Prettier
- [ ] Think about mocks for web, do we need them for easy development
- [ ] Move all shown strings, text to one file
- [ ] Improve cluster detection by patroni name without port (pros and cons)

### [Long Term Plans]

- [ ] Stability
  - [ ] add tests to backend, think about frontend
  - [ ] cover all edge cases (errors, loadings, etc)
- [x] Add password storage for patroni clusters and postgres nodes (for pgcompact table)
  - [x] add secret to encrypt passwords
  - [x] add passwords storage with interface
  - [x] use secret to show and do action (reinit, switchover, clean bloat, etc)
- [ ] Add backups that can be manage by Ivory, think how to restore it (pg_backrest - point in time, pg_dump - easier, pg_restore etc)
- [ ] Add some basics live time graphics to Node Overview (like memory, disk size, core, cpu, etc)
- [ ] Add opportunity to add some grafana graphs
- [x] Add opportunity to use tool [pg_compacttable](https://github.com/dataegret/pgcompacttable)
- [ ] Add opportunity to use tool as [pg_repack](https://github.com/reorg/pg_repack) or [pg_squeeze](https://github.com/cybertec-postgresql/pg_squeeze)
- [ ] Add opportunity to collect and see different [Postgres Checkups](https://gitlab.com/postgres-ai/postgres-checkup)
- [ ] Think about support other tools as Patroni - [PgPool 2](https://www.pgpool.net/), [PAF](http://clusterlabs.github.io/PAF/), [repmgr](https://repmgr.org/), [Stolon](https://github.com/sorintlab/stolon), etc
- [ ] Think about support sharding tools as [Citus](https://www.citusdata.com/), [Postgres-XL](https://www.postgres-xl.org/), etc (can we?)
