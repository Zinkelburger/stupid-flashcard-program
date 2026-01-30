SAS stands for Serial Attached SCSI

The Rule: You can plug a SATA drive into a SAS controller, and it will work.
The Trap: You cannot plug a SAS drive into a SATA controller. It will physically fit (sometimes), but the SATA controller cannot speak the "SCSI" language that SAS drives use.
