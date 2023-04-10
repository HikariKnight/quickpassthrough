package internal

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
	"github.com/HikariKnight/quickpassthrough/internal/logger"
	"github.com/charmbracelet/bubbles/list"
)

func getIOMMU(args ...string) []string {
	var stdout, stderr bytes.Buffer
	// Write to logger
	logger.Printf("Executing: utils/ls-iommu %s", strings.Join(args, " "))

	// Configure the ls-iommu command
	cmd := exec.Command("utils/ls-iommu", args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	// Execute the command
	err := cmd.Run()

	// If ls-iommu returns an error then IOMMU is disabled
	errorcheck.ErrorCheck(err, "IOMMU disabled in either UEFI/BIOS or in bootloader!")

	// Read the output
	var items []string
	output, _ := io.ReadAll(&stdout)

	// Write to logger
	logger.Printf("ls-iommu query returned\n%s", string(output))

	// Parse the output line by line
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		// Write the objects into the list
		items = append(items, scanner.Text())
	}

	// Return our list of items
	return items
}

func iommuList2ListItem(stringList []string) []list.Item {
	// Make the []list.Item struct
	items := []list.Item{}

	deviceID := regexp.MustCompile(`\[[a-f0-9]{4}:[a-f0-9]{4}\]\s+`)
	// Parse the output line by line
	for _, v := range stringList {
		// Get the current line and split by :
		objects := strings.Split(v, ": ")

		// Write the objects into the list
		items = append(items, item{title: deviceID.ReplaceAllString(objects[2], ""), desc: fmt.Sprintf("%s: %s: DeviceID: %s", objects[0], objects[1], objects[3])})
	}

	// Return our list of items
	return items
}
