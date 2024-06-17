package fileio

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"

	"github.com/HikariKnight/quickpassthrough/internal/common"
)

/*
 * This just implements repetetive tasks I have to do with files
 */

// AppendContent creates a file and appends the content to the file.
// (ending newline must be supplied with content string)
func AppendContent(content string, fileName string) {
	// Open the file
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)

	errorcheck.ErrorCheck(err, fmt.Sprintf("Error opening \"%s\" for writing", fileName))
	defer f.Close()

	// Write the content
	_, err = f.WriteString(content)
	errorcheck.ErrorCheck(err, fmt.Sprintf("Error writing to %s", fileName))
}

// ReadLines reads the file and returns a stringlist with each line.
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

// ReadFile reads a file and returns all the content as a string.
func ReadFile(fileName string) string {
	// Read the whole file
	content, err := os.ReadFile(fileName)
	if errors.Is(err, os.ErrPermission) {
		errorcheck.ErrorCheck(err, common.PermissionNotice)
		return "" // note: unreachable due to ErrorCheck calling fatal
	}
	errorcheck.ErrorCheck(err, fmt.Sprintf("Failed to ReadFile on %s", fileName))

	// Return all the lines as one string
	return string(content)

}

// FileExist checks if a file exists and returns a bool and any error that isn't os.ErrNotExist.
func FileExist(fileName string) (bool, error) {
	var exist bool

	// Check if the file exists
	_, err := os.Stat(fileName)
	switch {
	case err == nil:
		exist = true
	case errors.Is(err, os.ErrNotExist):
		// Set the value to true
		exist = false
		err = nil
	}
	// Return if the file exists
	return exist, err
}

// FileCopy copies a FILE from source to dest.
func FileCopy(sourceFile, destFile string) {
	// Get the file info
	filestat, err := os.Stat(sourceFile)
	errorcheck.ErrorCheck(err, "Error getting fileinfo of: %s", sourceFile)

	// If the file is a regular file
	if filestat.Mode().IsRegular() {
		// Open the source file for reading
		source, err := os.Open(sourceFile)
		errorcheck.ErrorCheck(err, "Error opening %s for copying", sourceFile)
		defer source.Close()

		// Create the destination file
		dest, err := os.Create(destFile)
		errorcheck.ErrorCheck(err, "Error creating %s", destFile)
		defer dest.Close()

		// Copy the contents of source to dest using io
		_, err = io.Copy(dest, source)
		errorcheck.ErrorCheck(err, "Failed to copy \"%s\" to \"%s\"", sourceFile, destFile)
	}
}
