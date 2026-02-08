# Generated Node Project Contract

**Feature**: 001-openapi-to-mcp | **Date**: 2025-02-06

## Purpose

The CLI generates a Node project that implements an MCP server using fastmcp. This contract describes the minimal structure and behavior so that integration tests and consumers can rely on it.

## Required files

| File | Description |
|------|-------------|
| `package.json` | Must include: `name`, `version`, `type` (module), `scripts.start`, dependency `fastmcp` (or equivalent). Must be valid for `npm install`. |
| Entry script | Default: `index.ts` or `index.js`. Must create a FastMCP server and register one MCP tool per OpenAPI operation. Must be runnable via `npm start`. |

## Runtime behavior

- **Start**: `npm install` then `npm start` (or `node index.js`) MUST start an MCP server (stdio or HTTP per fastmcp config).
- **Tools**: Each OpenAPI operation (path + method) MUST be exposed as one MCP tool. Tool name MUST be deterministic from operationId or path+method.
- **Input schema**: Each tool MUST accept arguments that reflect the OpenAPI parameters/requestBody (JSON Schema). Types: string, number, boolean, array, object as needed.

## Out of scope (v1)

- Actual HTTP call from generated MCP to the OpenAPI backend (optional; can be placeholder or user-configured later).
- MCP resources (URI-like); only tools are required for v1.
- Authentication (OpenAPI security); generated code may pass through or document; not required for contract.

## Validation

- **Contract tests** (MUST): After generation, verify: `package.json` exists and is valid JSON; required fields present (`name`, `version`, `scripts.start`, dependency `fastmcp`); entry script exists. Fail build if any check fails.
- **Integration tests** (MUST): After CLI generates into a temp dir: run `npm install` (must exit 0); run `npm start` with short timeout (process must start without immediate crash). When a test MCP client is available: list tools and invoke at least one; must succeed. Fail build if any step fails.
- **Guarantees doc**: See [generated-output-guarantees.md](./generated-output-guarantees.md) for full guarantee mechanisms (structure, libs, runnable, deterministic).
