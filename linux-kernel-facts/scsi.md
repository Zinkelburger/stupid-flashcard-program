SCSI Tape (st) driver

https://tldp.org/HOWTO/SCSI-2.4-HOWTO/

 The tape driver interface is documented in the file /usr/src/linux/drivers/scsi/README.st and on the st(4) man page (type man st). The file README.st also documents the different parameters and options of the driver together with the basic mechanisms used in the driver.

The tape driver is usually accessed via the mt command (see man mt). mtx is an associated program for controlling tape autoloaders (see mtx.sourceforge.net).

The st driver detects those SCSI devices whose peripheral device type is "Sequential-access" (code number 1) unless they appear on the driver's "reject_list". [Currently the OnStream tape drives (described in a following section) are the only entry in this reject_list.]

The st driver is capable of recognizing 32 tape drives. There are 8 device file names for each tape drive: a rewind and non-rewind variant for each of 4 modes (numbered 0 to 3). See the tape device file name examples in Section 3.2 on device names. Any number of tape drives (up to the overall limit of 32) can be added after the st driver is loaded.

ATAPI tape drives can be controlled by this driver with help from the ide-scsi pseudo adapter driver. The discussion in Section 9.2.4 also applies for ATAPI tape drives (and ATAPI floppies).

SCSI Level,Period,Characteristics
SCSI-1,Early 80s,"The original ""alphabet."" Very basic commands."
SCSI-2,Early 90s,"Added Command Queuing and established the ""Tape"" command set we still use today."
SCSI-3+,Late 90s-Present,Transitioned into modern SAS (Serial Attached SCSI) and LUN mapping.