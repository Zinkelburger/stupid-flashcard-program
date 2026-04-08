echo check > /sys/devices/virtual/block/md0/md/sync_action

This triggers a RAID consistency check (aka "scrub") on the md0 array via the md (multiple devices) sysfs interface.

What happens under the hood:

Writing check to sync_action kicks off a non-destructive read pass across all member devices of the array. The md driver walks every stripe (for RAID5/6) or mirror segment (for RAID1/10), reads the corresponding blocks from each component device, and recomputes parity or compares mirror copies. If a mismatch is found, it increments /sys/devices/virtual/block/md0/md/mismatch_cnt but does not repair — it's read-only verification. This is the key distinction from writing repair, which would additionally overwrite the inconsistent copy with the "correct" one (parity-reconstructed or primary-mirror data).

The I/O path: The resync thread (md0_resync kthread) issues sequential reads against the component rdevs, honoring sync_speed_min / sync_speed_max in /proc/sys/dev/raid/. It uses the md bitmap (write-intent bitmap, if configured) to know which regions are potentially dirty, though check ignores the bitmap and walks everything. The reads go through the block layer like normal I/O but are tagged so the I/O scheduler can deprioritize them (IOPRIO_CLASS_IDLE by default in newer kernels).

Why you'd do this: Silent bit rot. A drive can return stale/corrupt data without reporting a read error (URE — unrecoverable read error, or firmware bugs). Without periodic scrubbing, you won't discover the rot until a drive fails and rebuild reads the bad block from the surviving device, at which point you've lost data. Most production setups cron this weekly or monthly and alert if mismatch_cnt > 0.


State machine context: sync_action exposes the md array's current resync state. Valid values you can read back include idle, check, repair, resync (initial sync or post-failure rebuild), and recover. Writing to it transitions the state machine — but only certain transitions are legal (e.g., you can write idle to abort a running check).


Repair:

For RAID1/10 (mirroring), the md driver treats the first device in the member list as authoritative. If mirror copies disagree, it overwrites the others with data from the primary. There's no parity to arbitrate — it's a blind "first wins" policy. If the primary device is the one with rot, you just propagated corruption to the good copies.

For RAID5/6 (parity), the driver recomputes parity from the data strips and overwrites the parity block if it disagrees. The problem is: if a data strip is the rotten one, the recomputed parity will be "consistent" with the bad data. The md driver can detect that a mismatch exists but fundamentally cannot determine which strip is wrong — it just forces consistency. After repair, the array is self-consistent but potentially consistent around the wrong data.



Check mismatch_cnt — if it's a small number of sectors, it's often benign (unclean shutdown can leave partial stripes). Huge counts are more concerning.
Before blindly repairing, if you have the means, try to identify which device is bad. Tools like smartctl -a /dev/sdX — look at Reallocated_Sector_Ct, Current_Pending_Sector, Offline_Uncorrectable. The drive with growing counters is your likely culprit.
If you have backups or checksums (ZFS does this right with per-block checksums — md/RAID does not), verify the data at the application layer first. Run sha256sum against known-good copies, btrfs scrub if on btrfs-over-md, or application-level integrity checks.
Then repair. In practice, most single-mismatch events are parity drift from unclean shutdowns or write holes (the classic md write hole — a crash between writing data and writing parity). In those cases, repair does the right thing by resyncing parity to match the data on disk.

md RAID has no data checksums. It can detect inconsistency between redundant copies but cannot detect which copy is correct. This is the core design advantage of ZFS/btrfs — they checksum every block independently of redundancy, so they can identify the rotten copy deterministically and heal from the good one.

