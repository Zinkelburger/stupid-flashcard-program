/proc (procfs) and /sys (sysfs) — Virtual Filesystems

These are pseudo-filesystems that don't store data on disk. They are windows into the running kernel, generated on-the-fly when you read them. They are the primary interface for inspecting and tuning the kernel from userspace.

---

/proc (procfs):

Mounted at /proc. Primarily exposes per-process information and system-wide kernel statistics.

Per-process files (/proc/<pid>/):

    /proc/<pid>/status     — Human-readable process info (name, state, memory, threads)
    /proc/<pid>/cmdline    — The command line that started the process (null-separated)
    /proc/<pid>/environ    — Environment variables (null-separated)
    /proc/<pid>/fd/        — Directory of symlinks to open file descriptors
    /proc/<pid>/maps       — Memory mappings (virtual address ranges, permissions, files)
    /proc/<pid>/smaps      — Detailed memory map with per-mapping RSS, PSS, swap usage
    /proc/<pid>/stat       — Raw process stats (used by ps, top)
    /proc/<pid>/io         — I/O statistics (bytes read/written by the process)
    /proc/<pid>/cgroup     — Which cgroup the process belongs to
    /proc/<pid>/ns/        — Symlinks to namespace inodes (net, pid, mnt, etc.)
    /proc/<pid>/limits     — Current resource limits (ulimits)
    /proc/<pid>/oom_score  — Current OOM killer score
    /proc/<pid>/oom_score_adj — Adjustable OOM score bias (-1000 to 1000)
    /proc/<pid>/stack      — Kernel stack trace of the process (requires privileges)
    /proc/<pid>/wchan      — Function name where the process is sleeping

/proc/self is a symlink to /proc/<current pid> — always refers to the reading process.

System-wide files (/proc/):

    /proc/meminfo         — Memory usage breakdown (MemTotal, MemFree, Buffers, Cached, SwapTotal, etc.)
    /proc/cpuinfo         — CPU model, cores, flags (sse, avx, vmx, etc.)
    /proc/interrupts      — Per-CPU interrupt counts per IRQ line
    /proc/softirqs        — Per-CPU softirq counts
    /proc/loadavg         — 1/5/15-minute load averages and running/total processes
    /proc/uptime          — Seconds since boot, seconds idle
    /proc/version         — Kernel version string
    /proc/modules         — Loaded kernel modules (same as lsmod)
    /proc/partitions      — Block devices and their sizes
    /proc/mounts          — Mounted filesystems (same as mount output)
    /proc/filesystems     — Filesystem types the kernel supports
    /proc/devices         — Registered character and block device major numbers
    /proc/buddyinfo       — Free page counts per zone per order (buddy allocator state)
    /proc/slabinfo        — Slab cache statistics
    /proc/vmstat          — Detailed virtual memory statistics
    /proc/diskstats       — Per-disk I/O statistics (reads, writes, time in queue)
    /proc/net/            — Network statistics (tcp, udp, dev, arp, route, etc.)
    /proc/sys/            — Tunable kernel parameters (writable — this IS sysctl)

/proc/sys/ is the sysctl interface:

    /proc/sys/kernel/      — Kernel parameters (hostname, pid_max, shmmax, etc.)
    /proc/sys/vm/          — Virtual memory (swappiness, dirty_ratio, overcommit, etc.)
    /proc/sys/net/         — Networking (tcp_*, ip_forward, etc.)
    /proc/sys/fs/          — Filesystem (file-max, inotify, etc.)

Reading: cat /proc/sys/net/ipv4/ip_forward
Writing: echo 1 > /proc/sys/net/ipv4/ip_forward (same as sysctl -w net.ipv4.ip_forward=1)

---

/sys (sysfs):

Mounted at /sys. Exposes the kernel's device model — a hierarchical representation of all hardware devices, drivers, buses, and kernel subsystems.

sysfs was introduced to replace the ad-hoc device-related files that had accumulated in /proc. The rule: /proc is for processes and kernel stats; /sys is for devices and drivers.

Key directories:

    /sys/bus/             — Each bus type (pci, usb, scsi, i2c, etc.)
    /sys/bus/pci/devices/ — Every PCI device, by BDF (0000:03:00.0)
    /sys/class/           — Devices grouped by type (net, block, tty, input, etc.)
    /sys/class/net/       — Network interfaces (eth0, bond0, etc.)
    /sys/class/block/     — Block devices (sda, dm-0, etc.)
    /sys/devices/         — The actual device hierarchy (the physical topology)
    /sys/module/          — Loaded kernel modules and their parameters
    /sys/kernel/          — Kernel-wide settings (mm, debug, security, etc.)
    /sys/firmware/        — Firmware interfaces (ACPI, DMI, EFI, etc.)
    /sys/fs/              — Filesystem-specific controls (cgroup, ext4, etc.)
    /sys/power/           — Power management (suspend, hibernate)

Example: Exploring a network interface:

    /sys/class/net/eth0/
    ├── address           — MAC address
    ├── mtu               — Current MTU (readable and writable)
    ├── speed             — Link speed in Mbps
    ├── operstate         — up, down, unknown
    ├── statistics/       — Packet/byte counters (rx_bytes, tx_bytes, etc.)
    ├── device/           — Symlink to the PCI device in /sys/devices/
    ├── queues/           — Per-queue stats and tuning (RSS queues)
    └── ...

Example: Exploring a PCI device:

    /sys/bus/pci/devices/0000:03:00.0/
    ├── vendor            — PCI vendor ID (0x8086 = Intel)
    ├── device            — PCI device ID
    ├── class             — PCI class code
    ├── driver/           — Symlink to the bound driver
    ├── resource           — BAR (Base Address Register) addresses
    ├── sriov_numvfs      — Number of SR-IOV VFs enabled
    ├── sriov_totalvfs    — Maximum VFs supported
    └── net/              — Network interfaces provided by this device

Example: Module parameters:

    /sys/module/bonding/parameters/
    ├── mode              — Current bonding mode
    ├── miimon            — MII monitoring interval
    └── xmit_hash_policy  — Hash policy for load balancing

    cat /sys/module/bonding/parameters/mode
    echo "balance-rr" > /sys/module/bonding/parameters/mode

---

udev and sysfs:

When a new device appears in sysfs (hardware plugged in, driver loaded, VF created), the kernel sends a uevent. udev (the userspace device manager) receives this event and:
  1. Creates the device node in /dev/ (using major/minor numbers from sysfs).
  2. Applies udev rules (/etc/udev/rules.d/) — naming, permissions, symlinks, running scripts.
  3. Notifies other daemons (NetworkManager, systemd, etc.).

    udevadm info /sys/class/net/eth0    # Show all udev properties for a device
    udevadm monitor                      # Watch uevents in real time
    udevadm trigger                      # Re-trigger all uevents (rescan devices)

---

Key Differences:

    /proc                              | /sys
    -----------------------------------|-----------------------------------
    Mostly process info + kernel stats  | Devices, drivers, and subsystems
    Flat, somewhat disorganized         | Structured device model hierarchy
    Older (Linux 1.0)                   | Newer (Linux 2.6)
    sysctl via /proc/sys/              | Module params via /sys/module/
    Some device info (legacy)           | All device info (canonical source)

Ref: Documentation/filesystems/proc.rst
Ref: Documentation/filesystems/sysfs.rst
