# Tasks: OpenAPI-to-MCP CLI

**Input**: Design documents from `specs/001-openapi-to-mcp/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Constitution and FR-008 require unit tests (domain, parsers, generators) and integration + contract tests for generated output. Test tasks are included.

**Organization**: Tasks are grouped by user story so each story can be implemented and tested independently.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: User story (US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Repo root**: `cmd/`, `internal/`, `templates/`, `tests/` at repository root (see plan.md)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and Go module structure per plan.md

- [x] T001 Create project structure: cmd/openapi2mcp/, internal/domain/openapi/, internal/domain/mapping/, internal/domain/model/, internal/generator/node/, internal/cli/, internal/fetch/, templates/node-fastmcp/, tests/unit/, tests/integration/, tests/contract/
- [x] T002 Initialize Go module (go mod init) and add dependency kin-openapi for OpenAPI 3.x parsing at repo root
- [x] T003 [P] Configure linting and formatting: gofmt, golangci-lint or staticcheck in Makefile or CI (e.g. Makefile or .golangci.yml at repo root)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core domain, mapping, template, and generator. MUST be complete before any user story.

**Independent Test**: Unit tests for domain and generator pass; template and generator can produce a valid Node project from a list of MCPTools.

- [x] T004 Define shared structs in internal/domain/model/: Operation, MCPTool (Name, Description, InputSchema), and any types needed by openapi and mapping packages
- [x] T005 Implement OpenAPI 3.x parsing in internal/domain/openapi/: parse from reader (YAML/JSON), return in-memory model; reject OpenAPI 2.0 with clear error; use kin-openapi
- [x] T006 Implement mapping in internal/domain/mapping/: Operation → MCPTool (tool name from operationId or path+method, inputSchema from parameters/requestBody per data-model.md)
- [x] T007 Create Node project template in templates/node-fastmcp/: package.json (name, version, type module, scripts.start, dependency fastmcp with fixed version), entry script stub (e.g. index.ts or index.js) that creates FastMCP and registers tools
- [x] T008 [P] Unit tests for internal/domain/openapi: parse valid OpenAPI 3.x, reject OpenAPI 2.0, reject invalid JSON/YAML in tests/unit/
- [x] T009 [P] Unit tests for internal/domain/mapping: operation to tool name and inputSchema in tests/unit/
- [x] T010 Implement generator in internal/generator/node/: emit package.json and entry script from []MCPTool using template; write to given output directory (interface for filesystem)
- [x] T011 [P] Unit tests for internal/generator/node: generate produces valid package.json and entry script for given MCPTools in tests/unit/

**Checkpoint**: Domain and generator ready; unit tests pass. User story implementation can begin.

---

## Phase 3: User Story 1 – Generate MCP from OpenAPI in current directory (Priority: P1) – MVP

**Goal**: User runs CLI with OpenAPI file (or URL) and gets a runnable Node MCP project in current directory (or specified dir). One command → npm install && npm start works.

**Independent Test**: Run CLI with valid OpenAPI 3.x file; cwd (or -o dir) contains generated project; npm install && npm start succeeds; MCP client can list/invoke tools.

### Implementation for User Story 1

- [x] T012 Implement CLI flags and args in internal/cli/: positional openapi-input (file path), -o/--output (default cwd), -f/--force; parse and validate
- [x] T013 Implement read OpenAPI from file path in internal/cli or internal/fetch: open file, pass to domain/openapi parser; return error if file not found or unreadable
- [x] T014 Implement full flow in internal/cli: Run(input, outputDir, force) → read input, parse OpenAPI, check output dir empty (unless force), map operations to MCPTools, call generator/node to write files; exit 0 on success
- [x] T015 Wire cmd/openapi2mcp/main.go: parse os.Args via cli, call cli.Run, exit with code from cli (stderr for errors)
- [x] T016 [P] [US1] Contract tests in tests/contract/: after generating with fixture OpenAPI, assert package.json exists and is valid JSON, required fields (name, scripts.start, dependency fastmcp), entry script exists
- [x] T017 [US1] Integration test in tests/integration/: build CLI binary, run with fixture OpenAPI 3.x file into temp dir, run npm install in generated dir (exit 0), run npm start with short timeout (process starts); fail if any step fails

**Checkpoint**: User Story 1 complete. CLI generates runnable Node MCP from file input; contract and integration tests pass.

---

## Phase 4: User Story 2 – Specify input and output paths (Priority: P2)

**Goal**: User can pass -o/--output for target directory; user can pass URL as openapi-input and CLI fetches it.

**Independent Test**: CLI with -o /path/to/dir and file input generates only in /path/to/dir; CLI with URL as input fetches and generates.

- [ ] T018 [P] [US2] Implement fetch in internal/fetch: HTTP GET OpenAPI from URL with timeout (e.g. 30s), return body or error; injectable for tests
- [ ] T019 [US2] In internal/cli: when input looks like URL (e.g. http:// or https://), use fetch to get body then parse; otherwise read from file path; support -o/--output (already in T012, ensure used in Run)
- [ ] T020 [US2] Integration test in tests/integration/: CLI with -o pointing to temp dir and file input → generated files only in that dir; CLI with public OpenAPI URL as input → generation succeeds

**Checkpoint**: User Story 2 complete. Explicit output dir and URL input work.

---

## Phase 5: User Story 3 – Clear failure and validation feedback (Priority: P3)

**Goal**: Invalid or unsupported input yields non-zero exit and clear, actionable stderr message (invalid OpenAPI, OpenAPI 2.0, file not found, output dir not empty, no mappable operations).

**Independent Test**: CLI with invalid OpenAPI, OpenAPI 2.0 file, missing file, or non-empty dir without -f exits non-zero and stderr contains identifiable message.

- [ ] T021 [P] [US3] Ensure OpenAPI 2.0 detection in internal/domain/openapi returns distinct error type or message (e.g. "OpenAPI 2.0 is not supported; use OpenAPI 3.x")
- [ ] T022 [US3] In internal/cli: map domain/parser and generator errors to exit codes 1–4 and stderr messages per contracts/cli-interface.md (1 invalid/unsupported, 2 input not found/fetch failed, 3 output dir not empty, 4 no mappable operations)
- [ ] T023 [US3] Unit tests in tests/unit/: cli or domain returns appropriate errors for invalid OpenAPI, OpenAPI 2.0, missing file
- [ ] T024 [US3] Integration tests in tests/integration/: CLI with invalid OpenAPI file → exit non-zero, stderr contains "invalid" or similar; CLI with OpenAPI 2.0 file → exit non-zero, stderr contains "2.0" or "not supported"; CLI with missing file path → exit non-zero, stderr indicates file/URL not found; CLI with non-empty output dir and no -f → exit non-zero, stderr indicates dir not empty and --force

**Checkpoint**: User Story 3 complete. All error paths have clear exit codes and messages.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, quickstart validation, and final quality pass.

- [ ] T025 [P] Add README or docs with Node/npm minimum version (e.g. Node 18+, npm 9+) and link to specs/001-openapi-to-mcp/quickstart.md
- [ ] T026 Run quickstart.md flow manually or via script: build CLI, generate from fixture, npm install && npm start in generated dir; document any deviations
- [ ] T027 [P] Run gofmt and linter on entire codebase; fix any issues

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: No dependencies – start immediately.
- **Phase 2 (Foundational)**: Depends on Phase 1 – BLOCKS all user stories.
- **Phase 3 (US1)**: Depends on Phase 2 – MVP.
- **Phase 4 (US2)**: Depends on Phase 2; can start after or in parallel with US1 (needs T012/T014 for -o and fetch).
- **Phase 5 (US3)**: Depends on Phase 2; can start after or in parallel with US1/US2 (error handling and exit codes).
- **Phase 6 (Polish)**: Depends on Phases 3–5 being done (or at least US1).

### User Story Dependencies

- **US1 (P1)**: After Foundational – no dependency on US2/US3.
- **US2 (P2)**: After Foundational – extends CLI with -o and URL; independently testable.
- **US3 (P3)**: After Foundational – extends CLI with exit codes and messages; independently testable.

### Parallel Opportunities

- Phase 1: T003 [P] with T001/T002.
- Phase 2: T008, T009, T011 [P] can run in parallel after T004–T007, T010.
- Phase 3: T016 [P] contract tests can run in parallel with T017 once T014–T015 are done.
- Phase 4: T018 [P] fetch implementation in parallel with T019.
- Phase 5: T021 [P], T023 [P] unit tests in parallel.
- Phase 6: T025 [P], T027 [P] in parallel.

---

## Implementation Strategy

### MVP First (User Story 1 only)

1. Complete Phase 1 (Setup).
2. Complete Phase 2 (Foundational) – CRITICAL.
3. Complete Phase 3 (US1): CLI with file input, default cwd, generate → npm install && npm start.
4. **STOP and VALIDATE**: Contract + integration tests pass; quickstart flow works.
5. Demo: one command from OpenAPI file to runnable MCP.

### Incremental Delivery

1. Setup + Foundational → unit tests pass.
2. Add US1 → contract + integration tests pass → MVP.
3. Add US2 → -o and URL input → test independently.
4. Add US3 → exit codes and messages → test independently.
5. Polish → README, quickstart validation, lint.

### Task Count Summary

| Phase            | Task IDs   | Count |
|------------------|------------|-------|
| Phase 1 Setup   | T001–T003  | 3     |
| Phase 2 Foundational | T004–T011 | 8  |
| Phase 3 US1      | T012–T017  | 6     |
| Phase 4 US2      | T018–T020  | 3     |
| Phase 5 US3      | T021–T024  | 4     |
| Phase 6 Polish   | T025–T027  | 3     |
| **Total**        |            | **27**|

**Suggested MVP scope**: Phases 1 + 2 + 3 (T001–T017). Independent test for MVP: CLI + fixture OpenAPI → generated dir → npm install && npm start succeeds.
