BERT (Boot Error Record Table), MCE, GHES, and Hardware Error Reporting

When a fatal hardware error crashes the system, the OS and firmware may stop instantly — no time to write to disk. The hardware error reporting chain ensures this information survives the reboot and is available for diagnosis.

---

Machine Check Exceptions (MCE):

MCEs are CPU-generated exceptions triggered by the hardware when it detects an internal error: ECC memory corruption, cache parity errors, bus errors, CPU thermal throttling, or internal logic faults.

When an MCE fires:
  - Corrected Error (CE): The hardware fixed it automatically (e.g., ECC corrected a single-bit flip). The kernel logs it but continues running.
  - Uncorrected Error (UCE): The hardware cannot fix it. If the error affects the currently executing code or critical data, the kernel panics. If it affects a user process's memory, the kernel can sometimes kill just that process (via memory_failure() / "machine check recovery").

MCE data is stored in Machine Specific Registers (MSRs) — CPU registers that persist across the exception but not across a power cycle.

On Linux, the kernel's MCE handler (arch/x86/kernel/cpu/mce/core.c) reads these MSRs and logs them to dmesg and /dev/mcelog.

---

BERT (Boot Error Record Table):

BERT handles the case where the error was so catastrophic that the OS didn't survive to log it.

The Timeline:

1. Fatal Error: A massive hardware failure (multi-bit memory corruption, CPU die failure, etc.) crashes the system instantly. The CPU and chipset have "sticky" registers (Machine Check Banks, NVRAM) that trap the error code in the split second before the system goes down.

2. Reboot (POST): The BIOS/UEFI firmware initializes and checks those sticky registers. It finds error data from the previous crash, formats it into the BERT (an ACPI table), and places it in memory for the OS to find.

3. OS Boot: Linux reads the ACPI tables during early init. It finds the BERT, parses the error record inside, and logs it to dmesg:

    BERT: Error records from previous boot:
    [Hardware Error]: event severity: fatal
    [Hardware Error]: Error 0, type: fatal
    [Hardware Error]:   section_type: memory error
    [Hardware Error]:   physical_addr: 0x0000001234560000

BERT is for errors that killed the previous boot. It's read once, at the start of the next boot.

---

HEST (Hardware Error Source Table):

While BERT reports errors from the dead system, HEST is the ACPI table that configures the live error reporting infrastructure. HEST tells the OS:
  - Which error sources exist (MCE banks, PCIe AER, firmware-first errors).
  - How to poll or receive interrupts for them.
  - Whether the firmware wants to handle errors first ("firmware-first" mode) or let the OS handle them directly ("OS-native" mode).

---

GHES (Generic Hardware Error Source):

GHES is the kernel driver (drivers/acpi/apei/ghes.c) that implements the HEST-defined error sources. When firmware detects a hardware error, it:
  1. Writes the error record to a shared memory region (the "error status block").
  2. Signals the OS via SCI (System Control Interrupt), NMI, or a GPIO interrupt.
  3. The GHES driver reads the error record and logs it.

GHES handles the CPER (Common Platform Error Record, UEFI specification) format — a standardized way to describe hardware errors across vendors. Error sections include: memory errors (with DIMM location, bank, row, column), PCIe errors, processor errors, and platform-specific errors.

    dmesg | grep -i "hardware error"
    dmesg | grep -i GHES

---

rasdaemon:

rasdaemon is a userspace daemon that reads hardware error events from the kernel's tracepoint interface (trace_mc_event, trace_aer_event, etc.) and stores them in a SQLite database.

Why rasdaemon instead of just dmesg?
  - dmesg is a ring buffer — old messages get overwritten.
  - rasdaemon persists errors to /var/lib/rasdaemon/ras-mc_event.db.
  - It tracks error trends over time (e.g., a DIMM accumulating corrected errors before it fails).
  - It provides ras-mc-ctl --summary and ras-mc-ctl --errors for human-readable reports.

rasdaemon replaced the older mcelog tool (which was deprecated because it relied on /dev/mcelog, a deprecated kernel interface).

    sudo systemctl start rasdaemon
    sudo ras-mc-ctl --summary           # Summary of memory errors
    sudo ras-mc-ctl --errors            # Detailed error list
    sudo sqlite3 /var/lib/rasdaemon/ras-mc_event.db "SELECT * FROM mc_event;"

---

EDAC (Error Detection And Correction):

EDAC is the kernel subsystem that monitors ECC memory controllers. It exposes error counts per memory controller, CSROW (chip-select row), and channel via sysfs:

    /sys/devices/system/edac/mc/mc0/ce_count     # Corrected errors on controller 0
    /sys/devices/system/edac/mc/mc0/ue_count     # Uncorrected errors on controller 0
    /sys/devices/system/edac/mc/mc0/csrow0/      # Per-DIMM-slot error info

EDAC feeds into rasdaemon's trace events.

---

The Error Reporting Hierarchy (summary):

    Hardware (CPU/Memory/PCIe)
        ↓
    Firmware (UEFI/BIOS)  ← configures via HEST, reports via BERT (post-crash)
        ↓
    Kernel (MCE handler, GHES driver, EDAC, AER driver)  ← logs to dmesg, tracepoints
        ↓
    Userspace (rasdaemon, journalctl -k)  ← persists to DB, alerts

Ref: Documentation/firmware-guide/acpi/apei/, UEFI Specification (CPER format)
