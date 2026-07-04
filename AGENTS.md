# Repository Guidelines

## Project Structure & Module Organization
Ivory is split into a Go backend in `server/` and a Vite/React frontend in `app/`. Backend packages are organized by concern: `features/` for business logic, `clients/` for integrations, `storage/` for persistence, and `core/` for wiring. Frontend source lives in `app/`, with `features/` for domain logic, `core/` for main pages and widgets, `shared/` for shared components, hooks, and providers. Tests live in `server/**/_test.go` and `app/test/`. Docs and screenshots are under `.doc/`; local Docker setups are under `.docker/`.

**Mandatory Check**: When modifying frontend code, always run `cd app && npm run lint` and `cd app && npm run build` to ensure no linting or TypeScript compilation errors were introduced.

## Build, Test, and Development Commands
Backend:
- `make -C server build` builds `server/build/ivory`.
- `make -C server test` runs `go test ./...`.
- `make -C server start` runs the built backend binary.

Frontend:
- `cd app && npm install` installs dependencies.
- `cd app && npm start` starts the Vite dev server.
- `cd app && npm run build` runs TypeScript compile plus production build.
- `cd app && npm run lint` runs ESLint over `core/`, `features/`, `shared/` and `test/`.
- `cd app && npm test` starts Vitest; `npm run test:coverage` runs coverage.

Use `.docker/ivory-dev/` for the stack.

## Architectural Patterns
- **Vertical Consolidation**: Prefer grouping all management logic for a single entity into one feature. For example, the `node` feature manages both Platform access (metrics, deployment operations) and the Database/Keeper (HA, config).
- **Service Splitting**: To keep code manageable, split large services by concern using a `service_<domain>.go` naming convention (e.g., `service_db.go` and `service_host.go`).
- **Generic Clients**: Centralize transport-level logic in generic clients (like `clients/http` or `clients/ssh`) and keep domain-specific logic in wrappers or consumers.
- **Platform Integration Boundary**: Keep `plugins/platform` generic. Platform adapters model deployment and troubleshooting primitives for a target platform, not Ivory-specific database behavior. Database-specific interpretation belongs in features such as `features/cluster`.

## Feature API Boundary (Frontend)
Each frontend feature exposes its own types in `features/<name>/api/type.ts`. Types must live in the feature that owns them — never in a shared "plugins" or "common" file. Cross-feature imports are allowed only in the consuming direction (e.g., `cluster` may import from `node` and `query`; `query` and `node` must not import from `cluster`).

Current ownership:
- `KeeperPlugin` enum → `features/node/api/type.ts` (keeper plugin selector is a node-level concept)
- `DbPlugin` enum → `features/query/api/type.ts` (database engine type is part of the DB connection config)
- `DbConfig` interface → `features/query/api/type.ts` (database connection config belongs to the query feature)

Backend mirrors this: `server/features/node/model.go` defines `KeeperPlugin = keeper.Plugin` (alias); `server/features/query/model.go` defines `DbPlugin = database.Plugin` (alias) and `DbConfig` struct. Plugin packages (`server/plugins/*`) are internal implementation details — features expose their own types and use mapper functions to translate to/from plugin types at the boundary.

