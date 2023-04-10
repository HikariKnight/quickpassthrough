package logger

import (
	"fmt"
	"log"
)

// Formats our log output to \n%s\n\n for readability
func Printf(content string, v ...any) {
	content = fmt.Sprintf("\n%s\n\n", content)
	log.Printf(content, v...)
}
