package node

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"your-mcp/internal/domain/model"
)

// FS writes files to the filesystem (abstraction for tests).
type FS interface {
	WriteFile(name string, data []byte, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
}

// OsFS uses the real os package.
type OsFS struct{}

func (OsFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (OsFS) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// Generate writes a Node project to outDir with package.json and entry script
// that registers one MCP tool per tool in tools.
func Generate(outDir string, tools []*model.MCPTool, fs FS) error {
	if fs == nil {
		fs = OsFS{}
	}
	if err := fs.MkdirAll(outDir, 0755); err != nil {
		return err
	}
	pkg := packageJSON(tools)
	pkgPath := filepath.Join(outDir, "package.json")
	pkgBytes, _ := json.MarshalIndent(pkg, "", "  ")
	if err := fs.WriteFile(pkgPath, pkgBytes, 0644); err != nil {
		return err
	}
	entry := entryScript(tools)
	entryPath := filepath.Join(outDir, "index.js")
	if err := fs.WriteFile(entryPath, []byte(entry), 0755); err != nil {
		return err
	}
	return nil
}

func packageJSON(tools []*model.MCPTool) map[string]interface{} {
	return map[string]interface{}{
		"name":    "generated-mcp",
		"version": "1.0.0",
		"type":    "module",
		"scripts": map[string]string{"start": "node index.js"},
		"dependencies": map[string]string{
			"fastmcp": "^3.29.0",
			"zod":     "^3.23.0",
		},
	}
}

func entryScript(tools []*model.MCPTool) string {
	b := `import { FastMCP } from "fastmcp";
import { z } from "zod";

const server = new FastMCP({ name: "generated-mcp", version: "1.0.0" });
`
	for _, t := range tools {
		url := t.BaseURL + t.Path
		b += fmt.Sprintf(`server.addTool({
  name: %q,
  description: %q,
  parameters: z.object({}),
  execute: async () => {
    const res = await fetch(%q, { method: %q });
    const body = await res.text();
    return body;
  },
});
`, t.Name, t.Description, url, t.Method)
	}
	b += `
server.start({ transportType: "stdio" });
`
	return b
}
