### [All]
- Choose licence
- Add README.md
- Add docker container with configured (storage)
- Add versioning and push image to docker hub
- Add opportunity to use tools as pgcompacttable and pg_repack
- Add opportunity to collect and see different Postgres Checkups
- Think about support Citus, Postgres-XL, etc (can we?)
- Think about support other tools as Patroni (PAF, repmgr, Stolon, etc)

### [Backend]
- Improve architecture and package structure

### [Web]
- Refactor all nested function components <> to render() or extract to separate file (diff behave could cause a problems)
- Add remove button for empty elements in the center of nodes array
- Refactor all strings from `""` to `''` and props in jsx from `=""` to `={''}`
- Add more strict Eslint rules (imports, spaces near objects `{}`, etc), Hooks for git and Prettier
