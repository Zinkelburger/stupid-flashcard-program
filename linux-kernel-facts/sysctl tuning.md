Sysctl Tuning for High-Connection-Count Workloads

These are kernel parameters commonly tuned for benchmarks or servers handling tens of thousands of simultaneous connections. All are set via sysctl -w and are non-persistent (lost on reboot unless added to /etc/sysctl.d/).

---

fs.file-max = 2000000

The system-wide limit on the total number of open file descriptors across ALL processes. Every socket, pipe, and open file consumes one. The default (~1M on modern kernels) is usually fine, but stress tests opening hundreds of thousands of connections need this raised.

Note: This is different from ulimit -n, which sets the per-process limit.

---

net.core.somaxconn = 2000000

The maximum length of the listen() backlog queue. When a server calls listen(fd, backlog), the kernel caps the backlog at somaxconn. If this is too low, new connections get dropped with SYN cookies or ECONNREFUSED during burst traffic.

Default: 4096 (older kernels: 128).

---

net.core.netdev_max_backlog = 2000000

When the NIC receives packets faster than the kernel can process them (via softirq/NAPI), incoming packets queue up in a per-CPU backlog. This sets the max length of that queue. If it fills up, packets are silently dropped. You can see drops in /proc/net/softnet_stat (second column).

Default: 1000.

---

net.ipv4.tcp_max_syn_backlog = 2000000

The maximum number of half-open TCP connections (SYN received, waiting for final ACK of the three-way handshake). Under SYN flood attacks or legitimate burst traffic, this queue fills up fast. When full, new SYNs are dropped (or handled by SYN cookies if enabled).

Default: 1024-2048 depending on memory.

---

net.ipv4.ip_local_port_range = "1024 65535"

The range of ephemeral (source) ports the kernel can assign to outgoing connections. Each outgoing TCP connection needs a unique (source IP, source port, dest IP, dest port) tuple. By default (~32768-60999), you get about 28,000 ports. Expanding to 1024-65535 gives ~64,500.

If you exhaust this range you get "Cannot assign requested address" errors.

---

net.ipv4.tcp_fin_timeout = 15

How long a socket stays in the FIN_WAIT_2 state (seconds). After the local side sends FIN and receives ACK, it waits for the remote side's FIN. If the remote never sends it (buggy client, crashed process), the kernel keeps the socket open for this long. Lowering it from the default (60) frees resources faster.

Also affects how quickly TIME_WAIT sockets are recycled in practice.

---

net.ipv4.tcp_tw_reuse = 1

Allows the kernel to reuse sockets in TIME_WAIT state for new outgoing connections, as long as the new connection's timestamp is greater than the last seen timestamp on the old connection. This is safe with TCP timestamps enabled (the default).

Without this, each TIME_WAIT socket (which lasts 60 seconds by default) occupies an ephemeral port. Under heavy churn, you run out of ports.

Note: tcp_tw_recycle was removed in kernel 4.12 because it broke clients behind NAT.

---

net.ipv4.tcp_mem = "524288 2097152 4194304"

Global TCP memory limits in pages (not bytes). Three values: min, pressure, max.
  - Below min: TCP can freely allocate memory.
  - At pressure: TCP starts throttling — memory pressure mode. New allocations are restricted.
  - At max: TCP refuses to allocate more memory. Connections start failing.

These values are in pages (usually 4KB each), so "524288 2097152 4194304" = 2GB / 8GB / 16GB.

---

net.ipv4.tcp_rmem = "4096 87380 8388608"
net.ipv4.tcp_wmem = "4096 87380 8388608"

Per-socket TCP buffer sizes in bytes. Three values: min, default, max.
  - min (4096): Minimum buffer size. Even under memory pressure, each socket gets at least this.
  - default (87380 ≈ 85KB): The buffer size assigned when a socket is created.
  - max (8388608 = 8MB): The maximum the kernel's autotuning can grow the buffer to.

TCP autotuning dynamically adjusts each socket's buffer between default and max based on the measured bandwidth-delay product. For high-latency, high-bandwidth links (e.g., cross-continent 10Gbps), larger max values allow bigger TCP windows.

tcp_rmem controls the receive buffer. tcp_wmem controls the send buffer.

---

Checking current values: sysctl -a | grep net.ipv4.tcp
Making permanent: echo "net.ipv4.tcp_tw_reuse = 1" >> /etc/sysctl.d/99-tuning.conf && sysctl --system
