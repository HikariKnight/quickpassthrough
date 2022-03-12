#!/bin/bash

function insert_INITRAMFSTOOLS() {
    # Get the header and enabled modules separately from the /etc/modules file
    local MODULES_HEADER
    local MODULES_ENABLED
    local VENDOR_RESET
    MODULES_HEADER=$(head -n $1 "$2" | grep -P "^#" | grep -v "# Added by quickpassthrough")
    MODULES_ENABLED=$(cat "$2" | grep -vP "^#" | grep -v "vendor-reset")
    VENDOR_RESET=0
    
    # If vendor-reset is present
    if grep -q "vendor-reset" "$2" ;
    then
        VENDOR_RESET=1
    fi

    # Write header
    echo "$MODULES_HEADER" > "$2"
    
    # If vendor-reset existed from before
    if [ $VENDOR_RESET == 1 ];
    then
        # Write vendor-reset as the first module!
        echo "vendor-reset" >> "$2"
    fi
    
    # Append vfio 
    printf "
# Added by quickpassthrough #
vfio
vfio_iommu_type1
vfio_pci
vfio_virqfd
#############################
" >> "$2"

    # Write the previously enabled modules under vfio in the load order
    echo "$MODULES_ENABLED" >> "$2"
}

function set_INITRAMFSTOOLS () {
    # Get the config paths
    source "$SCRIPTDIR/lib/paths.sh"
    
    # Insert modules in the correct locations as early as possible without
    # conflicting with vendor-reset module if it is enabled
    insert_INITRAMFSTOOLS 4 "$SCRIPTDIR/$ETCMODULES"
    insert_INITRAMFSTOOLS 11 "$SCRIPTDIR/$INITRAMFS/modules"

    "$SCRIPTDIR/lib/set_VFIO.sh" "$1"
    "$SCRIPTDIR/lib/set_MODPROBE.sh" "$1"

    exec "$SCRIPTDIR/lib/get_USB_CTL.sh"
}


function main () {
    SCRIPTDIR=$(dirname "$(which $0)" | perl -pe "s/\/\.\.\/lib//" | perl -pe "s/\/lib$//")

    set_INITRAMFSTOOLS "$1"
}

main "$1"
