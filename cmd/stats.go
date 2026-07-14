package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/devops-dojo/cli/internal/colors"
	"github.com/devops-dojo/cli/internal/db"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "View your Dojo leaderboard rank and stats",
	Run: func(cmd *cobra.Command, args []string) {
		l, err := db.GetPlayerStats()
		if err != nil {
			fmt.Println("❌ Error loading leaderboard:", err)
			return
		}

		fmt.Println(colors.Colorize(colors.Magenta, "\n🏆 DevOps Dojo - Your Stats 🏆"))
		fmt.Println("==============================")
		fmt.Printf("Rank: %s\n", colors.Colorize(colors.Cyan, db.GetRank(l.TotalXP)))
		fmt.Printf("Total Points: %s\n", colors.Colorize(colors.Green, fmt.Sprintf("%d", l.TotalXP)))
		fmt.Printf("Incidents Resolved: %d\n", l.IncidentsResolved)
		fmt.Printf("Hints Used: %d\n", l.HintsUsed)
		fmt.Println("==============================\n")
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
