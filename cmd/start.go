package cmd

import (
	"fmt"
	"time"
	"github.com/spf13/cobra"
	"github.com/devops-dojo/cli/internal/session"
)

var startCmd = &cobra.Command{
	Use:   "start [scenario]",
	Short: "Start a predefined learning scenario in sandbox mode",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		scenario := args[0]
		mode := args[1]
		fmt.Printf("Starting sandbox scenario: %s in %s mode\n", scenario, mode)
		
		// TODO: Implement Sandbox Provisioner and Injector logic
		// For now, just save state so the session is active
		err := session.SaveState(&session.State{
			ActiveIncidentID:     scenario,
			StartTime:            time.Now(),
			VerificationAttempts: 0,
			HintLevel:            0,
			Mode:                 mode,
		})
		if err != nil {
			fmt.Printf("⚠️  Warning: Failed to save session state: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
