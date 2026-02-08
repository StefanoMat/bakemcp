# Feature Specification: OpenAPI-to-MCP CLI

**Feature Branch**: `001-openapi-to-mcp`  
**Created**: 2025-02-06  
**Status**: Draft  
**Input**: User description: "The project aims to be a Go CLI. The core feature is to receive an OpenAPI spec as input and output in the directory of who calls the CLI an MCP server using github.com/punkpeye/fastmcp."

## Clarifications

### Session 2025-02-06

- Q: When the output directory already contains files (e.g. previous generation), what should the CLI do? → A: Refuse if the directory is not empty unless the user passes a force/overwrite flag (e.g. `--force`).
- Q: Is OpenAPI 2.0 (Swagger) in scope for v1? → A: Out of scope for v1; support only OpenAPI 3.x; reject 2.0 with a clear "unsupported version" message.
- Q: What form should the generated artifact take? → A: A Node project (source) that uses the fastmcp library; the user runs it with Node (e.g. npm install && npm start).
- Q: How is the OpenAPI → MCP endpoints mapping done? → A: Deterministic code construction (rules + templates); no LLM or external AI.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Generate MCP from OpenAPI in current directory (Priority: P1)

A developer has an OpenAPI specification (file or URL) and wants to get a runnable MCP server in their current working directory so they can immediately use it with MCP clients (e.g. IDEs, agents) without manual wiring.

**Why this priority**: This is the core value: one command turns an existing API contract into an MCP that exposes it. Delivering this alone is a usable MVP.

**Independent Test**: Can be fully tested by running the CLI with a valid OpenAPI spec and verifying that the current directory contains a generated MCP that can be started and exposes the API operations. Delivers value as a standalone code generator.

**Acceptance Scenarios**:

1. **Given** a valid OpenAPI 3.x file on disk, **When** the user runs the CLI with that file as input and no output path, **Then** the current working directory contains a generated Node project (e.g. package.json, source using fastmcp) that reflects the OpenAPI operations and is runnable with Node (e.g. npm install && npm start).
2. **Given** the user is in an empty directory, **When** they run the CLI with an OpenAPI URL, **Then** the same directory contains the generated Node MCP project and the user can run it from there (e.g. npm install && npm start).
3. **Given** a generated MCP (Node project) in the output directory, **When** the user runs it (e.g. npm start), **Then** an MCP client can discover and invoke tools/resources that correspond to the OpenAPI operations.

---

### User Story 2 - Specify input and output paths (Priority: P2)

A developer wants to point the CLI at an OpenAPI spec (local path or URL) and optionally choose a target directory for the generated MCP instead of the current directory.

**Why this priority**: Supports automation and project layout; not required for the first working slice.

**Independent Test**: Can be tested by running the CLI with explicit input (file or URL) and output directory and confirming the MCP is generated only in the specified output path.

**Acceptance Scenarios**:

1. **Given** an OpenAPI file at a path and a target directory, **When** the user runs the CLI with that input path and output directory, **Then** the MCP is generated only in the target directory and the current directory is unchanged.
2. **Given** an OpenAPI spec available at a URL, **When** the user runs the CLI with that URL as input, **Then** the CLI fetches the spec and generates the MCP without requiring a local file.

---

### User Story 3 - Clear failure and validation feedback (Priority: P3)

A developer provides invalid or unsupported input (e.g. malformed OpenAPI, unsupported version) and expects a clear, actionable error message instead of a generic failure.

**Why this priority**: Improves usability and debuggability after the core generate-and-run flow works.

**Independent Test**: Can be tested by invoking the CLI with invalid input and checking that the process exits with a non-zero code and that stderr (or equivalent) contains a message that identifies the problem (e.g. invalid OpenAPI, unsupported version).

**Acceptance Scenarios**:

1. **Given** a file that is not valid OpenAPI, **When** the user runs the CLI with that file, **Then** the CLI exits with an error and reports that the input is not valid OpenAPI (or similar).
2. **Given** an OpenAPI 2.0 (Swagger) file, **When** the user runs the CLI with that file, **Then** the CLI exits with an error and reports that OpenAPI 2.0 is not supported (e.g. "use OpenAPI 3.x").
3. **Given** an input path that does not exist or is not readable, **When** the user runs the CLI, **Then** the CLI exits with an error and indicates the input file or URL could not be read.

---

### Edge Cases

