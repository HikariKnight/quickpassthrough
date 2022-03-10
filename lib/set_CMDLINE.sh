#!/bin/bash

# Function to configure systemd-boot using kernelstub
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
}

# Function to configure grub
function set_GRUB () {
    clear
    # Get the config paths
    source "$SCRIPTDIR/lib/paths.sh"

    local CMDLINE
    CMDLINE=$(cat "$SCRIPTDIR/config/kernel_args")

    # HIGHLY EXPERIMENTAL!
    GRUB_CMDLINE=$(cat "/etc/default/grub" | grep -P "^GRUB_CMDLINE_LINUX_DEFAULT" | perl -pe "s/GRUB_CMDLINE_LINUX_DEFAULT=\"(.+)\"/\1/" | perl -pe "s/iommu=(pt|on)|amd_iommu=on|vfio_pci.ids=.+|vfio_pci.disable_vga=\d{1}//g" | perl -pe "s/(^\s+|\s+$)//g")
    GRUB_CMDLINE_LINUX_DEFAULT=$(cat "/etc/default/grub" | grep -P "^GRUB_CMDLINE_LINUX_DEFAULT")

    # Update the GRUB_CMDLINE_LINUX_DEFAULT line
    perl -pi -e "s/${GRUB_CMDLINE_LINUX_DEFAULT}/GRUB_CMDLINE_LINUX_DEFAULT=\"${GRUB_CMDLINE} ${CMDLINE}\"/" "${SCRIPTDIR}/config/etc/default/grub"

    echo "The script will now replace your default grub file with a new one.
Then attempt to update grub and generate a new grub.cfg.
If generating the grub.cfg file fails, you can find a backup of your grub default file here:
$SCRIPTDIR/backup/etc/default/grub
"
    read -r -p "Press ENTER to continue"

    sudo cp -v "$SCRIPTDIR/config/etc/default/grub" "/etc/default/grub"
    sudo grub-mkconfig -o "/boot/grub/grub.cfg"

    echo ""
    read -r -p "Please verify there was no errors generating the grub.cfg file, then press ENTER"    
}

function show_FINISH () {
    clear
    # Get the config paths
    source "$SCRIPTDIR/lib/paths.sh"

    local CMDLINE
    CMDLINE=$(cat "$SCRIPTDIR/config/kernel_args")

    echo "Configuration is now complete!"

    if [ "$1" == 0 ];
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

    source "${SCRIPTDIR}/config/quickemu/qemu-vfio_vars.conf"

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
        # Configure kernelstub
        set_KERNELSTUB
        BOOTLOADER_AUTOCONFIG=1
    fi

    # If grub exists
    if which grub-mkconfig ;
    then
        # Configure grub
        set_GRUB
        BOOTLOADER_AUTOCONFIG=1
    fi

    show_FINISH $BOOTLOADER_AUTOCONFIG
}


function main () {
    SCRIPTDIR=$(dirname "$(which $0)" | perl -pe "s/\/\.\.\/lib//" | perl -pe "s/\/lib$//")
    set_CMDLINE
}

main
