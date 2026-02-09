package mapping

import (
	"fmt"
	"regexp"
	"strings"

	"bakemcp/internal/domain/model"
)

// OperationToMCPTool converts one OpenAPI operation to one MCP tool.
// Tool name: operationId if present, else sanitized path+method (e.g. get_users).
// InputSchema: properties from parameters + requestBody; required array.
func OperationToMCPTool(op *model.Operation, baseURL string) *model.MCPTool {
	name := toolName(op)
	desc := op.Summary
	if desc == "" {
		desc = fmt.Sprintf("%s %s", op.Method, op.Path)
	}
	schema := buildInputSchema(op)

	var params []model.MCPToolParam
	for _, p := range op.Parameters {
		params = append(params, model.MCPToolParam{
			Name:     p.Name,
			In:       p.In,
			Required: p.Required,
			Schema:   p.Schema,
		})
	}

	var body *model.MCPToolBody
	if op.RequestBody != nil {
		body = &model.MCPToolBody{
			Required: op.RequestBody.Required,
			Schema:   op.RequestBody.Schema,
		}
	}

	return &model.MCPTool{
		Name:        name,
		Description: desc,
		InputSchema: schema,
		Params:      params,
		Body:        body,
		Method:      strings.ToUpper(op.Method),
		Path:        op.Path,
		BaseURL:     baseURL,
	}
}

// OperationsToMCPTools maps each operation to one MCP tool, ensuring unique
// and descriptive tool names. It detects auto-generated numeric suffixes
// (e.g. create_1, updateById_1) and name collisions, falling back to
// path-based naming for disambiguation.
func OperationsToMCPTools(ops []*model.Operation, baseURL string) []*model.MCPTool {
	tools := make([]*model.MCPTool, 0, len(ops))

	// Pass 1: generate all tools with default naming.
	for _, op := range ops {
		tools = append(tools, OperationToMCPTool(op, baseURL))
	}

	// Pass 2: detect bad names (numeric suffix or collision) and fix them.
	nameCount := make(map[string]int)
	for _, t := range tools {
		nameCount[t.Name]++
	}
	for i, t := range tools {
		if nameCount[t.Name] > 1 || numericSuffix.MatchString(t.Name) {
			tools[i].Name = pathBasedName(ops[i])
		}
	}

	// Pass 3: final dedup — if collisions remain, append _2, _3, etc.
	seen := make(map[string]int)
	for i, t := range tools {
		seen[t.Name]++
		if seen[t.Name] > 1 {
			tools[i].Name = fmt.Sprintf("%s_%d", t.Name, seen[t.Name])
		}
	}

	return tools
}

var nonID = regexp.MustCompile(`[^a-zA-Z0-9_]+`)

// numericSuffix detects auto-generated suffixes like _1, _2 appended by
// frameworks (e.g. SpringDoc, Swagger Codegen) when operationIds collide.
var numericSuffix = regexp.MustCompile(`_\d+$`)

// pathBasedName generates a tool name from the HTTP method and path,
// ignoring the operationId. Used as fallback for bad operationIds.
func pathBasedName(op *model.Operation) string {
	pathPart := strings.Trim(pathToName(op.Path), "_")
	methodPart := strings.ToLower(op.Method)
	if pathPart == "" {
		return methodPart
	}
	return methodPart + "_" + pathPart
}

func toolName(op *model.Operation) string {
	if op.OperationID != "" {
		return sanitizeName(op.OperationID)
	}
	// path + method: e.g. /users -> get_users
	pathPart := strings.Trim(pathToName(op.Path), "_")
	methodPart := strings.ToLower(op.Method)
	if pathPart == "" {
		return methodPart
	}
	return methodPart + "_" + pathPart
}

func pathToName(path string) string {
	path = strings.Trim(path, "/")
	if path == "" {
		return ""
	}
	path = strings.ReplaceAll(path, "/", "_")
	path = strings.ReplaceAll(path, "{", "")
	path = strings.ReplaceAll(path, "}", "")
	return sanitizeName(path)
}

// camelBoundary matches transitions like aB (lowercase→uppercase) or ABc (acronym end).
var camelBoundary = regexp.MustCompile(`([a-z0-9])([A-Z])`)

func sanitizeName(s string) string {
	// Insert underscore at camelCase boundaries: listProducts -> list_Products
	s = camelBoundary.ReplaceAllString(s, "${1}_${2}")
	s = nonID.ReplaceAllString(s, "_")
	s = strings.Trim(s, "_")
	s = strings.ToLower(s)
	if s == "" {
		s = "op"
	}
	return s
}

func buildInputSchema(op *model.Operation) map[string]interface{} {
	props := make(map[string]interface{})
	var required []string
	for _, p := range op.Parameters {
		props[p.Name] = p.Schema
		if p.Required {
			required = append(required, p.Name)
		}
	}
	if op.RequestBody != nil {
		props["body"] = op.RequestBody.Schema
		if op.RequestBody.Required {
			required = append(required, "body")
		}
	}
	schema := map[string]interface{}{
		"type":       "object",
		"properties": props,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}
