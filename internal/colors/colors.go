package colors

import "fmt"

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Bold    = "\033[1m"
)

// Colorize wraps a string with the given ANSI color code and resets it at the end.
func Colorize(color string, text string) string {
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}
