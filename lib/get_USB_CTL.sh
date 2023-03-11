#!/bin/bash

function get_USB_CTL () {
    clear
    printf "THIS STEP IS OPTIONAL IF YOU DO NOT PLAN TO USE ANYTHING OTHER THAN MOUSE AND KEYBOARD!
The USB Controller you want to passthrough cannot be in a group with other devices.
Passing through a whole USB Controller (a set of hardwired 1-4 usb ports on the motherboard)
is only needed if you intend to use other devices than just mouse and keyboard with the VFIO enabled VM.

"
    echo "#------------------------------------------#"
    exec "$SCRIPTDIR/utils/ls-iommu" | grep -i "usb controller" | cut -d " " -f 1-5,9-
    echo "#------------------------------------------#"
    printf "
Press q to quit
"

    read -r -p "Which group number do you want to check?: " IOMMU_GROUP

   case "${IOMMU_GROUP}" in
       [1-9]*)
            exec "$SCRIPTDIR/lib/get_USB_CTL_GROUP.sh" "$IOMMU_GROUP"
        ;;
       [Qq]*)
            exec "$SCRIPTDIR/lib/apply_CHANGES.sh"
        ;;
       *)
            exec "$SCRIPTDIR/lib/apply_CHANGES.sh"
        ;;
   esac
}

function main () {
    SCRIPTDIR=$(dirname "$(realpath "$0")" | perl -pe "s/\/\.\.\/lib//" | perl -pe "s/\/lib$//")

    get_USB_CTL
}

main