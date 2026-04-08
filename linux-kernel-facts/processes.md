Processes and the Process Lifecycle

Everything running on a Linux system — your shell, a web server, a kernel thread — is represented by a process (or thread) with its own task_struct in the kernel.

---

task_struct:

The task_struct (defined in include/linux/sched.h) is the central data structure for every process and thread. It contains everything the kernel needs to manage a task:

  - Process state (running, sleeping, stopped, zombie, dead)
  - PID and TGID (thread group ID — the "process ID" seen by userspace)
  - Pointer to the memory descriptor (mm_struct) — page tables, VMAs
  - Pointer to the filesystem info (fs_struct) — root dir, current dir
  - Pointer to the open file table (files_struct) — file descriptor array
  - Scheduling info (priority, nice value, CPU affinity, scheduler class)
  - Signal handlers and pending signals
  - Credentials (uid, gid, capabilities, SELinux context)
  - Parent/child/sibling pointers (the process tree)
  - Namespaces and cgroup membership

The task_struct is allocated from a dedicated slab cache. On x86_64, it's roughly 6-8KB.

The macro current always points to the task_struct of the process currently executing on the local CPU.

---

Process States:

TASK_RUNNING (R): The process is either executing on a CPU or waiting in the run queue ready to execute. This is the only state that means "actively runnable."

TASK_INTERRUPTIBLE (S): The process is sleeping, waiting for something (I/O completion, a mutex, a signal). It can be woken by delivering a signal. This is the most common sleep state — what you see as "S" in ps/top.

TASK_UNINTERRUPTIBLE (D): The process is sleeping and CANNOT be woken by signals. Used for short, critical waits (typically I/O operations where interrupting mid-operation would corrupt data). If a process is stuck in D state, it usually means the I/O is stuck (e.g., NFS server unreachable, hung block device). You cannot kill a D-state process — not even with SIGKILL.

TASK_STOPPED (T): The process has been stopped by SIGSTOP or SIGTSTP (Ctrl+Z). It does nothing until SIGCONT resumes it.

TASK_TRACED (t): The process is stopped by a debugger (ptrace). Similar to STOPPED but only a debugger (or the tracer) can resume it.

EXIT_ZOMBIE (Z): The process has exited, but its parent hasn't called wait() to read the exit status yet. The task_struct still exists (holding the exit code and resource usage stats), but all other resources (memory, files) have been freed.

EXIT_DEAD (X): Final state. The parent has collected the exit status. The task_struct is being freed. You never see this in ps.

---

fork() — Creating a Process:

fork() creates a new process by duplicating the calling process. The child is a near-identical copy of the parent:
  - Same code, same open files, same memory layout.
  - Different PID.
  - fork() returns the child's PID to the parent and 0 to the child.

Copy-on-Write (CoW): fork() does NOT immediately copy the parent's memory. Instead, both parent and child share the same physical pages, marked read-only. Only when one of them writes to a page does the kernel actually copy that page. This makes fork() extremely fast even for processes with gigabytes of memory.

In the kernel, fork() is implemented by the clone() syscall with specific flags.

clone(): A more flexible version of fork(). The flags argument controls what is shared:
  - CLONE_VM: Share memory space (this is how threads are created).
  - CLONE_FILES: Share the file descriptor table.
  - CLONE_FS: Share filesystem info (root, cwd).
  - CLONE_THREAD: Create a thread (same TGID, different TID).
  - CLONE_NEWPID, CLONE_NEWNS, etc.: Create new namespaces (containers).

pthread_create() in glibc uses clone(CLONE_VM | CLONE_FILES | CLONE_FS | CLONE_THREAD | ...) to create threads — processes that share everything except their stack and thread-local storage.

---

exec() — Replacing the Program:

exec() (actually execve()) replaces the current process's code and data with a new program:
  1. The kernel reads the ELF binary header.
  2. Destroys the old memory mappings (code, data, heap, stack).
  3. Maps the new program's segments into memory.
  4. Sets up the new stack with argc, argv, envp.
  5. Sets the instruction pointer to the new program's entry point.
  6. The PID does not change. Open file descriptors are preserved (unless marked FD_CLOEXEC).

The classic pattern is fork-then-exec: the parent forks, and the child immediately execs a new program. The shell does this for every command you run.

---

wait() — Reaping Children:

When a child process exits, it becomes a zombie — it's dead but its exit status is preserved in its task_struct until the parent calls wait() (or waitpid(), waitid()).

wait() does two things:
  1. Returns the child's exit status to the parent.
  2. Frees the child's task_struct (removing the zombie).

If the parent exits before calling wait(), the orphaned children are re-parented to init (PID 1) or a subreaper. init periodically calls wait() to clean up orphaned zombies.

Zombie accumulation: If a parent creates many children but never calls wait(), zombies pile up. Each zombie is small (just a task_struct), but they consume PID numbers. A system can run out of PIDs if zombies aren't reaped. Fix: signal(SIGCHLD, SIG_IGN) tells the kernel to auto-reap children.

---

The Process Tree:

Every process except PID 0 (the idle task / swapper) has a parent. This forms a tree:

    PID 0: idle (swapper) — per-CPU kernel thread, never assigned to user processes
     └── PID 1: init / systemd — first userspace process
          ├── PID 500: sshd
          │    └── PID 501: bash (your shell)
          │         └── PID 502: vim
          ├── PID 600: httpd
          └── ...
     └── PID 2: kthreadd — parent of all kernel threads
          ├── kworker/0:0
          ├── ksoftirqd/0
          ├── migration/0
          └── ...

    pstree -p       # Show the full process tree with PIDs
    ps auxf         # Show process forest (BSD syntax)

---

Kernel Threads:

Kernel threads are processes that run entirely in kernel space — they have no userspace memory (mm_struct is NULL). They are created with kthread_create() and run kernel functions directly.

Examples: kworker (workqueues), ksoftirqd (softirq processing), migration (CPU migration), kswapd (page reclaim), jbd2 (ext4 journaling).

All kernel threads are children of kthreadd (PID 2).

---

Signals:

Signals are the primary IPC mechanism for process control. Common signals:

    Signal   | Number | Default Action | Description
    ---------|--------|----------------|-----------------------------------
    SIGHUP   | 1      | Terminate      | Hangup (terminal closed)
    SIGINT   | 2      | Terminate      | Interrupt (Ctrl+C)
    SIGQUIT  | 3      | Core dump      | Quit (Ctrl+\)
    SIGKILL  | 9      | Terminate      | Kill (cannot be caught or ignored)
    SIGSEGV  | 11     | Core dump      | Segmentation fault (invalid memory access)
    SIGTERM  | 15     | Terminate      | Graceful termination request
    SIGSTOP  | 19     | Stop           | Stop process (cannot be caught)
    SIGTSTP  | 20     | Stop           | Terminal stop (Ctrl+Z)
    SIGCONT  | 18     | Continue       | Resume stopped process
    SIGCHLD  | 17     | Ignore         | Child exited (sent to parent)

SIGKILL and SIGSTOP cannot be caught, blocked, or ignored. Everything else can be handled by the process with signal() or sigaction().

    kill -SIGTERM <pid>     # Send SIGTERM (graceful shutdown)
    kill -9 <pid>           # Send SIGKILL (force kill — last resort)
    kill -l                 # List all signal names and numbers

Ref: man 7 signal, include/linux/sched.h
