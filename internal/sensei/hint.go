package sensei

import (
	"fmt"
	"os"
	
	"github.com/devops-dojo/cli/internal/colors"
	"github.com/devops-dojo/cli/internal/session"
)

func ProvideHint() {
	state, err := session.LoadState()
	if err != nil || state == nil {
		fmt.Println(colors.Colorize(colors.Yellow, "💡 Sensei says: 'You must start a Dojo session before asking for hints!'"))
		return
	}

	fmt.Println(colors.Colorize(colors.Cyan, "Sensei is analyzing your progress..."))

	// Escalate hint level
	state.HintLevel++
	session.SaveState(state)

	// AI Mentor Mode
	if os.Getenv("DOJO_API_KEY") != "" {
		fmt.Println(colors.Colorize(colors.Blue, "🤖 Consulting AI Sensei..."))
		hint, err := AskAI(state)
		if err == nil {
			fmt.Printf("\n%s\n%s\n", colors.Colorize(colors.Yellow, "💡 AI Sensei says:"), colors.Colorize(colors.Cyan, hint))
			return
		}
		fmt.Printf("%s\n\n", colors.Colorize(colors.Red, fmt.Sprintf("⚠️  AI API failed (%v). Falling back to standard hints.", err)))
	}
	
	// In a real system, these would be loaded from a map or struct per incident
	if state.ActiveIncidentID == "medium-break" || state.ActiveIncidentID == "k8s-oomkilled" {
		switch state.HintLevel {
		case 1:
			fmt.Println("💡 Hint 1 (Vague): Look closely at the pod events. Why is it restarting?")
		case 2:
			fmt.Println("💡 Hint 2 (Direction): Check the `resources.limits.memory` in the deployment.")
		case 3:
			fmt.Println("💡 Hint 3 (Root Cause): The memory limit is set to 1Mi, which is too low for most applications.")
		default:
			fmt.Println("💡 Solution: Remove or increase the memory limit in the pod spec to fix OOMKilled.")
		}
	} else {
		// Generic fallback
		switch state.HintLevel {
		case 1:
			fmt.Println("💡 Hint 1 (Vague): Look at recent logs or events to find the failure.")
		case 2:
			fmt.Println("💡 Hint 2 (Direction): Check for typos or missing keys in your configuration files.")
		case 3:
			fmt.Println("💡 Hint 3 (Command): Run tools like `docker logs` or `kubectl describe`.")
		default:
			fmt.Println("💡 Solution: Review the git diff against your backup or previous commit to find the exact change.")
		}
	}

	session.SaveState(state)
}
