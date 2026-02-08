package openapi_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"bakemcp/internal/domain/openapi"
)

func TestParse_ValidOpenAPI3(t *testing.T) {
	spec := `{"openapi":"3.0.3","info":{"title":"x","version":"1.0"},"paths":{"/ping":{"get":{"operationId":"ping","responses":{"200":{}}}}}}`
	result, err := openapi.Parse(bytes.NewReader([]byte(spec)))
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if len(result.Operations) != 1 {
		t.Fatalf("expected 1 operation, got %d", len(result.Operations))
	}
	if result.Operations[0].OperationID != "ping" {
		t.Errorf("operationId: got %q", result.Operations[0].OperationID)
	}
}

func TestParse_RejectOpenAPI2(t *testing.T) {
	spec := `{"swagger":"2.0","info":{"title":"x","version":"1.0"},"paths":{}}`
	_, err := openapi.Parse(bytes.NewReader([]byte(spec)))
	if err == nil {
		t.Fatal("expected error for OpenAPI 2.0")
	}
	if !errors.Is(err, openapi.ErrOpenAPI2Unsupported) {
		t.Errorf("expected ErrOpenAPI2Unsupported, got %v", err)
	}
}

func TestParse_RejectInvalidJSON(t *testing.T) {
	_, err := openapi.Parse(strings.NewReader("not json"))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
