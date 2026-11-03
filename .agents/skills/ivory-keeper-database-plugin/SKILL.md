---
name: ivory-keeper-database-plugin
description: Add, fix, or audit Ivory keeper/database adapters against integrity, healthiness, informativeness, and supported actions. Use for adapter work in this repository.
---

# Ivory keeper and database plugins

Read the repository's AGENTS.md, including its plugin checklist, and the shared backend skill with its Go profile. Then read [the workflow](references/workflow.md) for the adapter contract, audit/add/fix modes, and report requirements.

Use audit mode for a review and add/fix mode for authorized implementation. This workflow is specific to Ivory; it is not a generic backend architecture rule. Preserve the existing checks and live-environment boundaries. Loading this skill does not authorize launching a live stack.
