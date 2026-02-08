package contract_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"your-mcp/internal/domain/mapping"
	"your-mcp/internal/domain/openapi"
	"your-mcp/internal/generator/node"
)

// Contract: after generating with fixture OpenAPI, package.json exists and is valid JSON;
// required fields (name, scripts.start, dependency fastmcp); entry script exists.
func TestGeneratedProject_Contract(t *testing.T) {
	specPath := filepath.Join("..", "fixtures", "openapi3-minimal.json")
	data, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	result, err := openapi.Parse(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("parse OpenAPI: %v", err)
	}
	if len(result.Operations) == 0 {
		t.Fatal("fixture has no operations")
	}
	tools := mapping.OperationsToMCPTools(result.Operations, result.BaseURL)
	dir := t.TempDir()
	if err := node.Generate(dir, tools, nil); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	// package.json exists and valid JSON
	pkgPath := filepath.Join(dir, "package.json")
	pkgData, err := os.ReadFile(pkgPath)
	if err != nil {
		t.Fatalf("package.json missing or unreadable: %v", err)
	}
	var pkg map[string]interface{}
	if err := json.Unmarshal(pkgData, &pkg); err != nil {
		t.Fatalf("package.json invalid JSON: %v", err)
	}
	if pkg["name"] == nil || pkg["name"] == "" {
		t.Error("package.json must have name")
	}
	if pkg["version"] == nil {
		t.Error("package.json must have version")
	}
	if pkg["type"] != "module" {
		t.Errorf("package.json type must be module, got %v", pkg["type"])
	}
	scripts, _ := pkg["scripts"].(map[string]interface{})
	if scripts == nil || scripts["start"] == nil || scripts["start"] == "" {
		t.Error("package.json must have scripts.start")
	}
	deps, _ := pkg["dependencies"].(map[string]interface{})
	if deps == nil || deps["fastmcp"] == nil {
		t.Error("package.json must have dependency fastmcp")
	}

	// entry script exists
	entryPath := filepath.Join(dir, "index.js")
	if _, err := os.Stat(entryPath); err != nil {
		t.Fatalf("entry script index.js missing: %v", err)
	}
}
