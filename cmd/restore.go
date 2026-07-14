package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/briandowns/spinner"
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore your project to its original state",
	Run: func(cmd *cobra.Command, args []string) {
		s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
		s.Suffix = " Restoring files from backup..."
		s.Start()
		time.Sleep(1 * time.Second)

		if _, err := os.Stat(".dojo_backup"); os.IsNotExist(err) {
			s.Stop()
			fmt.Println("ℹ️ No active backup found (.dojo_backup does not exist). Your project is safe!")
			return
		}

		// Restore files by walking the backup directory
		restoredCount := 0
		filepath.Walk(".dojo_backup", func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				// Rel represents the original path
				rel, _ := filepath.Rel(".dojo_backup", path)
				
				content, err := os.ReadFile(path)
				if err == nil {
					os.WriteFile(rel, content, 0644)
					restoredCount++
				}
			}
			return nil
		})

		// Cleanup
		os.RemoveAll(".dojo_backup")
		os.Remove(".dojo_state.json")

		s.Stop()
		fmt.Printf("✅ Project fully restored! (%d files recovered)\n", restoredCount)
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}
