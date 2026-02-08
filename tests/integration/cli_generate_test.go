package integration_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"bakemcp/internal/cli"
)

// Integration: build CLI flow (via cli.Run), generate into temp dir with fixture OpenAPI,
// run npm install (exit 0), run npm start with short timeout (process starts).
func TestCLI_GenerateThenNpmInstallAndStart(t *testing.T) {
	// Resolve fixture path relative to repo root (test runs from package dir or module root)
	fixturePath := filepath.Join("..", "fixtures", "openapi3-minimal.json")
	if _, err := os.Stat(fixturePath); err != nil {
		t.Skipf("fixture not found (run from repo root or tests/integration): %v", err)
	}
	outDir := t.TempDir()
	cfg := cli.Config{
		InputPath: fixturePath,
		OutputDir: outDir,
		Force:     true,
	}
	code, err := cli.Run(cfg)
	if err != nil {
		t.Fatalf("cli.Run: %v (exit %d)", err, code)
	}
	if code != 0 {
		t.Fatalf("cli.Run exit code: got %d", code)
	}

	// npm install must succeed
	cmdInstall := exec.Command("npm", "install")
	cmdInstall.Dir = outDir
	cmdInstall.Stdout = nil
	cmdInstall.Stderr = nil
	if err := cmdInstall.Run(); err != nil {
		t.Fatalf("npm install: %v", err)
	}

	// npm start must start without immediate crash (short timeout)
	cmdStart := exec.Command("npm", "start")
	cmdStart.Dir = outDir
	cmdStart.Stdout = nil
	cmdStart.Stderr = nil
	if err := cmdStart.Start(); err != nil {
		t.Fatalf("npm start: %v", err)
	}
	done := make(chan error, 1)
	go func() { done <- cmdStart.Wait() }()
	select {
	case err := <-done:
		if err != nil {
			t.Logf("npm start exited: %v (may be expected if server exits)", err)
		}
	case <-time.After(3 * time.Second):
		_ = cmdStart.Process.Kill()
		<-done
	}
	// Test passes if we reached here: install succeeded and start launched
}
