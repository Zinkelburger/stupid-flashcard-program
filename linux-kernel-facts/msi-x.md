MSI and MSI-X (Message Signaled Interrupts)

Ref: Documentation/PCI/msi-howto.rst, https://www.kernel.org/doc/html/latest/PCI/msi-howto.html

---

Legacy Pin-Based Interrupts (the old way):

In the original PCI design, devices shared physical interrupt lines (wires). A single IRQ line might be shared by a NIC, a SATA controller, and a USB controller. When any of them triggered an interrupt, the kernel had to call every handler registered on that line to figure out which device actually needed attention. This is slow and wasteful.

Also, pin-based interrupts have a subtle ordering issue: a device might write data to memory via DMA, then raise the interrupt pin. But the interrupt can arrive at the CPU before the DMA data lands in memory (especially with PCI-PCI bridges). The interrupt handler must read a device register to force all preceding DMA writes to complete — adding latency.

---

MSI (Message Signaled Interrupts):

First defined in PCI 2.2. Instead of toggling a physical wire, the device performs a memory write to a special address that the interrupt controller watches. This memory write IS the interrupt.

Key properties:
  - A device can allocate 1, 2, 4, 8, 16, or 32 MSI interrupts (must be a power of two).
  - Each interrupt is identified by a 16-bit data word written to the configured address.
  - MSI interrupts are never shared — each vector belongs to exactly one device.
  - The DMA ordering problem goes away: the interrupt write cannot pass preceding DMA writes (PCI ordering rules guarantee this).

---

MSI-X (Extended MSI):

First defined in PCI 3.0. Removes the limitations of MSI:

  - Up to 2048 interrupt vectors per device (vs 32 for MSI).
  - Each vector gets its own independent target address and data word, stored in a per-device MSI-X Table in BAR (Base Address Register) space.
  - Each vector can be independently masked.
  - Each vector can target a different CPU — critical for multi-queue devices.

A device may support both MSI and MSI-X, but only one can be active at a time. The kernel prefers MSI-X when available.

---

Why MSI-X Matters for Network Cards:

Modern NICs (10GbE, 25GbE, 100GbE) have multiple hardware TX/RX queues. Each queue needs its own interrupt so that:
  1. Each queue's interrupt can be pinned to a specific CPU core.
  2. Multiple cores can process packets in parallel without lock contention.
  3. The CPU that handles the interrupt is the same one that processes the packet (cache locality).

Without MSI-X, a NIC with 16 queues would have to multiplex all queue events through 1-2 shared interrupts, and the interrupt handler would have to check every queue — defeating the purpose of multi-queue.

With MSI-X, each queue gets its own vector → its own IRQ number → its own CPU affinity.

---

Per-Queue Interrupts and RSS:

RSS (Receive Side Scaling) is a NIC feature that hashes incoming packet headers (src/dst IP + ports) and distributes packets across multiple RX queues. Each RX queue has its own MSI-X interrupt. Combined with IRQ affinity, this spreads network processing across all CPU cores.

The hash → queue → MSI-X vector → CPU core pipeline is the foundation of high-performance Linux networking.

---

IRQ Affinity:

Each MSI-X vector gets a Linux IRQ number visible in /proc/interrupts. You can control which CPU handles each interrupt:

    cat /proc/interrupts | grep eth0
    # Shows IRQ numbers and per-CPU interrupt counts for each queue

    cat /proc/irq/75/smp_affinity         # Bitmask of allowed CPUs for IRQ 75
    echo 4 > /proc/irq/75/smp_affinity    # Pin IRQ 75 to CPU 2 (bitmask 0x4 = bit 2)

    cat /proc/irq/75/smp_affinity_list    # Human-readable CPU list
    echo 2 > /proc/irq/75/smp_affinity_list  # Same as above but using CPU number

irqbalance is the daemon that automatically distributes interrupts across CPUs. It usually does a reasonable job, but for latency-critical applications (DPDK, financial trading), manual pinning is preferred.

---

Kernel API for Drivers:

Modern drivers use a single function to allocate MSI/MSI-X vectors:

    int pci_alloc_irq_vectors(struct pci_dev *dev,
                              unsigned int min_vecs,
                              unsigned int max_vecs,
                              unsigned int flags);

flags can be PCI_IRQ_MSI, PCI_IRQ_MSIX, PCI_IRQ_INTX, or PCI_IRQ_ALL_TYPES.
If PCI_IRQ_AFFINITY is set, the kernel automatically spreads vectors across CPUs.

To get the Linux IRQ number for a vector:

    int irq = pci_irq_vector(dev, vector_number);
    request_irq(irq, my_handler, 0, "my_driver", data);

To free:

    pci_free_irq_vectors(dev);

Ref: https://www.kernel.org/doc/html/latest/PCI/msi-howto.html
