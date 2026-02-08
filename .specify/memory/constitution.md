<!--
  Sync Impact Report
  Version change: (none) → 1.0.0
  Modified principles: (initial creation)
  Added sections: Core Principles (5), Quality Standards, Development Workflow, Governance
  Removed sections: (none)
  Templates: plan-template.md ✅ (Constitution Check unchanged, gates derived from this file); spec-template.md ✅ (no change); tasks-template.md ✅ (task types align with testing principles); commands ✅ (no agent-specific references)
  Follow-up TODOs: (none)
-->

# OpenAPI-to-MCP CLI Constitution

## Core Principles

### I. Clean Code

Code MUST be readable, self-explanatory, and maintainable. Names (variables, functions, types) MUST reflect intent. Functions MUST be small and do one thing. Duplication MUST be removed in favor of reusable abstractions. Complexity MUST be justified; when in doubt, prefer the simpler solution. Comments explain "why" when the code cannot; avoid comments that restate "what". This principle ensures long-term maintainability and reduces defect rate.

### II. Clean Architecture

The codebase MUST follow clean architecture boundaries: dependencies point inward (e.g. domain does not depend on delivery or infrastructure). Core business rules (parsing OpenAPI, generating MCP) MUST live in a domain/core layer; CLI and I/O (files, network) MUST be in outer layers. External dependencies (libraries, frameworks) MUST be injectable or behind interfaces so that tests and future changes do not require rewrites. This principle protects product stability and enables testability.

### III. Product, Stability & Quality

Decisions MUST consider the product as a whole: user value, reliability, and long-term quality over short-term speed. Stability is non-negotiable: avoid breaking changes without a clear migration path; failures MUST be predictable and reported clearly (exit codes, messages). Quality is measured by correctness, test coverage, and operational behavior; technical debt MUST be tracked and reduced, not accumulated without plan. This principle keeps the CLI trustworthy and evolvable.

### IV. Unit Testing

Unit tests are REQUIRED for domain logic, parsers, generators, and any code that contains branching or business rules. Tests MUST be deterministic, fast, and isolated (no real I/O unless explicitly integration). Coverage targets SHOULD be defined per module; critical paths MUST be covered. Tests are written as part of implementation (test-first or test-along); code that cannot be unit-tested MUST be justified and minimized. This principle guards against regressions and enables safe refactoring.

### V. Integration Testing

Integration tests are REQUIRED for: CLI entrypoints (invoking the binary with args), OpenAPI parsing + MCP generation end-to-end, and any contract between components (e.g. generated Node project runnable by Node). Integration tests MAY use real filesystem and subprocesses; they MUST be clearly separated from unit tests and run in a dedicated phase or suite. New features that change external behavior or generated output MUST add or update integration tests. This principle ensures the product works as a whole and that generated artifacts remain valid.

## Quality Standards

- Code review MUST verify compliance with Clean Code and Clean Architecture before merge.
- New code that adds or changes behavior MUST include or update unit and/or integration tests as appropriate.
- Linting and formatting (e.g. gofmt, static analysis) MUST pass; configuration is versioned in the repo.
- Breaking changes to CLI flags, output format, or generated project structure MUST be documented and versioned (e.g. changelog, semver).

## Development Workflow

- Work is driven by specs and plans; implementation follows tasks derived from user stories.
- Before marking a task done: implementation exists, relevant tests exist and pass, and Constitution Check (plan phase) is satisfied.
- Refactors that touch multiple layers SHOULD be preceded by or accompanied by tests to preserve behavior.
- Use the constitution as the source of truth for "how we build"; when in doubt, prefer the principle that favors stability and quality.

## Governance

- This constitution supersedes ad-hoc practices for this project. All PRs and reviews MUST verify compliance with the principles above.
- Amendments require: a proposed change, rationale, and impact on existing specs/plans/tasks. Version MUST be bumped (MAJOR for principle removal/redefinition, MINOR for new principles/sections, PATCH for wording/clarifications).
- Complexity or exceptions (e.g. skipping tests in a narrow case) MUST be documented and justified in the code or design docs.
- Last ratification is the date this version was adopted; Last Amended is updated on every constitution change.

**Version**: 1.0.0 | **Ratified**: 2025-02-06 | **Last Amended**: 2025-02-06
