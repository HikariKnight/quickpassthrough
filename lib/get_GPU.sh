#!/bin/bash

function get_GPU () {
    clear
    printf "These are your graphic cards, they have to be in separate groups.
The graphic card you want to passthrough cannot be in a group with other devices that
does not belong to itself. Both cards must also have unique hardware ids [xxxx:yyyy]!:

"
    echo "#------------------------------------------#"
    exec "$SCRIPTDIR/utils/ls-iommu" | grep -i "vga" | cut -d " " -f 1-4,9-
    echo "#------------------------------------------#"

    printf "
Press q to quit
"

    read -p "Which group number do you want to check?: " IOMMU_GROUP

   case "${IOMMU_GROUP}" in
       [1-9]*)
            exec "$SCRIPTDIR/lib/get_GPU_GROUP.sh" $IOMMU_GROUP
        ;;
       [Qq]*)
            echo ""
        ;;
       *)
            echo ""
        ;;
   esac
}

function main () {
    SCRIPTDIR=$(dirname `which $0` | perl -pe "s/\/\.\.\/lib//")
    SCRIPTDIR="$SCRIPTDIR/.."
    
    get_GPU
}

main