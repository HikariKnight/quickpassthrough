package command

import (
	"bytes"
	"encoding/base64"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
)

// Run a command and return STDOUT
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

	// Get the output
	outputs := []string{}
	outputs = append(outputs, string(output))

	// Return our list of items
	return outputs, err
}

// This function is just like command.Run() but also returns STDERR
func RunErr(binary string, args ...string) ([]string, []string, error) {
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

	// Get the output
	var outputs, outerrs []string
	outputs = append(outputs, string(output))
	outerrs = append(outerrs, string(outerr))

	// Return our list of items
	return outputs, outerrs, err
}

// This functions runs the command "sudo -Sk -- echo", this forces sudo
// to re-authenticate and lets us enter the password to STDIN
// giving us the ability to run sudo commands
func Elevate(password string) {
	// Do a simple sudo command to just authenticate with sudo
	cmd := exec.Command("sudo", "-S", "--", "echo")

	// Wait for 500ms, if the password is correct, sudo will return immediately
	cmd.WaitDelay = 1000 * time.Millisecond

	// Open STDIN
	stdin, err := cmd.StdinPipe()
	errorcheck.ErrorCheck(err, "\nFailed to get sudo STDIN")

	// Start the authentication
	cmd.Start()

	// Get the passed password
	pw, _ := base64.StdEncoding.DecodeString(password)
	_, err = stdin.Write([]byte(string(pw) + "\n"))
	errorcheck.ErrorCheck(err, "\nFailed at typing to STDIN")
	// Clear the password
	pw = nil
	password = ""

	stdin.Close()

	// Wait for the sudo prompt (If the correct password was given, it will not stay behind)
	err = cmd.Wait()
	errorcheck.ErrorCheck(err, "\nError, password given was wrong")
}

// Function to just clear the terminal
func Clear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
}
