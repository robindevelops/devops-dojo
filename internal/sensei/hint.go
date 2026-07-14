package sensei

import (
	"fmt"
)

// ProvideHint returns a helpful suggestion for the current incident
func ProvideHint() {
	// In a real scenario, this would check a state file for the currently active incident and level
	fmt.Println("💡 Sensei says: 'When investigating Kubernetes pods that fail to start, checking the `kubectl describe` events is your best friend!'")
	fmt.Println("💡 Hint: Look for typos in critical keys like 'apiVersion' or resource limits.")
}
