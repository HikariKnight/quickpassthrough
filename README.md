# QuickPassthrough

A project to simplify setting up GPU passthrough on your Linux host for libvirt/virt-manager
![Quickpassthrough Preview GIF](https://raw.githubusercontent.com/HikariKnight/quickpassthrough/main/quickpassthrough_preview.gif)

You can use it by simply downloading the latest [release](https://github.com/HikariKnight/quickpassthrough/releases/) and run it inside a terminal or by downloading and compiling it yourself with the commands below.

This project is aimed at **systems with 2 GPUs** and **headless servers**, where the only GPU is not needed.

**Note:** Quickpassthrough is not designed to be installed by a package manager! As you would usually have to only run it once, unless you change the GPU.

## Features
* Show general warning to user and inform about making a backup and general expectations
* Detect if user has an amd or intel CPU and provide the correct IOMMU kernel args based on that
* Configure your 2nd GPU for GPU Passthrough (1 for host, 1 for VM)
* Use [ls-iommu](https://github.com/HikariKnight/ls-iommu) to find PCI devices like graphic cards, usb controllers, etc and see what IOMMU group they are in
* Enable and configure vfio modules (initramfs-tools, dracut, modprobe and mkinitcpio)
* Generate the correct kernel arguments for grub and systemd-boot
* Generate a script to use for dumping the vbios rom (as some cards require a romfile for passthrough to work), however no rom patching support planned.
* Configure kernel arguments for systemd-boot (using kernelstub or grubby)
* Configure kernel arguments for grub2 (editing /etc/default/grub or using grubby)
* A menu system you can navigate through.
* Make sure [vendor-reset](https://github.com/gnif/vendor-reset) module is loaded before vfio, check the repository for the list of cards that require it!
* Provides you with the correct kernel arguments to add to your bootloader entry if a supported bootloader is not found

## Features handled by [ls-iommu](https://github.com/HikariKnight/ls-iommu)
* Automatically handle GPUs where parts of it might be in separate IOMMU groups (ex: RX6600XT)
* Fetch the ID and PCI Address of devices
* Locate the vbios rom path on the system
* Tell the user to enable IOMMU (VT-d/AMD-v) on their motherboard and bootloader
* Get a list of devices, their IOMMU groups and various other information

## Build with current dependencies
```bash
git clone https://github.com/HikariKnight/quickpassthrough.git
cd quickpassthrough
go mod download
CGO_ENABLED=0 go build -ldflags="-X github.com/HikariKnight/quickpassthrough/internal/version.Version=$(git rev-parse --short HEAD)" -o quickpassthrough cmd/main.go
```

## Build with newest dependencies (may break)
```bash
git clone https://github.com/HikariKnight/quickpassthrough.git
cd quickpassthrough
go mod download
go get -u ./cmd
CGO_ENABLED=0 go build -ldflags="-X github.com/HikariKnight/quickpassthrough/internal/version.Version=$(git rev-parse --short HEAD)" -o quickpassthrough cmd/main.go
```

## Does this work on atomic or immutable systems?
Currently no, however [Bazzite](https://bazzite.gg), [Bluefin](https://projectbluefin.io) and [Aurora](https://getaurora.dev) has an `ujust` command that does a very similar job.

In Bazzite you run `ujust setup-virtualization` and follow the prompts to `Enable Virtualization` and `Enable VFIO drivers`.<br>
In Bluefin and Aurora you run `ujust setup-vfio` and follow the prompts.

## How do I undo the changes?
There is a `backup/` folder generated on the first run that will have a copy of all your files (and their paths) from before we edited anything.
Compare that folder with the `config/` folder to see which files you need to delete in addition to copying the files from `backup/` to your system before rebuilding your initramfs and updating your bootloader config.

## How do I just disable vfio for 1 boot?
Remove the vfio kernel arguments from your bootloader by pressing E on the boot menu. The kernel arguments added to the bootloader can be found in the `config/kernel_args` file. <br>
NOTE: You can also just remove them from your bootloader permanently and update your bootloader if you want to keep the config files on your system.

## What this project does NOT do
* Setup or configure your Virtual Machine (that is your job)
* Optimize your Virtual Machine for Passthrough (again this is your job)
* Optimize your host machine for Passthrough or Virtualization (out of this projects scope)
* Setup and configure GPU Passthrough on systems with 1 graphic card (iGPU counts as 1 Graphic Card by itself, so iGPU with another GPU will work)
* Does not configure passthrough of 3D controllers, [as it will not work](https://lantian.pub/en/article/modify-computer/laptop-intel-nvidia-optimus-passthrough.lantian/) (this is most gaming laptops so do not even think about it). If you try run this on a laptop with a 3D controller, the "2nd GPU" will not show up.

## Contributing
<img src="https://user-images.githubusercontent.com/2557889/156038229-4e70352f-9182-4474-8e32-d14d3ad67566.png" width="250px">

This project originally started out as a bash only project, upon completing the proof of concept it became very clear that bash would become very messy with all the weird quirks and regex and inline editing of files. <br>
So the project moved over to golang, this lets us utilize TUI toolkits to build a proper menu system for the project. <br>

If you know golang, passthrough or qemu, you are welcome to help! Just make a pull request!<br>
Just remember to add comments to document the work and explain it for people who are less familiar with the golang syntax or anything else you use. 😄

Also if you know English, you can help by just proof reading. English is not my native language, plus I have dyslexia so I often make spelling mistakes.
Proof reading is still contribution!

----

### Why GO?

I wanted to learn it, while also using a language that would potentially not create system dependencies. GO fits this criteria when you compile using CGO_ENABLED=0 as this will statically link the libraries and still produce a fairly small (when compressed with upx) binary.
