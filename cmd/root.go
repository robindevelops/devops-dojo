package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dojo",
	Short: "DevOps Dojo - Chaos Engineering for Learning",
	Long: `DevOps Dojo is a CLI tool that injects production-grade failures into your pipeline or infrastructure for training purposes.
You can use it in sandbox mode or point it at your own project (BYOP) to break your own configurations.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
