#!/bin/bash

function set_CMDLINE () {
    clear
    # Get the config paths
    source "$SCRIPTDIR/lib/paths.sh"

    CMDLINE=$(cat "$SCRIPTDIR/config/kernel_args")

    printf "Configuration is now complete, however no changes have been done to your system.
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
    SCRIPTDIR=$(dirname `which $0` | perl -pe "s/\/\.\.\/lib//" | perl -pe "s/\/lib$//")
    set_CMDLINE
}

main