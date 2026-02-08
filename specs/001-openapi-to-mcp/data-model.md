# Data Model: OpenAPI-to-MCP CLI

**Feature**: 001-openapi-to-mcp | **Date**: 2025-02-06

## Purpose

In-memory and generated-artifact models used by the CLI. No persistent storage; CLI reads OpenAPI, produces a Node project on disk.

---

## 1. OpenAPI (input) — parsed

| Concept | Description | Source |
|--------|-------------|--------|
| **OpenAPIDoc** | Root: info, servers, paths, components | Parsed from YAML/JSON (file or URL) |
| **PathItem** | Path string → operations (get, post, etc.) | OpenAPI `paths` |
| **Operation** | Method, operationId, parameters, requestBody, responses | OpenAPI path item |
| **Parameter** | name, in (path/query/header), schema (type, format) | Operation.parameters |
| **RequestBody** | content (e.g. application/json), schema | Operation.requestBody |
| **Schema** | type (string, number, boolean, array, object), format, properties | components.schemas or inline |

**Validation rules**: OpenAPI 3.x only; reject 2.0. Valid parse required before mapping.

---

## 2. Mapping (internal) — OpenAPI → MCP

| Concept | Description | Rules |
|--------|-------------|--------|
| **MCPTool** | Name, description, inputSchema (JSON Schema) | One per Operation |
| **Tool name** | operationId if present, else sanitize path + method (e.g. `get_users`, `post_items`) | Deterministic, unique per operation |
| **Input schema** | JSON Schema: properties from parameters + requestBody; required array | Types: string, number, boolean, array, object; no $ref resolution for v1 (inline or flatten) |

**Relationships**: One Operation → one MCPTool. Many MCPTools → one GeneratedProject.

---

## 3. Generated Node project (output)

| Concept | Description | Validation |
|--------|-------------|------------|
| **GeneratedProject** | Root of output dir: package.json, entry script, optional files | Must be runnable with `npm install && npm start` |
| **package.json** | name, version, type (module), scripts.start, dependencies (fastmcp, etc.) | Valid npm package |
| **Entry script** | Creates FastMCP server; registers one tool per MCPTool; each handler receives args, returns placeholder or call result | Executable by Node |
| **Tool registration** | .tool(name, description, inputSchema, handler) | Matches MCP protocol |

**State**: No runtime state in CLI; generated project is static files. Generated project may call external API (OpenAPI servers) at runtime — out of scope for CLI data model.

---

## 4. CLI invocation (boundary)

| Concept | Description | Validation |
|--------|-------------|------------|
| **Input** | openapi_input (file path or URL), output_dir (default: cwd), force (bool) | Path/URL exists and readable; output_dir writable; if not empty, force required |
| **Exit codes** | 0 success; non-zero with clear stderr message (invalid OpenAPI, unsupported version, fetch failure, output dir not empty, etc.) | Constitution: predictable failures |

---

## 5. Entity relationship (summary)

```text
OpenAPIDoc
  └── paths → PathItem[]
        └── Operation[] (get, post, put, delete, etc.)
              └── parameters, requestBody, responses
                    ↓ (mapping rules)
              MCPTool (name, description, inputSchema)
                    ↓ (codegen)
              GeneratedProject (package.json + entry script with registered tools)
```

No identity/lifecycle beyond a single CLI run; no DB.
