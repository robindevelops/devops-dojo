package engine

import (
	"fmt"
	"os"
	"strings"
	"path/filepath"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/devops-dojo/cli/internal/colors"
	"github.com/devops-dojo/cli/internal/engine/scenarios"
	"github.com/devops-dojo/cli/internal/project"
	"github.com/devops-dojo/cli/internal/session"
)

type Injector struct {
	Stack *project.Stack
}

func NewInjector(stack *project.Stack) *Injector {
	return &Injector{Stack: stack}
}

func (i *Injector) InjectFailure(level scenarios.Difficulty) error {
	var targetFile string
	var injectionFunc func(string) error

	// Determine injection target based on stack
	if i.Stack.HasKubernetes && len(i.Stack.K8sManifests) > 0 {
		targetFile = i.Stack.K8sManifests[0]
		if level == scenarios.Easy {
			injectionFunc = injectK8sTypo
		} else {
			injectionFunc = injectOOMKilled
		}
	} else if i.Stack.HasDocker && len(i.Stack.Dockerfiles) > 0 {
		targetFile = i.Stack.Dockerfiles[0]
		injectionFunc = injectDockerFailure
	} else if i.Stack.HasTerraform && len(i.Stack.TerraformFiles) > 0 {
		targetFile = i.Stack.TerraformFiles[0]
		injectionFunc = injectTerraformFailure
	} else {
		return fmt.Errorf("no suitable targets found for failure injection in this project")
	}

	// Blast Radius Confirmation
	fmt.Println(colors.Colorize(colors.Red, "\n⚠️  BLAST RADIUS WARNING ⚠️"))
	fmt.Println(colors.Colorize(colors.Yellow, "Dojo will mutate the following file in your project:"))
	fmt.Printf(" -> %s\n\n", colors.Colorize(colors.Cyan, targetFile))

	prompt := promptui.Prompt{
		Label:     "Do you want to proceed and inject this failure",
		IsConfirm: true,
	}
	
	result, err := prompt.Run()
	if err != nil || (result != "y" && result != "Y" && result != "") { // empty string is often accepted as 'y' in Confirm but promptui usually returns err on abort
		return fmt.Errorf("injection cancelled by user")
	}

	// Create strict backup before mutating
	fmt.Println("Creating snapshot...")
	err = i.backupFiles([]string{targetFile})
	if err != nil {
		return fmt.Errorf("failed to create backup (aborting injection for safety): %w", err)
	}

	err = injectionFunc(targetFile)
	if err != nil {
		return err
	}

	// Save session state
	incidentID := string(level) + "-break" // Need a better ID mapper later
	err = session.SaveState(&session.State{
		ActiveIncidentID:     incidentID,
		StartTime:            time.Now(),
		VerificationAttempts: 0,
		HintLevel:            0,
	})
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to save session state: %v\n", err)
	}
	return nil
}

func (i *Injector) backupFiles(targets []string) error {
	for _, f := range targets {
		dst := filepath.Join(".dojo_backup", f)
		os.MkdirAll(filepath.Dir(dst), 0755)
		err := copyFile(f, dst)
		if err != nil {
			return err
		}
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

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}

func injectTerraformFailure(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	
	// Inject invalid syntax at the top
	broken := "resource \"aws_s3_bucket\" \"broken\" {\n  bucket = \"my-bucket\"\n  # DOJO INJECTED FAILURE (Missing closing brace)\n\n" + string(content)
	return os.WriteFile(filePath, []byte(broken), 0644)
}
