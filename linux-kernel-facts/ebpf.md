eBPF (extended Berkeley Packet Filter)

eBPF is a technology that lets you run sandboxed programs inside the Linux kernel without modifying kernel source code or loading kernel modules. It's used for observability (tracing, profiling), networking (firewalling, load balancing), and security (syscall filtering).

Think of it as: "safely inject custom code into the running kernel."

Ref: Documentation/bpf/

---

How eBPF Works:

1. Write: You write a small program in C (or Rust), compiled to eBPF bytecode.
2. Verify: The kernel's eBPF verifier statically analyzes the bytecode to guarantee safety:
   - No unbounded loops (must provably terminate).
   - No out-of-bounds memory access.
   - No null pointer dereferences.
   - All code paths return a value.
   - Stack depth is bounded (512 bytes max).
   If the verifier rejects the program, it is never executed.
3. JIT: The bytecode is Just-In-Time compiled to native machine code (x86_64, ARM64, etc.) for near-native performance.
4. Attach: The compiled program is attached to a kernel hook point (a tracepoint, kprobe, network event, etc.).
5. Execute: Every time the hook fires, the eBPF program runs. It can read kernel data, make decisions, and communicate results to userspace via maps.

---

eBPF Program Types (the most important ones):

Tracing:
  - kprobe/kretprobe: Attach to any kernel function entry/exit. Trace arguments and return values.
  - tracepoint: Attach to stable, documented kernel tracepoints (e.g., sched:sched_switch, net:net_dev_xmit).
  - uprobe/uretprobe: Attach to userspace function entry/exit (e.g., trace every malloc() call).
  - perf_event: Attach to hardware performance counters (cache misses, branch mispredictions, cycles).

Networking:
  - XDP (eXpress Data Path): Runs at the earliest possible point — before the kernel's network stack even allocates an sk_buff. Can DROP, PASS, TX (redirect), or REDIRECT packets at line rate. Used for DDoS mitigation, load balancing.
  - tc (Traffic Control): Runs in the TC (traffic control) layer. More features than XDP (can access sk_buff) but slower.
  - socket filter: Attached to sockets. Can filter packets before they reach the application.
  - cgroup/sock: Attached to cgroups. Can control socket operations for all processes in a cgroup (connection tracking, rate limiting).

Security:
  - LSM (Linux Security Module): eBPF programs that implement MAC (mandatory access control) checks.
  - seccomp: Syscall filtering (restrict which syscalls a process can make). Widely used by container runtimes and Chrome/Firefox sandboxes.

---

eBPF Maps:

Maps are key-value data structures shared between the eBPF program (kernel side) and the userspace control program. They are how eBPF communicates results to you.

Common map types:
  - BPF_MAP_TYPE_HASH: Hash table. General-purpose key-value store.
  - BPF_MAP_TYPE_ARRAY: Fixed-size array. Indexed by integer.
  - BPF_MAP_TYPE_PERCPU_HASH / PERCPU_ARRAY: Per-CPU variants (no lock contention, each CPU has its own copy).
  - BPF_MAP_TYPE_RINGBUF: Efficient ring buffer for streaming events to userspace.
  - BPF_MAP_TYPE_LRU_HASH: Hash with automatic LRU eviction (useful for connection tracking).
  - BPF_MAP_TYPE_PROG_ARRAY: Array of eBPF programs for tail calls (chaining programs).

Example: An eBPF program attached to a tracepoint increments a counter in a hash map every time a syscall fires. A userspace program reads the map and displays "top 10 syscalls per second."

---

The eBPF Verifier:

The verifier is the safety guarantee. Before any eBPF program runs, the verifier:
  1. Walks every possible code path (up to 1 million instructions verified).
  2. Tracks register types and pointer bounds for every instruction.
  3. Ensures no unbounded loops (bounded loops with known iteration count are allowed since kernel 5.3).
  4. Ensures all memory accesses are within safe bounds.
  5. Checks that helper function calls are valid for the program type.

If the verifier rejects a program, you get a detailed error explaining which instruction is unsafe and why.

---

BTF (BPF Type Format):

BTF is type information embedded in the kernel and eBPF programs. It makes eBPF programs portable across kernel versions — a program compiled on kernel 6.1 can run on kernel 6.8 because BTF describes the layout of kernel structures, and the CO-RE (Compile Once, Run Everywhere) loader adjusts field offsets at load time.

Without BTF/CO-RE, eBPF programs broke whenever a struct changed between kernel versions.

    ls /sys/kernel/btf/vmlinux       # BTF data for the running kernel

---

Userspace Tools and Libraries:

bpftool: The Swiss-army knife for eBPF inspection.
    bpftool prog list               # List all loaded eBPF programs
    bpftool map list                # List all maps
    bpftool prog dump xlated id N   # Disassemble a loaded program
    bpftool net list                # Show eBPF programs attached to network interfaces

libbpf: The standard C library for loading and attaching eBPF programs. Handles ELF parsing, map creation, program loading, BTF relocation, and CO-RE.

bcc (BPF Compiler Collection): A toolkit of pre-built tracing tools:
    opensnoop      — Trace every file open() across the entire system
    execsnoop      — Trace every execve() (every new process)
    biolatency     — Histogram of block I/O latencies
    tcplife        — Trace TCP connection lifetimes and throughput
    funccount      — Count kernel function calls
    trace          — Trace arbitrary kernel/user functions with custom filters

bpftrace: A high-level tracing language (like awk for eBPF):
    bpftrace -e 'tracepoint:syscalls:sys_enter_openat { printf("%s %s\n", comm, str(args.filename)); }'
    # Prints every file opened and by which process

    bpftrace -e 'kprobe:tcp_sendmsg { @bytes[comm] = sum(arg2); }'
    # Sum TCP bytes sent per process name

---

XDP (eXpress Data Path):

XDP is the highest-performance eBPF hook for networking. The program runs inside the NIC driver's NAPI poll loop, before any sk_buff allocation. It can process millions of packets per second per core.

XDP actions:
  XDP_DROP:     Drop the packet. Used for DDoS mitigation (drop bad packets before they consume kernel resources).
  XDP_PASS:     Pass the packet up to the normal network stack.
  XDP_TX:       Transmit the packet back out the same interface (hairpin).
  XDP_REDIRECT: Redirect to another interface, CPU, or socket.

XDP modes:
  - Native (driver mode): Runs in the NIC driver. Fastest. Requires driver support.
  - Generic (skb mode): Runs later in the stack. Works with any NIC but slower.
  - Offloaded: Runs on the NIC hardware itself (supported by some Netronome/Mellanox cards).

    ip link set dev eth0 xdpgeneric obj my_xdp_prog.o   # Attach XDP program
    ip link set dev eth0 xdpgeneric off                   # Detach

---

Real-World eBPF Users:

  - Facebook/Meta: Katran (XDP-based load balancer handling all Facebook traffic).
  - Cloudflare: DDoS mitigation at the edge using XDP.
  - Cilium: Kubernetes networking and security entirely in eBPF (replaces iptables).
  - Android: eBPF-based network traffic accounting and firewall.
  - systemd: Uses eBPF cgroup programs for IP accounting.
  - Falco: Runtime security monitoring via eBPF syscall tracing.

Ref: Documentation/bpf/
Ref: https://ebpf.io/
Ref: https://docs.cilium.io/en/latest/bpf/
