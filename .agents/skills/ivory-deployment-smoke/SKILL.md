---
name: ivory-deployment-smoke
description: Run an explicitly requested live verification of Ivory's shipped single-host or multi-host deployment templates through its API, and produce the shared HTML report.
---

# Ivory deployment template smoke verification

Read AGENTS.md and select the requested mode:

- [Single host](references/single-host.md): shipped templates against one disposable Docker host.
- [Multiple hosts](references/multi-host.md): shipped templates against the existing three-host lab.

Read the selected procedure completely before execution, including its prerequisites, authorization boundaries, stopping conditions, cleanup policy and reporting contract. If both modes were requested, run both procedures and distinguish their results. A source-code change alone is not authorization to start a live deployment test.

Both modes use the [report template](assets/deploy-template-smoke-report.template.html) and [illustrative example](assets/deploy-template-smoke-report.example.html). The example is invented, not evidence of a successful run. Report as-shipped behavior before workarounds. Do not edit application source during these smoke workflows.

The procedures are retained from the former Claude agent definitions; they no longer require a particular model or provider. The root AGENTS.md and user's current task scope govern their use.
