package lsiommu

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/HikariKnight/quickpassthrough/internal/logger"
	"github.com/klauspost/cpuid/v2"
)

func GetIOMMU(args ...string) []string {
	var stdout, stderr bytes.Buffer
	// Write to logger
	logger.Printf("Executing: utils/ls-iommu %s\n", strings.Join(args, " "))

	// Configure the ls-iommu command
	cmd := exec.Command("utils/ls-iommu", args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	// Execute the command
	err := cmd.Run()

	// Generate the correct iommu string for the system
	var iommu_args string
	cpuinfo := cpuid.CPU
	// Write the argument based on which cpu the user got
	switch cpuinfo.VendorString {
	case "AuthenticAMD":
		iommu_args = "iommu=pt amd_iommu=on"
	case "GenuineIntel":
		iommu_args = "iommu=pt intel_iommu=on"
	}

	// If ls-iommu returns an error then IOMMU is disabled
	if err != nil {
		fmt.Printf(
			"IOMMU disabled in either UEFI/BIOS or in bootloader, or run inside container!\n"+
				"For your bootloader, make sure you have added the kernel arguments:\n"+
				"%s\n",
			iommu_args,
		)
		os.Exit(1)
	}

	// Read the output
	var items []string
	output, _ := io.ReadAll(&stdout)

	// Write to logger
	logger.Printf("ls-iommu query returned\n%s", string(output))

	// Make regex to shorten vendor names
	shortenVendor := regexp.MustCompile(` Corporation:| Technology Inc.:| Electronics Co Ltd:|Advanced Micro Devices, Inc\. \[(AMD)(|\/ATI)\]:`)

	// Parse the output line by line
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		// Write the objects into the list
		items = append(items, shortenVendor.ReplaceAllString(scanner.Text(), "${1}"))
	}

	// Return our list of items
	return items
}
