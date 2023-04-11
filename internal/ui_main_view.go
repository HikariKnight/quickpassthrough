package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) View() string {
	if m.width != 0 {
		title := ""
		view := ""
		switch m.focused {
		case INTRO:
			title = dialogStyle.Render(
				fmt.Sprint(
					titleStyle.MarginLeft(0).Render("Welcome to QuickPassthrough!"),
					"\n\n",
					"This script is meant to make it easier to setup GPU passthrough for\n",
					"Qemu based systems.\n",
					"However due to the complexity of GPU passthrough\n",
					"This script assumes you know how to do (or have done) the following.\n\n",
					"* You have already enabled IOMMU, VT-d, SVM and/or AMD-v\n  inside your UEFI/BIOS advanced settings.\n",
					"* Know how to edit your bootloader\n",
					"* Have a bootloader timeout of at least 3 seconds to access the menu\n",
					"* Enable & Configure kernel modules\n",
					"* Have a backup/snapshot of your system in case the script causes your\n  system to be unbootable\n\n",
					"By continuing you accept that I am not liable if your system\n",
					"becomes unbootable, as you will be asked to verify the files generated",
				),
			)

			view = listStyle.Render(m.lists[m.focused].View())

		case GPUS:
			title = titleStyle.MarginLeft(2).Render(
				"Select a GPU to check the IOMMU groups of",
			)

			view = listStyle.Render(m.lists[m.focused].View())

		case GPU_GROUP:
			title = titleStyle.Render(
				fmt.Sprint(
					"Press ENTER/RETURN to set up all these devices for passthrough.\n",
					"This list should only contain items related to your GPU.",
				),
			)

			view = listStyle.Render(m.lists[m.focused].View())

		case USB:
			title = titleStyle.Render(
				"[OPTIONAL]: Select a USB Controller to check the IOMMU groups of",
			)

			view = listStyle.Render(m.lists[m.focused].View())

		case USB_GROUP:
			title = titleStyle.Render(
				fmt.Sprint(
					"Press ENTER/RETURN to set up all these devices for passthrough.\n",
					"This list should only contain the USB controller you want to use.",
				),
			)

			view = listStyle.Render(m.lists[m.focused].View())

		case VBIOS:
			// Get the program directory
			exe, _ := os.Executable()
			scriptdir := filepath.Dir(exe)

			// If we are using go run use the working directory instead
			if strings.Contains(scriptdir, "/tmp/go-build") {
				scriptdir, _ = os.Getwd()
			}

			text := dialogStyle.Render(
				fmt.Sprint(
					"Based on your GPU selection, a vbios extraction script has been generated for your convenience.\n",
					"Passing a VBIOS rom to the card used for passthrough is required for some cards, but not all.\n",
					"Some cards also requires you to patch your VBIOS romfile, check online if this is neccessary for your card!\n",
					"The VBIOS will be read from:\n",
					fmt.Sprintf(
						"%s\n\n",
						m.vbios_path,
					),
					"The script to extract the vbios has to be run as sudo and without a displaymanager running for proper dumping!\n",
					"\n",
					"You can run the script with:\n",
					fmt.Sprintf(
						"%s/utils/dump_vbios.sh",
						scriptdir,
					),
				),
			)

			title = fmt.Sprintf(text, m.vbios_path, scriptdir)

			view = listStyle.Render(m.lists[m.focused].View())

		case VIDEO:
			title = dialogStyle.Render(
				fmt.Sprint(
					"Disabling video output in Linux for the card you want to use in a VM\n",
					"will make it easier to successfully do the passthrough without issues.\n",
					"\n",
					"Do you want to force disable video output in linux on this card?",
				),
			)

			view = listStyle.Render(m.lists[m.focused].View())

		case INSTALL:
			title = dialogStyle.Render(
				fmt.Sprint(
					"The configuration files have been generated and are\n",
					"located inside the \"config\" folder\n",
					"\n",
					"* The \"kernel_args\" file contains kernel arguments that your bootloader needs\n",
					"* The \"quickemu\" folder contains files that might be\n  useable for quickemu in the future\n",
					"* The files inside the \"etc\" folder must be copied to your system.\n",
					"  NOTE: Verify that these files are correctly formated/edited!\n",
					"* Once all files have been copied, you need to update your bootloader and rebuild\n",
					"  your initramfs using the tools to do so by your system.\n",
					"\n",
					"This program can do this for you, however the program will have to\n",
					"type your password to sudo using STDIN, to avoid using STDIN press CTRL+C\n",
					"and copy the files, update your bootloader and rebuild your initramfs manually.\n",
					"If you want to go back and change something, press CTRL+Z\n",
					"\nNOTE: A backup of the original files from the first run can be found in the backup folder",
				),
			)

			view = m.authDialog.View()

		case DONE:
			title = titleStyle.Render("Applying configurations!")
			view = dialogStyle.Render(strings.Join(m.installOutput, "\n"))
		}
		//return listStyle.SetString(fmt.Sprintf("%s\n\n", title)).Render(m.lists[m.focused].View())
		return lipgloss.JoinVertical(lipgloss.Left, fmt.Sprintf("%s\n%s\n", title, view))
	} else {
		return "Loading..."
	}
}
