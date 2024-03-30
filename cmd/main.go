package main

import (
	"fmt"

	internal "github.com/HikariKnight/quickpassthrough/internal"
	downloader "github.com/HikariKnight/quickpassthrough/internal/ls_iommu_downloader"
	"github.com/HikariKnight/quickpassthrough/internal/params"
	version "github.com/HikariKnight/quickpassthrough/internal/version"
)

func main() {
	// Get all our arguments in 1 neat struct
	pArg := params.NewParams()

	// Display the version
	if p.Arg.Flag["version"] {
		fmt.Printf("QuickPassthrough Version %s\n", version.Version)
	}

	if !pArg.Flag["gui"] {
		downloader.CheckLsIOMMU()
		internal.Tui()
	}
}
