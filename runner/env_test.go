package runner

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp .env file: %v", err)
	}
	return path
}

func TestLoadEnvFile_Basic(t *testing.T) {
	path := writeTempEnvFile(t, "PORT=3000\nDEBUG=true\n")
	env, err := LoadEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["PORT"] != "3000" {
		t.Errorf("expected PORT=3000, got %q", env["PORT"])
	}
	if env["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true, got %q", env["DEBUG"])
	}
}

func TestLoadEnvFile_IgnoresComments(t *testing.T) {
	path := writeTempEnvFile(t, "# this is a comment\nKEY=val\n")
	env, err := LoadEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := env["# this is a comment"]; ok {
		t.Error("comment line should not be parsed as key")
	}
	if env["KEY"] != "val" {
		t.Errorf("expected KEY=val, got %q", env["KEY"])
	}
}

func TestLoadEnvFile_StripQuotes(t *testing.T) {
	path := writeTempEnvFile(t, `NAME="hello world"\nTITLE='dev'\n`)
	env, err := LoadEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["NAME"] != "hello world" {
		t.Errorf("expected NAME=hello world, got %q", env["NAME"])
	}
}

func TestLoadEnvFile_MissingFile(t *testing.T) {
	env, err := LoadEnvFile("/nonexistent/.env")
	if err != nil {
		t.Fatalf("expected no error for missing file, got: %v", err)
	}
	if len(env) != 0 {
		t.Errorf("expected empty map for missing file, got %v", env)
	}
}

func TestLoadEnvFile_InlineComment(t *testing.T) {
	path := writeTempEnvFile(t, "HOST=localhost # the host\n")
	env, err := LoadEnvFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %q", env["HOST"])
	}
}

func TestEnvMapToSlice(t *testing.T) {
	env := map[string]string{"A": "1", "B": "2"}
	slice := EnvMapToSlice(env)
	if len(slice) != 2 {
		t.Errorf("expected 2 entries, got %d", len(slice))
	}
}
