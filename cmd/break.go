package cmd

import (
	"fmt"
	"time"
	"github.com/spf13/cobra"
	"github.com/devops-dojo/cli/internal/project"
	"github.com/devops-dojo/cli/internal/engine"
	"github.com/devops-dojo/cli/internal/engine/scenarios"
	"github.com/briandowns/spinner"
)

var level string

var breakCmd = &cobra.Command{
	Use:   "break",
	Short: "Analyze the current project and inject a failure",
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Analyzing current project for injection targets..."
		s.Start()

		// Simulate deep scan for effect
		time.Sleep(1 * time.Second)
		
		stack, err := project.Analyze(".")
		s.Stop()

		if err != nil {
			fmt.Printf("❌ Error analyzing project: %v\n", err)
			return
		}

		fmt.Printf("🔍 Detected Stack:\n")
		fmt.Printf("- Docker: %v (Found %d Dockerfiles)\n", stack.HasDocker, len(stack.Dockerfiles))
		fmt.Printf("- Kubernetes: %v (Found %d Manifests)\n", stack.HasKubernetes, len(stack.K8sManifests))
		fmt.Printf("- GitHub Actions: %v\n", stack.HasGitHubActions)

		fmt.Printf("\n⚙️  Preparing to inject a %s level failure...\n", level)
		s.Suffix = " Injecting failure..."
		s.Start()
		time.Sleep(1 * time.Second) // Simulate complex injection
		
		// Call Injector
		inj := engine.NewInjector(stack)
		err = inj.InjectFailure(scenarios.Difficulty(level))
		s.Stop()

		if err != nil {
			fmt.Printf("❌ Failed to inject failure: %v\n", err)
			return
		}

		fmt.Println("✅ Failure injected successfully! Your training begins now. Use your standard tools to investigate.")
	},
}

func init() {
	rootCmd.AddCommand(breakCmd)
	breakCmd.Flags().StringVarP(&level, "level", "l", "medium", "Difficulty level (easy, medium, hard, extreme)")
}
