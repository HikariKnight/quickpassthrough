package fileio

import (
	"bufio"
	"fmt"
	"os"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
)

// This just implements repetetive tasks I have to do with files

func AppendContent(content string, fileName string) {
	// Open the file
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	errorcheck.ErrorCheck(err, fmt.Sprintf("Error opening \"%s\" for writing", fileName))
	defer f.Close()

	// Write the content
	_, err = f.WriteString(content)
	errorcheck.ErrorCheck(err, fmt.Sprintf("Error writing to %s", fileName))
}

func ReadLines(fileName string) []string {
	content, err := os.Open(fileName)
	errorcheck.ErrorCheck(err, fmt.Sprintf("Error reading file %s", fileName))
	defer content.Close()

	// Make a list of lines
	var lines []string

	// Read the file line by line
	scanner := bufio.NewScanner(content)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Return all the lines
	return lines

}

func ReadFile(fileName string) string {
	// Read the whole file
	content, err := os.ReadFile(fileName)
	errorcheck.ErrorCheck(err, fmt.Sprintf("Failed to ReadFile on %s", fileName))

	// Return all the lines as one string
	return string(content)

}
