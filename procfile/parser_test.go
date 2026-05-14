package procfile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/procfile-run/procfile"
)

func writeTempProcfile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "Procfile")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp Procfile: %v", err)
	}
	return path
}

func TestParse_ValidProcfile(t *testing.T) {
	path := writeTempProcfile(t, `
# This is a comment
web: go run ./cmd/server
worker: go run ./cmd/worker --queue default
`)

	processes, err := procfile.Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(processes) != 2 {
		t.Fatalf("expected 2 processes, got %d", len(processes))
	}
	if processes[0].Name != "web" || processes[0].Command != "go run ./cmd/server" {
		t.Errorf("unexpected first process: %+v", processes[0])
	}
	if processes[1].Name != "worker" || processes[1].Command != "go run ./cmd/worker --queue default" {
		t.Errorf("unexpected second process: %+v", processes[1])
	}
}

func TestParse_EmptyFile(t *testing.T) {
	path := writeTempProcfile(t, "# only comments\n\n")
	_, err := procfile.Parse(path)
	if err == nil {
		t.Fatal("expected error for empty Procfile, got nil")
	}
}

func TestParse_InvalidSyntax(t *testing.T) {
	path := writeTempProcfile(t, "web go run ./cmd/server\n")
	_, err := procfile.Parse(path)
	if err == nil {
		t.Fatal("expected error for invalid syntax, got nil")
	}
}

func TestParse_MissingFile(t *testing.T) {
	_, err := procfile.Parse("/nonexistent/Procfile")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestParse_EmptyCommand(t *testing.T) {
	path := writeTempProcfile(t, "web:   \n")
	_, err := procfile.Parse(path)
	if err == nil {
		t.Fatal("expected error for empty command, got nil")
	}
}
