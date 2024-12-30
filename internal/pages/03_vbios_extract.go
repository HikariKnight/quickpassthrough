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
	"github.com/gookit/color"
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

	// Search for a vbios path and generate the vbios dumping script if found
	vbios_paths := lsiommu.GetIOMMU("-g", "-i", config.Gpu_Group, "--rom")
	if len(vbios_paths) != 0 {
		configs.GenerateVBIOSDumper(vbios_paths[0])
	}

	// Make the qemu config folder
	os.Mkdir(fmt.Sprintf("%s/%s", scriptdir, config.Path.QEMU), os.ModePerm)

	// Generate a dummy rom (1MB rom of zeroes) for use with AMD RX 7000 series cards by recommendation from Gnif
	// Source: https://forum.level1techs.com/t/the-state-of-amd-rx-7000-series-vfio-passthrough-april-2024/210242
	command.Run("dd", "if=/dev/zero", fmt.Sprintf("of=%s/%s/dummy.rom", scriptdir, config.Path.QEMU), "bs=1M", "count=1")

	// Write a title
	title := color.New(color.BgHiBlue, color.White, color.Bold)
	title.Println("VBIOS roms for Passthrough")

	// Tell users about the VBIOS dumper script and dummy rom for RX 7000 series cards
	fmt.Print(
		"If you have an RX 7000 series (and possibly newer AMD cards) GPUs, please use the dummy.rom file\n",
		fmt.Sprintf("%s/%s/dummy.rom\n", scriptdir, config.Path.QEMU),
		"Or disable ROM BAR for the card in qemu/libvirt\n",
		"\n",
		"For some other GPUs, you will need to instead dump the VBIOS (and possibly patch it) and pass the\n",
		"rom to the VM along with the card in order to get a functional passthrough.\n",
		"In many cases you can find your vbios at https://www.techpowerup.com/vgabios/\n",
		"\n",
		"If we found a romfile for your GPU you can also attempt to dump your own vbios from TTY using the script in\n",
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
