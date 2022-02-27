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
To use these devices for passthrough please type in ALL their device ids in the format (without brackets or quotes) --> \"xxxx:yyyy,xxxx:yyyy\"
NOTE: The device ID is the part inside the last [] brackets, example: [1002:aaf0]

To return to the previous page just press ENTER without typing in any ids
"
read -p "Enter the ids for all devices you want to passthrough: " GPU_DEVID

if [[ $GPU_DEVID =~ : ]];
then
    # Get the PCI ids
    local PCI_ID=$($SCRIPTDIR/utils/ls-iommu | grep -i "group $1" | cut -d " " -f 4 | perl -pe "s/\n/ /" | perl -pe "s/\s$//")

    # Write the currently unused vfio config for quickemu
    printf "GPU_PCI_ID=($PCI_ID)
USB_CTL_ID=\"\"" > "$SCRIPTDIR/$QUICKEMU/qemu-vfio_vars.conf"

    # Get the PCI_ID
    local ROM_PCI_ID=$($SCRIPTDIR/utils/ls-iommu | grep -i "vga" | grep -i "group $1" | cut -d " " -f 4)
    "$SCRIPTDIR/lib/get_GPU_ROM.sh" "$ROM_PCI_ID"

    exec "$SCRIPTDIR/lib/set_MODULES.sh" $GPU_DEVID
else
    exec "$SCRIPTDIR/lib/get_GPU.sh"
fi

}

function main () {
    SCRIPTDIR=$(dirname `which $0` | perl -pe "s/\/\.\.\/lib//")
    SCRIPTDIR="$SCRIPTDIR/.."
    
    get_GROUP $1
}

main $1