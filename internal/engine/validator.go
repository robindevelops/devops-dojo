package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"bytes"

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
	fmt.Println("Analyzing project state and comparing with backups...")

	allFixed := true

	for _, f := range v.Stack.K8sManifests {
		backupPath := filepath.Join(".dojo_backup", filepath.Base(f))
		if !compareFiles(f, backupPath) {
			fmt.Printf("❌ Kubernetes manifest %s still has issues.\n", f)
			allFixed = false
		}
	}

	for _, f := range v.Stack.Dockerfiles {
		backupPath := filepath.Join(".dojo_backup", filepath.Base(f))
		if !compareFiles(f, backupPath) {
			fmt.Printf("❌ Dockerfile %s still has issues.\n", f)
			allFixed = false
		}
	}

	return allFixed, nil
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
