# React, TypeScript, and MUI profile

This is the selected style of these projects, not a universal React requirement.

- Use PascalCase source filenames; feature API files are `<Name>Type.ts`, `<Name>Router.ts`, and `<Name>Hook.ts`. Root non-colliding files need no feature prefix.
- Declare hooks and state before `return`. Below it, order function declarations as `render*`, action `handle*`, hook implementations `handleEffect*`/`handleMemo*`/`handleCallback*`, then `get*`. Extract a non-trivial hook callback into its named function; keep the hook call unconditional in the component body. Do not place executable initialization after a return.
- Extract multi-line JSX in props into `render*` helpers. Short inline expressions are fine.
- Each component owns one `SX` style constant, typed with the project's style helper where present. Do not share `SX` across components; extract the repeated UI into a reusable component. Use theme tokens instead of hard-coded colors.
- Application absence is `undefined`. Translate external JSON `null` at the boundary when needed; this convention does not change a catalogue file's documented JSON schema.
- Where the project supplies a shared text-input component, use it for fields and select/autocomplete rendering, preserving documented special controls. Otherwise follow its existing MUI input composition; do not add a wrapper merely to match another project's component names. In a new project, implement any selected wrapper before using it.
- Use the existing mutation adapter for error handling, optional success messages, and invalidation. Preserve intentional query-key semantics; investigate a dependency warning before changing keys or adding a narrowly justified suppression supported by the configured linter.
- A dialog shell owns trigger/title/close/back navigation. Each screen owns both its scrolling content and sibling actions area. One discriminated step union and one render switch keep content, actions, and back behavior together. Initialize screen state from props when mounting with ready data; reset on actual close handlers, not effects that react to freshly allocated step objects.
- Existing shell names differ: Ivory uses `DialogButton`, Veembly uses `DialogShell`. Consult local components rather than renaming them to match a sample.
- Rstest/Testing Library tests are colocated as `*.test.ts(x)`. Helpers, hooks, providers, and pure UI calculations are behavior worth testing. These projects normally omit presentational component tests unless requested; preserve existing tests.
- Let the repository's lint configuration enforce double quotes, no semicolons, JSX braces, object spacing, and import order. Do not duplicate its full rule set here.
