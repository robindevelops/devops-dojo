package project

import (
	"os"
	"path/filepath"
	"strings"
)

// Stack represents the detected technologies in the user's project
type Stack struct {
	HasDocker        bool
	HasKubernetes    bool
	HasGitHubActions bool
	Dockerfiles      []string
	K8sManifests     []string
}

// Analyze scans the given directory path to detect the project stack
func Analyze(dirPath string) (*Stack, error) {
	stack := &Stack{
		Dockerfiles:  []string{},
		K8sManifests: []string{},
	}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .git directory to speed up scanning
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			// Detect Docker
			if info.Name() == "Dockerfile" || strings.HasSuffix(info.Name(), ".Dockerfile") {
				stack.HasDocker = true
				stack.Dockerfiles = append(stack.Dockerfiles, path)
			}

			// Detect Kubernetes Manifests (YAML/YML)
			if strings.HasSuffix(info.Name(), ".yaml") || strings.HasSuffix(info.Name(), ".yml") {
				// Simple heuristic: check if it's in a k8s dir or might be a k8s manifest
				// In a real implementation, we'd parse the YAML to check for apiVersion/kind
				content, err := os.ReadFile(path)
				if err == nil {
					if strings.Contains(string(content), "apiVersion:") && strings.Contains(string(content), "kind:") {
						stack.HasKubernetes = true
						stack.K8sManifests = append(stack.K8sManifests, path)
					}
				}
			}

			// Detect GitHub Actions
			if strings.Contains(path, ".github/workflows") {
				stack.HasGitHubActions = true
			}
		}

		return nil
	})

	return stack, err
}
