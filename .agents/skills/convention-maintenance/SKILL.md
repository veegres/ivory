---
name: convention-maintenance
description: Resolve gaps or conflicts in adopted engineering conventions and propose or implement narrowly scoped updates when real project work establishes a new pattern.
---

# Maintain conventions from evidence

Conventions are editable engineering decisions. They should help complete the user's task, not prevent it because a case has not been documented.

1. Identify the actual gap, contradiction, or repeated workaround. Compare AGENTS.md, the applicable skill/profile, configuration, and nearby implementation. An unfamiliar feature alone is not evidence of a new convention.
2. Continue authorized work using the closest established pattern and language/framework idiom. Ask only when an unresolved decision materially changes product behavior, compatibility, or scope; absence of a written rule is not itself a reason to stop.
3. Propose a concrete change when the gap affected a decision: the observed case, suggested wording, owning file, rationale, and projects it should affect. Include this briefly in the final report. Do not repeat settled suggestions or recommend changes after tasks that exposed no gap.
4. Classify the rule: product facts and exceptions belong in project AGENTS.md or a project skill; cross-language responsibilities in architecture core; language/framework behavior in the relevant reference; repeatable new workflows in a focused skill; mechanically enforceable rules in tooling.
5. Explicit user direction and established authorization govern updates. If asked to maintain conventions as part of the task, implement the supported change and summarize it. Otherwise suggest a broader policy change rather than silently imposing it on every project. A user's explicit correction can be documented at its stated scope without another confirmation round.

## Shared source and installed copies

Projects vendor selected skills under `.agents/skills/`; `.agents/agent-skills.lock.json` records the source and file hashes. Edit shared guidance in the `agent-skills` source repository, then use its sync script for projects included in the task. Never silently edit a managed copy and leave its manifest stale. If the source is unavailable, give a proposed patch or record a local exception instead of assuming a sibling checkout exists.

Keep a single source for each rule. When a core rule changes, inspect its profiles and affected project exceptions for conflicts. Prefer replacing obsolete wording over appending another universal prohibition. Test concrete scenarios in both an existing consumer and any new profile; structural validation alone cannot prove the instruction makes good decisions.

Project-specific workflows live in that project's `.agents/skills/<name>/SKILL.md`. Promote one into this package only when its reusable boundary is clear. Explain skill locations in AGENTS.md rather than creating provider-specific skill directories or adapters. Use AGENTS.md directly without provider-specific instruction aliases; the neutral skills remain authoritative.
