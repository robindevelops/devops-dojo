package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [scenario]",
	Short: "Start a predefined learning scenario in sandbox mode",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		scenario := args[0]
		fmt.Printf("Starting sandbox scenario: %s\n", scenario)
		// TODO: Implement Sandbox Provisioner and Injector logic
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
