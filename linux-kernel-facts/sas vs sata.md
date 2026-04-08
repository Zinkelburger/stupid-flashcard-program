SAS vs SATA

SAS stands for Serial Attached SCSI. SATA stands for Serial ATA.

Both are serial point-to-point interfaces for connecting storage devices, but they target very different markets.

---

Physical Compatibility:

The Rule: You can plug a SATA drive into a SAS controller (backplane), and it will work. The SAS connector is a superset of the SATA connector.

The Trap: You cannot plug a SAS drive into a SATA controller. Even if it physically fits, the SATA controller cannot speak the SCSI protocol that SAS drives use. SATA controllers only understand the ATA command set.

---

Speed:

Generation         | SAS           | SATA
-------------------|---------------|--------
Gen 1              | 3 Gbps        | 1.5 Gbps
Gen 2              | 6 Gbps        | 3 Gbps
Gen 3 (current)    | 12 Gbps       | 6 Gbps
Gen 4 (SAS only)   | 22.5 Gbps     | N/A

SAS has always been approximately 2x faster than the contemporary SATA generation. SATA topped out at 6 Gbps and has no further roadmap (NVMe replaced it for high-speed storage). SAS continues to evolve.

---

Duplex:

SAS is full-duplex — it can send and receive data simultaneously on separate wire pairs. The effective bandwidth between two SAS devices is up to 2x the link speed.

SATA is half-duplex — it can send OR receive at any given moment, never both at once.

---

Dual-Port SAS (High Availability):

SAS drives have two physical ports (two sets of signal pins). Each port can connect to a different SAS controller or SAS expander. If one controller or cable fails, the other path takes over transparently. This is critical in enterprise storage where a single drive failure path cannot be tolerated.

SATA drives have a single port — one path, no failover.

In Linux, multipath I/O (dm-multipath or the kernel's built-in SCSI multipath via scsi_dh) manages the two paths automatically.

---

SAS Expanders:

A SAS expander is like an Ethernet switch for SAS. It lets you connect many more drives than the controller has physical ports.

A single SAS HBA with 8 ports can connect to an expander, which fans out to 128+ drives. Expanders can be daisy-chained for even larger topologies. This is how large JBODs (Just a Bunch Of Disks) and enterprise storage shelves work.

SATA has no equivalent. Each SATA port connects to exactly one drive.

---

Command Queuing:

SAS uses TCQ (Tagged Command Queuing) with a queue depth of up to 254 commands per device. The drive and controller can reorder and optimize I/O.

SATA uses NCQ (Native Command Queuing) with a queue depth of up to 32. Sufficient for consumer workloads but limiting under heavy enterprise I/O.

---

Error Handling:

SAS drives have much stronger error recovery built into the protocol (SCSI error handling, automatic retries, detailed sense data). SAS SSDs and HDDs are also rated for significantly higher write endurance and longer MTBF (2 million hours vs 1 million for SATA). The firmware is tuned for 24/7 operation.

---

When to use which:

SATA: Consumer desktops, NAS boxes, cold storage, budget builds. Cost per TB is significantly lower.
SAS: Enterprise servers, database storage, anything requiring high availability, dual-pathing, or sustained random I/O under heavy load.

In Linux, both SAS and SATA drives appear as /dev/sdX and use the sd (SCSI disk) driver. The kernel's libata layer translates ATA commands for SATA drives into the SCSI framework, so from the kernel's perspective, everything looks like SCSI.
