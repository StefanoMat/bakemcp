# bakemcp

A CLI that turns any OpenAPI 3.x spec into a ready-to-run MCP server.

Give it an OpenAPI file, get a Node.js project powered by [fastmcp](https://github.com/punkpeye/fastmcp) — each API operation becomes an MCP tool that calls the real endpoint.

## Quick Start

```bash
bakemcp openapi.yaml
cd generated-mcp
npm install
npm start
```

That's it. Your MCP server is running.

## Installation

### Homebrew (recommended)

```bash
brew install stefanoMat/tap/bakemcp
```

No Go required — installs a pre-built binary.

### Download binary

Grab the latest release for your platform from [GitHub Releases](https://github.com/stefanoMat/bakemcp/releases), extract, and move to your PATH:

```bash
# Example for macOS arm64
tar -xzf bakemcp_darwin_arm64.tar.gz
sudo mv bakemcp /usr/local/bin/
```

### From source

Requires Go 1.21+.

```bash
git clone https://github.com/stefanoMat/bakemcp.git
cd bakemcp
make install
```

This builds the binary and copies it to `/usr/local/bin/bakemcp`.

## Usage

```
Usage: bakemcp [options] <openapi-input>
  openapi-input  path to OpenAPI 3.x file (JSON or YAML)
  -o string      output directory (default: current directory)
  -f             overwrite non-empty output directory
```

### Examples

```bash
# Generate in current directory
bakemcp api.yaml

# Generate in a specific directory
bakemcp -o ./my-mcp api.yaml

# Overwrite existing files
bakemcp -f api.yaml
```

## What it generates

Given an OpenAPI spec like:

```yaml
openapi: 3.0.3
info:
  title: My API
  version: 1.0.0
servers:
  - url: http://localhost:8080
paths:
  /ping:
    get:
      operationId: getPing
      summary: Ping
```

bakemcp generates a Node.js project with:

**`package.json`** — dependencies (`fastmcp`, `zod`), start script

**`index.js`** — MCP server with one tool per operation:

```javascript
import { FastMCP } from "fastmcp";
import { z } from "zod";

const server = new FastMCP({ name: "generated-mcp", version: "1.0.0" });
server.addTool({
  name: "getping",
  description: "Ping",
  parameters: z.object({}),
  execute: async () => {
    const res = await fetch("http://localhost:8080/ping", { method: "GET" });
    const body = await res.text();
    return body;
  },
});

server.start({ transportType: "stdio" });
```

## Using with Cursor / Claude Desktop

Add to your MCP config:

```json
{
  "mcpServers": {
    "my-api": {
      "command": "node",
      "args": ["/path/to/generated/index.js"]
    }
  }
}
```

## How it works

1. Parses the OpenAPI 3.x spec (JSON or YAML) using [kin-openapi](https://github.com/getkin/kin-openapi)
2. Maps each operation to an MCP tool (name from `operationId` or `method_path`)
3. Extracts the base URL from `servers[0]`
4. Generates a Node.js project where each tool calls the real API endpoint via `fetch`

## Constraints

- **OpenAPI 3.x only** — OpenAPI 2.0 (Swagger) is rejected with a clear error
- **Output directory must be empty** unless `-f` is used
- Tools use `stdio` transport by default

## Development

```bash
make build     # build binary
make test      # run all tests (unit, contract, integration)
make fmt       # format code
make lint      # vet + format
make install   # build + install to /usr/local/bin
```
## Next steps
 - HTTP Stream
 - Authentication
## License

MIT
