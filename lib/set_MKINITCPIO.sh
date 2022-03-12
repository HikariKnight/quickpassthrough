#!/bin/bash
function set_MKINITCPIO () {
    # Get the config paths
    source "$SCRIPTDIR/lib/paths.sh"

    # Grab the current modules but exclude vfio and vendor-reset
    CURRENTMODULES=$(grep -P "^MODULES" "$SCRIPTDIR/$MKINITCPIO" | perl -pe "s/MODULES=\((.+)\)/\1/")
    MODULES="$(grep -P "^MODULES" "$SCRIPTDIR/$MKINITCPIO" | perl -pe "s/MODULES=\((.+)\)/\1/" | perl -pe "s/\s?(vfio_iommu_type1|vfio_pci|vfio_virqfd|vfio|vendor-reset)\s?//g")"

    if [[ $CURRENTMODULES =~ "vendor-reset" ]];
    then
        perl -pi -e "s/MODULES=\(${CURRENTMODULES}\)/MODULES=\(vendor-reset vfio vfio_iommu_type1 vfio_pci vfio_virqfd ${MODULES}\)/" "$SCRIPTDIR/$MKINITCPIO"
    else
        perl -pi -e "s/MODULES=\(${CURRENTMODULES}\)/MODULES=\(vfio vfio_iommu_type1 vfio_pci vfio_virqfd ${MODULES}\)/" "$SCRIPTDIR/$MKINITCPIO"
    fi

}

function main () {
    SCRIPTDIR=$(dirname "$(realpath "$0")" | perl -pe "s/\/\.\.\/lib//" | perl -pe "s/\/lib$//")

    set_MKINITCPIO "$1"
}

main "$1"