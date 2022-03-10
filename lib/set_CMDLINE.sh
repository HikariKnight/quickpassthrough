#!/bin/bash

function set_CMDLINE () {
    clear
    # Get the config paths
    source "${SCRIPTDIR}/lib/paths.sh"

    local CMDLINE
    CMDLINE=$(cat "${SCRIPTDIR}/config/kernel_args")

    echo "Configuration is now complete!

For VFIO to work properly you need to make sure these kernel parameters are in your bootloader entry:
#-----------------------------------------------#
$CMDLINE
#-----------------------------------------------#

A backup the files we replaced on your system can be found inside
$SCRIPTDIR/backup/
In order to restore these files just copy them back to your system and run
\"sudo update-initramfs -u\"

You can remove the the vfio_pci kernel arguments from the linux line in your bootloader to disable/unbind the graphic card from the vfio driver on boot.

Finally reboot your system and run \"$SCRIPTDIR/vfio-verify\" to check if your GPU is properly set up.
NOTE: Some AMD GPUs will require the vendor-reset kernel module from https://github.com/gnif/vendor-reset to be installed!

The files inside \"$SCRIPTDIR/$QUICKEMU\" are currently unused files, however they provide
the required information that the QuickEMU project can hook into and use to add support for VFIO enabled VMs.
"
}


function main () {
    SCRIPTDIR=$(dirname "$(which $0)" | perl -pe "s/\/\.\.\/lib//" | perl -pe "s/\/lib$//")
    set_CMDLINE
}

main
