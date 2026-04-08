Linux Namespaces

Namespaces provide process-level isolation of global system resources. A process inside a namespace sees its own private copy of a resource (PIDs, network stack, mount table, etc.) that is invisible to processes in other namespaces.

Namespaces are the second foundational technology behind Linux containers (the first being cgroups). Together they provide: cgroups = "how much" you can use; namespaces = "what you can see."

---

The Seven Namespace Types:

1. Mount Namespace (CLONE_NEWNS, since Linux 2.4.19):

Each mount namespace has its own mount table. Mounting or unmounting a filesystem in one namespace doesn't affect others.

This is how containers have their own root filesystem (overlayfs) without interfering with the host's mounts.

    unshare --mount -- bash          # Enter a new mount namespace
    mount --bind /tmp /mnt           # Only visible in this namespace

Was the first namespace type added (2002). The flag is named CLONE_NEWNS (not CLONE_NEWMNT) because at the time, it was the only namespace.

2. PID Namespace (CLONE_NEWPID, since Linux 3.8):

Processes inside a PID namespace see their own set of PIDs starting from 1. The first process in the namespace becomes PID 1 (the "init" of that namespace).

From outside, the same process has a different PID in the parent namespace. This is how a container's process tree looks self-contained (PID 1 = the container's entrypoint).

    unshare --pid --fork -- bash     # New PID namespace
    echo $$                          # Shows 1 (you are PID 1 inside)

PID namespaces are hierarchical — the parent can see all children's PIDs, but children can't see the parent's PIDs.

3. Network Namespace (CLONE_NEWNET, since Linux 2.6.29):

Each network namespace has its own network interfaces, routing tables, firewall rules (iptables/nftables), and /proc/net.

You already use this: ip netns add n1 creates a network namespace. Containers get their own network namespace with a virtual Ethernet pair (veth) connecting them to the host's bridge.

    ip netns add test
    ip link add veth0 type veth peer name veth1
    ip link set veth1 netns test     # Move one end into the namespace
    ip netns exec test ip addr show  # See interfaces inside

4. UTS Namespace (CLONE_NEWUTS, since Linux 2.6.19):

Isolates the hostname and NIS domain name. Each container can have its own hostname without affecting the host.

    unshare --uts -- bash
    hostname container-1             # Only changes hostname inside this namespace

UTS stands for "UNIX Time-Sharing" — a historical name from the uname() syscall data structure.

5. IPC Namespace (CLONE_NEWIPC, since Linux 2.6.19):

Isolates System V IPC objects (shared memory segments, message queues, semaphores) and POSIX message queues. Processes in different IPC namespaces cannot see each other's IPC objects.

6. User Namespace (CLONE_NEWUSER, since Linux 3.8):

Maps UIDs and GIDs inside the namespace to different UIDs outside. A process can be root (UID 0) inside the namespace but an unprivileged user (UID 1000) outside. This is the foundation of rootless containers.

    unshare --user --map-root-user -- bash
    id     # Shows uid=0(root) — but you're still unprivileged on the host

The mapping is defined in /proc/<pid>/uid_map and /proc/<pid>/gid_map.

User namespaces also affect capability checks — a process has full capabilities (CAP_SYS_ADMIN, CAP_NET_ADMIN, etc.) inside its user namespace but no capabilities outside it.

7. Cgroup Namespace (CLONE_NEWCGROUP, since Linux 4.6):

Virtualizes the view of /proc/<pid>/cgroup so that processes inside see their own cgroup as the root. A container process sees itself in "/" instead of "/docker/abc123def/...".

---

How Containers Combine Namespaces:

When Docker/Podman creates a container, it creates a new instance of (usually) all seven namespaces:

    clone(CLONE_NEWPID | CLONE_NEWNS | CLONE_NEWNET | CLONE_NEWUTS |
          CLONE_NEWIPC | CLONE_NEWUSER | CLONE_NEWCGROUP)

The container sees:
  - Its own PID 1 (PID namespace)
  - Its own root filesystem (mount namespace + overlayfs)
  - Its own network interfaces and routing (network namespace)
  - Its own hostname (UTS namespace)
  - Its own IPC objects (IPC namespace)
  - UID 0 mapped to an unprivileged host UID (user namespace, rootless mode)
  - Its cgroup as "/" (cgroup namespace)

---

Syscalls for Namespaces:

    clone()   — Create a new process in new namespaces (used by container runtimes).
    unshare() — Move the calling process into new namespaces (used by unshare command).
    setns()   — Join an existing namespace (used by nsenter / docker exec).

---

Inspecting Namespaces:

    ls -la /proc/<pid>/ns/          # Show namespace inodes for a process
    lsns                             # List all namespaces on the system
    nsenter -t <pid> -m -p -n -- bash  # Enter another process's mount/PID/net namespaces
    unshare --pid --fork -- bash     # Create and enter a new PID namespace

Each namespace is identified by an inode number in /proc/<pid>/ns/. Two processes in the same namespace have the same inode number for that type:

    readlink /proc/1/ns/net          # e.g., net:[4026531840]
    readlink /proc/500/ns/net        # Same number = same namespace

Ref: man 7 namespaces, man 2 clone, man 1 unshare, man 1 nsenter
