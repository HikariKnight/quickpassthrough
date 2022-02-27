# VFIO-enabler
A project to simplify setting up GPU passthrough for QuickEMU and libvirt

Currently this project does NO MODIFICATIONS to your system, all it does is generate the files needed for testing inside `./config/

## Features
* General warning and info about what you will be needing
* Enable and configure vfio modules
* Configure 2nd GPU for GPU Passthrough
* Dump the selected GPU rom (as some cards require a romfile for passthrough to work), however no rom patching support planned.
* Enable and configure the correct kernel modules

## TODO
* ~~Everything~~
* ~~Show general warning to user and inform about making a backup and general expectations~~
* ~~Detect if user has an amd or intel CPU and provide the correct IOMMU kernel args based on that~~
* ~~Tell user to enable IOMMU (VT-d/AMD-v) on their motherboard and bootloader~~
* ~~Integrate ls-iommu and locate graphic cards and see what IOMMU group they are in~~
~~* Enable and configure vfio modules~~
* Fetch the ID for the GPUs and generate the correct kernel arguments for grub and systemd-boot
* ~~Dump the GPU rom, just in case it will be needed for passthrough~~ (no rom patching planned due to complexity)
* A non hacky menu system? (I will need help by some bash wizards for this)
* Blacklist drivers? (some cards require blacklisting as softdep is not enough)
* Install vendor_reset kernel module? (maybe far future)
* Setup looking-glass? (far future maybe)