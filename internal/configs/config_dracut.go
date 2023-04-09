package configs

import (
	"fmt"

	"github.com/HikariKnight/quickpassthrough/pkg/fileio"
)

func Set_Dracut() {
	config := GetConfig()

	// Write the dracut config file
	fileio.AppendContent("add_drivers+=\" vfio_pci vfio vfio_iommu_type1 vfio_virqfd \"\n", fmt.Sprintf("%s/vfio.conf", config.Path.DRACUT))

	// Add to our kernel arguments file that vfio_pci should load early (dracut does this using kernel arguments)
	fileio.AppendContent(" rd.driver.pre=vfio_pci", config.Path.CMDLINE)
}
