#!/bin/bash

function get_GROUP () {
    clear
    # Get the config paths
    source "$SCRIPTDIR/lib/paths.sh"

    printf "
For this card to be passthrough-able, it must contain only:
* The GPU/Graphic card
* The GPU Audio Controller

Optionally it may also include:
* GPU USB Host Controller
* GPU Serial Port
* GPU USB Type-C UCSI Controller

"
    echo "#------------------------------------------#"
    exec "$SCRIPTDIR/utils/ls-iommu" | grep -i "group $1" | cut -d " " -f 1-4,8- | perl -pe "s/\[[0-9a-f]{4}\]: //"
    echo "#------------------------------------------#"

    printf "
To use any of these devices for passthrough ALL of them has to be passed through to the VMs\

To return to the previous page just press ENTER without typing in anything.
"
    read -p "Do you want to use these devices for passthrough? [y/N]: " YESNO

    case "${YESNO}" in
        [Yy]*)
            # Get the hardware ids from the selected group
            local GPU_DEVID=$($SCRIPTDIR/utils/ls-iommu | grep -i "group $1" | perl -pe "s/.+\[([0-9a-f]{4}:[0-9a-f]{4})\].+/\1/" | perl -pe "s/\n/,/" | perl -pe "s/,$/\n/")

            # Get the PCI ids
            local PCI_ID=$($SCRIPTDIR/utils/ls-iommu | grep -i "group $1" | cut -d " " -f 4 | perl -pe "s/\n/ /" | perl -pe "s/\s$//")

            # Write the GPU_PCI_IDs to the config that quickemu might make use of in the future
            printf "GPU_PCI_ID=($PCI_ID)
USB_CTL_ID=()
" > "$SCRIPTDIR/$QUICKEMU/qemu-vfio_vars.conf"

            # Get the PCI_ID
            local ROM_PCI_ID=$($SCRIPTDIR/utils/ls-iommu | grep -i "vga" | grep -i "group $1" | cut -d " " -f 4)

            # Get the GPU ROM
            "$SCRIPTDIR/lib/get_GPU_ROM.sh" "$ROM_PCI_ID"
            
            # Start setting up modules
            exec "$SCRIPTDIR/lib/set_MODULES.sh" $GPU_DEVID
        ;;
        *)
            exec "$SCRIPTDIR/lib/get_GPU.sh"
        ;;
    esac
}

function main () {
    SCRIPTDIR=$(dirname `which $0` | perl -pe "s/\/\.\.\/lib//")
    SCRIPTDIR="$SCRIPTDIR/.."
    
    get_GROUP $1
}

main $1