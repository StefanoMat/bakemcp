package cli

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"your-mcp/internal/domain/mapping"
	"your-mcp/internal/domain/openapi"
	"your-mcp/internal/generator/node"
)

// Config holds parsed CLI arguments.
type Config struct {
	InputPath  string
	OutputDir  string
	Force      bool
}

// Run executes the full flow: read input, parse OpenAPI, check output dir, map operations to tools, generate Node project.
// Returns exit code (0 = success) and error message for stderr.
func Run(cfg Config) (exitCode int, err error) {
	if cfg.OutputDir == "" {
		cfg.OutputDir, _ = os.Getwd()
	}

	// Read input (file only for now; URL in US2)
	data, err := os.ReadFile(cfg.InputPath)
	if err != nil {
		if os.IsNotExist(err) {
			return 2, fmt.Errorf("input file not found: %s", cfg.InputPath)
		}
		return 2, fmt.Errorf("cannot read input: %w", err)
	}

	// Parse OpenAPI 3.x
	result, err := openapi.Parse(bytes.NewReader(data))
	if err != nil {
		if errors.Is(err, openapi.ErrOpenAPI2Unsupported) {
			return 1, err
		}
		return 1, fmt.Errorf("invalid OpenAPI: %w", err)
	}

	if len(result.Operations) == 0 {
		return 4, fmt.Errorf("no mappable operations found in OpenAPI spec")
	}

	// Check output dir empty unless --force (create if missing)
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return 1, fmt.Errorf("cannot create output directory: %w", err)
	}
	entries, _ := os.ReadDir(cfg.OutputDir)
	if len(entries) > 0 && !cfg.Force {
		return 3, fmt.Errorf("output directory is not empty; use --force to overwrite")
	}

	// Map to MCP tools
	tools := mapping.OperationsToMCPTools(result.Operations, result.BaseURL)

	// Generate Node project
	if err := node.Generate(cfg.OutputDir, tools, nil); err != nil {
		return 1, fmt.Errorf("generation failed: %w", err)
	}

	return 0, nil
}
