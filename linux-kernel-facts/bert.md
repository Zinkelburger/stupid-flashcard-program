BERT (Boot Error Record Table) is read on the next boot

When a fatal hardware error happens (like a CPU failure or massive memory corruption), the Operating System and Firmware might stop running instantly. They don't have time to write to the hard drive.

However, the Hardware itself (the CPU and Chipset) has a tiny amount of "sticky" storage (Machine Check Registers or NVRAM) that stays powered for a split second or persists across a reboot. The hardware "traps" the error code into these registers just as the system goes down.
2. The Reboot (POST)

The system restarts. The BIOS/UEFI Firmware initializes (Power-On Self-Test).

    The Firmware checks those "sticky" registers.

    It sees the error data from the previous crash.

    It formats this data into the BERT (Boot Error Record Table).

    It hands this table off to the OS memory area.


  3. The OS Loads (The Reader)

Linux boots up. During its initialization, it reads the ACPI tables provided by the BIOS.

    It finds the BERT.

    It reads the error entry inside it.

    It prints it to dmesg or /var/log/messages.

BERT is for errors that kill the system. HEST (Hardware Error Source Table) is what sets up the communication line for the live errors
