package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/manifoldco/promptui"
	"github.com/devops-dojo/cli/internal/engine/scenarios"
)

const dojoLogo = `
 ____             ___                   ____        _       
|  _ \  _____   _/ _ \ _ __  ___       |  _ \  ___ (_) ___  
| | | |/ _ \ \ / / | | | '_ \/ __|      | | | |/ _ \| |/ _ \ 
| |_| |  __/\ V /| |_| | |_) \__ \      | |_| | (_) | | (_) |
|____/ \___| \_/  \___/| .__/|___/      |____/ \___// |\___/ 
                       |_|                        |__/       
`

var rootCmd = &cobra.Command{
	Use:   "dojo",
	Short: "DevOps Dojo - Chaos Engineering for Learning",
	Long: dojoLogo + `
DevOps Dojo is a CLI tool that injects production-grade failures into your pipeline or infrastructure for training purposes.
You can use it in sandbox mode or point it at your own project (BYOP) to break your own configurations.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(dojoLogo)
		fmt.Println("Welcome to DevOps Dojo! Choose an action to begin:")
		
		prompt := promptui.Select{
			Label: "Select Action",
			Items: []string{
				"Start Sandbox Scenario",
				"Break Current Project (BYOP Mode)",
				"Verify Fix",
				"Get a Hint",
				"Show Help & Commands",
				"Exit",
			},
		}

		_, result, err := prompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		switch result {
		case "Start Sandbox Scenario":
			// Fetch the catalog dynamically
			incidents := scenarios.GetAvailableIncidents()
			var items []string
			for _, inc := range incidents {
				items = append(items, fmt.Sprintf("%s (%s)", inc.Name, inc.Difficulty))
			}
			items = append(items, "Cancel")

			scenarioPrompt := promptui.Select{
				Label: "Select Scenario",
				Items: items,
				Size:  10,
			}
			idx, scenario, _ := scenarioPrompt.Run()
			if scenario != "Cancel" {
				startCmd.Run(startCmd, []string{incidents[idx].ID})
			}
		case "Break Current Project (BYOP Mode)":
			levelPrompt := promptui.Select{
				Label: "Select Difficulty Level",
				Items: []string{"easy", "medium", "hard", "extreme"},
			}
			_, levelResult, _ := levelPrompt.Run()
			// Override level for break command
			level = levelResult
			breakCmd.Run(breakCmd, []string{})
		case "Verify Fix":
			verifyCmd.Run(verifyCmd, []string{})
		case "Get a Hint":
			hintCmd.Run(hintCmd, []string{})
		case "Show Help & Commands":
			cmd.Help()
		case "Exit":
			fmt.Println("Goodbye, Chaos Master!")
			os.Exit(0)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
