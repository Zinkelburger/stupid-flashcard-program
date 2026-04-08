Device Mapper (DM)

Device Mapper is a kernel framework that maps virtual block devices onto underlying physical block devices. It lives between the filesystem layer and the physical disk drivers. LVM, LUKS encryption, dm-multipath, and software RAID (dm-raid) are all built on top of Device Mapper.

The user creates a virtual device by defining a "table" of target mappings. Each mapping says: "virtual sectors X through Y are handled by target type T, with these parameters."

Ref: Documentation/admin-guide/device-mapper/

---

Core Concepts:

Mapped Device: The virtual block device (e.g., /dev/dm-0, or its symlink /dev/mapper/vg0-lv_root). Applications read/write to this like a normal block device.

Table: A list of target rules that define how I/O to the mapped device is translated to I/O on physical devices. Each row specifies: start sector, length, target type, and target-specific arguments.

Target: A kernel module that implements a specific I/O transformation. The kernel ships with many targets.

---

Common DM Targets:

linear — Maps a contiguous range of virtual sectors to a contiguous range on a physical device. The simplest target. This is what LVM logical volumes use for basic allocations.

    0 1048576 linear /dev/sda2 0
    (virtual sectors 0-1048575 map to /dev/sda2 starting at sector 0)

striped — Stripes data across multiple physical devices (like RAID-0). Spreads I/O for performance.

    0 2097152 striped 2 128 /dev/sda2 0 /dev/sdb2 0
    (2 devices, 128-sector stripe size)

mirror — Mirrors data to multiple devices (like RAID-1). Writes go to all mirrors; reads come from any one.

snapshot / snapshot-origin — Creates copy-on-write snapshots of a volume. The origin continues normal I/O; the snapshot records only changed blocks. Used by LVM snapshots.

thin-pool / thin — Thin provisioning. Allocates storage on demand rather than upfront. A 1TB thin volume might only consume 100GB of actual disk until data is written. Supports snapshots natively. Used by Docker's devicemapper storage driver and LVM thin volumes.

crypt — Encrypts/decrypts I/O using the kernel's crypto API. This is what LUKS (dm-crypt/cryptsetup) uses.

    0 1048576 crypt aes-xts-plain64 <key> 0 /dev/sda2 0

raid — Software RAID implemented in Device Mapper (alternative to md). Supports RAID 0/1/4/5/6/10. Used by LVM RAID.

multipath — Routes I/O across multiple physical paths to the same storage device. If one path fails (cable, HBA, switch), I/O continues on the surviving path. Essential for SAN storage and dual-port SAS.

delay — Adds artificial latency to I/O. Used for testing.

flakey — Returns I/O errors intermittently. Used for testing error handling in filesystems and applications.

zero — Returns zeros for all reads, discards all writes. Like /dev/zero but as a block device.

verity — Read-only target that verifies data integrity using a Merkle tree of hashes. Used by Android's verified boot and dm-verity for integrity-protected images.

---

LVM and Device Mapper:

LVM (Logical Volume Manager) is a userspace tool that creates Device Mapper tables. When you run:

    lvcreate -L 10G -n myvolume vg0

LVM tells the kernel to create a DM device with a linear (or thin) target mapping that points to the allocated physical extents in vg0. The result is /dev/mapper/vg0-myvolume, which is a symlink to /dev/dm-N.

LVM concepts map directly to DM:
  PV (Physical Volume) → the underlying block device
  VG (Volume Group) → a pool of PVs
  LV (Logical Volume) → a DM mapped device with a table pointing into the VG

---

dmsetup — The Low-Level DM Tool:

dmsetup talks directly to the Device Mapper ioctl interface, bypassing LVM.

    dmsetup table                     # Show the mapping table for all DM devices
    dmsetup table /dev/dm-0           # Table for a specific device
    dmsetup status                    # Show runtime status of all DM targets
    dmsetup ls --tree                 # Show DM device dependency tree
    dmsetup info                      # Show info (name, UUID, state, open count)
    dmsetup create mydev --table '0 2048 linear /dev/sdb 0'   # Create a DM device
    dmsetup remove mydev              # Remove a DM device
    dmsetup suspend mydev             # Pause I/O (for table reload)
    dmsetup resume mydev              # Resume I/O after table reload

---

Live Table Reload (how LVM resizes volumes online):

1. dmsetup suspend mydev          — Queues all new I/O.
2. dmsetup load mydev --table '...' — Loads a new table (e.g., extended size).
3. dmsetup resume mydev           — Resumes I/O with the new table.

In-flight I/O completes against the old table. New I/O uses the new table. This is atomic from the application's perspective.

In the kernel, this is managed via dm_get_live_table() / dm_put_live_table() with RCU-protected table swaps (see put-table.md).

---

Zone Reports (for zoned block devices):

For SMR (Shingled Magnetic Recording) or ZNS (Zoned Namespace NVMe) drives, DM needs to know which zones are full, empty, or partially written. The kernel can request a "zone report" from the drive — a list of zones and their write pointer positions.

The comment "For internal zone reports bypassing the top BIO submission path" in the kernel means: instead of sending a standard BIO request through the full block layer stack (which might be queued, throttled, or modified), the DM subsystem calls the drive directly to get the metadata it needs for mapping decisions.

Ref: https://www.kernel.org/doc/html/latest/admin-guide/device-mapper/index.html