## Platform Architecture
- **Platform Package**: `server/plugins/platform` is the backend integration point for deployment targets. Existing implementations live under `server/plugins/platform/<adapter>/`, for example `linux/`. Future adapters such as Kubernetes or OpenShift should be added as sibling packages.
- **Platform vs Transport**: Platform is the product-level abstraction. SSH is only a transport/client detail for the Linux adapter and credential flow. Keep names like `clients/ssh`, `sshPort`, `SSH_KEY`, and SSH credential UI labels when they describe the actual transport or stored credential type.
- **Linux Adapter**: `platform/linux` implements platform deployment operations for on-prem Linux hosts by executing Docker commands over SSH. It is split into exactly two implementation files by interface: `vm_manager.go` for `VmManager` methods (host metrics, logs, process listing, SSH key provisioning) and `container_manager.go` for `ContainerManager` methods (container lifecycle and container metrics). Docker command helpers and tests may still use Docker-specific names internally when they are truly Docker CLI details.
- **Generic Platform Methods**: Platform adapters expose generic deployment operations: `KeeperNodeList`, `KeeperDeploy`, `KeeperStop`, `KeeperDelete`, and `KeeperLogs`, returning `OperationResult`. Do not name these methods after Docker containers, Kubernetes pods, or databases.
- **OperationResult**: Use `OperationResult` for command/API execution output with `stdout`, `stderr`, and `exitCode`. It is intentionally not called Container, Pod, Workload, or Database because it is only an operation result.
- **Node Platform Boundary**: The node feature uses prefixed service methods such as `PlatformContainerUp`, `PlatformContainerStop`, `PlatformContainerStart`, `PlatformContainerDown`, `PlatformContainerList`, and `PlatformContainerLogs` to avoid collisions with existing Keeper methods like `KeeperNodeList`. Frontend node APIs mirror this as `NodeApi.deployment.*` and `useRouterNodePlatform*`.
- **Routes and Features**: Platform routes live under `/node/platform/...`, currently `/node/platform/metrics`, `/node/platform/logs`, `/node/platform/copy-id`, and `/node/platform/container/...`. Permissions use `view.node.platform.*` and `manage.node.platform.*`.
- **Domain-Specific Logic**: `features/cluster` may keep database-focused names such as `normalizeDatabaseOptions` and cluster `Deploy`, because that feature converts Ivory database cluster configuration into generic platform deployment requests.
- **Keeper Metadata**: `keeper.Metadata` is a separate interface from `keeper.Adapter`: Adapter covers operations against a running keeper, Metadata covers plugin self-description — `SupportedFeatures()` and `DeploymentSpec()` (platform-agnostic image requirements: env, ports, volumes, default image). The split lets plugins without a management API yet (plain postgres, `plugins/keeper/postgres`) still declare features and deployment defaults. `database` and `tools` plugins keep `SupportedFeatures` on their Adapter for now.
- **Deployment Defaults**: Default deploy options are split across two plugin vocabularies to avoid an n×m matrix: keeper plugins declare requirements via `keeper.Metadata.DeploymentSpec()`, and platform adapters implement `RenderOptions(platform.DeploySpec) string` (native, user-editable options text — docker flags for `linux`). `features/node` maps between them (`mapKeeperDeploymentToPlatformSpec`) and serves the result via `GET /node/platform/container/deploy-options`. Keeper plugins must never know about platforms and vice versa.
- **Request Naming**: Use neutral platform request/response names at the node/platform boundary, such as `PlatformUpRequest`, `PlatformLogsRequest`, `PlatformVaultConnection`, and `PlatformResponse`. Request fields like `name` should identify the target deployment generically; adapters map that to a Docker container name, Kubernetes resource name, pod name, etc.

## Nomenclature
- **Node**: Refers to a single server entity (Hardware + Software).
- **Keeper**: Refers to the HA management tool (e.g., Patroni).
- **VM**: Refers specifically to the virtual machine/host level.
- **Platform**: Refers to the deployment/troubleshooting target that Ivory integrates with, such as on-prem Docker-over-SSH now and Kubernetes/OpenShift in the future.
- **Deployment**: Refers to the generic platform-managed thing Ivory deploys, stops, deletes, lists, or reads logs from. Avoid using `Container`, `Pod`, or `DatabaseRegistry` for generic platform APIs.
- **Job**: Refers to a generic command execution entity managed by a background goroutine, supporting event-driven status updates, console logs persistence, and live streaming.

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
- **Frontend Tests**: Frontend tests use Vitest and Testing Library. Tests are colocated with the source files they test, following the pattern `[filename].test.[ext]`. No coverage threshold is defined, but new logic should ship with targeted tests.

## Commit & Pull Request Guidelines
Recent commits use short, imperative summaries such as `add tooltips for refresh buttons` and `fix problem with different column widths in list`. Keep commit subjects concise and specific; add a body when the change needs context. Pull requests should describe behavior changes clearly, link relevant issues, and include screenshots or GIFs for UI work. Before opening a PR, run the relevant backend tests, frontend lint, and affected frontend tests.
