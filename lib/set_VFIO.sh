#!/bin/bash
function set_VFIO () {
    # Get the config paths
    source "$SCRIPTDIR/lib/paths.sh"

    # Assign the GPU device ids to a variable
    GPU_DEVID="$1"

    # Get the kernel_args file content
    CMDLINE=$(cat "$SCRIPTDIR/config/kernel_args")

    # Ask if we shall disable video output on this card
    echo "
Disabling video output in Linux for the card you want to use in a VM
will make it easier to successfully do the passthrough without issues."
    read -p "Do you want to force disable video output in linux on this card? [Y/n]: " DISABLE_VGA
    case "${DISABLE_VGA}" in
    [Yy]*)
        # Update kernel_args file
        echo "${CMDLINE} vfio_pci.ids=${GPU_DEVID} vfio_pci.disable_vga=1" > "$SCRIPTDIR/config/kernel_args"

        # Update GPU_DEVID
        GPU_DEVID="$GPU_DEVID disable_vga=1"
    ;;
    [Nn]*)
        # Update kernel_args file
        echo "${CMDLINE} vfio_pci.ids=${GPU_DEVID}" > "$SCRIPTDIR/config/kernel_args"
    ;;
    *)
        # Update kernel_args file
        echo "${CMDLINE} vfio_pci.ids=${GPU_DEVID} vfio_pci.disable_vga=1" > "$SCRIPTDIR/config/kernel_args"

        # Update GPU_DEVID
        GPU_DEVID="$GPU_DEVID disable_vga=1"
    ;;
    esac
}

function main () {
    SCRIPTDIR=$(dirname "$(realpath $0)" | perl -pe "s/\/\.\.\/lib//" | perl -pe "s/\/lib$//")

    set_VFIO "$1"
}

main "$1"