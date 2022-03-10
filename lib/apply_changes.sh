#!/bin/bash

function apply_CHANGES () {
	clear
    # Get the config paths
    source "$SCRIPTDIR/lib/paths.sh"
	
	read -p "Do you want to proceed with the installation of the files? [y/N]: " YESNO

    case "${YESNO}" in
        [Yy]*)
            echo ""
        ;;
        *)
            exit 1
        ;;
    esac
}


function main () {
    SCRIPTDIR=$(dirname "$(which $0)" | perl -pe "s/\/\.\.\/lib//" | perl -pe "s/\/lib$//")
    apply_CHANGES
}

main
