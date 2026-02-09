package node

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"bakemcp/internal/domain/model"
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

// ---------------------------------------------------------------------------
// Entry script generation
// ---------------------------------------------------------------------------

func entryScript(tools []*model.MCPTool) string {
	// Extract the base URL from the first tool (all share the same base).
	defaultBaseURL := ""
	if len(tools) > 0 {
		defaultBaseURL = tools[0].BaseURL
	}

	var b strings.Builder
	b.WriteString(`import { FastMCP } from "fastmcp";
import { z } from "zod";

`)
	b.WriteString(fmt.Sprintf("const BASE_URL = process.env.BASE_URL || %q;\n", defaultBaseURL))
	b.WriteString(`
const server = new FastMCP({ name: "generated-mcp", version: "1.0.0" });
`)
	for _, t := range tools {
		b.WriteString(toolBlock(t))
	}
	b.WriteString(`
server.start({ transportType: "stdio" });
`)
	return b.String()
}

func toolBlock(t *model.MCPTool) string {
	zodSchema := buildZodSchema(t)
	executeFn := buildExecuteFn(t)
	return fmt.Sprintf(`server.addTool({
  name: %q,
  description: %q,
  parameters: %s,
  execute: %s,
});
`, t.Name, t.Description, zodSchema, executeFn)
}

// ---------------------------------------------------------------------------
// Zod schema generation
// ---------------------------------------------------------------------------

type zodField struct {
	name     string
	zod      string
	required bool
}

func buildZodSchema(t *model.MCPTool) string {
	var fields []zodField

	// Add params (path, query, header)
	for _, p := range t.Params {
		fields = append(fields, zodField{
			name:     p.Name,
			zod:      schemaToZod(p.Schema, 6),
			required: p.Required,
		})
	}

	// Flatten body properties into the same z.object
	if t.Body != nil && t.Body.Schema != nil {
		props, _ := t.Body.Schema["properties"].(map[string]interface{})
		reqSet := requiredSet(t.Body.Schema)
		for name, propRaw := range props {
			prop, _ := propRaw.(map[string]interface{})
			fields = append(fields, zodField{
				name:     name,
				zod:      schemaToZod(prop, 6),
				required: reqSet[name],
			})
		}
	}

	if len(fields) == 0 {
		return "z.object({})"
	}

	// Sort for deterministic output
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].name < fields[j].name
	})

	var lines []string
	for _, f := range fields {
		zod := f.zod
		if !f.required {
			zod += ".optional()"
		}
		lines = append(lines, fmt.Sprintf("    %s: %s,", f.name, zod))
	}
	return fmt.Sprintf("z.object({\n%s\n  })", strings.Join(lines, "\n"))
}

func schemaToZod(schema map[string]interface{}, innerIndent int) string {
	if schema == nil {
		return "z.any()"
	}

	// Check for enum first (can apply to any type)
	if enumVals, ok := schema["enum"]; ok {
		if arr, ok := enumVals.([]interface{}); ok && len(arr) > 0 {
			var vals []string
			for _, v := range arr {
				vals = append(vals, fmt.Sprintf("%q", fmt.Sprint(v)))
			}
			return fmt.Sprintf("z.enum([%s])", strings.Join(vals, ", "))
		}
	}

	typ, _ := schema["type"].(string)
	switch typ {
	case "string":
		return "z.string()"
	case "number":
		return "z.number()"
	case "integer":
		return "z.number()"
	case "boolean":
		return "z.boolean()"
	case "array":
		if items, ok := schema["items"].(map[string]interface{}); ok {
			return fmt.Sprintf("z.array(%s)", schemaToZod(items, innerIndent))
		}
		return "z.array(z.any())"
	case "object":
		return objectToZod(schema, innerIndent)
	default:
		return "z.any()"
	}
}

func objectToZod(schema map[string]interface{}, innerIndent int) string {
	props, ok := schema["properties"].(map[string]interface{})
	if !ok || len(props) == 0 {
		if _, hasAP := schema["additionalProperties"]; hasAP {
			return "z.record(z.any())"
		}
		return "z.object({})"
	}

	reqSet := requiredSet(schema)

	var names []string
	for name := range props {
		names = append(names, name)
	}
	sort.Strings(names)

	pad := strings.Repeat(" ", innerIndent)
	closePad := strings.Repeat(" ", innerIndent-2)

	var lines []string
	for _, name := range names {
		propSchema, _ := props[name].(map[string]interface{})
		zod := schemaToZod(propSchema, innerIndent+2)
		if !reqSet[name] {
			zod += ".optional()"
		}
		lines = append(lines, fmt.Sprintf("%s%s: %s,", pad, name, zod))
	}

	return fmt.Sprintf("z.object({\n%s\n%s})", strings.Join(lines, "\n"), closePad)
}

