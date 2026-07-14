package cmd

import (
	"fmt"
	"time"
	"github.com/spf13/cobra"
	"github.com/devops-dojo/cli/internal/project"
	"github.com/devops-dojo/cli/internal/engine"
	"github.com/briandowns/spinner"
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify if you have successfully fixed the injected failure",
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Verifying environment state and analyzing fixes..."
		s.Start()
		
		time.Sleep(1500 * time.Millisecond) // Simulate deep verification
		
		stack, err := project.Analyze(".")
		if err != nil {
			s.Stop()
			fmt.Printf("❌ Error analyzing project: %v\n", err)
			return
		}

		v := engine.NewValidator(stack)
		fixed, err := v.Verify()
		s.Stop()
		
		if err != nil {
			fmt.Printf("❌ Error during verification: %v\n", err)
			return
		}

		if fixed {
			fmt.Println("🎉 Congratulations! You have successfully resolved the incident. You are a Chaos Master!")
		} else {
			fmt.Println("❌ The issue is not fully resolved yet. Keep investigating! Type 'dojo hint' if you get stuck.")
		}
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
