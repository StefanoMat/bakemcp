# Research: OpenAPI-to-MCP CLI

**Feature**: 001-openapi-to-mcp | **Date**: 2025-02-06

## 1. Go OpenAPI 3.x parsing (read-only)

**Decision**: Use a Go library that parses OpenAPI 3.x YAML/JSON into an in-memory struct (paths, methods, parameters, response schemas). No code generation from OpenAPI; only read for our own mapping and Node codegen.

**Rationale**: We need deterministic, offline parsing. No need for full validation beyond "is this valid OpenAPI 3.x"; we need paths, operations, parameters, and (optionally) schemas to drive MCP tool generation.

**Alternatives considered**:
- **kin-openapi**: Popular, supports OpenAPI 3.x, read-only parse. Chosen for maturity and YAML/JSON support.
- **go-openapi/swag**: More codegen-oriented; heavier. Rejected for read-only use.
- **Manual JSON/YAML + structs**: Possible but error-prone; prefer a maintained parser.

**References**: [kin-openapi](https://github.com/getkin/kin-openapi) (OpenAPI 3.x).

---

## 2. Mapping OpenAPI operations to MCP tools

**Decision**: One OpenAPI operation (path + HTTP method) maps to one MCP tool. Tool name: derived from operationId if present, else from path + method (e.g. `get_users`, `post_items`). Parameters: from OpenAPI parameters (path, query, header) and requestBody; types mapped to MCP-friendly JSON schema. No LLM; rules are fixed (e.g. string, number, boolean, array, object).

**Rationale**: Spec requires deterministic, reproducible mapping. OperationId is the standard way to name operations; fallback to path+method keeps tools identifiable. MCP tools accept JSON arguments; we map OpenAPI parameter schemas to that.

**Alternatives considered**:
- **One tool per path (all methods)**: Would require tool to take "method" as arg; less ergonomic for MCP clients. Rejected.
- **Resources instead of tools for GET**: MCP has resources (URI-like) and tools (callable). For v1, mapping all operations to tools is simpler; resources can be added later if needed.

---

## 3. Node project structure and fastmcp usage (generated output)

**Decision**: Generated output is a Node project with `package.json` (fastmcp as dependency), an entry script (e.g. `index.ts` or `index.js`) that creates a FastMCP server and registers one tool per OpenAPI operation. Each tool handler receives args (from OpenAPI params/body) and can call an optional backend URL (from OpenAPI `servers`) or leave that to the user. For v1, generated tools return a placeholder or the OpenAPI description; actual HTTP call to the API can be a follow-up.

**Rationale**: Spec requires Node + fastmcp; runnable with `npm install && npm start`. FastMCP (TypeScript) exposes `.tool()` for registration; we generate one tool per operation with a deterministic name and parameter schema.

**Alternatives considered**:
- **JavaScript only**: Simpler for users without TypeScript; fastmcp supports both. Can generate `.js` or `.ts`; TypeScript preferred for generated code clarity.
- **Single file vs multiple**: Single entry file with all tools is simpler for v1; splitting per-path can come later if needed.

**References**: [fastmcp](https://github.com/punkpeye/fastmcp), [MCP tools](https://modelcontextprotocol.io/docs/concepts/tools).

---

## 4. Unsupported operations and edge cases

**Decision**: Document supported patterns in the CLI (e.g. path params, query, requestBody with application/json). Operations that cannot be mapped with current rules (e.g. custom extensions, unsupported auth) are skipped with a warning written to stderr; generation continues. If no operations remain mappable, exit with non-zero and clear message.

**Rationale**: Spec edge case: "skip unsupported with warning or fail with guidance." Skip-with-warning keeps partial generation useful; fail when zero tools generated avoids a useless artifact.

---

## 5. URL fetch for OpenAPI input

**Decision**: Support OpenAPI input via URL (HTTP/HTTPS). Use Go standard library or a small HTTP client; timeout (e.g. 30s); no auth in v1 (public URLs only). On failure (timeout, non-2xx, invalid body), exit with non-zero and message (e.g. "failed to fetch OpenAPI from URL: ...").

**Rationale**: Spec FR-001 and User Story 2 require URL as input. Deterministic and offline after fetch; no LLM.

**Alternatives considered**: Only file path for v1 — rejected; spec explicitly requires URL support.
