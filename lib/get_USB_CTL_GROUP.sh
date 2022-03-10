#!/bin/bash

function get_USB_CTL_GROUP () {
    clear
    # Get the config paths
    source "$SCRIPTDIR/lib/paths.sh"

    printf "
For this USB controller device to be passthrough-able, it must be the ONLY device in this group!
Passing through more than just the USB controller can in some cases cause system issues
if you do not know what you are doing.

"
    echo "#------------------------------------------#"
    exec "$SCRIPTDIR/utils/ls-iommu" | grep -i "group $1" | cut -d " " -f 1-4,8- | perl -pe "s/\[[0-9a-f]{4}\]: //"
    echo "#------------------------------------------#"

    printf "
To use any of the devices shown for passthrough, all of them have to be passed through

To return to the previous page just press ENTER.
"
    read -p "Do you want to use the displayed devices for passthrough? [y/N]: " YESNO

    case "${YESNO}" in
        [Yy]*)
            # Get the PCI ids
            local PCI_ID
            PCI_ID=$($SCRIPTDIR/utils/ls-iommu | grep -i "group $1" | cut -d " " -f 4 | perl -pe "s/([0-9a-f]{2}:[0-9a-f]{2}.[0-9a-f]{1})\n/\"\1\" /" | perl -pe "s/\s$//")

            # Replace the blank USB_CTL_ID with the PCI_ID for the usb controller the user wants to pass through
            perl -pi -e "s/USB_CTL_ID=\(\)/USB_CTL_ID=\($PCI_ID\)/" "$SCRIPTDIR/$QUICKEMU/qemu-vfio_vars.conf"
            exec "$SCRIPTDIR/lib/apply_CHANGES.sh"
        ;;
        *)
            exec "$SCRIPTDIR/lib/get_USB_CTL.sh"
        ;;
    esac
}

function main () {
    SCRIPTDIR=$(dirname "$(which $0)" | perl -pe "s/\/\.\.\/lib//" | perl -pe "s/\/lib$//")

    get_USB_CTL_GROUP $1
}

main $1
