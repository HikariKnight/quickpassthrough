package pages

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/HikariKnight/quickpassthrough/internal/configs"
	lsiommu "github.com/HikariKnight/quickpassthrough/internal/lsiommu"
	"github.com/HikariKnight/quickpassthrough/pkg/command"
	"github.com/HikariKnight/quickpassthrough/pkg/menu"
)

func genVBIOS_dumper(config *configs.Config) {
	// Clear the scren
	command.Clear()

	// Get the program directory
	exe, _ := os.Executable()
	scriptdir := filepath.Dir(exe)

	// If we are using go run use the working directory instead
	if strings.Contains(scriptdir, "/tmp/go-build") {
		scriptdir, _ = os.Getwd()
	}

	// Get the vbios path and generate the vbios dumping script
	vbios_path := lsiommu.GetIOMMU("-g", "-i", config.Gpu_Group, "--rom")[0]
	configs.GenerateVBIOSDumper(vbios_path)

	// Tell users about the VBIOS dumper script
	fmt.Print(
		"For some GPUs, you will need to dump the VBIOS and pass the\n",
		"rom to the VM along with the card in order to get a functional passthrough.\n",
		"In many cases you can find your vbios at https://www.techpowerup.com/vgabios/\n",
		"\n",
		"You can also attempt to dump your own vbios using the script in\n",
		fmt.Sprintf("%s/utils/dump_vbios.sh\n", scriptdir),
		"\n",
	)

	// Get the OK press
	choice := menu.OkBack("Make sure you run the script with the display-manager stopped using ssh or tty!")

	// Parse choice
	switch choice {
	case "next":
		disableVideo(config)

	case "back":
		SelectGPU(config)

	case "":
		fmt.Println("")
		os.Exit(0)
	}
}
