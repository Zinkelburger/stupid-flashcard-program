cgroups (Control Groups)

cgroups are a kernel mechanism for organizing processes into hierarchical groups and applying resource limits, accounting, and isolation to those groups. They are one of the two foundational technologies behind Linux containers (the other being namespaces).

Ref: Documentation/admin-guide/cgroup-v2.rst

---

What cgroups control:

A cgroup limits how much of a system resource a group of processes can consume. Without cgroups, any process can use 100% of CPU, all available RAM, or all disk I/O — potentially starving every other process.

With cgroups, you can say: "This group of processes can use at most 2 CPU cores, 4GB of RAM, and 100MB/s of disk I/O."

---

cgroup v1 vs cgroup v2:

cgroup v1 (legacy, 2008):
  - Each resource controller (cpu, memory, blkio, etc.) has its own independent hierarchy.
  - A process can be in different cgroups for different controllers.
  - Complex, inconsistent, hard to manage — different controllers had different interfaces.
  - Still supported but deprecated for new usage.

cgroup v2 (unified, default since kernel 5.2+):
  - A single unified hierarchy for ALL controllers.
  - A process belongs to exactly one cgroup, and all controllers apply to that cgroup.
  - Cleaner, more consistent interface.
  - Mandatory for systemd (>= 248), Kubernetes (>= 1.25), and modern container runtimes.

The cgroup filesystem is mounted at:
  v1: /sys/fs/cgroup/<controller>/  (separate mount per controller)
  v2: /sys/fs/cgroup/               (single unified mount)

    mount | grep cgroup              # See what's mounted
    cat /sys/fs/cgroup/cgroup.controllers   # Available controllers (v2)

---

Key Controllers (cgroup v2):

cpu:
  - cpu.max: "quota period" format. "200000 100000" means 200ms of CPU time every 100ms = 2 CPU cores max.
  - cpu.weight: Proportional CPU sharing (1-10000, default 100). Higher weight = more CPU when contended.
  - Replaces v1's cpu.cfs_quota_us and cpu.shares.

memory:
  - memory.max: Hard limit in bytes. If exceeded, the OOM killer targets processes in this cgroup.
  - memory.high: Soft limit. Exceeding it triggers aggressive reclaim (slows the process) but doesn't kill it.
  - memory.current: Current memory usage.
  - memory.swap.max: Limit on swap usage.
  - Replaces v1's memory.limit_in_bytes.

io (disk I/O):
  - io.max: Rate limit per block device. "8:0 rbps=10485760 wbps=10485760" = 10MB/s read and write on /dev/sda.
  - io.weight: Proportional I/O weight (1-10000).
  - Replaces v1's blkio.

pids:
  - pids.max: Maximum number of processes in this cgroup. Prevents fork bombs.

cpuset:
  - cpuset.cpus: Which CPU cores the cgroup can use ("0-3" or "0,2,4").
  - cpuset.mems: Which NUMA nodes.

---

How Containers Use cgroups:

When Docker (or Podman, containerd) creates a container, it:
  1. Creates a new cgroup under /sys/fs/cgroup/system.slice/docker-<id>.scope/ (or similar).
  2. Writes resource limits to the cgroup's control files.
  3. Moves the container's PID into the cgroup's cgroup.procs file.

    docker run --cpus="1.5" --memory="512m" ...
    # Docker creates a cgroup with:
    #   cpu.max = "150000 100000"    (1.5 CPUs)
    #   memory.max = 536870912       (512MB)

---

systemd and cgroups:

systemd uses cgroups extensively. Every systemd unit (service, slice, scope) gets its own cgroup:

    /sys/fs/cgroup/
    ├── system.slice/           # System services
    │   ├── sshd.service/       # sshd's cgroup
    │   ├── httpd.service/      # httpd's cgroup
    │   └── ...
    ├── user.slice/             # User sessions
    │   └── user-1000.slice/    # UID 1000's cgroup
    └── init.scope/             # PID 1 (systemd itself)

You can set resource limits on any service:

    systemctl set-property httpd.service MemoryMax=1G
    systemctl set-property httpd.service CPUQuota=200%   # 2 cores

These are persisted in /etc/systemd/system.control/ and applied on every restart.

    systemd-cgtop           # Live view of cgroup resource usage (like top for cgroups)
    systemd-cgls            # Tree view of the cgroup hierarchy

---

Manual cgroup v2 Usage:

Creating a cgroup:
    mkdir /sys/fs/cgroup/my_group

Enabling controllers:
    echo "+cpu +memory" > /sys/fs/cgroup/cgroup.subtree_control

Setting limits:
    echo "200000 100000" > /sys/fs/cgroup/my_group/cpu.max     # 2 CPUs
    echo 1073741824 > /sys/fs/cgroup/my_group/memory.max        # 1GB

Moving a process into the cgroup:
    echo <pid> > /sys/fs/cgroup/my_group/cgroup.procs

Checking membership:
    cat /proc/<pid>/cgroup     # Shows which cgroup a process belongs to

---

cgroup v2 Delegation:

Unprivileged users can manage subtrees of the cgroup hierarchy if the parent cgroup delegates control. This is how rootless containers (Podman rootless, user namespaces) work — the user gets a cgroup subtree they own and can subdivide.

Ref: Documentation/admin-guide/cgroup-v2.rst
Ref: man 7 cgroups
