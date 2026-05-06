# Repository Guidelines

## Project Structure & Module Organization
Ivory is split into a Go backend in `service/` and a Vite/React frontend in `web/`. Backend packages live under `service/src/`, organized by concern: `features/` for business logic, `clients/` for integrations, `storage/` for persistence, and `app/` for wiring. Frontend source lives in `web/src/`, with `api/` for client adapters, `component/` for UI, `provider/` and `hook/` for shared state and behavior, and `style/` for CSS modules. Tests live in `service/src/**/_test.go` and `web/test/`. Docs and screenshots are under `doc/`; local Docker setups are under `docker/`.

**Mandatory Check**: When modifying frontend code, always run `cd web && npm run lint` and `cd web && npm run build` to ensure no linting or TypeScript compilation errors were introduced.

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

## Architectural Patterns
- **Vertical Consolidation**: Prefer grouping all management logic for a single entity into one feature. For example, the `node` feature manages both the Host/VM (SSH, Docker, metrics) and the Database/Keeper (HA, config).
- **Service Splitting**: To keep code manageable, split large services by concern using a `service_<domain>.go` naming convention (e.g., `service_db.go` and `service_host.go`).
- **Generic Clients**: Centralize transport-level logic in generic clients (like `clients/http` or `clients/ssh`) and keep domain-specific logic in wrappers or consumers.

## Nomenclature
- **Node**: Refers to a single server entity (Hardware + Software).
- **Keeper**: Refers to the HA management tool (e.g., Patroni).
- **VM**: Refers specifically to the virtual machine/host level.

## Coding & UI Standards
- **Surgical Renaming**: When asked to rename components or fields, ONLY update the terminology. DO NOT refactor implementation logic, move methods to props, or change the component structure unless explicitly directed.
- **Data Synchronization**: When syncing frontend types with backend changes, perform the MINIMAL necessary updates to ensure compilation. DO NOT redesign or refactor UI components during a data sync task. If a type change appears to require a major UI redesign, ASK for clarification first.
- **Component Style**: Follow the established pattern of using dedicated helper methods (e.g., `handleAction`, `renderContent`, `handleEffect`) within functional components to keep the JSX clean.
  - **Hook Convention**: If a React hook (like `useEffect`, `useMemo`, `useCallback`) is not a simple one-liner, move its implementation to a dedicated function below the `return` statement and name it `handleEffect***`, `handleMemo***`, etc.
  - **JSX in Props**: Never write multi-line JSX directly in component props. Instead, create a dedicated `render***` function (e.g., `renderDescription`) below the `return` statement. One-liner JSX is allowed.
  - **Styles Naming**: Always name the constant for component-specific styles as `SX`. Do not use descriptive names like `VAULT_NEW_SX`.
  - **Shared Styles**: Never share `SX` styles between components. Sharing styles is a signal to create a separate reusable component (e.g., instead of sharing a `code` style, create a `Code` component in `view/box`).
  - **Member Order**: Maintain the following order within functional components:
    1. `return` (JSX)
    2. `render***` methods
    3. `handle***` methods
    4. `handleEffect***`, `handleMemo***`, etc. (specific React hooks)
    5. `get***` methods
- **Cache Management**: Be cautious when updating React Query keys. If a linter warns about a missing dependency that was intentionally omitted to control cache invalidation, use `// eslint-disable-next-line` instead of changing the key signature.
- **TypeScript Nullability**: Never use `null` in TypeScript code; always use `undefined` for representing missing or optional values.

## Backward Compatibility & Persistence
- **Backup Models are Sacred**: Structs used for backups (e.g., `Backup`, `backupCluster`) must remain unchanged in their nomenclature and JSON tags to ensure long-term compatibility. If changes are necessary, a new backup version must be introduced.
- **Internal DB Migration**: Internal database models (e.g., `Cluster` in `features/cluster/model.go`) can be refactored and their JSON tags updated. The backup tool is the primary mechanism for users to migrate data between versions if the internal schema breaks.
- **Data Integrity**: When renaming internal models, always verify if they are part of the backup/export logic before changing their JSON representation.

## Testing Guidelines
- **Zero Deletion Policy**: NEVER delete existing tests during refactoring or type synchronization. If types change, UPDATE the tests to match the new structures. A "refactor" that results in fewer tests is a failure.
- **Backend Tests**: Backend tests use Go's standard `testing` package with table-driven tests and `t.Run()` subtests. Focus on storage adapters, client mappers, and service logic; avoid thin routers and external network calls.
- **Frontend Tests**: Frontend tests use Vitest and Testing Library. Mirror `web/src/` inside `web/test/`, for example `web/src/hook/Debounce.ts` -> `web/test/hook/Debounce.test.ts`. No coverage threshold is defined, but new logic should ship with targeted tests.

## Commit & Pull Request Guidelines
Recent commits use short, imperative summaries such as `add tooltips for refresh buttons` and `fix problem with different column widths in list`. Keep commit subjects concise and specific; add a body when the change needs context. Pull requests should describe behavior changes clearly, link relevant issues, and include screenshots or GIFs for UI work. Before opening a PR, run the relevant backend tests, frontend lint, and affected frontend tests.
