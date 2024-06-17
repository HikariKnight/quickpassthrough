package menu

import (
	"fmt"
	"strings"

	"github.com/gookit/color"

	"github.com/HikariKnight/quickpassthrough/internal/common"
)

func ManualInput(msg string, format string) []string {
	// Print the title
	color.Bold.Println(msg)

	// Tell user the format to use
	color.Bold.Printf("The format is %s\n", format)

	// Get the user input
	var input string
	_, err := fmt.Scan(&input)
	common.ErrorCheck(err)

	input_list := strings.Split(input, ",")

	// Return the input
	return input_list
}