- When the output directory already contains files (e.g. a previous generation), the CLI MUST refuse to write and exit with an error unless the user passes an explicit force/overwrite flag (e.g. `--force`). The error message MUST state that the output directory is not empty and that the flag may be used to overwrite.
- When the input is OpenAPI 2.0 (Swagger), the CLI MUST reject it with a clear unsupported-version message and non-zero exit.
- How does the system handle very large OpenAPI specs or deeply nested schemas? The CLI MUST complete or fail with a clear resource or complexity message rather than hanging or crashing without explanation.
- How does the system handle OpenAPI operations that are not trivially mappable to MCP tools/resources (e.g. ambiguous or custom extensions)? The CLI MUST document supported patterns and either skip unsupported operations with a warning or fail with guidance.

## Assumptions

- OpenAPI 3.x is the only supported input format for v1. OpenAPI 2.0 (Swagger) is out of scope and MUST be rejected with a clear unsupported-version message.
- The caller has a writable directory for output and, to run the generated MCP, has Node.js and npm (or compatible) available.
- "Output an MCP" means generating a Node project (source) that implements an MCP server using github.com/punkpeye/fastmcp; the user runs it in the output directory (e.g. npm install && npm start).
- The mapping from OpenAPI operations to MCP tools/resources is done by deterministic rules and code generation (templates); no LLM or external AI service is used.
- Default output directory is the current working directory when not specified.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST provide a CLI that accepts an OpenAPI specification as input (local file path or URL).
- **FR-002**: The system MUST generate a Node project (source) in the caller’s chosen directory (default: current working directory) that implements an MCP server using github.com/punkpeye/fastmcp and is runnable with Node (e.g. package.json, npm install && npm start).
- **FR-003**: The generated MCP MUST expose the OpenAPI operations (paths and methods) as MCP tools or resources so that MCP clients can invoke them.
- **FR-003a**: The system MUST map OpenAPI to MCP tools/resources using deterministic rules and code generation (templates); MUST NOT use LLM or external AI for this mapping. Generation MUST be reproducible for the same OpenAPI input.
- **FR-004**: The system MUST support specifying an explicit output directory so that the generated artifact is written only there.
- **FR-005**: The system MUST validate or parse the OpenAPI input and MUST exit with a non-zero status and an actionable error message when the input is invalid or unsupported.
- **FR-005a**: The system MUST reject OpenAPI 2.0 (Swagger) input with a non-zero exit and a clear message that the version is unsupported (e.g. "OpenAPI 2.0 is not supported; use OpenAPI 3.x").
- **FR-006**: The system MUST refuse to generate into a non-empty output directory unless the user passes a force/overwrite flag (e.g. `--force`). When refusing, the CLI MUST exit with a non-zero status and an error message indicating the directory is not empty and that the flag may be used to allow overwrite.
- **FR-007**: The CLI MUST be invokable non-interactively for the minimal case (e.g. input path + default output directory) so that it can be used in scripts and automation.
- **FR-008**: The system MUST guarantee that the generated Node project has solid structure, correct dependencies, and works: (1) contract tests MUST validate generated structure (package.json valid, required files present, required dependencies); (2) integration tests MUST run `npm install` and `npm start` in the generated directory and MUST succeed (and, when feasible, list/invoke MCP tools). Any failure in these tests MUST fail the build/CI.

### Key Entities

- **OpenAPI specification**: The input contract (file or URL); describes paths, methods, parameters, and response schemas. The CLI consumes it to drive code generation via deterministic rules and templates (no LLM).
- **Generated MCP server**: A Node project (source) in the output directory that implements the MCP protocol using fastmcp and maps OpenAPI operations to MCP tools/resources. Generated by rule-based code and templates; the user runs it with Node (e.g. npm install && npm start).
- **Caller / user**: The developer or script that invokes the CLI; provides input spec and optional output path, and expects a runnable MCP in the chosen directory.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A user with a valid OpenAPI 3.x spec can produce a runnable MCP server in their chosen directory with a single CLI invocation.
- **SC-002**: The time from CLI invocation to completion of generation is under 30 seconds for a typical OpenAPI spec (e.g. under 50 paths).
- **SC-003**: When input is invalid or unsupported, the CLI exits with a non-zero code and reports a clear, actionable reason within one sentence or a short list.
- **SC-004**: The generated MCP can be started by the user and successfully listed/invoked by at least one reference MCP client for the operations defined in the input OpenAPI spec.
