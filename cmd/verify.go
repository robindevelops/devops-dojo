package cmd

import (
	"fmt"
	"time"
	"github.com/spf13/cobra"
	"github.com/devops-dojo/cli/internal/project"
	"github.com/devops-dojo/cli/internal/engine"
	"github.com/devops-dojo/cli/internal/session"
	"github.com/devops-dojo/cli/internal/colors"
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
		
		stack, err := project.Analyze(".", "")
		if err != nil {
			s.Stop()
			fmt.Printf("❌ Error analyzing project: %v\n", err)
			return
		}

		v := engine.NewValidator(stack)
		fixed, err := v.Verify()
		s.Stop()
		
		if err != nil {
			fmt.Println(colors.Colorize(colors.Red, fmt.Sprintf("❌ Error during verification: %v", err)))
			return
		}

		if fixed {
			fmt.Println(colors.Colorize(colors.Green, "\n🎉 Congratulations! You have successfully resolved the incident. You are a Chaos Master!"))
			
			// Scoring logic
			state, err := session.LoadState()
			if err == nil && state != nil {
				timeElapsed := time.Since(state.StartTime).Minutes()
				baseScore := 100
				
				// Deduct points based on time
				timePenalty := int(timeElapsed) * 2
				hintPenalty := state.HintLevel * 10
				
				finalScore := baseScore - timePenalty - hintPenalty
				if finalScore < 10 {
					finalScore = 10 // minimum points
				}

				fmt.Printf("⏱️  Time to resolve: %.1f minutes\n", timeElapsed)
				fmt.Printf("💡 Hints used: %d (-%d points)\n", state.HintLevel, hintPenalty)
				fmt.Printf("🏆 Points Earned: %d\n", finalScore)

				session.AddScore(finalScore, state.HintLevel)
				session.ClearState() // End the active session
			}
		} else {
			// Increment verification attempts
			state, err := session.LoadState()
			if err == nil && state != nil {
				state.VerificationAttempts++
				session.SaveState(state)
			}
			fmt.Println(colors.Colorize(colors.Red, "\n❌ The issue is not fully resolved yet. Keep investigating! Type 'dojo hint' if you get stuck."))
		}
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
