package menu

import (
	"fmt"
	"strings"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
	"github.com/gookit/color"
)

func ManualInput(msg string, format string) []string {
	// Print the title
	color.Bold.Println(msg)

	// Tell user the format to use
	color.Bold.Printf("The format is %s\n", format)

	// Get the user input
	var input string
	_, err := fmt.Scan(&input)
	errorcheck.ErrorCheck(err)

	input_list := strings.Split(input, ",")

	// Return the input
	return input_list
}
