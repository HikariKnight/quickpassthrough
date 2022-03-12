#!/bin/bash

function make_BACKUP () {
    local BACKUPDIR
    BACKUPDIR="$SCRIPTDIR/backup"

    if [ ! -d "$BACKUPDIR" ];
    then
        # Make the backup directories and backup the files  
        if [ -d "/etc/initramfs-tools" ];
        then
            mkdir -p "$BACKUPDIR/etc/initramfs-tools"
            cp -v "/etc/initramfs-tools/modules" "$BACKUPDIR/etc/initramfs-tools/modules"
            cp -v "/etc/modules" "$BACKUPDIR/etc/modules"

        elif [ -d "/etc/dracut.conf" ];
        then
            mkdir -p "$BACKUPDIR/etc/dracut.conf.d"
            if [ -f "/etc/dracut.conf.d/10-vfio.conf" ];
            then
                cp -v "/etc/dracut.conf.d/10-vfio.conf" "$BACKUPDIR/etc/dracut.conf.d/10-vfio.conf"

            fi

        elif [ -f "/etc/mkinitcpio.conf" ];
        then
            mkdir -p "$BACKUPDIR/etc"
            cp -v "/etc/mkinitcpio.conf" "$BACKUPDIR/etc/mkinitcpio.conf"

        fi

        if [ -f "/etc/default/grub" ];
        then
            mkdir -p "$BACKUPDIR/etc/default"
            cp -v "/etc/default/grub" "$BACKUPDIR/etc/default/grub"

        fi

        mkdir -p "$BACKUPDIR/etc/modprobe.d"

        # If a vfio.conf file exists, backup that too
        if [ -f "/etc/modprobe.d/vfio.conf" ];
        then
            cp -v "/etc/modprobe.d/vfio.conf" "$BACKUPDIR/etc/modprobe.d/vfio.conf"

        fi
        
        printf "Backup completed!\n"

    else
        echo "
A backup already exists!
backup skipped.
"
    fi
}

function copy_FILES () {
    echo "Starting copying files to the system!"

    sudo cp -v "$SCRIPTDIR/$MODPROBE/vfio.conf" "/etc/modprobe.d/vfio.conf"

    if [ -d "/etc/initramfs-tools" ];
    then
        sudo cp -v "$SCRIPTDIR/$ETCMODULES" "/etc/modules"
        sudo cp -v "$SCRIPTDIR/$INITRAMFS/modules" "/etc/initramfs-tools/modules"
        echo "
Rebuilding initramfs"
        sudo update-initramfs -u

    elif [ -d "/etc/dracut.conf" ];
    then
        cp -v "$SCRIPTDIR/$DRACUT/10-vfio.conf" "/etc/dracut.conf.d/10-vfio.conf"
        echo "
Rebuilding initramfs"
        sudo dracut -f -kver "$(uname -r)"

    else
        echo "
Unsupported initramfs infrastructure
In order to make vfio work, please add these modules to your
initramfs and make them load early, then rebuild initramfs.

vfio
vfio_iommu_type1
vfio_pci
vfio_virqfd
"
        
    fi

    
}

function apply_CHANGES () {
	clear
    # Get the config paths
    source "$SCRIPTDIR/lib/paths.sh"

    echo "Configuration is now complete and these files have been generated for your system:
$SCRIPTDIR/$ETCMODULES
$SCRIPTDIR/$INITRAMFS/modules
$SCRIPTDIR/$MODPROBE/vfio.conf

By proceeding, a backup of your system's version of these files will be placed in
$SCRIPTDIR/backup
unless a backup already exist.

Then the files above will be copied to your system followed by running followed by updating your
initramfs and then attempt adding new kernel arguments to your bootloader."
	
	read -p "Do you want to proceed with the installation of the files? (no=quit) [Y/n]: " YESNO

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
    SCRIPTDIR=$(dirname "$(realpath "$0")" | perl -pe "s/\/\.\.\/lib//" | perl -pe "s/\/lib$//")

    apply_CHANGES
}

main
