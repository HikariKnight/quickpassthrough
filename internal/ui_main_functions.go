package internal

import (
	"os"
	"regexp"

	"github.com/HikariKnight/quickpassthrough/internal/configs"
)

// This function processes the enter event
func (m *model) processSelection() {
	switch m.focused {
	case GPUS:
		configs.InitConfigs()

		// Gets the selected item
		selectedItem := m.lists[m.focused].SelectedItem()

		// Gets the IOMMU group of the selected item
		iommu_group_regex := regexp.MustCompile(`(\d{1,3})`)
		iommu_group := iommu_group_regex.FindString(selectedItem.(item).desc)

		// Add the gpu group to our model (this is so we can grab the vbios details later)
		m.gpu_group = iommu_group

		// Get all the gpu devices and related devices (same device id or in the same group)
		items := iommuList2ListItem(getIOMMU("-grr", "-i", m.gpu_group, "-F", "vendor:,prod_name,optional_revision:,device_id"))

		// Add the devices to the list
		m.lists[GPU_GROUP].SetItems(items)

		// Change focus to next index
		m.focused++

	case GPU_GROUP:
		// Generate the VBIOS dumper script once the user has selected a GPU
		generateVBIOSDumper(*m)
		m.focused++

	case USB:
		// Gets the selected item
		selectedItem := m.lists[m.focused].SelectedItem()

		// Gets the IOMMU group of the selected item
		iommu_group_regex := regexp.MustCompile(`(\d{1,3})`)
		iommu_group := iommu_group_regex.FindString(selectedItem.(item).desc)

		// Get the USB controllers in the selected iommu group
		items := iommuList2ListItem(getIOMMU("-ur", "-i", iommu_group, "-F", "vendor:,prod_name,optional_revision:,device_id"))

		// Add the items to the list
		m.lists[USB_GROUP].SetItems(items)

		// Change focus to next index
		m.focused++

	case USB_GROUP:
		m.focused++

	case VBIOS:
		m.focused++

	case VIDEO:
		m.focused++

	case INTRO:
		m.focused++

	case DONE:
		os.Exit(0)
	}
}
