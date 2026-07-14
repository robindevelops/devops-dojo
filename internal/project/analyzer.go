package project

import (
	"os"
	"path/filepath"
	"strings"
)

// Stack represents the detected technologies in the user's project
type Stack struct {
	HasDocker        bool
	Dockerfiles      []string
	HasKubernetes    bool
	K8sManifests     []string
	HasGitHubActions bool
	HasTerraform     bool
	TerraformFiles   []string
	HasCompose       bool
	ComposeFiles     []string
}

// Analyze scans the given directory path to detect the project stack
func Analyze(dirPath string, focus string) (*Stack, error) {
	stack := &Stack{
		Dockerfiles:  []string{},
		K8sManifests: []string{},
		TerraformFiles: []string{},
		ComposeFiles: []string{},
	}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			// Skip typical ignore dirs
			if info.Name() == ".git" || info.Name() == "node_modules" || info.Name() == ".dojo_backup" || info.Name() == ".terraform" {
				return filepath.SkipDir
			}
			return nil
		}

		baseName := filepath.Base(path)
		
		if strings.Contains(strings.ToLower(baseName), "dockerfile") {
			stack.HasDocker = true
			stack.Dockerfiles = append(stack.Dockerfiles, path)
		} else if strings.HasSuffix(baseName, ".yaml") || strings.HasSuffix(baseName, ".yml") {
			if strings.Contains(strings.ToLower(baseName), "docker-compose") {
				stack.HasCompose = true
				stack.ComposeFiles = append(stack.ComposeFiles, path)
			} else if !strings.Contains(path, ".github/workflows") {
				stack.HasKubernetes = true
				stack.K8sManifests = append(stack.K8sManifests, path)
			} else {
				stack.HasGitHubActions = true
			}
		} else if strings.HasSuffix(baseName, ".tf") {
			stack.HasTerraform = true
			stack.TerraformFiles = append(stack.TerraformFiles, path)
		}
		return nil
	})

	// Apply focus filter if provided
	if focus != "" {
		focusLower := strings.ToLower(focus)
		if focusLower != "kubernetes" && focusLower != "k8s" {
			stack.HasKubernetes = false
			stack.K8sManifests = nil
		}
		if focusLower != "docker" {
			stack.HasDocker = false
			stack.Dockerfiles = nil
		}
		if focusLower != "github-actions" && focusLower != "ci" {
			stack.HasGitHubActions = false
		}
		if focusLower != "terraform" && focusLower != "tf" {
			stack.HasTerraform = false
			stack.TerraformFiles = nil
		}
		if focusLower != "compose" && focusLower != "docker-compose" {
			stack.HasCompose = false
			stack.ComposeFiles = nil
		}
	}

	return stack, err
}
