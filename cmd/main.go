package main

import (
	"fmt"
	"os"

	internal "github.com/HikariKnight/quickpassthrough/internal"
	downloader "github.com/HikariKnight/quickpassthrough/internal/ls_iommu_downloader"
	"github.com/HikariKnight/quickpassthrough/internal/params"
	"github.com/HikariKnight/quickpassthrough/internal/version"
)

func main() {
	// Get all our arguments in 1 neat struct
	pArg := params.NewParams()

	// Display the version
	if pArg.Flag["version"] {
		fmt.Printf("QuickPassthrough Version %s\n", version.Version)
		os.Exit(0)
	} else {
		downloader.CheckLsIOMMU()
		internal.Tui()
	}
}
