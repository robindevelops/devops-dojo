package cmd

import (
	"fmt"
	"time"
	"github.com/spf13/cobra"
	"github.com/devops-dojo/cli/internal/project"
	"github.com/devops-dojo/cli/internal/engine"
	"github.com/devops-dojo/cli/internal/session"
	"github.com/devops-dojo/cli/internal/colors"
	"github.com/devops-dojo/cli/internal/db"
	"github.com/devops-dojo/cli/internal/sensei"
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
			state, err := session.LoadState()
			if err == nil && state != nil {
				if state.Mode == "timed" && time.Since(state.StartTime) > 5*time.Minute {
					fmt.Println(colors.Colorize(colors.Red, "\n⏰ Time's up! You fixed it, but took longer than 5 minutes. Session failed."))
					
					db.SaveSession(db.SessionRecord{
						IncidentID:       state.ActiveIncidentID,
						Difficulty:       "unknown",
						Mode:             state.Mode,
						StartTime:        state.StartTime,
						EndTime:          time.Now(),
						Status:           "failed",
						HintsUsed:        state.HintLevel,
						TimeTakenSeconds: int(time.Since(state.StartTime).Seconds()),
						XPGained:         0,
					})
					session.ClearState()
					return
				}
			}

			fmt.Println(colors.Colorize(colors.Green, "\n🎉 Congratulations! You have successfully resolved the incident. You are a Chaos Master!"))
			
			// Scoring logic
			state, err = session.LoadState()
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

				db.SaveSession(db.SessionRecord{
					IncidentID:       state.ActiveIncidentID,
					Difficulty:       "unknown", // To be updated if we track difficulty in state
					Mode:             state.Mode,
					StartTime:        state.StartTime,
					EndTime:          time.Now(),
					Status:           "resolved",
					HintsUsed:        state.HintLevel,
					TimeTakenSeconds: int(time.Since(state.StartTime).Seconds()),
					XPGained:         finalScore,
				})

				// Post-Mortem AI Review
				sessionData := map[string]interface{}{
					"incident": state.ActiveIncidentID,
					"mode": state.Mode,
					"time_taken_minutes": timeElapsed,
					"hints_used": state.HintLevel,
					"score": finalScore,
				}
				fmt.Println(colors.Colorize(colors.Blue, "\n🤖 Generating Post-Mortem Review with Sensei AI..."))
				
				// Need to import sensei in verify.go if not already imported
				// Let's assume it's imported or I will add it
				pm, pmErr := sensei.GeneratePostMortem(sessionData)
				if pmErr == nil {
					fmt.Printf("\n%s\n%s\n", colors.Colorize(colors.Magenta, "📝 Post-Mortem:"), colors.Colorize(colors.Cyan, pm))
				} else {
					fmt.Printf(colors.Colorize(colors.Yellow, "⚠️ Could not generate post-mortem: %v\n"), pmErr)
				}

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
