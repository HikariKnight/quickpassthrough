package configs

import (
	"github.com/HikariKnight/quickpassthrough/pkg/command"
)

func getBootloader(config *Config) {
	// Check what bootloader handler we are using
	// Check for grub-mkconfig
	_, err := command.Run("which", "grub-mkconfig")
	if err == nil {
		// Mark bootloader as grub2
		config.bootloader = "grub2"
	}

	// Check for grubby (used by fedora)
	_, err = command.Run("which", "grubby")
	if err == nil {
		// Mark it as unknown as i do not support it yet
		config.bootloader = "unknown"
	}

	// Check for kernelstub (used by pop os)
	_, err = command.Run("which", "kernelstub")
	if err == nil {
		config.bootloader = "kernelstub"
	}
}
