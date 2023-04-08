package command

import (
	"bytes"
	"io"
	"os/exec"
)

func Run(binary string, args ...string) ([]string, error) {
	var stdout, stderr bytes.Buffer

	// Configure the ls-iommu c--ommand
	cmd := exec.Command(binary, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	// Execute the command
	err := cmd.Run()

	// Read the output
	output, _ := io.ReadAll(&stdout)
	outerr, _ := io.ReadAll(&stderr)

	outputs := []string{}
	outputs = append(outputs, string(output))
	outputs = append(outputs, string(outerr))

	// Return our list of items
	return outputs, err
}