func requiredSet(schema map[string]interface{}) map[string]bool {
	req := make(map[string]bool)
	if reqArr, ok := schema["required"].([]interface{}); ok {
		for _, r := range reqArr {
			if s, ok := r.(string); ok {
				req[s] = true
			}
		}
	}
	return req
}

// ---------------------------------------------------------------------------
// Execute function generation
// ---------------------------------------------------------------------------

var pathParamRe = regexp.MustCompile(`\{([^}]+)\}`)

func buildExecuteFn(t *model.MCPTool) string {
	pathParams := filterByIn(t.Params, "path")
	queryParams := filterByIn(t.Params, "query")
	hasBody := t.Body != nil
	hasPathParams := len(pathParams) > 0
	hasQueryParams := len(queryParams) > 0
	hasNonBodyParams := hasPathParams || hasQueryParams
	needsArgs := hasBody || hasPathParams || hasQueryParams

	var lines []string

	// Destructure path/query params from body when both are present
	useDestructuring := hasBody && hasNonBodyParams
	if useDestructuring {
		var names []string
		for _, p := range pathParams {
			names = append(names, p.Name)
		}
		for _, p := range queryParams {
			names = append(names, p.Name)
		}
		lines = append(lines, fmt.Sprintf("    const { %s, ...bodyArgs } = args;", strings.Join(names, ", ")))
	}

	// URL construction
	paramPrefix := "args."
	if useDestructuring {
		paramPrefix = ""
	}
	urlExpr := buildURLExpr(t.Path, pathParams, paramPrefix)

	if hasQueryParams || hasPathParams {
		lines = append(lines, fmt.Sprintf("    let url = %s;", urlExpr))
	}

	// Query params
	if hasQueryParams {
		lines = append(lines, "    const qp = new URLSearchParams();")
		for _, q := range queryParams {
			ref := paramPrefix + q.Name
			lines = append(lines, fmt.Sprintf("    if (%s !== undefined) qp.append(%q, String(%s));", ref, q.Name, ref))
		}
		lines = append(lines, "    const qs = qp.toString();")
		lines = append(lines, `    if (qs) url += "?" + qs;`)
	}

	// Fetch call
	var fetchURL string
	if hasQueryParams || hasPathParams {
		fetchURL = "url"
	} else {
		fetchURL = urlExpr
	}

	if hasBody {
		bodyRef := "args"
		if useDestructuring {
			bodyRef = "bodyArgs"
		}
		lines = append(lines, fmt.Sprintf("    const res = await fetch(%s, {", fetchURL))
		lines = append(lines, fmt.Sprintf("      method: %q,", t.Method))
		lines = append(lines, `      headers: { "Content-Type": "application/json" },`)
		lines = append(lines, fmt.Sprintf("      body: JSON.stringify(%s),", bodyRef))
		lines = append(lines, "    });")
	} else {
		lines = append(lines, fmt.Sprintf("    const res = await fetch(%s, { method: %q });", fetchURL, t.Method))
	}

	lines = append(lines, "    const body = await res.text();")
	lines = append(lines, `    if (!res.ok) throw new Error("HTTP " + res.status + ": " + body);`)
	lines = append(lines, "    return body;")

	argsStr := "()"
	if needsArgs {
		argsStr = "(args)"
	}
	return fmt.Sprintf("async %s => {\n%s\n  }", argsStr, strings.Join(lines, "\n"))
}

func buildURLExpr(path string, pathParams []model.MCPToolParam, paramPrefix string) string {
	if len(pathParams) == 0 {
		// No path params → simple string concatenation: BASE_URL + "/path"
		return fmt.Sprintf("BASE_URL + %q", path)
	}
	// Convert {param} to ${encodeURIComponent(prefix.param)} in template literal
	result := pathParamRe.ReplaceAllStringFunc(path, func(match string) string {
		name := match[1 : len(match)-1] // strip { and }
		return fmt.Sprintf("${encodeURIComponent(%s%s)}", paramPrefix, name)
	})
	return fmt.Sprintf("`${BASE_URL}%s`", result)
}

func filterByIn(params []model.MCPToolParam, in string) []model.MCPToolParam {
	var out []model.MCPToolParam
	for _, p := range params {
		if p.In == in {
			out = append(out, p)
		}
	}
	return out
}
