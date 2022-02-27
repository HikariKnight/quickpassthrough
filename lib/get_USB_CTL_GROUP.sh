#!/bin/bash

function get_USB_CTL_GROUP () {
    clear
    # Get the config paths
    source "$SCRIPTDIR/lib/paths.sh"

    printf "
For this USB controller device to be passthrough-able, it must be the ONLY device in this group!

"
    echo "#------------------------------------------#"
    exec "$SCRIPTDIR/utils/ls-iommu" | grep -i "group $1" | cut -d " " -f 1-4,8- | perl -pe "s/\[[0-9a-f]{4}\]: //"
    echo "#------------------------------------------#"
    
    printf "
To use this device for passthrough please type in the device id in the format (without brackets or quotes) --> \"xxxx:yyyy\"
NOTE: The device ID is the part inside the last [] brackets, example: [1002:aaf0]

To return to the previous page just press ENTER without typing in any ids
"
read -p "Enter the id for the device you want to passthrough: " USB_CTL_DEVID

if [[ $USB_CTL_DEVID =~ : ]];
then
    # Get the PCI ids
    PCI_ID=$($SCRIPTDIR/utils/ls-iommu | grep -i "group $1" | cut -d " " -f 4)
    
    exec perl -pi -e "s/USB_CTL_ID=\"\"/USB_CTL_ID=\"$PCI_ID\"/" "$SCRIPTDIR/$QUICKEMU/qemu-vfio_vars.conf"
else
    exec "$SCRIPTDIR/lib/get_USB_CTL.sh"
fi

}

function main () {
    SCRIPTDIR=$(dirname `which $0` | perl -pe "s/\/\.\.\/lib//")
    SCRIPTDIR="$SCRIPTDIR/.."
    
    get_USB_CTL_GROUP $1
}

main $1