package engine

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/devops-dojo/cli/internal/project"
)

type Validator struct {
	Stack *project.Stack
}

func NewValidator(stack *project.Stack) *Validator {
	return &Validator{Stack: stack}
}

// Verify checks if the current state matches the backup state (issue resolved)
func (v *Validator) Verify() (bool, error) {
	fmt.Println("Running objective actionable validations...")

	// Objective Validation: Docker
	for _, f := range v.Stack.Dockerfiles {
		fmt.Printf("Building %s to verify fix...\n", f)
		cmd := exec.Command("docker", "build", "-t", "dojo-test-build", "-f", f, filepath.Dir(f))
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("❌ Docker build failed for %s:\n%s\n", f, string(output))
			return false, nil // Still broken
		}
	}

	// Objective Validation: Kubernetes
	for _, f := range v.Stack.K8sManifests {
		fmt.Printf("Validating K8s manifest %s...\n", f)
		cmd := exec.Command("kubectl", "apply", "--dry-run=client", "-f", f)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("❌ Kubernetes validation failed for %s:\n%s\n", f, string(output))
			return false, nil
		}
		
		// For scenarios where dry-run succeeds but business logic fails (like OOM limits)
		content, err := os.ReadFile(f)
		if err == nil {
			if strings.Contains(string(content), "memory: 1Mi") {
				fmt.Printf("❌ Pod memory limit is still 1Mi. This will cause OOMKilled under load.\n")
				return false, nil 
			}
		}
	}

	// Objective Validation: Terraform
	if len(v.Stack.TerraformFiles) > 0 {
		fmt.Printf("Validating Terraform syntax...\n")
		// Assume running in root of tf project
		tfDir := filepath.Dir(v.Stack.TerraformFiles[0])
		
		cmd := exec.Command("terraform", "validate")
		cmd.Dir = tfDir
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("❌ Terraform validation failed:\n%s\n", string(output))
			return false, nil
		}
	}

	return true, nil
}

func compareFiles(file1, file2 string) bool {
	content1, err1 := os.ReadFile(file1)
	content2, err2 := os.ReadFile(file2)

	if err1 != nil || err2 != nil {
		// If backup doesn't exist, assume no issues were injected there
		if os.IsNotExist(err2) {
			return true
		}
		return false
	}

	return bytes.Equal(content1, content2)
}
