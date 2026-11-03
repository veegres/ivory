---
name: frontend-architecture
description: Design or change a feature-oriented frontend, including feature ownership, API boundaries, shared UI, routing, and state. Select framework-specific details from the project's profile.
---

# Frontend architecture

Read AGENTS.md for the framework, design system, routing, and exceptions. For React/TypeScript/MUI projects, also read [the React profile](references/react-typescript-mui.md). Other frameworks retain the boundaries below but use their native component and state idioms; do not introduce React conventions into them.

- `app/core/` owns application composition, routed pages, and widgets. `app/features/<name>/` owns a business capability's API and UI. `app/shared/` holds reusable components, helpers, hooks, providers, tests, and assets. Add a standalone `tools/` area only when the product has independent tools.
- Feature-owned types belong with that feature's API. Cross-feature dependencies follow the consuming direction. Shared code does not own unrelated domain models just because two features use them.
- Keep network adaptation, transport errors, and mutation/cache behavior in API helpers and the project's established adapters. Presentation consumes those contracts rather than inventing a second request mechanism.
- Use existing components and theme tokens for common behavior. Shared visual patterns should become components; a global collection of copied style fragments is not a component system.
- Specify URL ownership, navigation/history behavior, persisted preferences, and ephemeral UI state deliberately. Preserve the project's route semantics rather than assuming every filter must use the same history behavior.
- Keep loading, empty results, missing resources, and errors distinguishable. Public-page metadata must not interpret network failures as missing content.
- When changing API types or terminology, keep UI updates scoped to the requested behavior; do not redesign the interface as a side effect.
- Verify helpers, state, hooks, and providers at their behavioral boundaries. Use the project's policy for presentational tests, lint, build, and accessibility verification. Keep configured lint directory lists in sync with new source roots.

For an actual missing convention, use [convention-maintenance](../convention-maintenance/SKILL.md) and propose a core, framework-profile, or project-local addition. Keep product branding and component names that differ by project in AGENTS.md.
