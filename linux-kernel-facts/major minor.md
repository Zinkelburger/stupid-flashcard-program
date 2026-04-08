Major and Minor Device Numbers

Every device node in /dev/ is identified by a pair of numbers: a major number and a minor number. Together, they tell the kernel which driver to invoke and which specific device instance to talk to.

Major number: Identifies the driver (or class of device). All SCSI disks share major 8. All virtual terminals share major 4.

Minor number: Identifies a specific device managed by that driver. For SCSI disks (major 8), minor 0 = /dev/sda, minor 1 = /dev/sda1 (first partition), minor 16 = /dev/sdb (second disk), minor 17 = /dev/sdb1, etc.

---

Common Examples:

    ls -l /dev/sda
    brw-rw---- 1 root disk 8, 0 Mar 30 12:00 /dev/sda
                                 ^  ^
                                 |  minor = 0 (first SCSI disk, whole device)
                                 major = 8 (SCSI disk driver)

Major | Minor  | Device
------|--------|------------------
1     | 3      | /dev/null
1     | 5      | /dev/zero
1     | 8      | /dev/random
4     | 0-63   | /dev/tty0 - tty63 (virtual consoles)
8     | 0      | /dev/sda (first SCSI/SATA/USB disk)
8     | 1      | /dev/sda1 (first partition of sda)
8     | 16     | /dev/sdb (second disk — each disk gets 16 minors)
8     | 32     | /dev/sdc
253   | 0+     | device-mapper (LVM volumes, dm-0, dm-1, ...)

For SCSI disks, the formula is: minor = (disk_number * 16) + partition_number. This limits each disk to 15 partitions (minors 1-15 within each block of 16).

---

Block vs Character Devices:

The first character of ls -l output tells you the type:
  b = block device (disks, partitions — random access, buffered I/O)
  c = character device (terminals, serial ports, /dev/null — sequential byte stream)

The major/minor system works the same for both, but they use separate number spaces. Major 8 for block devices (sd) is different from major 8 for character devices (if one existed).

---

Looking Up Registered Majors:

    cat /proc/devices

This shows two lists: "Character devices" and "Block devices", each showing major numbers and their associated driver names.

---

Creating Device Nodes Manually (mknod):

    mknod /dev/mydevice b 8 0
          ^             ^ ^ ^
          path          | | minor
                        | major
                        b=block, c=char

In practice, device nodes are created automatically by udev based on events from the kernel (kobject uevents). You rarely need mknod on modern systems.

---

Dynamic vs Static Major Numbers:

Historically, major numbers were statically assigned in Documentation/admin-guide/devices.txt (e.g., 8 = sd, 3 = hd, 9 = md). The file is maintained by LANANA.

Modern practice is dynamic allocation: the driver calls register_blkdev(0, "mydriver") with major=0, and the kernel assigns the next available major number. This avoids conflicts as the number of drivers grows. Device-mapper (usually 253) and NVMe (259) use dynamically assigned majors.

Static allocation is still used for the classic devices (sd=8, md=9, etc.) for backward compatibility.

---

In the Kernel:

Drivers register their major/minor range with the kernel:

  Block devices:   register_blkdev() / unregister_blkdev()
  Character devices: register_chrdev_region() / alloc_chrdev_region()

When userspace opens /dev/sda, the VFS looks up the major number in a hash table, finds the registered driver, and calls its open() function. The minor number is passed along so the driver knows which specific device is being accessed.

Ref: Documentation/admin-guide/devices.txt
