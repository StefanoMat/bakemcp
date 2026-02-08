# Quickstart: OpenAPI-to-MCP CLI

**Feature**: 001-openapi-to-mcp | **Date**: 2025-02-06

## Prerequisites

- **Go 1.21+** (to build and run the CLI)
- **Node.js and npm** (to run the generated MCP project)
- An **OpenAPI 3.x** spec (YAML or JSON), local file or URL

## Build the CLI (from repo root)

```bash
go build -o openapi2mcp ./cmd/openapi2mcp
```

Or install locally:

```bash
go install ./cmd/openapi2mcp
# Ensure $GOPATH/bin or $GOBIN is in PATH
```

## Generate an MCP from OpenAPI

1. **Empty output directory** (recommended for first run):

   ```bash
   mkdir my-mcp && cd my-mcp
   openapi2mcp /path/to/openapi.yaml
   # Or from URL:
   openapi2mcp https://api.example.com/openapi.json
   ```

2. **Or specify output directory**:

   ```bash
   openapi2mcp -o ./my-mcp /path/to/openapi.yaml
   ```

3. **If output directory already has files**, use `--force` to overwrite:

   ```bash
   openapi2mcp --force -o ./my-mcp /path/to/openapi.yaml
   ```

## Run the generated MCP

```bash
cd my-mcp   # or the output directory you used
npm install
npm start
```

The MCP server starts (stdio or HTTP depending on generated config). Use an MCP client (e.g. IDE integration, CLI client) to list tools and invoke them; each tool corresponds to one OpenAPI operation.

## Exit codes and errors

- **0**: Success.
- **1**: Invalid or unsupported OpenAPI (e.g. OpenAPI 2.0, malformed). Check stderr for message.
- **2**: Input file/URL not found or fetch failed.
- **3**: Output directory not empty; use `--force` to overwrite or choose another directory.
- **4**: No operations could be mapped (all skipped). Check stderr for warnings.

All errors are written to stderr with a short, actionable message.

## Next steps

- Run **unit tests**: `go test ./internal/...`
- Run **integration tests**: see `tests/integration/` (CLI + generated project runnable).
- Customize generated Node project (e.g. add backend base URL) by editing the generated files after first run.
