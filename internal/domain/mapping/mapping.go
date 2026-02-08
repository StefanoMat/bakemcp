package mapping

import (
	"fmt"
	"regexp"
	"strings"

	"your-mcp/internal/domain/model"
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
	return &model.MCPTool{
		Name:        name,
		Description: desc,
		InputSchema: schema,
		Method:      strings.ToUpper(op.Method),
		Path:        op.Path,
		BaseURL:     baseURL,
	}
}

// OperationsToMCPTools maps each operation to one MCP tool.
func OperationsToMCPTools(ops []*model.Operation, baseURL string) []*model.MCPTool {
	tools := make([]*model.MCPTool, 0, len(ops))
	for _, op := range ops {
		tools = append(tools, OperationToMCPTool(op, baseURL))
	}
	return tools
}

var nonID = regexp.MustCompile(`[^a-zA-Z0-9_]+`)

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

func sanitizeName(s string) string {
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
