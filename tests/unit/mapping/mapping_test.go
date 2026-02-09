package mapping_test

import (
	"testing"

	"bakemcp/internal/domain/mapping"
	"bakemcp/internal/domain/model"
)

func TestOperationToMCPTool_NameFromOperationId(t *testing.T) {
	op := &model.Operation{
		OperationID: "getUser",
		Path:        "/users/{id}",
		Method:      "GET",
		Summary:     "Get user",
	}
	tool := mapping.OperationToMCPTool(op, "http://localhost:8080")
	if tool.Name != "get_user" {
		t.Errorf("name: got %q (sanitized from operationId)", tool.Name)
	}
	if tool.Description != "Get user" {
		t.Errorf("description: got %q", tool.Description)
	}
}

func TestOperationToMCPTool_NameFromPathMethod(t *testing.T) {
	op := &model.Operation{
		Path:    "/items",
		Method:  "POST",
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

// ── Numeric suffix detection ──────────────────────────────────────────────

func TestOperationsToMCPTools_NumericSuffix_FallsBackToPath(t *testing.T) {
	ops := []*model.Operation{
		{OperationID: "create_1", Path: "/deeplink/domain", Method: "POST", Summary: "Create domain param"},
		{OperationID: "getAll_1", Path: "/deeplink/domain", Method: "GET", Summary: "List domain params"},
	}
	tools := mapping.OperationsToMCPTools(ops, "http://localhost:8080")

	expected := map[int]string{
		0: "post_deeplink_domain",
		1: "get_deeplink_domain",
	}
	for i, want := range expected {
		if tools[i].Name != want {
			t.Errorf("tools[%d].Name = %q, want %q", i, tools[i].Name, want)
		}
	}
}

func TestOperationsToMCPTools_NumericSuffix_UpdateById(t *testing.T) {
	ops := []*model.Operation{
		{OperationID: "updateById", Path: "/deeplinks/{id}", Method: "PUT", Summary: "Update deeplink"},
		{OperationID: "updateById_1", Path: "/deeplink/domain/update/{id}", Method: "PUT", Summary: "Update domain param"},
	}
	tools := mapping.OperationsToMCPTools(ops, "http://localhost:8080")

	// updateById has no numeric suffix → keeps operationId-based name
	if tools[0].Name != "update_by_id" {
		t.Errorf("tools[0].Name = %q, want %q", tools[0].Name, "update_by_id")
	}
	// updateById_1 has numeric suffix → falls back to path-based name
	if tools[1].Name != "put_deeplink_domain_update_id" {
		t.Errorf("tools[1].Name = %q, want %q", tools[1].Name, "put_deeplink_domain_update_id")
	}
}

// ── Collision detection ───────────────────────────────────────────────────

func TestOperationsToMCPTools_Collision_FallsBackToPath(t *testing.T) {
	ops := []*model.Operation{
		{OperationID: "delete", Path: "/products/{id}", Method: "DELETE", Summary: "Delete product"},
		{OperationID: "delete", Path: "/orders/{id}", Method: "DELETE", Summary: "Delete order"},
	}
	tools := mapping.OperationsToMCPTools(ops, "http://localhost:8080")

	if tools[0].Name != "delete_products_id" {
		t.Errorf("tools[0].Name = %q, want %q", tools[0].Name, "delete_products_id")
	}
	if tools[1].Name != "delete_orders_id" {
		t.Errorf("tools[1].Name = %q, want %q", tools[1].Name, "delete_orders_id")
	}
}

func TestOperationsToMCPTools_Collision_SamePath_FinalDedup(t *testing.T) {
	// Edge case: same path AND same method (shouldn't happen in valid OpenAPI, but handle gracefully)
	ops := []*model.Operation{
		{OperationID: "doThing", Path: "/things", Method: "POST", Summary: "Do thing A"},
		{OperationID: "doThing", Path: "/things", Method: "POST", Summary: "Do thing B"},
	}
	tools := mapping.OperationsToMCPTools(ops, "")

	// Both fall back to path-based → post_things, then dedup → post_things, post_things_2
	if tools[0].Name != "post_things" {
		t.Errorf("tools[0].Name = %q, want %q", tools[0].Name, "post_things")
	}
	if tools[1].Name != "post_things_2" {
		t.Errorf("tools[1].Name = %q, want %q", tools[1].Name, "post_things_2")
	}
}

// ── No false positives ───────────────────────────────────────────────────

func TestOperationsToMCPTools_GoodNames_Unchanged(t *testing.T) {
	ops := []*model.Operation{
		{OperationID: "listProducts", Path: "/products", Method: "GET"},
		{OperationID: "createProduct", Path: "/products", Method: "POST"},
		{OperationID: "getProductById", Path: "/products/{id}", Method: "GET"},
	}
	tools := mapping.OperationsToMCPTools(ops, "")

	expected := []string{"list_products", "create_product", "get_product_by_id"}
	for i, want := range expected {
		if tools[i].Name != want {
			t.Errorf("tools[%d].Name = %q, want %q", i, tools[i].Name, want)
		}
	}
}

// ── Full scenario matching the user's real deeplink API ──────────────────

func TestOperationsToMCPTools_DeeplinkAPI_FullScenario(t *testing.T) {
	ops := []*model.Operation{
		{OperationID: "updateById", Path: "/deeplinks/{id}", Method: "PUT", Summary: "Update a specific deeplink"},
		{OperationID: "create_1", Path: "/deeplink/domain", Method: "POST", Summary: "Post Deeplink Params domain"},
		{OperationID: "getAll_1", Path: "/deeplink/domain", Method: "GET", Summary: "Get All Deeplinks Params domain"},
		{OperationID: "deleteById", Path: "/deeplink/domain/delete/{id}", Method: "DELETE", Summary: "Delete domain param"},
		{OperationID: "updateById_1", Path: "/deeplink/domain/update/{id}", Method: "PUT", Summary: "Update domain param"},
	}
	tools := mapping.OperationsToMCPTools(ops, "http://localhost:8080")

	expected := []struct {
		name string
		desc string
	}{
		{"update_by_id", "Update a specific deeplink"},
		{"post_deeplink_domain", "Post Deeplink Params domain"},
		{"get_deeplink_domain", "Get All Deeplinks Params domain"},
		{"delete_by_id", "Delete domain param"},
		{"put_deeplink_domain_update_id", "Update domain param"},
	}
	for i, want := range expected {
		if tools[i].Name != want.name {
			t.Errorf("tools[%d].Name = %q, want %q", i, tools[i].Name, want.name)
		}
	}
}
