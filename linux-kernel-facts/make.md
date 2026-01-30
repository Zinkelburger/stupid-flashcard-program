make menuconfig: Opens a blue text-mode GUI to toggle kernel features. This is the most popular way to customize a build.

make defconfig: Generates a "default" config for your current architecture (x86_64).

make -j$(nproc): The main build command. -j$(nproc) tells it to use every CPU core you have, making it much faster.

make binrpm-pkg: Crucial for RHEL. This compiles the kernel and packages it into an .rpm file. This is the "proper" way to install a custom kernel on RHEL so dnf can track it

sudo make modules_install: Installs the compiled driver modules to /lib/modules/.

sudo make install: Installs the kernel image itself to /boot/ and updates your GRUB bootloader automatically.

make clean: Deletes most generated files but keeps your .config.

make mrproper: The "nuclear" option. Deletes everything, including your .config and all generated files.

you usually want to run cp /boot/config-$(uname -r) .config first to ensure your new kernel is based on the exact settings your system is currently running.
