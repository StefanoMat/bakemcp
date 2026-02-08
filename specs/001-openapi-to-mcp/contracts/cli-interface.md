# CLI Contract: openapi2mcp

**Feature**: 001-openapi-to-mcp | **Date**: 2025-02-06

## Invocation

```text
openapi2mcp [OPTIONS] <openapi-input>
```

- **openapi-input**: Required. Local file path or URL (HTTP/HTTPS) to an OpenAPI 3.x document (YAML or JSON).

## Options

| Option | Short | Description | Default |
|--------|-------|-------------|---------|
| `--output` | `-o` | Output directory for generated Node MCP project | Current working directory |
| `--force` | `-f` | Allow generation into non-empty output directory (overwrite) | false |

## Exit codes

| Code | Meaning |
|------|--------|
| 0 | Success: Node project generated in output directory |
| 1 | Invalid or unsupported OpenAPI (e.g. malformed, OpenAPI 2.0) |
| 2 | Input file or URL not found / not readable / fetch failed |
| 3 | Output directory not empty and `--force` not set |
| 4 | No mappable operations (all skipped with warnings) |

## Output

- **stdout**: Optional success message (e.g. "Generated MCP at <path>").
- **stderr**: All errors and warnings (e.g. unsupported version, skipped operations, non-empty dir). Exit code non-zero on failure.

## Examples

```bash
# Generate in current directory (must be empty or use --force)
openapi2mcp ./openapi.yaml

# Generate into specific directory
openapi2mcp -o ./my-mcp https://api.example.com/openapi.json

# Overwrite existing content in output directory
openapi2mcp --force -o ./out ./spec.yaml
```

## Non-interactive

CLI MUST NOT prompt for input when invoked with valid args (input + optional -o, -f). Scriptable and CI-friendly.
