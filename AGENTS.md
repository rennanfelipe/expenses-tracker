# AGENTS.md

## Project Context

Expenses Tracker is a Go application for importing Brazilian bank statements, managing expense categories, and generating reports.

Consult `README.md` for the project overview and setup instructions.

## Mandatory Instructions

Before implementing changes, read:

- [Architecture Rules](ai/architecture-rules.md)
- [Coding Standards](ai/coding-standards.md)
- [Testing Guidelines](ai/testing-guidelines.md)

These documents define the project's architectural constraints, coding conventions, and testing requirements.

## Project Documentation

Consult the relevant documentation before making changes:

- `docs/business.md` — Business context, requirements, and rules.
- `docs/architecture.md` — System architecture and structural documentation.
- `docs/data-model.md` — Database entities and relationships.
- `docs/decisions/` — Architecture Decision Records (ADRs).

Database definitions and deployment configuration are maintained in `db-schemas/`.

## Working Instructions

- Inspect relevant source code and documentation before implementing changes.
- Implement only requested or documented functionality.
- Preserve existing behavior unless explicitly instructed otherwise.
- Ask for clarification when requirements are ambiguous and affect implementation.
- Do not introduce architectural layers, frameworks, dependencies, or major abstractions without justification and approval.
- Do not overwrite unrelated user changes.
- Do not commit or push changes unless explicitly requested.

## Completion

After implementation, report:

- What was changed.
- Which files were affected.
- Which tests and verification commands were executed and their results.
- Any remaining limitations or unresolved issues.

Never claim that a command or test succeeded unless it was actually executed.