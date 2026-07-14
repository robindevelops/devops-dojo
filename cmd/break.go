package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/devops-dojo/cli/internal/project"
	"github.com/devops-dojo/cli/internal/engine"
	"github.com/devops-dojo/cli/internal/engine/scenarios"
)

var level string

var breakCmd = &cobra.Command{
	Use:   "break",
	Short: "Analyze the current project and inject a failure",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Analyzing current project and injecting failure at level: %s\n", level)
		
		stack, err := project.Analyze(".")
		if err != nil {
			fmt.Printf("Error analyzing project: %v\n", err)
			return
		}

		fmt.Printf("Detected Stack:\n")
		fmt.Printf("- Docker: %v (Found %d Dockerfiles)\n", stack.HasDocker, len(stack.Dockerfiles))
		fmt.Printf("- Kubernetes: %v (Found %d Manifests)\n", stack.HasKubernetes, len(stack.K8sManifests))
		fmt.Printf("- GitHub Actions: %v\n", stack.HasGitHubActions)

		// Call Injector
		inj := engine.NewInjector(stack)
		err = inj.InjectFailure(scenarios.Difficulty(level))
		if err != nil {
			fmt.Printf("Failed to inject failure: %v\n", err)
			return
		}

		fmt.Println("✅ Failure injected successfully! Your training begins now. Use your standard tools to investigate.")
	},
}

func init() {
	rootCmd.AddCommand(breakCmd)
	breakCmd.Flags().StringVarP(&level, "level", "l", "medium", "Difficulty level (easy, medium, hard, extreme)")
}
