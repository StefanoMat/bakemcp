package node_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"bakemcp/internal/domain/model"
	"bakemcp/internal/generator/node"
)

func TestGenerate_ProducesValidPackageJSON(t *testing.T) {
	dir := t.TempDir()
	tools := []*model.MCPTool{
		{Name: "ping", Description: "Ping", InputSchema: map[string]interface{}{"type": "object"}},
	}
	err := node.Generate(dir, tools, nil)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	pkgPath := filepath.Join(dir, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		t.Fatalf("ReadFile package.json: %v", err)
	}
	var pkg map[string]interface{}
	if err := json.Unmarshal(data, &pkg); err != nil {
		t.Fatalf("package.json invalid JSON: %v", err)
	}
	if pkg["name"] != "generated-mcp" {
		t.Errorf("name: got %v", pkg["name"])
	}
	if pkg["type"] != "module" {
		t.Errorf("type: got %v", pkg["type"])
	}
	scripts, _ := pkg["scripts"].(map[string]interface{})
	if scripts == nil || scripts["start"] != "node index.js" {
		t.Errorf("scripts.start: got %v", pkg["scripts"])
	}
	deps, _ := pkg["dependencies"].(map[string]interface{})
	if deps == nil || deps["fastmcp"] == nil {
		t.Errorf("dependencies.fastmcp: got %v", pkg["dependencies"])
	}
}

func TestGenerate_ProducesEntryScript(t *testing.T) {
	dir := t.TempDir()
	tools := []*model.MCPTool{
		{Name: "ping", Description: "Ping", InputSchema: map[string]interface{}{"type": "object"}},
	}
	err := node.Generate(dir, tools, nil)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	entryPath := filepath.Join(dir, "index.js")
	data, err := os.ReadFile(entryPath)
	if err != nil {
		t.Fatalf("ReadFile index.js: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "FastMCP") {
		t.Error("entry script should reference FastMCP")
	}
	if !strings.Contains(content, "ping") {
		t.Error("entry script should register tool ping")
	}
}
