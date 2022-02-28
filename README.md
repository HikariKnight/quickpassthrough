# VFIO-enabler (name pending)
A project to simplify setting up GPU passthrough for [QuickEMU](https://github.com/quickemu-project/quickemu) and libvirt

Currently this project does NO MODIFICATIONS to your system, all it does is generate the files needed for testing inside `./config/`
In a future version it will ask if you want changes applied to your system, however I am not enabling that until I have confirmation that this generates a working configuration for other systems than my own.

You can use it by simply running
```bash
git clone https://github.com/HikariKnight/VFIO-enabler.git
cd VFIO-enabler
chmod +x ./vfio-setup
./vfio-setup
```

## Features
* General warning and info about what you will be needing
* Enable and configure vfio modules
* Configure 2nd GPU for GPU Passthrough
* Dump the selected GPU rom (as some cards require a romfile for passthrough to work), however no rom patching support planned.
* Enable and configure the correct kernel modules
* Provides you with the correct kernel arguments to add to your bootloader entry

## Contributing
<img src="https://user-images.githubusercontent.com/2557889/156038229-4e70352f-9182-4474-8e32-d14d3ad67566.png" width="250px"><br>
I know my bash skills are not great, so help is always welcome! And help is wanted here.<br>
If you know bash well, you will be able to help! Just make a pull request with your changes!<br>
Just remember to add comments to document the work and explain it for people who are less familiar with the bash syntax or anything else you use. 😄<br>
<br>
Also if you know english, you can just help proof reading. English is not my native language, plus I have dyslexia so I often make spelling mistakes.<br>
Proof reading is still contribution!


## TODO
* ~~Everything~~
* ~~Show general warning to user and inform about making a backup and general expectations~~
* ~~Detect if user has an amd or intel CPU and provide the correct IOMMU kernel args based on that~~
* ~~Tell user to enable IOMMU (VT-d/AMD-v) on their motherboard and bootloader~~
* ~~Integrate ls-iommu and locate graphic cards and see what IOMMU group they are in~~
* ~~Enable and configure vfio modules~~
* ~~Fetch the ID for the GPUs and generate the correct kernel arguments for grub and systemd-boot~~
* ~~Dump the GPU rom, just in case it will be needed for passthrough~~ (no rom patching planned due to complexity)
* Get help to actually make the scripts better
* A non hacky menu system? (I will need help by some bash wizards for this)
* Colored highlight/text for important information?
* Blacklist drivers? (some cards require blacklisting as softdep is not enough)
* Install vendor_reset kernel module? (maybe far future)
* Setup looking-glass? (far future maybe)

<br>

### Why bash?
I wanted the dependencies to be minimal without the need for compilation and not have a potential breaking change in the future (like with the transition from python2 to python3).

I know enough bash to make things work, but I am in no way a professional in writing bash scripts as I usually write python and golang.

There is also quite a lot of perl usage as I am quite familiar with the perl regex format over something like sed
