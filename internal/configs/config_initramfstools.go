package configs

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
	"github.com/HikariKnight/quickpassthrough/pkg/uname"
)

func readHeader(lines int, fileName string) string {
	// Open the file
	f, err := os.Open(fileName)
	errorcheck.ErrorCheck(err, fmt.Sprintf("Error opening %s", fileName))
	defer f.Close()

	header_re := regexp.MustCompile(`^#`)
	var header []string

	// Make a new scanner
	scanner := bufio.NewScanner(f)

	// Read the first N lines
	for i := 0; i < lines; i++ {
		scanner.Scan()
		if header_re.MatchString(scanner.Text()) {
			header = append(header, scanner.Text())
		}
	}

	// Return the header
	return fmt.Sprintf("%s\n", strings.Join(header, "\n"))
}

func addModules(conffile string) {
	// Make a regex to get the system path instead of the config path
	syspath_re := regexp.MustCompile(`^config`)

	// Make a regex to skip specific modules and comments
	skipmodules_re := regexp.MustCompile(`(^#|vendor-reset)`)

	// Get the syspath
	syspath := syspath_re.ReplaceAllString(conffile, "")

	// Open the system file for reading
	sysfile, err := os.Open(syspath)
	errorcheck.ErrorCheck(err, fmt.Sprintf("Error opening file for reading %s", syspath))
	defer sysfile.Close()

	// Open config file for writing
	out, err := os.OpenFile(conffile, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	errorcheck.ErrorCheck(err, fmt.Sprintf("Error opening file for writing %s", conffile))
	defer out.Close()

	// Make a list of modules
	//var modules []string

	// Scan the file line by line
	scanner := bufio.NewScanner(sysfile)
	for scanner.Scan() {
		if scanner.Text() == "vendor-reset" {
			out.WriteString(scanner.Text())
		} else if !skipmodules_re.MatchString(scanner.Text()) {
			writeContent(fmt.Sprintf("%s\n", scanner.Text()), conffile)
			sysinfo := uname.New()
			writeContent(fmt.Sprintf("%s\n%s\n%s\n%s\n", sysinfo.Nodename, sysinfo.Sysname, sysinfo.Domainname, sysinfo.Machine), conffile)
		}
	}
}

func writeContent(content string, fileName string) {
	// Open the file
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	errorcheck.ErrorCheck(err, fmt.Sprintf("Error opening %s", fileName))
	defer f.Close()

	// Make a new scanner
	_, err = f.WriteString(content)
	errorcheck.ErrorCheck(err, fmt.Sprintf("Error writing to %s", fileName))
}
