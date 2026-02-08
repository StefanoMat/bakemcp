package mapping_test

import (
	"testing"

	"your-mcp/internal/domain/mapping"
	"your-mcp/internal/domain/model"
)

func TestOperationToMCPTool_NameFromOperationId(t *testing.T) {
	op := &model.Operation{
		OperationID: "getUser",
		Path:        "/users/{id}",
		Method:      "GET",
		Summary:     "Get user",
	}
	tool := mapping.OperationToMCPTool(op, "http://localhost:8080")
	if tool.Name != "getuser" {
		t.Errorf("name: got %q (sanitized from operationId)", tool.Name)
	}
	if tool.Description != "Get user" {
		t.Errorf("description: got %q", tool.Description)
	}
}

func TestOperationToMCPTool_NameFromPathMethod(t *testing.T) {
	op := &model.Operation{
		Path:   "/items",
		Method: "POST",
		Summary: "Create item",
	}
	tool := mapping.OperationToMCPTool(op, "")
	if tool.Name != "post_items" {
		t.Errorf("name: got %q (expected path+method derived)", tool.Name)
	}
}

func TestOperationsToMCPTools_InputSchema(t *testing.T) {
	op := &model.Operation{
		OperationID: "ping",
		Parameters: []model.Parameter{
			{Name: "q", In: "query", Required: true, Schema: map[string]interface{}{"type": "string"}},
		},
	}
	tools := mapping.OperationsToMCPTools([]*model.Operation{op}, "")
	if len(tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(tools))
	}
	props, _ := tools[0].InputSchema["properties"].(map[string]interface{})
	if props == nil {
		t.Fatal("expected properties in input schema")
	}
	if _, ok := props["q"]; !ok {
		t.Error("expected property q in schema")
	}
}
