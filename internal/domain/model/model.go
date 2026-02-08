package model

// Operation represents an OpenAPI operation (path + method) for mapping.
type Operation struct {
	Path        string
	Method      string
	OperationID string
	Summary     string
	Parameters  []Parameter
	RequestBody *RequestBody
}

// Parameter represents an OpenAPI parameter (path, query, header).
type Parameter struct {
	Name     string
	In       string // path, query, header
	Required bool
	Schema   map[string]interface{}
}

// RequestBody represents OpenAPI requestBody (e.g. application/json schema).
type RequestBody struct {
	Required bool
	Schema   map[string]interface{}
}

// MCPTool represents an MCP tool derived from an OpenAPI operation.
type MCPTool struct {
	Name        string
	Description string
	InputSchema map[string]interface{} // JSON Schema for tool arguments
	Method      string                 // HTTP method (GET, POST, etc.)
	Path        string                 // API path (e.g. /ping)
	BaseURL     string                 // Base URL from OpenAPI servers (e.g. http://localhost:8080)
}
