// Generated MCP entry - tools are registered by the generator
import { FastMCP } from "fastmcp";

const mcp = new FastMCP({ name: "generated-mcp" });

// Tools will be registered here by the Go generator
// Example: mcp.tool("toolName", "description", { schema }, async (args) => ({ result }));

//example tool registration
// server.addTool({
//   name: "ping",
//   description: "Chama o endpoint GET /ping da Go API e retorna a resposta (esperado: 'pong')",
//   annotations: {
//     readOnlyHint: true,
//     openWorldHint: true,
//   },
//   execute: async () => {
//     const response = await fetch(`${GO_API_BASE_URL}/ping`);

//     if (!response.ok) {
//       throw new Error(`Go API respondeu com status ${response.status}: ${response.statusText}`);
//     }

//     const text = await response.text();
//     return text;
//   },
// });

export default mcp;
