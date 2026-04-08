System Calls (Syscalls)

A system call is the mechanism by which userspace programs request services from the kernel. Every meaningful interaction with the OS — opening a file, reading a socket, forking a process, allocating memory — goes through a syscall.

---

The User/Kernel Boundary:

The CPU runs in two privilege levels:
  - Ring 3 (User mode): Your application runs here. It cannot access hardware, kernel memory, or other processes' memory directly.
  - Ring 0 (Kernel mode): The kernel runs here. It has unrestricted access to everything.

A system call is a controlled, safe transition from Ring 3 to Ring 0. The CPU has a special instruction for this:
  - x86_64: SYSCALL instruction (fast, introduced with AMD64; replaces the older INT 0x80).
  - ARM64: SVC (Supervisor Call) instruction.

---

How a System Call Works (x86_64):

1. Setup: The userspace program loads the syscall number into the RAX register and the arguments into RDI, RSI, RDX, R10, R8, R9 (up to 6 arguments).

2. Transition: The program executes the SYSCALL instruction. The CPU:
   - Saves the user's RIP (instruction pointer) and RFLAGS.
   - Switches to Ring 0.
   - Jumps to the kernel's syscall entry point (entry_SYSCALL_64 in arch/x86/entry/entry_64.S).

3. Dispatch: The kernel looks up the syscall number (RAX) in the syscall table (sys_call_table) — an array of function pointers. Each entry points to the kernel function that implements that syscall.

4. Execution: The kernel function runs (e.g., sys_read, sys_write, sys_open). It can access hardware, kernel data structures, and perform privileged operations.

5. Return: The kernel places the return value in RAX and executes SYSRET (or IRET), which switches back to Ring 3 and resumes the user program.

The whole transition (user→kernel→user) takes roughly 50-100ns on modern hardware.

---

The Syscall Table:

Each syscall has a unique number. The table is defined in:
  - x86_64: arch/x86/entry/syscalls/syscall_64.tbl
  - ARM64: include/uapi/asm-generic/unistd.h

Common syscall numbers (x86_64):

    Number | Name
    -------|----------
    0      | read
    1      | write
    2      | open
    3      | close
    9      | mmap
    56     | clone (fork)
    57     | fork
    59     | execve
    60     | exit
    62     | kill
    202    | futex
    231    | exit_group
    257    | openat
    435    | clone3

New syscalls are added with increasing numbers. Syscall numbers are NEVER reused or changed — this is the kernel's ABI stability guarantee.

---

The C Library (glibc) Wrapper:

Application code almost never invokes SYSCALL directly. Instead, glibc provides wrapper functions:

    write(fd, buf, count);   // C library function
        → sets RAX = 1 (sys_write), RDI = fd, RSI = buf, RDX = count
        → executes SYSCALL instruction
        → checks return value, sets errno if negative
        → returns result to caller

Some glibc wrappers add extra logic. For example, printf() buffers output in userspace and only calls write() when the buffer is full.

---

Error Handling:

Syscalls return a negative value on error. The kernel returns -errno (e.g., -ENOENT, -EPERM, -EINVAL). The glibc wrapper:
  1. Checks if the raw return value is negative (between -1 and -4095).
  2. If so, negates it and stores it in the thread-local errno variable.
  3. Returns -1 to the caller.

This is why C code checks if (ret == -1) and then reads errno.

---

Tracing Syscalls:

strace is the primary tool for observing syscalls in real time:

    strace ls /tmp                  # Trace all syscalls made by "ls"
    strace -c ls /tmp               # Summary: count, time per syscall
    strace -e trace=open,read ls    # Only trace specific syscalls
    strace -p <pid>                 # Attach to a running process
    strace -f ./myapp               # Follow child processes (after fork)

Each line shows: syscall_name(args...) = return_value

    openat(AT_FDCWD, "/tmp", O_RDONLY|O_DIRECTORY) = 3
    getdents64(3, /* 5 entries */, 32768) = 152
    write(1, "file1.txt\nfile2.txt\n", 19) = 19

ltrace is similar but traces library calls (glibc functions) instead of syscalls.

---

vDSO (Virtual Dynamic Shared Object):

Some syscalls (gettimeofday, clock_gettime, getcpu) are called extremely frequently but only need to read kernel data, not modify it. Going through the full SYSCALL instruction for these is wasteful.

The vDSO is a small shared library mapped by the kernel into every process's address space. It contains optimized implementations of these read-only syscalls that execute entirely in userspace — reading data from a kernel-mapped shared page — without ever entering Ring 0.

    gettimeofday() via vDSO: ~20ns
    gettimeofday() via SYSCALL: ~80ns

---

Adding a New Syscall (kernel development):

1. Add an entry to the syscall table (syscall_64.tbl).
2. Define the syscall function with SYSCALL_DEFINE macros:
       SYSCALL_DEFINE3(my_syscall, int, arg1, char __user *, arg2, size_t, arg3)
   The macro handles type checking and the __user annotation (marks pointers as coming from userspace — must be accessed via copy_from_user/copy_to_user).
3. Add a prototype to include/linux/syscalls.h.
4. Syscall numbers are never removed or reused — only appended.

Ref: Documentation/process/adding-syscalls.rst
