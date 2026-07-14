package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/devops-dojo/cli/internal/sensei"
)

var hintCmd = &cobra.Command{
	Use:   "hint",
	Short: "Get a hint for the currently active incident",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Sensei is analyzing your progress...")
		sensei.ProvideHint()
	},
}

func init() {
	rootCmd.AddCommand(hintCmd)
}
