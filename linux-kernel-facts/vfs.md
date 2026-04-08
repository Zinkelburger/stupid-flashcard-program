VFS (Virtual File System)

The VFS is an abstraction layer in the kernel that provides a uniform interface for all filesystems. Whether you're reading from ext4 on a local disk, NFS over the network, procfs in memory, or a FUSE filesystem in userspace, the same open/read/write/close system calls work. The VFS translates these calls into the correct filesystem-specific operations.

Ref: Documentation/filesystems/vfs.rst

---

The Four Core VFS Objects:

1. Superblock (struct super_block):

Represents a mounted filesystem instance. It stores:
  - The filesystem type (ext4, xfs, btrfs, tmpfs, etc.)
  - Block size
  - Mount flags (read-only, noexec, nosuid, etc.)
  - A pointer to the root dentry of this filesystem
  - Filesystem-specific private data

Created when a filesystem is mounted (mount /dev/sda1 /mnt). Destroyed on unmount.

Each mounted filesystem has exactly one superblock.

2. Inode (struct inode):

Represents a single file, directory, symlink, device, pipe, or socket on disk. The inode contains all metadata about a file except its name:
  - File type and permissions (mode)
  - Owner (uid/gid)
  - Size
  - Timestamps (atime, mtime, ctime)
  - Link count (how many directory entries point to this inode)
  - Pointers to data blocks (or extent tree in ext4)
  - The inode number (unique within a filesystem)

Key insight: The inode does NOT contain the filename. Filenames are stored in directory entries (dentries). This is why hard links work — multiple names can point to the same inode.

    ls -i filename     # Show inode number
    stat filename      # Show all inode metadata

3. Dentry (struct dentry, "directory entry"):

Maps a filename component to an inode. The path /home/user/file.txt involves three dentries:
  - "home" → inode of /home
  - "user" → inode of /home/user
  - "file.txt" → inode of the file

Dentries are cached in the dentry cache (dcache) — a kernel hash table that avoids re-reading directory data from disk on every path lookup. The dcache is one of the most performance-critical caches in the kernel.

Negative dentries: The kernel also caches "this file does NOT exist" results. If you stat("/tmp/nonexistent"), the kernel remembers that lookup failed and returns ENOENT immediately next time without hitting disk.

4. File (struct file):

Represents an open file — the association between a process and an inode. Created by open(), destroyed by close(). Contains:
  - Current file position (the seek offset)
  - Open flags (O_RDONLY, O_WRONLY, O_APPEND, etc.)
  - A pointer to the file_operations table (the function pointers for read, write, ioctl, etc.)
  - A pointer to the dentry (and through it, the inode)

Multiple processes can have different struct file objects pointing to the same inode (same physical file, different positions and flags).

---

File Operations (struct file_operations):

Every filesystem (and device driver) provides a file_operations table — a struct full of function pointers that implement the standard file operations:

    struct file_operations {
        ssize_t (*read)(struct file *, char __user *, size_t, loff_t *);
        ssize_t (*write)(struct file *, const char __user *, size_t, loff_t *);
        int (*open)(struct inode *, struct file *);
        int (*release)(struct inode *, struct file *);
        long (*unlocked_ioctl)(struct file *, unsigned int, unsigned long);
        int (*mmap)(struct file *, struct vm_area_struct *);
        loff_t (*llseek)(struct file *, loff_t, int);
        // ... many more
    };

When you call read(fd, buf, 100) from userspace, the VFS:
  1. Looks up the struct file from the process's file descriptor table.
  2. Finds the file_operations table.
  3. Calls fops->read(file, buf, 100, &pos).
  4. The filesystem-specific code reads the data from disk (or wherever).

This is how /proc files work: procfs registers file_operations with custom read() functions that generate data on-the-fly from kernel data structures instead of reading from a disk.

---

Path Lookup (namei):

When you open("/home/user/file.txt", O_RDONLY), the VFS must resolve the path:

1. Start at the root dentry (/) or current directory (for relative paths).
2. Look up "home" in the dcache → find its inode → verify execute permission (you need +x on a directory to traverse it).
3. Look up "user" in "home"'s directory → same process.
4. Look up "file.txt" in "user"'s directory → find the target inode.
5. Check permissions on the final inode for the requested access mode (O_RDONLY).
6. Allocate a struct file, assign it a file descriptor, return the fd to userspace.

Each component lookup checks the dcache first. Only on a cache miss does the VFS call the filesystem's lookup() inode operation to read from disk.

Symlinks: If any component is a symlink, the VFS follows it (up to a recursion limit of 40 to prevent symlink loops).

Mount points: If a component is a mount point, the VFS transparently crosses into the mounted filesystem's root dentry.

---

Common Filesystem Types in Linux:

    ext4     — Default on most Linux distributions. Journaling, extents, fast.
    xfs      — Default on RHEL. Excellent for large files and parallel I/O.
    btrfs    — Copy-on-write, snapshots, checksums, subvolumes. Modern alternative.
    tmpfs    — RAM-backed filesystem (/tmp, /run). Data lost on reboot.
    procfs   — /proc. Virtual filesystem exposing kernel/process data.
    sysfs    — /sys. Virtual filesystem exposing device/driver/subsystem data.
    devtmpfs — /dev. Automatically populated device nodes.
    nfs      — Network File System. Access remote files over the network.
    overlayfs — Union filesystem. Layers multiple directories. Used by Docker.
    FUSE     — Filesystem in Userspace. Allows non-privileged users to create filesystems.

---

Useful Commands:

    mount                         # Show all mounted filesystems
    findmnt                       # Tree view of mount hierarchy
    df -hT                        # Disk usage with filesystem types
    lsof /path/to/file            # Show which processes have a file open
    cat /proc/filesystems         # Filesystems the kernel supports
    echo 2 > /proc/sys/vm/drop_caches  # Drop dentries and inodes from cache

Ref: Documentation/filesystems/vfs.rst
