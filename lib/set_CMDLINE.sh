#!/bin/bash

function set_KERNELSTUB () {
    clear
    # Tell what we are going to do
    echo "Adding vfio kernel arguments to systemd-boot using kernelstub"

    # Get the config paths
    source "${SCRIPTDIR}/lib/paths.sh"

    # Check if systemd-boot already has vfio parameters from before
    KERNELSTUB_TEST=$(sudo kernelstub -p 2>&1 | grep "Kernel Boot Options" | perl -pe "s/.+Kernel Boot Options:\..+(vfio_pci.ids=.+ ).+/\1/")
    
    # If there are already vfio_pci parameters in kernelstub
    if [[ "$KERNELSTUB_TEST" =~ vfio_pci.ids ]] ;
    then
        # Remove the old parameters
        sudo kernelstub -d "$KERNELSTUB_TEST"
        sudo kernelstub -d "vfio_pci.disable_vga=1"
        sudo kernelstub -d "vfio_pci.disable_vga=0"
    fi

    # Apply new parameters
    CMDLINE=$(cat "${SCRIPTDIR}/config/kernel_args")
    sudo kernelstub -a "$CMDLINE"

    show_FINISH
}

function show_FINISH () {
    clear
    # Get the config paths
    source "${SCRIPTDIR}/lib/paths.sh"

    local CMDLINE
    CMDLINE=$(cat "${SCRIPTDIR}/config/kernel_args")

    echo "Configuration is now complete!"

    if [ $1 == 0 ];
    then
        printf "For VFIO to work properly you need to make sure these kernel parameters are in your bootloader entry:
#-----------------------------------------------#
%s
#-----------------------------------------------#

" "$CMDLINE"
    fi

    echo "Restart your system and run 
\"$SCRIPTDIR/vfio-verify\"
to check if your GPU is properly set up.

If the graphic card is bound to vfio-pci then you can
proceed to add it to your virtual machines.

A backup the files we replaced on your system can be found inside
$SCRIPTDIR/backup/
In order to restore these files just copy them back to your system and run
\"sudo update-initramfs -u\"

You can remove the the vfio_pci kernel arguments from the linux line in your bootloader
to disable/unbind the graphic card from the vfio driver on boot.

The files inside \"$SCRIPTDIR/$QUICKEMU\" are currently unused files, however they provide
the required information that the QuickEMU project can hook into and use to add support for VFIO enabled VMs.

The PCI Devices with these IDs are what you should add to your VMs:
NOTE: Some AMD GPUs will require the vendor-reset kernel module from https://github.com/gnif/vendor-reset to be installed!"

    source "$SCRIPTDIR/config/quickemu/qemu-vfio_vars.conf"

    for dev in "${GPU_PCI_ID[@]}"
    do
        echo "* $dev"
    done
    for dev in "${USB_CTL_ID[@]}"
    do
        echo "* $dev"
    done
}

function set_CMDLINE () {
    # Make a variable to tell if 
    local BOOTLOADER_AUTOCONFIG
    BOOTLOADER_AUTOCONFIG=0
    
    # If kernelstub is detected (program to manage systemd-boot)
    if which kernelstub ;
    then
        # Configure kernelstub then exit
        set_KERNELSTUB
        BOOTLOADER_AUTOCONFIG=1
    fi

    show_FINISH $BOOTLOADER_AUTOCONFIG
}


function main () {
    SCRIPTDIR=$(dirname "$(which $0)" | perl -pe "s/\/\.\.\/lib//" | perl -pe "s/\/lib$//")
    set_CMDLINE
}

main
