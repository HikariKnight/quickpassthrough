package main

import (
	"fmt"

	internal "github.com/HikariKnight/quickpassthrough/internal"
	downloader "github.com/HikariKnight/quickpassthrough/internal/ls_iommu_downloader"
	"github.com/HikariKnight/quickpassthrough/internal/params"
	"github.com/HikariKnight/quickpassthrough/internal/version"
)

func main() {
	// Get all our arguments in 1 neat struct
	pArg := params.NewParams()

	if pArg.Flag["version"] {
		fmt.Printf("Quickpassthrough version: %s\n", version.Version)
	} else {
		downloader.CheckLsIOMMU()
		internal.Tui()
	}
}
