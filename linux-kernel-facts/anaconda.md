Section 1: The Linux Boot Process (Pre-Installer & General)

    Flashcard: UEFI vs. BIOS Boot Handoff

        Concept: How the firmware hands execution to the bootloader.

        Details: In Legacy BIOS, the MBR (first 512 bytes) is executed, loading GRUB stage 1. In UEFI, the firmware reads the ESP (EFI System Partition) and directly executes an EFI binary (e.g., grubx64.efi or the kernel's own EFISTUB).

    Flashcard: The Kernel Command Line

        Concept: How the bootloader passes parameters to the kernel and initramfs.

        Details: The linux or linuxefi line in GRUB defines parameters. Parameters without a prefix are parsed by the kernel (e.g., console=ttyS0, panic=0). Parameters prefixed with rd. are parsed by Dracut (the initramfs generator). Parameters prefixed with inst. are parsed specifically by Anaconda.

    Flashcard: PID 1 & switch_root

        Concept: The transition from the temporary boot environment to the real OS.

        Details: The kernel unpacks the initramfs into rootfs (a tmpfs) and executes /init (usually a symlink to systemd). Once the real root filesystem (or the installer squashfs) is mounted, systemd uses switch_root to move the mount point to /, clear the initramfs from memory, and execute the real system's /sbin/init as the new PID 1.

Section 2: Dracut & The Initramfs

    Flashcard: What is Dracut?

        Concept: The event-driven initramfs infrastructure.

        Details: Dracut replaces the old mkinitrd. It relies heavily on udev to detect hardware and only loads the modules necessary to mount the root filesystem.

    Flashcard: Dracut Installer Hooks

        Concept: How Dracut fetches the Anaconda environment.

        Details: When booting an installer, Dracut uses the livenet or network modules to bring up the NICs (e.g., fetching DHCP based on the BOOTIF MAC address passed by PXE). It then fetches the stage2 image (usually a squashfs file named LiveOS/squashfs.img containing Anaconda), mounts it via loopback, and switch_roots into it.

    Flashcard: rd.break Debugging

        Concept: Interrupting the early boot process.

        Details: Passing rd.break to the kernel drops you into an emergency shell inside the initramfs before switch_root. You can use rd.break=pre-udev, pre-mount, mount, or pre-pivot to stop at specific stages. This is critical for debugging missing storage drivers or failed stage2 network downloads.

Section 3: Anaconda Installer Architecture

    Flashcard: Anaconda Framework

        Concept: The overall structure of the Red Hat/Fedora installer.

        Details: Anaconda is not a single script; it's a modular Python application running under systemd. It communicates via D-Bus. It has a UI tier (GUI via GTK, TUI via simpleline) and a headless payload tier.

    Flashcard: Storage Subsystem (Blivet)

        Concept: How Anaconda manages disks.

        Details: Anaconda relies on the Blivet Python library. Blivet builds a device tree in memory by querying udev and sysfs. It orchestrates partitioning, LVM, mdraid, and LUKS encryption. If a storage controller driver is missing from the initramfs, Blivet will never see the block device.

    Flashcard: Anaconda Boot Arguments (inst.*)

        Concept: Directing the installer's behavior via the kernel cmdline.

        Details: * inst.stage2=: URL/path to the installer squashfs image.

            inst.ks=: URL/path to the kickstart file.

            inst.repo=: The package payload repository.

            (Note: As of modern RHEL versions, arguments intended for Anaconda must have the inst. prefix. For example, inst.selinux=0 applies to the installer environment, while selinux=0 applies to the installed system's kernel).

Section 4: Kickstart Execution & Environments

    Flashcard: Kickstart Parsing

        Concept: How the automated installation file is processed.

        Details: Kickstart is parsed by the pykickstart library. It validates the syntax before Anaconda executes the steps. If there is a syntax error, Anaconda usually halts early.

    Flashcard: %pre Script Environment

        Concept: Execution state before the installation begins.

        Details: %pre scripts run after the kickstart is downloaded but before Anaconda evaluates the storage layout. They run in the volatile installer RAM environment. They are commonly used to dynamically generate storage configurations (e.g., dynamically writing partitioning lines to an /tmp/storage.ks file and %include-ing it).

    Flashcard: %post (Chroot vs. No-Chroot)

        Concept: Execution state after packages are installed.

        Details: * %post (default): Runs chrooted inside the newly installed system located at /mnt/sysroot. The kernel is still the installer's kernel, but the binaries and libraries are the newly installed OS.

            %post --nochroot: Runs in the installer's RAM environment. The new system filesystem is mounted at /mnt/sysroot. This is highly useful for copying logs (/tmp/anaconda.log) from the volatile installer space into the persistent disk (/mnt/sysroot/var/log/).

Section 5: Logs and Troubleshooting

    Flashcard: Anaconda Log Locations

        Concept: Where to look when things crash.

        Details: During installation, logs are kept in a tmpfs at /tmp.

            /tmp/anaconda.log: Main installer logic and UI routing.

            /tmp/storage.log: Blivet disk discovery and partitioning (critical for block-layer debugging).

            /tmp/program.log: Output of external commands (like %pre/%post scripts or subprocess calls).

            /tmp/packaging.log: DNF/RPM transaction details.

        After a successful install, these are copied to /var/log/anaconda/ on the new system.

    Flashcard: tmux Installer Shells

        Concept: Accessing the console during a run.

        Details: Anaconda runs inside tmux on the virtual consoles. If the system hangs, you can switch to TTY2 (Ctrl+Alt+F2). The first pane is the main shell; you can cycle panes (Ctrl+b, then 2, 3, etc.) to view live tails of the storage.log or program.log while the UI is frozen on TTY1/TTY6.