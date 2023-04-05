#!/bin/bash
# shellcheck disable=SC2034
MODPROBE="config/etc/modprobe.d"
INITRAMFS="config/etc/initramfs-tools"
ETCMODULES="config/etc/modules"
DEFAULT="config/etc/default"
QUICKEMU="config/quickemu"
DRACUT="config/etc/dracut.conf.d"
MKINITCPIO="config/etc/mkinitcpio.conf"

READAPI="wget -O-"
DOWNLOAD="wget -0 \"$SCRIPTDIR/utils/ls-iommu.tar.gz\""
# Get the tool to use for downloading
if [ -f "/usr/bin/curl" ];
then
    READAPI="curl"
    DOWNLOAD="curl -JLo ls-iommu.tar.gz"
fi
