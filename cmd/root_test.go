package cmd

import (
	"os"
	"path/filepath"
	"testing"
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

func TestRun_MissingProcfile(t *testing.T) {
	rootCmd.ResetFlags()
	init()

	rootCmd.SetArgs([]string{"--procfile", "/nonexistent/Procfile"})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for missing Procfile, got nil")
	}
}

func TestRun_EmptyProcfile(t *testing.T) {
	path := writeTempProcfile(t, "")

	rootCmd.ResetFlags()
	init()

	rootCmd.SetArgs([]string{"--procfile", path})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for empty Procfile, got nil")
	}
}

func TestRun_InvalidProcfileSyntax(t *testing.T) {
	path := writeTempProcfile(t, "this is not valid")

	rootCmd.ResetFlags()
	init()

	rootCmd.SetArgs([]string{"--procfile", path})
	err := rootCmd.Execute()
	if err == nil {
		t.Error("expected error for invalid Procfile syntax, got nil")
	}
}

func TestExecute_DefaultFlags(t *testing.T) {
	rootCmd.ResetFlags()
	init()

	if rootCmd.Flag("procfile") == nil {
		t.Error("expected --procfile flag to be registered")
	}
	if rootCmd.Flag("env") == nil {
		t.Error("expected --env flag to be registered")
	}
}
