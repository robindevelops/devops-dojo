package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/manifoldco/promptui"
	"github.com/devops-dojo/cli/internal/colors"
	"github.com/devops-dojo/cli/internal/engine/scenarios"
	"github.com/devops-dojo/cli/internal/db"
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
		fmt.Println(colors.Colorize(colors.Cyan, dojoLogo))
		fmt.Println(colors.Colorize(colors.Green, "Welcome to DevOps Dojo! Choose an action to begin:"))
		
		for {
			prompt := promptui.Select{
				Label: "Select Action",
				Items: []string{
					"Start Sandbox Scenario",
					"Break Current Project (BYOP Mode)",
					"Verify Fix",
					"Get a Hint",
					"View Stats",
					"Restore Project (Undo Break)",
					"Show Help & Commands",
					"Exit",
				},
				Size: 10,
			}

			_, result, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			switch result {
			case "Start Sandbox Scenario":
				// Prompt for Practice Mode
				modePrompt := promptui.Select{
					Label: "Select Practice Mode",
					Items: []string{
						"Normal (Hints allowed)",
						"Timed Challenge (5 Minutes, no hints)",
						"Blind Debugging (No hints)",
						"Interview Mode (AI Interviewer)",
						"Back",
					},
				}
				_, modeChoice, _ := modePrompt.Run()
				if modeChoice == "Back" {
					continue
				}

				modeStr := "normal"
				if modeChoice == "Timed Challenge (5 Minutes, no hints)" {
					modeStr = "timed"
				} else if modeChoice == "Blind Debugging (No hints)" {
					modeStr = "blind"
				} else if modeChoice == "Interview Mode (AI Interviewer)" {
					modeStr = "interview"
				}

				// Fetch the catalog dynamically
				incidents := scenarios.GetAvailableIncidents()
				var items []string
				for _, inc := range incidents {
					items = append(items, fmt.Sprintf("%s (%s)", inc.Name, inc.Difficulty))
				}
				items = append(items, "Back")

				scenarioPrompt := promptui.Select{
					Label: "Select Scenario",
					Items: items,
					Size:  10,
				}
				idx, scenario, _ := scenarioPrompt.Run()
				if scenario == "Back" {
					continue
				}
				startCmd.Run(startCmd, []string{incidents[idx].ID, modeStr})
				return
			case "Break Current Project (BYOP Mode)":
				levelPrompt := promptui.Select{
					Label: "Select Difficulty Level",
					Items: []string{"easy", "medium", "hard", "extreme", "Back"},
				}
				_, levelResult, _ := levelPrompt.Run()
				if levelResult == "Back" {
					continue
				}
				// Override level for break command
				level = levelResult
				breakCmd.Run(breakCmd, []string{})
				return
			case "Verify Fix":
				verifyCmd.Run(verifyCmd, []string{})
				return
			case "Get a Hint":
				hintCmd.Run(hintCmd, []string{})
				return
			case "View Stats":
				statsCmd.Run(statsCmd, []string{})
				// Do not return here; let the user read stats and pick another action
			case "Restore Project (Undo Break)":
				restoreCmd.Run(restoreCmd, []string{})
				return
			case "Show Help & Commands":
				cmd.Help()
				// Do not return here; let the user read help and pick another action
			case "Exit":
				fmt.Println("Goodbye, Chaos Master!")
				os.Exit(0)
			}
		}
	},
}

func Execute() {
	if err := db.InitDB(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize database: %v\n", err)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
