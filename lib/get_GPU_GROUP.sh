#!/bin/bash
# shellcheck disable=SC1091

function get_GPU_GROUP () {
    clear
    # Get the config paths
    source "$SCRIPTDIR/lib/paths.sh"

    printf "For this card to be passthrough-able, it must contain only:
* The GPU/Graphic card
* The GPU Audio Controller

Optionally it may also include:
* GPU USB Host Controller
* GPU Serial Port
* GPU USB Type-C UCSI Controller
* PCI Bridge (if they are in their own IOMMU groups)

"
    echo "#------------------------------------------#"
    exec "$SCRIPTDIR/utils/ls-iommu" -i "$1" -r | cut -d " " -f 1-5,6- | perl -pe "s/\[[0-9a-f]{4}\]: //"
    echo "#------------------------------------------#"

    printf "
To use any of these devices for passthrough ALL of them has to be passed through to the VMs\

To return to the previous page just press ENTER without typing in anything.
"
    read -r -p "Do you want to use these devices for passthrough? [y/N]: " YESNO

    case "${YESNO}" in
        [Yy]*)
            # Get the hardware ids from the selected group
            local GPU_DEVID
            GPU_DEVID=$("$SCRIPTDIR/utils/ls-iommu" -i "$1" -r --id | perl -pe "s/\n/,/" | perl -pe "s/,$/\n/")

            # Get the PCI ids
            local PCI_ID
            PCI_ID=$("$SCRIPTDIR/utils/ls-iommu" -i "$1" -r --pciaddr | perl -pe "s/([0-9a-f]{2}:[0-9a-f]{2}.[0-9a-f]{1})\n/\"\1\" /" | perl -pe "s/\s$//")

            # Write the GPU_PCI_IDs to the config that quickemu might make use of in the future
            echo "GPU_PCI_ID=($PCI_ID)
USB_CTL_ID=()" > "$SCRIPTDIR/$QUICKEMU/qemu-vfio_vars.conf"

            # Get the GPU ROM
            "$SCRIPTDIR/lib/get_GPU_ROM.sh" "$1"

            # Start setting up modules
            if [ -d "/etc/initramfs-tools" ];
            then
                exec "$SCRIPTDIR/lib/set_INITRAMFSTOOLS.sh" "$GPU_DEVID"
            
            elif [ -d "/etc/dracut.conf" ];
            then
                exec "$SCRIPTDIR/lib/set_DRACUT.sh" "$GPU_DEVID"
            
            elif [ -f "/etc/mkinitcpio.conf" ];
            then
                exec "$SCRIPTDIR/lib/set_MKINITCPIO.sh" "$GPU_DEVID"
            else
                # Bind GPU to VFIO
                "$SCRIPTDIR/lib/set_VFIO.sh" "$GPU_DEVID"

                # Configure modprobe
                "$SCRIPTDIR/lib/set_MODPROBE.sh" "$GPU_DEVID"
            fi
        ;;
        *)
            exec "$SCRIPTDIR/lib/get_GPU.sh"
        ;;
    esac
}

function main () {
    SCRIPTDIR=$(dirname "$(realpath "$0")" | perl -pe "s/\/\.\.\/lib//" | perl -pe "s/\/lib$//")

    get_GPU_GROUP "$1"
}

main "$1"
