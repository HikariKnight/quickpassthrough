package main

import (
	internal "github.com/HikariKnight/quickpassthrough/internal"
	downloader "github.com/HikariKnight/quickpassthrough/internal/ls_iommu_downloader"
	"github.com/HikariKnight/quickpassthrough/internal/params"
)

func main() {
	// Get all our arguments in 1 neat struct
	pArg := params.NewParams()

	if !pArg.Flag["gui"] {
		downloader.CheckLsIOMMU()
		internal.Tui()
	}
}
