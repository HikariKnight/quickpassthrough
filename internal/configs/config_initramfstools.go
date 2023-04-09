package configs

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
	"github.com/HikariKnight/quickpassthrough/pkg/fileio"
)

// Special function to read the header of a file (reads the first N lines)
func initramfs_readHeader(lines int, fileName string) string {
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

// Reads the system file and copies over the content while inserting the vfio modules
// Takes the config file as argument
func initramfs_addModules(conffile string) {
	// Make a regex to get the system path instead of the config path
	syspath_re := regexp.MustCompile(`^config`)

	// Make a regex to skip specific modules and comments
	skipmodules_re := regexp.MustCompile(`(^#|vendor-reset|vfio|vfio_pci|vfio_iommu_type1|vfio_virqfd)`)

	// Get the syspath
	syspath := syspath_re.ReplaceAllString(conffile, "")

	// Open the system file for reading
	sysfile, err := os.Open(syspath)
	errorcheck.ErrorCheck(err, fmt.Sprintf("Error opening file for reading %s", syspath))
	defer sysfile.Close()

	// Check if user has vendor-reset installed/enabled and make sure that is first
	content := fileio.ReadFile(syspath)
	if strings.Contains(content, "vendor-reset") {
		fileio.AppendContent("vendor-reset\n", conffile)
	}

	// Write the vfio modules
	fileio.AppendContent(
		fmt.Sprint(
			"# Added by quickpassthrough #\n",
			fmt.Sprintf(
				"%s\n",
				strings.Join(vfio_modules(), "\n"),
			),
			"#############################\n",
		),
		conffile,
	)

	// Scan the system file line by line
	scanner := bufio.NewScanner(sysfile)
	for scanner.Scan() {
		// If this is not a line we skip then
		if !skipmodules_re.MatchString(scanner.Text()) {
			// Add the module to our config
			fileio.AppendContent(fmt.Sprintf("%s\n", scanner.Text()), conffile)
		}
	}
}
