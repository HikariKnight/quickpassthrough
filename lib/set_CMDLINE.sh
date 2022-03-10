#!/bin/bash

function set_CMDLINE () {
    clear
    # Get the config paths
    source "${SCRIPTDIR}/lib/paths.sh"

    local CMDLINE
    CMDLINE=$(cat "${SCRIPTDIR}/config/kernel_args")

    echo "Configuration is now complete, however no changes have been done to your system.
The files needed have just been written to $SCRIPTDIR/config/etc

At this point if you know what you are doing, you can use these files to enable VFIO on your system
and achieve GPU Passthrough with your VMs by moving them to the correct location and then updating your
initramfs and bootloader.

For VFIO to work properly you need to make sure these kernel parameters are in your bootloader entry:
#-----------------------------------------------#
$CMDLINE
#-----------------------------------------------#

Make sure that the files inside \"$SCRIPTDIR/config/etc\" are copied to your /etc
AND PLEASE MAKE A BACKUP FIRST!

Then run \"sudo update-initramfs -u\", that way you can boot an older kernel without vfio if needed, before commiting fully.
You can remove the the vfio_pci kernel arguments from the linux line in your bootloader to disable/unbind the graphic card from the vfio driver on boot.

Finally reboot your system and run \"$SCRIPTDIR/vfio-verify\" to check if your GPU is properly set up.
NOTE: Some AMD GPUs will require the vendor-reset kernel module from https://github.com/gnif/vendor-reset to be installed!

The files inside \"$SCRIPTDIR/$QUICKEMU\" are currently unused files, however they provide
the required information that the QuickEMU project can hook into and use to add support for VFIO enabled VMs.

######################################################################
####  In the future, when I have enough confirmation that this script works for other people.
####  This page will get replaced with a prompt asking if you want to apply the changes and make backups
####  of your current system config.
######################################################################

"
}


function main () {
    SCRIPTDIR=$(dirname "$(which $0)" | perl -pe "s/\/\.\.\/lib//" | perl -pe "s/\/lib$//")
    set_CMDLINE
}

main
