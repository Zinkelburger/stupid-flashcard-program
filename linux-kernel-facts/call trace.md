Kernel Call Traces (Stack Traces in dmesg)

When the kernel hits a bug (WARN, BUG, Oops, panic), it dumps a call trace — a snapshot of the function call stack at the moment of the error. Reading these is a core kernel debugging skill.

Ref: https://gist.github.com/doughgle/735229c34c52f9006ca92a2cf24da990

---

Anatomy of a Call Trace Entry:

    [  844.569823][ T7073]  ? uncore_pmu_event_add+0xa25/0x16d0

Breaking this down:

    [ 844.569823]     — Timestamp (seconds since boot).
    [ T7073]          — Thread/task ID that was running when the trace was captured.
    ?                 — Reliability indicator:
                          (no symbol)  = The unwinder is confident this frame is real.
                          ?            = The unwinder is NOT confident. This entry
                                        might be a stale value left on the stack
                                        from a previous function call. Treat with suspicion.
    uncore_pmu_event_add — The function name (kernel symbol).
    +0xa25             — Offset in bytes from the start of the function where
                         execution was at the time of the trace.
    /0x16d0            — Total size of the function in bytes.

The offset+size gives you a rough idea of "where in the function" the problem is. +0xa25/0x16d0 means about 63% of the way through a ~5.8KB function.

---

Reliable vs Unreliable Entries ("?" marks):

The kernel's stack unwinder walks the stack frame by frame. On x86_64 with ORC unwinder (the default since kernel 4.14), the unwinder uses a side table (.orc_unwind) to trace frames accurately. When it encounters a stack value that looks like a return address but doesn't match the ORC data, it marks it with "?".

Common causes of "?" entries:
  - Stale return addresses from functions that already returned (left on the stack).
  - Interrupt frames that interleave with the call chain.
  - Functions that don't have ORC metadata (assembly, inline, or very old code).

Rule of thumb: focus on entries WITHOUT the "?" first. They form the reliable call chain.

---

Translating Addresses to Source Lines:

The raw call trace only shows function names and offsets — not source file lines. To get line-level information, you need debug symbols.

Method 1: decode_stacktrace.sh (preferred)

The kernel source tree includes scripts/decode_stacktrace.sh. It takes a raw dmesg call trace and resolves function+offset to file:line using the vmlinux with debug symbols.

    ./scripts/decode_stacktrace.sh /path/to/vmlinux < dmesg_output.txt

This converts:
    uncore_pmu_event_add+0xa25/0x16d0
to something like:
    uncore_pmu_event_add (arch/x86/events/intel/uncore.c:382)

Method 2: addr2line (manual per-address)

    addr2line -e /path/to/vmlinux -fi ffffffff81234567

The -f flag shows the function name, -i shows inlined functions. You need the absolute address, not the offset — for module addresses, you must add the module's base address (from /proc/modules or /sys/module/NAME/sections/.text).

Method 3: gdb

    gdb /path/to/vmlinux
    (gdb) list *(uncore_pmu_event_add+0xa25)

---

Getting Debug Symbols:

Kernel images ship without debug symbols by default to keep the image small and boot fast.

Without debug symbols: Typical uncompressed vmlinux is ~40-50MB.
With debug symbols (debuginfo): The vmlinux can be 500MB-1GB+.

On RHEL/CentOS/Fedora:
    sudo dnf debuginfo-install kernel-$(uname -r)
    # Installs to /usr/lib/debug/lib/modules/$(uname -r)/vmlinux

On Ubuntu/Debian:
    sudo apt install linux-image-$(uname -r)-dbgsym
    # Or from the ddebs repository

For a custom-built kernel:
    make menuconfig → Kernel hacking → Compile-time checks → Debug information → DWARF v5
    # Or: CONFIG_DEBUG_INFO=y, CONFIG_DEBUG_INFO_DWARF5=y

---

Types of Kernel Errors (from least to most severe):

WARN_ON / WARN_ON_ONCE:
  - Prints a call trace and a "WARNING" banner to dmesg.
  - The kernel continues running. This flags an "unexpected but survivable" condition.
  - See warn_on_once.md for details.

BUG_ON:
  - Prints a call trace and kills the current process/thread.
  - The kernel continues for other tasks, but the offending task is dead.
  - More severe than WARN — the developer believes this state is unrecoverable for the task.

Oops:
  - A kernel exception (NULL pointer dereference, page fault in kernel space, etc.).
  - Prints a full register dump, call trace, and loaded module list.
  - The offending task is killed. The kernel may continue but is in a "tainted" state.
  - If the Oops happens in an interrupt handler or with locks held, the kernel may hang.

Panic:
  - The kernel considers the system unrecoverable and halts (or reboots, if panic= is set).
  - Causes: Oops in the idle thread, double fault, explicit panic() call,
    inability to mount root filesystem, init (PID 1) dying.

---

Useful dmesg Commands:

    dmesg -T                     # Human-readable timestamps
    dmesg --level=err,warn       # Only errors and warnings
    dmesg -w                     # Follow (like tail -f)
    journalctl -k                # Same as dmesg but persistent across reboots
    journalctl -k -b -1          # Kernel messages from the previous boot
