package command

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/HikariKnight/quickpassthrough/internal/common"
	"github.com/HikariKnight/quickpassthrough/internal/logger"
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
	outputs := make([]string, 0, 1)
	outputs = append(outputs, string(output))

	// Return our list of items
	return outputs, err
}

// RunErr is just like command.Run() but also returns STDERR
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

func RunErrSudo(isRoot bool, binary string, args ...string) ([]string, []string, error) {
	if !isRoot && binary != "sudo" {
		args = append([]string{binary}, args...)
		binary = "sudo"
	}
	logger.Printf("Executing (elevated): %s %s\n", binary, strings.Join(args, " "))
	fmt.Printf("Executing (elevated): %s %s\n", binary, strings.Join(args, " "))
	return RunErr(binary, args...)
}

// Elevate elevates this functions runs the command "sudo -Sk -- echo",
// this forces sudo to re-authenticate and lets us enter the password to STDIN
// giving us the ability to run sudo commands
func Elevate(password string) {
	// Do a simple sudo command to just authenticate with sudo
	cmd := exec.Command("sudo", "-S", "--", "echo")

	// Wait for 500ms, if the password is correct, sudo will return immediately
	cmd.WaitDelay = 1000 * time.Millisecond

	// Open STDIN
	stdin, err := cmd.StdinPipe()
	common.ErrorCheck(err, "\nFailed to get sudo STDIN")

	// Start the authentication
	err = cmd.Start()
	common.ErrorCheck(err, "\nFailed to start sudo command")

	// Get the passed password
	pw, _ := base64.StdEncoding.DecodeString(password)
	_, err = stdin.Write([]byte(string(pw) + "\n"))
	common.ErrorCheck(err, "\nFailed at typing to STDIN")
	// Clear the password
	pw = nil
	password = ""

	_ = stdin.Close()

	// Wait for the sudo prompt (If the correct password was given, it will not stay behind)
	err = cmd.Wait()
	common.ErrorCheck(err, "\nError, password given was wrong")
}

// Clear clears the terminal.
func Clear() {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	_ = c.Run()
}

func processCmdString(cmd string) (string, []string) {
	// handle quoted arguments
	args := strings.Fields(cmd)
	cmdBin := args[0]
	args = args[1:]
	for i, arg := range args {
		if !strings.HasPrefix(arg, "\"") {
			continue
		}
		// find the end of the quoted argument
		for j, a := range args[i:] {
			if strings.HasSuffix(a, "\"") {
				args[i] = strings.Join(args[i:i+j+1], " ")
				args = append(args[:i+1], args[i+j+1:]...)
				break
			}
		}
	}

	return cmdBin, args
}

// ExecAndLogSudo executes an elevated command and logs the output.
//
// * if we're root, the command is executed directly
// * if we're not root, the command is prefixed with "sudo"
//
//   - noisy determines if we should print the command to the user
//     noisy isn't set to true by our copy caller, as it logs differently,
//     but other callers set it.
func ExecAndLogSudo(isRoot, noisy bool, cmd string) error {
	if !isRoot && !strings.HasPrefix(cmd, "sudo") {
		cmd = fmt.Sprintf("sudo %s", cmd)
	}

	// Write to logger
	logger.Printf("Executing (elevated): %s\n", cmd)

	if noisy {
		// Print to the user
		fmt.Printf("Executing (elevated): %s\nSee debug.log for detailed output\n", cmd)
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	cmdBin, args := processCmdString(cmd)
	r := exec.Command(cmdBin, args...)
	r.Dir = wd

	cmdCombinedOut, err := r.CombinedOutput()
	outStr := string(cmdCombinedOut)

	// Write to logger, tabulate output
	// tabulation denotes it's hierarchy as a child of the command
	outStr = strings.ReplaceAll(outStr, "\n", "\n\t")
	logger.Printf("\t" + string(cmdCombinedOut) + "\n")
	if noisy {
		// Print to the user
		fmt.Printf("%s\n", outStr)
	}

	if err != nil {
		err = fmt.Errorf("failed to execute %s: %w\n%s", cmd, err, outStr)
	}

	return err
}
