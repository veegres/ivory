---
name: project-instructions
description: Create or update a repository's AGENTS.md, install selected shared skills, place project-specific skills, and add requested agent-client compatibility links without duplicating instructions.
---

# Project instructions and skill installation

Inspect existing instructions, source layout, package/build files, and working-tree changes first. Preserve domain constraints and local exceptions; never replace an existing AGENTS.md wholesale with a generic template.

## Locations and composition

- Shared skills are authored in the independent `agent-skills` repository under `skills/<name>/SKILL.md`.
- Install selected skills into a consumer's `.agents/skills/` with `python3 /path/to/agent-skills/scripts/sync.py /path/to/project --skills <names...>`. The script includes required companion skills, records a content digest and available Git revision, and refuses to overwrite local edits. Use `--check` to verify without writing. Installation never requires a sibling clone at runtime.
- Check in the installed files and `.agents/agent-skills.lock.json` with the project's instruction changes. Use a reviewed source commit/tag for later updates; this initial package may be installed before its first commit, so its content hash is authoritative until then.
- New project-only workflows belong in `.agents/skills/<name>/SKILL.md`, outside the managed names. Give each a concise `name` and trigger `description`; put long conditional procedures in `references/`, output templates in `assets/`, and reusable deterministic automation in `scripts/`.
- References to another skill are explicit reading instructions, not automatic dependency execution. Ensure referenced companion skills are installed. Do not tell a model to read every skill on every task.

## AGENTS.md content

Use [the template](assets/AGENTS.template.md) as a checklist, replacing every placeholder with verified project facts and omitting irrelevant sections. Include:

1. Product purpose and the root layout actually present.
2. Selected backend language, frontend framework, persistence profile, and project-specific contracts.
3. Links to the installed architecture skills and precise triggers for reading each.
4. Real commands, their working directories, required checks, and runtime prerequisites. Keep mandatory checks here so they do not depend on skill selection.
5. Domain invariants, exceptions, and project-only workflows. Link long details when appropriate; do not copy another project's facts or an entire shared skill.
6. A brief instruction to suggest a concrete convention improvement when real work exposes a gap, with a link to convention-maintenance.

## Client compatibility

AGENTS.md is the canonical project instruction source. Explain there that skills live only in `.agents/skills/<name>/SKILL.md`, and instruct the model to read the relevant entrypoint directly even when its client does not discover that directory automatically. Include task-specific links and triggers. Do not create `.claude/skills`, `.claude/agents`, or other provider-specific skill directories or delegation adapters.

Use AGENTS.md directly. Do not create provider-specific instruction aliases or duplicate instruction files. When removing legacy aliases, preserve any unique instructions in the canonical file first.

When migrating legacy provider-specific workflows, move their bodies/assets into neutral project skills, update the AGENTS.md links, and remove the replaced definitions after verifying preservation. Keep unrelated local client settings untouched. Moving a workflow does not authorize executing it.

Validate local Markdown links, symlink targets, frontmatter, selected profiles, and sync drift. Verify a fresh project checkout needs neither a global skill installation nor another application checkout. Follow [convention-maintenance](../convention-maintenance/SKILL.md) when the instruction design itself needs a new rule.
