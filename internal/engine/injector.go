package engine

import (
	"fmt"
	"os"
	"strings"
	"path/filepath"

	"github.com/devops-dojo/cli/internal/engine/scenarios"
	"github.com/devops-dojo/cli/internal/project"
)

type Injector struct {
	Stack *project.Stack
}

func NewInjector(stack *project.Stack) *Injector {
	return &Injector{Stack: stack}
}

func (i *Injector) InjectFailure(level scenarios.Difficulty) error {
	fmt.Printf("Preparing to inject a %s level failure...\n", level)
	
	// Create backup
	err := i.backupFiles()
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// For demonstration, if we have K8s manifests, inject a typo for Easy level
	if i.Stack.HasKubernetes && len(i.Stack.K8sManifests) > 0 {
		manifest := i.Stack.K8sManifests[0]
		fmt.Printf("Targeting Kubernetes manifest: %s\n", manifest)
		
		if level == scenarios.Easy {
			return injectK8sTypo(manifest)
		} else if level == scenarios.Medium {
			return injectOOMKilled(manifest)
		}
	}

	// Fallback to Docker
	if i.Stack.HasDocker && len(i.Stack.Dockerfiles) > 0 {
		dockerfile := i.Stack.Dockerfiles[0]
		fmt.Printf("Targeting Dockerfile: %s\n", dockerfile)
		return injectDockerFailure(dockerfile)
	}

	return fmt.Errorf("no suitable targets found for failure injection in this project")
}

func (i *Injector) backupFiles() error {
	// Simple backup: Create a .dojo_backup dir
	os.Mkdir(".dojo_backup", 0755)
	
	for _, f := range i.Stack.K8sManifests {
		copyFile(f, filepath.Join(".dojo_backup", filepath.Base(f)))
	}
	for _, f := range i.Stack.Dockerfiles {
		copyFile(f, filepath.Join(".dojo_backup", filepath.Base(f)))
	}
	return nil
}

func injectK8sTypo(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	broken := strings.Replace(string(content), "apiVersion", "apiVersoin", 1)
	return os.WriteFile(filePath, []byte(broken), 0644)
}

func injectOOMKilled(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	// Very naive injection for demo purposes
	broken := string(content) + "\n        resources:\n          limits:\n            memory: 1Mi"
	return os.WriteFile(filePath, []byte(broken), 0644)
}

func injectDockerFailure(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	broken := string(content) + "\n# DOJO INJECTED FAILURE\nRUN false\n"
	return os.WriteFile(filePath, []byte(broken), 0644)
}

func copyFile(src, dst string) {
	input, _ := os.ReadFile(src)
	os.WriteFile(dst, input, 0644)
}
