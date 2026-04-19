# Repository Guidelines

## Project Structure & Module Organization
Ivory is split into a Go backend in `service/` and a Vite/React frontend in `web/`. Backend packages live under `service/src/`, organized by concern: `features/` for business logic, `clients/` for integrations, `storage/` for persistence, and `app/` for wiring. Frontend source lives in `web/src/`, with `api/` for client adapters, `component/` for UI, `provider/` and `hook/` for shared state and behavior, and `style/` for CSS modules. Tests live in `service/src/**/_test.go` and `web/test/`. Docs and screenshots are under `doc/`; local Docker setups are under `docker/`.

## Build, Test, and Development Commands
Backend:
- `make -C service build` builds `service/build/ivory`.
- `make -C service test` runs `go test -C src ./...`.
- `make -C service start` runs the built backend binary.

Frontend:
- `cd web && npm install` installs dependencies.
- `cd web && npm start` starts the Vite dev server.
- `cd web && npm run build` runs TypeScript compile plus production build.
- `cd web && npm run lint` runs ESLint over `src/` and `test/`.
- `cd web && npm test` starts Vitest; `npm run test:coverage` runs coverage.

Use `docker/ivory-dev/` for the stack.

## Coding Style & Naming Conventions
Follow existing file naming: Go packages are lowercase; frontend files use PascalCase for components (`QueryTable.tsx`) and descriptive camel-case names for hooks/providers. In TypeScript, prefer `interface` for shared types, `type` for component props, and function declarations for named React components. Frontend lint rules enforce double quotes, no semicolons, sorted imports, and React Hook discipline. Keep styling in MUI `sx` objects, CSS modules, or top-level constants rather than inline JSX clutter.

## Testing Guidelines
Backend tests use Go's standard `testing` package with table-driven tests and `t.Run()` subtests. Focus on storage adapters, client mappers, and service logic; avoid thin routers and external network calls. Frontend tests use Vitest and Testing Library. Mirror `web/src/` inside `web/test/`, for example `web/src/hook/Debounce.ts` -> `web/test/hook/Debounce.test.ts`. No coverage threshold is defined, but new logic should ship with targeted tests.

## Commit & Pull Request Guidelines
Recent commits use short, imperative summaries such as `add tooltips for refresh buttons` and `fix problem with different column widths in list`. Keep commit subjects concise and specific; add a body when the change needs context. Pull requests should describe behavior changes clearly, link relevant issues, and include screenshots or GIFs for UI work. Before opening a PR, run the relevant backend tests, frontend lint, and affected frontend tests.
