package engine

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devops-dojo/cli/internal/project"
)

func TestInjector_BackupFiles(t *testing.T) {
	// Create a temp dir
	tmpDir, err := os.MkdirTemp("", "dojo-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create dummy test file
	testFile := "Dockerfile"
	err = os.WriteFile(filepath.Join(tmpDir, testFile), []byte("FROM alpine\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Change working dir so relative paths work like production
	origWd, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origWd)

	inj := NewInjector(&project.Stack{})
	err = inj.backupFiles([]string{testFile})
	if err != nil {
		t.Fatalf("Failed to backup: %v", err)
	}

	// Verify backup exists
	if _, err := os.Stat(filepath.Join(".dojo_backup", testFile)); os.IsNotExist(err) {
		t.Errorf("Backup file was not created in .dojo_backup")
	}
}

func TestInjectDockerFailure(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "dojo-test-docker")
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "Dockerfile")
	os.WriteFile(testFile, []byte("FROM ubuntu:latest\nRUN apt-get update\n"), 0644)

	err := injectDockerFailure(testFile)
	if err != nil {
		t.Fatalf("injectDockerFailure returned error: %v", err)
	}

	content, _ := os.ReadFile(testFile)
	if !strings.Contains(string(content), "RUN false") {
		t.Errorf("Failure was not injected. Content: %s", string(content))
	}
}

func TestInjectK8sTypo(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "dojo-test-k8s")
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "deploy.yaml")
	os.WriteFile(testFile, []byte("apiVersion: apps/v1\nkind: Deployment\n"), 0644)

	err := injectK8sTypo(testFile)
	if err != nil {
		t.Fatalf("injectK8sTypo returned error: %v", err)
	}

	content, _ := os.ReadFile(testFile)
	if !strings.Contains(string(content), "apiVersoin") {
		t.Errorf("Typo failure was not injected. Content: %s", string(content))
	}
}
