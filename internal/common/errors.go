package common

import (
	"os"
	"time"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
	"github.com/gookit/color"
)

const PermissionNotice = `
<yellowB>Permissions error occured during file operations.</>

<blue_b>Hint</>:

	If you initially ran QuickPassthrough as root or using sudo,
	but are now running it as a normal user, this is expected behavior.

	<us>Try running QuickPassthrough as root or using sudo if so.</>

	If this does not work, double check your filesystem's permissions,
	and be sure to check the debug log for more information.`

// ErrorCheck serves as a wrapper for HikariKnight/ls-iommu/pkg/common.ErrorCheck that allows for visibile error messages
func ErrorCheck(err error, msg ...string) {
	_, _ = os.Stdout.WriteString("\033[H\033[2J") // clear the screen
	oneMsg := ""
	if err != nil {
		if len(msg) < 1 {
			oneMsg = ""
		} else {
			for _, v := range msg {
				oneMsg += v + "\n"
			}
		}
		color.Printf("\n<red_b>FATAL</>: %s\n%s\nAborting", err.Error(), oneMsg)
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			print(".")
		}
		print("\n")
		errorcheck.ErrorCheck(err, msg...)
	}
}
