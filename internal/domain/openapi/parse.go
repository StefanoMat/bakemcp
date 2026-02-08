package openapi

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"

	"your-mcp/internal/domain/model"
)

// ErrOpenAPI2Unsupported is returned when the spec is OpenAPI 2.0 (Swagger).
var ErrOpenAPI2Unsupported = &UnsupportedVersionError{Msg: "OpenAPI 2.0 is not supported; use OpenAPI 3.x"}

// UnsupportedVersionError represents an unsupported OpenAPI version.
type UnsupportedVersionError struct {
	Msg string
}

func (e *UnsupportedVersionError) Error() string {
	return e.Msg
}

// ParseResult holds the result of parsing an OpenAPI spec.
type ParseResult struct {
	Operations []*model.Operation
	BaseURL    string // First server URL, if present
}

// Parse reads an OpenAPI 3.x document from r (YAML or JSON) and returns
// a list of operations and the base URL. Rejects OpenAPI 2.0 with ErrOpenAPI2Unsupported.
func Parse(r io.Reader) (*ParseResult, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(data)
	if err != nil {
		return nil, err
	}
	if doc.OpenAPI == "" {
		return nil, ErrOpenAPI2Unsupported
	}
	// Reject 2.x
	if len(doc.OpenAPI) >= 1 && doc.OpenAPI[0] == '2' {
		return nil, ErrOpenAPI2Unsupported
	}
	baseURL := ""
	if doc.Servers != nil && len(doc.Servers) > 0 {
		baseURL = strings.TrimRight(doc.Servers[0].URL, "/")
	}
	return &ParseResult{
		Operations: extractOperations(doc),
		BaseURL:    baseURL,
	}, nil
}

func extractOperations(doc *openapi3.T) []*model.Operation {
	var out []*model.Operation
	for path, pathItem := range doc.Paths.Map() {
		if pathItem == nil {
			continue
		}
		for method, op := range pathItem.Operations() {
			if op == nil {
				continue
			}
			m := &model.Operation{
				Path:        path,
				Method:      method,
				OperationID: op.OperationID,
				Summary:     op.Summary,
			}
			for _, p := range op.Parameters {
				if p == nil || p.Value == nil {
					continue
				}
				var schema map[string]interface{}
				if p.Value.Schema != nil && p.Value.Schema.Value != nil {
					schema = schemaToMap(p.Value.Schema.Value)
				}
				m.Parameters = append(m.Parameters, model.Parameter{
					Name:     p.Value.Name,
					In:       p.Value.In,
					Required: p.Value.Required,
					Schema:   schema,
				})
			}
			if op.RequestBody != nil && op.RequestBody.Value != nil {
				ct := op.RequestBody.Value.Content.Get("application/json")
				if ct != nil && ct.Schema != nil && ct.Schema.Value != nil {
					m.RequestBody = &model.RequestBody{
						Required: op.RequestBody.Value.Required,
						Schema:   schemaToMap(ct.Schema.Value),
					}
				}
			}
			out = append(out, m)
		}
	}
	return out
}

func schemaToMap(s *openapi3.Schema) map[string]interface{} {
	if s == nil {
		return nil
	}
	b, _ := json.Marshal(s)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	return m
}
