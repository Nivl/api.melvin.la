package logger

import "fmt"
import "strings"

func formatError(staticData []string, msg string) string {
	static := strings.Join(staticData, ", ")
	if static != "" {
		static = ", " + static
	}
	return fmt.Sprintf(`level: "ERROR"%s, %s"`, static, msg)
}
