package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_PathTraversal(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	err := os.WriteFile(configPath, []byte("networks: [\"127.0.0.1/32\"]"), 0644)
	if err != nil {
		t.Fatalf("failed to create temp config: %v", err)
	}

	// Create another file outside the "expected" directory (though LoadConfig doesn't have an expected dir yet)
	secretPath := filepath.Join(tmpDir, "secret.txt")
	err = os.WriteFile(secretPath, []byte("secret-data"), 0644)
	if err != nil {
		t.Fatalf("failed to create secret file: %v", err)
	}

	// Test with path traversal
	traversalPath := filepath.Join(tmpDir, "subdir", "..", "secret.txt")

	// Currently LoadConfig will just read it.
	// After fix, it will still read it if we only use filepath.Clean.
	// But filepath.Clean will resolve the "subdir/.." part.

	_, err = LoadConfig(traversalPath)
	// It should fail to parse as YAML, but it should succeed in reading if traversal works.
	if err != nil && err.Error() == "failed to read config file: open "+traversalPath+": no such file or directory" {
		t.Errorf("Expected to be able to read the file (even if it fails parsing): %v", err)
	}
}
