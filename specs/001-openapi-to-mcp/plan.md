# Implementation Plan: OpenAPI-to-MCP CLI

**Branch**: `001-openapi-to-mcp` | **Date**: 2025-02-06 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `specs/001-openapi-to-mcp/spec.md`

## Summary

Build a Go CLI that accepts an OpenAPI 3.x spec (file or URL) and generates a Node project in the caller’s directory that implements an MCP server using fastmcp. Mapping from OpenAPI operations to MCP tools/resources is deterministic (rules + templates); no LLM. Output directory must be empty unless `--force`; OpenAPI 2.0 is rejected with a clear message. Constitution requires clean code, clean architecture, unit tests for domain/generators, and integration tests for CLI and generated artifact.

## Technical Context

**Language/Version**: Go 1.21+ (CLI); generated output is Node (TypeScript/JavaScript) with fastmcp  
**Primary Dependencies**: Go: OpenAPI 3.x parser (e.g. kin-openapi or go-openapi for read-only parse); Node (generated): fastmcp, package.json  
**Storage**: N/A (CLI is stateless; reads OpenAPI, writes files to output dir)  
**Testing**: Go: `go test`; unit tests for parser, mapper, generator; integration tests for CLI binary + generated Node project runnable  
**Target Platform**: OS-agnostic CLI (Linux, macOS, Windows); generated MCP runs on Node  
**Project Type**: Single project (CLI binary + internal packages)  
**Performance Goals**: Generation completes in &lt;30s for typical spec (e.g. &lt;50 paths)  
**Constraints**: No LLM/external AI; deterministic, reproducible output; offline-capable CLI  
**Scale/Scope**: Single binary; one OpenAPI input → one Node project output per invocation  

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Gate | Status |
|-----------|------|--------|
| I. Clean Code | Names reflect intent; small functions; no duplication; complexity justified | Pass (plan enforces structure) |
| II. Clean Architecture | Domain (parse, map, generate) independent of CLI/I/O; dependencies inward | Pass (see Project Structure) |
| III. Product, Stability & Quality | Predictable failures; clear exit codes/messages; no breaking changes without migration | Pass (spec FR-005, FR-006) |
| IV. Unit Testing | Domain logic, parsers, generators covered; deterministic, isolated tests | Pass (tests/unit/) |
| V. Integration Testing | CLI entrypoint + OpenAPI→MCP e2e; generated Node project runnable | Pass (tests/integration/) |

## Project Structure

### Documentation (this feature)

```text
specs/001-openapi-to-mcp/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (CLI contract, generated output contract)
└── tasks.md             # Phase 2 output (/speckit.tasks - not created by plan)
```

### Source Code (repository root)

```text
cmd/
└── openapi2mcp/         # CLI entrypoint (main.go)
    └── main.go

internal/
├── domain/              # Core: no I/O, no CLI
│   ├── openapi/         # OpenAPI 3.x parse (in-memory model)
│   ├── mapping/         # OpenAPI operation → MCP tool/resource (rules)
│   └── model/           # Shared structs (Operation, MCPTool, etc.)
├── generator/           # Node project generation (templates + code)
│   └── node/            # fastmcp Node project emitter
├── cli/                 # Flags, args, exit codes, stderr
│   └── cmd.go
└── fetch/               # Optional: URL fetch for OpenAPI (injectable)

templates/               # Embedded or file-based Node project templates
└── node-fastmcp/        # package.json, index.ts, tool stubs

tests/
├── unit/                # internal/domain, internal/generator, internal/mapping
├── integration/         # CLI → generate → npm install && npm start; optional: MCP client list/invoke
└── contract/            # Generated output: package.json valid, required files, required deps (see contracts/generated-output-guarantees.md)
```

**Guarantees for generated Node project**: Structure, libs, and “works” are guaranteed by (1) versioned template + fixed deps in generated package.json; (2) contract tests (structure + package.json + deps); (3) integration tests (npm install + npm start + optional MCP list/invoke). See `specs/001-openapi-to-mcp/contracts/generated-output-guarantees.md`.

**Structure Decision**: Single Go module; `cmd/openapi2mcp` for entrypoint; `internal/` for all logic (domain, generator, cli, fetch). Clean architecture: `domain` and `mapping` have no I/O; `cli` and `fetch` are adapters. `templates/` holds Node project skeleton for deterministic codegen. Tests mirror constitution: unit for domain/generator, integration for CLI and generated MCP.

## Complexity Tracking

*(No violations; structure aligns with constitution.)*
