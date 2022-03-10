#!/bin/bash

function make_BACKUP () {
    local BACKUPDIR
    BACKUPDIR="$SCRIPTDIR/config/backup"

    if [ ! -d "$BACKUPDIR" ];
    then
        mkdir -p "$BACKUPDIR/etc/initramfs-tools"
        mkdir -p "$BACKUPDIR/etc/modprobe.d"

        sudo cp -v "/etc/modules" "$BACKUPDIR/etc/modules"
        sudo cp -v "/etc/initramfs-tools/modules" "$BACKUPDIR/etc/initramfs-rools/modules"

        if [ -f "/etc/modprobe.d/vfio.conf" ];
        then
            sudo cp -v "/etc/modprobe.d/vfio.conf" "$BACKUPDIR/etc/modprobe.d/vfio.conf"
        fi
        
        echo "Backup completed!"

    else
        echo "A backup already exists!
backup skipped."
    fi
}

function copy_FILES () {
    echo "Starting copying files to the system!"
    sudo cp -v "$SCRIPTDIR/$MODULES" "/etc/modules"
    sudo cp -v "$SCRIPTDIR/$INITRAMFS/modules" "/etc/initramfs-tools/modules"
    sudo cp -v "$SCRIPTDIR/$MODPROBE/vfio.conf" "/etc/modprobe.d/vfio.conf"

    echo ""
    echo "Rebuilding initramfs"
    sudo update-initramfs -u
}

function apply_CHANGES () {
	clear
    # Get the config paths
    source "$SCRIPTDIR/lib/paths.sh"

    echo "Configuration is now complete and these files have been generated for your system:
$SCRIPTDIR/$MODULES
$SCRIPTDIR/$INITRAMFS/modules
$SCRIPTDIR/$MODPROBE/vfio.conf

By proceeding, a backup of your system's version of these files will be placed in
$SCRIPTDIR/config/backup
unless a backup already exist.

Then the files above will be copied to your system followed by running \"update-initramfs -u\"
to build your new initrd image (all of this will require sudo permissions!)"
	
	read -p "Do you want to proceed with the installation of the files? [Y/n]: " YESNO

    case "${YESNO}" in
        [Yy]*)
            make_BACKUP
            copy_FILES
            exec "$SCRIPTDIR/lib/set_CMDLINE.sh"
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
