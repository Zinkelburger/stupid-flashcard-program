iavf — Intel Adaptive Virtual Function Driver

iavf is the VF (Virtual Function) driver for Intel Ethernet NICs. It only loads when SR-IOV Virtual Functions exist — either inside a VM with a passed-through VF, or on the host after VFs have been created from a PF.

You cannot request device driver="iavf" in systems like Beaker because no machine has an iavf device in its PCI inventory at boot time. The iavf driver only loads after you create VFs via:

    echo N > /sys/bus/pci/devices/<PF BDF>/sriov_numvfs

---

PF Drivers (the "parent" drivers that create VFs):

NIC Family            | PF Driver | VF Driver
----------------------|-----------|----------
Intel E810 (100GbE)   | ice       | iavf
Intel X710/XL710      | i40e      | iavf
Intel 82599 (10GbE)   | ixgbe     | ixgbevf (older, NOT iavf)
Intel I350 (1GbE)     | igb       | igbvf (older, NOT iavf)

iavf is the "unified" VF driver for newer Intel NICs. Older NICs had per-family VF drivers (ixgbevf, igbvf). The trend is for all new Intel NICs to use iavf as the single VF driver.

---

What Can a VF Do vs a PF?

VF (iavf) capabilities:
  - Send and receive packets on its allocated queues
  - Basic offloads: checksum, TSO (TCP Segmentation Offload), RSS (within its queue set)
  - Request MAC address and VLAN filter changes (subject to PF approval)
  - Interrupt handling via MSI-X vectors assigned by the PF

VF (iavf) limitations:
  - Cannot change physical link speed or view link-level statistics
  - Cannot modify firmware or NVM (non-volatile memory)
  - Cannot see or affect other VFs' traffic
  - Queue count and interrupt vectors are capped by PF configuration
  - Cannot enter promiscuous mode unless the PF sets "trust on" for that VF

---

PF-VF Communication: The Mailbox / Admin Queue

The VF cannot directly access PF registers. Communication happens through a mailbox mechanism implemented in the NIC hardware:

1. The VF writes a virtchnl (Virtual Channel) message into a shared mailbox region.
2. The PF driver receives an interrupt indicating a VF mailbox message.
3. The PF driver reads the message, validates it (is this VF allowed to do this?), and writes back a response.
4. The VF receives an interrupt with the response.

The message opcodes are defined in virtchnl.h (e.g., VIRTCHNL_OP_CONFIG_VSI_QUEUES, VIRTCHNL_OP_ADD_ETH_ADDR, VIRTCHNL_OP_ENABLE_QUEUES).

This is how the VF negotiates its capabilities at initialization:
  1. VF sends VIRTCHNL_OP_GET_VF_RESOURCES → PF responds with queue count, supported offloads, MAC address.
  2. VF sends VIRTCHNL_OP_CONFIG_VSI_QUEUES → PF programs the hardware queues.
  3. VF sends VIRTCHNL_OP_ENABLE_QUEUES → PF enables the queues. Traffic can now flow.

---

Why iavf Exists (instead of just using the PF driver in a VM):

Performance: With VF passthrough (VFIO), the VM's iavf driver talks directly to hardware registers — no hypervisor in the data path. Near-native I/O performance.

Isolation: Each VF has its own queue set, MAC address, and VLAN. The IOMMU restricts each VF's DMA to only the VM's memory. One VM cannot snoop another's traffic.

Density: A single 100GbE NIC can present 128+ VFs, each assigned to a different VM. Far more efficient than software-emulated NICs.

---

Useful debugging:

    ethtool -i <vf_interface>            # Shows driver (iavf), firmware version
    ip link show <pf_interface>          # Shows VF list with MAC, VLAN, spoofchk, trust
    dmesg | grep iavf                    # VF init messages, mailbox errors
    cat /proc/interrupts | grep iavf     # MSI-X vector assignments

Ref: Documentation/networking/device_drivers/ethernet/intel/iavf.rst
