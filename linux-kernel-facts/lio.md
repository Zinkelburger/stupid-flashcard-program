Rather, LIO is a specific software component in Linux used for creating storage targets (servers), and iSCSI just happens to be the most common way people use it.

LIO (Linux-IO) is the standard, open-source SCSI target subsystem inside the Linux kernel.

In storage terminology, there are "Initiators" (the client asking for data) and "Targets" (the server providing the data). LIO is the software that allows a Linux machine to act as a Target. It takes local storage (like a physical hard drive, a RAID array, or a simple file) and presents it over a network or cable so that other computers think it's a local disk.

iSCSI (Internet Small Computer System Interface) is a transport protocol. It takes standard SCSI storage commands, wraps them up in TCP/IP network packets, and sends them over a standard Ethernet network.

If you want to turn a Linux server into an iSCSI storage server (like a DIY NAS or SAN), you have to use a target software stack to present those disks to the network. In modern Linux, LIO is that software stack. Therefore, when administrators are setting up, troubleshooting, or "testing" an iSCSI target on Linux, they are intrinsically interacting with and testing LIO.
