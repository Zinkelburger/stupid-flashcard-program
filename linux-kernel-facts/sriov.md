SR-IOV (Single Root I/O Virtualization)

SR-IOV is a PCI Express extended capability that makes one physical device appear as multiple virtual devices. It is defined in the PCI-SIG SR-IOV specification and documented in the kernel at Documentation/PCI/pci-iov-howto.rst.

Physical Function (PF): The real, full-featured PCIe device. It has complete control over the hardware and can configure SR-IOV. The PF driver (e.g., ice, ixgbe, i40e for Intel NICs, or mlx5_core for Mellanox) manages the physical hardware and creates/destroys VFs.

Virtual Function (VF): A lightweight PCIe function that shares the physical hardware with the PF but has its own Bus/Device/Function number, its own PCI config space, and its own memory-mapped register set. From the OS's perspective, each VF looks like a separate PCI device. The VF driver in Linux for Intel NICs is iavf (Intel Adaptive Virtual Function).

---

Enabling SR-IOV:

By default, SR-IOV is disabled. VFs are created at runtime, not at boot. There are two methods:

Method 1 (sysfs — recommended for new devices):

    echo 4 > /sys/bus/pci/devices/0000:03:00.0/sriov_numvfs

This creates 4 VFs under PF 0000:03:00.0. Each VF appears as a new PCI device (e.g., 0000:03:02.0, 0000:03:02.1, etc.) and is automatically probed by the iavf driver.

    echo 0 > /sys/bus/pci/devices/0000:03:00.0/sriov_numvfs

This destroys all VFs.

Method 2 (PF driver module parameter — older approach):

    modprobe ixgbe max_vfs=4

The PF driver calls pci_enable_sriov(dev, nr_virtfn) internally. On removal, it calls pci_disable_sriov(dev).

---

Why SR-IOV exists:

Without SR-IOV, VMs access the network through a software bridge (virtio-net or an emulated e1000). Every packet goes: VM → hypervisor → host kernel → NIC. This means the hypervisor must copy data and context-switch for every packet.

With SR-IOV, a VF is passed directly into the VM using VFIO/IOMMU. The VM's driver talks directly to the VF hardware registers — the hypervisor is completely out of the packet path. This gives near-native I/O performance.

The IOMMU (Intel VT-d / AMD-Vi) is required for safe VF passthrough. It restricts each VF's DMA to only the memory regions assigned to that VM, preventing one VM from reading another's memory.

---

PF vs VF capabilities:

The PF has full control: it can change link speed, update firmware, configure RSS, set the switch mode, manage MAC/VLAN filters for all VFs, and enable/disable spoofcheck.

A VF is intentionally limited. It can:
  - Send and receive packets on its assigned queues
  - Request MAC address changes (subject to PF approval)
  - Request VLAN filters
  - Use basic offloads (checksum, TSO, RSS within its allocated queues)

A VF cannot:
  - Change the physical link settings
  - See or modify other VFs' traffic
  - Update firmware
  - Exceed its allocated queue/interrupt budget

---

PF-VF communication: The Mailbox

The VF cannot access PF registers directly. Instead, they communicate through a mailbox mechanism — a shared region of memory in the NIC. The VF writes a message (e.g., "please set my MAC to aa:bb:cc:dd:ee:ff"), the PF driver reads it, decides whether to allow or deny, and writes back a response. This is asynchronous and interrupt-driven.

On Intel NICs, this is handled by the Admin Queue (AQ). Mailbox messages have opcodes (VIRTCHNL_OP_*) defined in the virtchnl.h header.

---

Useful commands:

    lspci | grep "Virtual Function"         # List VFs on the system
    ip link show dev eth0                    # Shows VF info if eth0 is a PF
    ip link set eth0 vf 0 mac aa:bb:cc:dd:ee:ff   # Set VF 0's MAC from the PF
    ip link set eth0 vf 0 vlan 100          # Assign VLAN 100 to VF 0
    ip link set eth0 vf 0 spoofchk on       # Enable spoofcheck on VF 0
    ip link set eth0 vf 0 trust on          # Trust VF 0 (allows promisc)
    cat /sys/bus/pci/devices/0000:03:00.0/sriov_totalvfs   # Max VFs supported

Ref: https://www.kernel.org/doc/html/latest/PCI/pci-iov-howto.html
