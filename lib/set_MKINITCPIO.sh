#!/bin/bash
function set_MKINITCPIO () {
    # Get the config paths
    source "$SCRIPTDIR/lib/paths.sh"

    # Grab the current modules but exclude vfio and vendor-reset
    CURRENTMODULES=$(grep -P "^MODULES" "$SCRIPTDIR/$MKINITCPIO" | perl -pe "s/MODULES=\((.+)\)/\1/")
    MODULES="$(grep -P "^MODULES" "$SCRIPTDIR/$MKINITCPIO" | perl -pe "s/MODULES=\((.+)\)/\1/" | perl -pe "s/\s?(vfio_iommu_type1|vfio_pci|vfio_virqfd|vfio|vendor-reset)\s?//g")"

    # Check if vendor-reset is present
    if [[ $CURRENTMODULES =~ "vendor-reset" ]];
    then
        # Inject vfio modules with vendor-reset
        perl -pi -e "s/MODULES=\(${CURRENTMODULES}\)/MODULES=\(vendor-reset vfio vfio_iommu_type1 vfio_pci vfio_virqfd ${MODULES}\)/" "$SCRIPTDIR/$MKINITCPIO"
    else
        # Inject vfio modules
        perl -pi -e "s/MODULES=\(${CURRENTMODULES}\)/MODULES=\(vfio vfio_iommu_type1 vfio_pci vfio_virqfd ${MODULES}\)/" "$SCRIPTDIR/$MKINITCPIO"
    fi

}

function main () {
    SCRIPTDIR=$(dirname "$(realpath "$0")" | perl -pe "s/\/\.\.\/lib//" | perl -pe "s/\/lib$//")

    set_MKINITCPIO "$1"
}

main "$1"