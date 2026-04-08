Tape Devices in Linux

Tape drives have two device nodes in Linux:

/dev/st0 (Rewinding): Every time a program finishes using the tape and closes the file descriptor, the drive automatically rewinds the tape to the beginning.

/dev/nst0 (Non-rewinding, the "n" prefix): The tape head stays exactly where it stopped. This is essential for writing multiple archives (files) sequentially on a single tape.

The number (0, 1, 2...) corresponds to the physical tape drive. Each drive also has four "modes" (0-3) for density/compression settings, giving device names like /dev/st0, /dev/st0l (mode 1), /dev/st0m (mode 2), /dev/st0a (mode 3).

---

Why Tapes Still Exist:

Cost: LTO tape is dramatically cheaper per TB than any disk or SSD for cold storage (~$5/TB vs ~$15-25/TB for HDD).
Longevity: A properly stored LTO tape lasts 30+ years.
Air gap: A tape on a shelf is immune to ransomware, network attacks, and accidental deletion.
Compliance: Many regulated industries (finance, healthcare, government) mandate tape backups.

---

LTO (Linear Tape-Open) Generations:

Generation | Capacity (native) | Compressed (~2.5:1) | Transfer Speed
-----------|-------------------|---------------------|---------------
LTO-5      | 1.5 TB            | 3 TB                | 140 MB/s
LTO-6      | 2.5 TB            | 6.25 TB             | 160 MB/s
LTO-7      | 6 TB              | 15 TB               | 300 MB/s
LTO-8      | 12 TB             | 30 TB               | 360 MB/s
LTO-9      | 18 TB             | 45 TB               | 400 MB/s

LTO drives can read tapes from one generation back and write tapes from the current generation only (LTO-9 reads LTO-8 tapes, but can only write LTO-9 format).

---

The mt Command (Magnetic Tape control):

mt is the primary tool for positioning and controlling tape drives. Always use the non-rewinding device (/dev/nst0) to avoid losing your position.

    mt -f /dev/nst0 status       # Show drive status (position, block size, density)
    mt -f /dev/nst0 rewind       # Rewind to beginning of tape
    mt -f /dev/nst0 fsf 1        # Forward-space one filemark (skip to next archive)
    mt -f /dev/nst0 bsf 1        # Back-space one filemark
    mt -f /dev/nst0 eod          # Seek to end of data (to append)
    mt -f /dev/nst0 erase        # Erase the entire tape (slow — writes zeros)
    mt -f /dev/nst0 offline      # Rewind and eject the tape
    mt -f /dev/nst0 setblk 0     # Set variable block size (required for some backup tools)

---

Filemarks:

A filemark is a special marker written to the tape that separates archives. When tar finishes writing an archive, it writes a filemark. When you mt fsf 1, the drive fast-forwards to the next filemark. This is how you navigate between multiple backups on a single tape.

---

Writing and Reading with tar:

Writing the first archive:
    tar cf /dev/nst0 /home/user/documents

Writing a second archive on the same tape (appending):
    mt -f /dev/nst0 eod          # Move to end of data
    tar cf /dev/nst0 /var/log     # Write second archive after the first

Reading the first archive:
    mt -f /dev/nst0 rewind       # Go to beginning
    tar xf /dev/nst0             # Extract first archive

Reading the second archive:
    mt -f /dev/nst0 rewind       # Go to beginning
    mt -f /dev/nst0 fsf 1        # Skip past first archive's filemark
    tar xf /dev/nst0             # Extract second archive

Listing contents without extracting:
    tar tf /dev/nst0

---

The SCSI Tape Driver (st):

In the kernel, tape drives are managed by the st (SCSI Tape) driver. It detects any SCSI device with peripheral device type 1 (Sequential-access). The driver supports up to 32 tape drives.

ATAPI tape drives (IDE interface) can also use the st driver via the ide-scsi compatibility layer (mostly historical).

For tape autoloaders and libraries (robots that automatically swap tapes), the mtx command controls the medium changer device (/dev/sch0, SCSI device type 8).

Ref: Documentation/scsi/st.rst, man st(4)
