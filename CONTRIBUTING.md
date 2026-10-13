# How to contribute

I'm really glad you're reading this, because we need volunteer developers to help this project come to fruition.

If you haven't already, you can write to us ([#ivory](https://postgresteam.slack.com/archives/C0304HP1ZSM) on postgres slack). We want you working on things you're excited about.

Here are some important resources:

  * [Enhancements](https://github.com/veegres/ivory/issues)
  * [Good for newcomers](https://github.com/veegres/ivory/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)
  * [Supported versions](SECURITY.md)
  * [Setup Local Environment](.docker/ivory-dev/README.md)
  * [Build Frontend](app/README.md)
  * [Build Backend](server/README.md)

## Coding conventions

Ivory has a strict coding convention, to prevent disputes in PRs and keep the code readable. The complete,
up-to-date rulebook is **[AGENTS.md](AGENTS.md)** — architecture, naming, file layout, model and service
rules, plugin contracts, backup compatibility and testing. It is the single source of truth: this file only
summarises it, and where the two disagree, AGENTS.md wins.

It is long on purpose, and it is written to be read by AI coding assistants as much as by people. Claude Code
(through `CLAUDE.md`, a link to it), Codex, Cursor and most other agents load it automatically. **If you use
one, ask it to follow AGENTS.md and to check your change against it before you open a PR** — that is the
fastest way to get the conventions right. Reviewers will hold the change to it either way.

The rules that come up most often:

- **Structure** — Go backend in `server/` (`features/` business logic, `clients/` integrations, `plugins/`
  keeper/database/platform adapters, `storage/` persistence, `core/` wiring); React frontend in `app/`
  (`features/` domain logic, `core/` pages and widgets, `shared/` reusable components, hooks and providers).
- **Go file names** are prefixed with their package: `node/node_service_keeper.go`, `cluster_router_http.go`.
  Domain types and mappers live in `<feature>_model.go`; split large services by concern as
  `<feature>_service_<domain>.go`.
- **Frontend file names** are PascalCase, `.ts` and `.tsx` alike; a feature's types live in
  `features/<name>/api/<Name>Type.ts`, owned by the feature they belong to.
- **Components** are functions with a fixed member order: `return` JSX, then `render***`, `handle***`,
  `handleEffect***`/`handleMemo***`, `get***`. Styles go in one constant named `SX` above the component —
  never inline `sx={{...}}`, never shared between components. `Props` is a `type` above the component.
- **TypeScript** uses `undefined`, never `null`.
- **Comments** only where the *why* is not obvious; no multi-line blocks.
- **Backups** are sacred: a released backup version's shape never changes, a change means a new version.

## Testing

Every change should come with tests and pass the same checks CI runs:

- **Backend** — `make -C server test`. Table-driven Go tests with `t.Run()` subtests, in a `_test.go` file
  that mirrors the file it tests. Services and repositories are tested; routers are not.
- **Frontend** — `cd app && npm run lint && npm run build && npm test`. Rstest + Testing Library, colocated
  as `[filename].test.[ext]`. Test helpers, hooks and providers; skip presentational components.

Never delete an existing test during a refactor — update it to the new shape. For anything touching
deployment or a plugin, also try it by hand against the [local environment](.docker/ivory-dev/README.md).

## Submitting changes

Please send a [GitHub Pull Request to Ivory](https://github.com/veegres/ivory/compare) with a clear list of what you've done (read more about [pull requests](http://help.github.com/pull-requests/)). Describe behaviour changes, link the related issue, and add screenshots or GIFs for UI work.

Always write a clear log message for your commits. Short imperative subjects (`add tooltips for refresh buttons`) are fine for small changes, but bigger changes should look like this:

    $ git commit -m "A brief summary of the commit
    > 
    > A paragraph describing what changed and its impact."

All new proposals about changing the coding convention should be discussed among contributors, and land as an
edit to AGENTS.md.


Thanks,
Andrei Sergeev
