# Repository guidelines

<Product purpose.>

## Architecture and shared skills

Backend: <language/framework>; persistence: <engine and profile>.
Frontend: <language/framework/design system>.

For backend changes, read [backend-architecture](.agents/skills/backend-architecture/SKILL.md) and its selected language profile.
For frontend changes, read [frontend-architecture](.agents/skills/frontend-architecture/SKILL.md) and its selected framework profile.
For changes to these instructions or skill installation, read [project-instructions](.agents/skills/project-instructions/SKILL.md).
When work exposes a convention gap, use [convention-maintenance](.agents/skills/convention-maintenance/SKILL.md): suggest a concrete rule at the appropriate scope; continue authorized work using established patterns.

Skills live only in `.agents/skills/<name>/SKILL.md`. Read the applicable entrypoints directly even if your client does not automatically discover them. Do not create provider-specific skill directories or adapters. Use AGENTS.md directly without provider-specific instruction aliases.

Shared rules live in the source package identified by `.agents/agent-skills.lock.json`. Project-specific requirements below refine them. Do not edit managed skill copies independently.

## Layout

<Actual directories and responsibilities; include only what exists.>

## Commands and required checks

<Exact commands, working directories, prerequisites, and which changes require each check.>

## Project contracts and exceptions

<API, authentication, routing, persistence/backup, domain invariants, and local style exceptions.>

## Project workflows

<Links and triggers for project-only skills, if any.>
