DMA (Direct Memory Access)

In the old days (PIO or Programmed I/O), if you wanted to load a file, the CPU had to manually read every byte from the hard drive port and write it into RAM. This meant the CPU was 100% busy just moving data.

DMA is a feature that allows hardware (like your NVMe SSD controller or Network Card) to write data directly into system RAM without bothering the CPU.

    CPU says: "Disk, please put the data from Sector 100 into RAM Address 0x5000."
    CPU: Goes back to doing other work.
    Disk: Writes data directly to RAM.
    Disk: Interrupts the CPU only when finished.

Ref: Documentation/core-api/dma-api.rst, Documentation/core-api/dma-api-howto.rst

---

DMA Addresses vs Physical Addresses:

A DMA address (dma_addr_t) is NOT necessarily the same as a physical address. Between the device and physical memory there may be an IOMMU that translates DMA addresses to physical addresses — just like the MMU translates virtual addresses to physical for the CPU.

The driver must never hand a raw physical address to a device. It must use the DMA API, which returns a dma_addr_t that the device can actually use.

---

Coherent (Consistent) DMA vs Streaming DMA:

These are the two fundamental allocation strategies:

Coherent DMA (dma_alloc_coherent):
  - Allocates a buffer that is simultaneously accessible by both the CPU and the device with no caching issues.
  - The hardware ensures coherence automatically (usually by marking the memory uncacheable or using hardware cache snooping).
  - Used for long-lived, bidirectional buffers: DMA descriptor rings, command queues, mailbox regions.
  - Expensive — the uncacheable memory is slower for CPU access, and the minimum allocation is often a full page.
  - Must be freed with dma_free_coherent().

Streaming DMA (dma_map_single / dma_map_sg):
  - Maps an existing kernel buffer for DMA access. The CPU and device take turns owning the buffer — they cannot both access it simultaneously.
  - Requires explicit ownership transfer:
      dma_sync_single_for_device() — CPU is done writing, device may now DMA.
      dma_sync_single_for_cpu()    — Device is done DMA'ing, CPU may now read.
  - Used for per-packet buffers in network drivers, per-I/O buffers in storage.
  - Much cheaper than coherent DMA because the memory stays cacheable.
  - Direction must be specified: DMA_TO_DEVICE (CPU→device), DMA_FROM_DEVICE (device→CPU), DMA_BIDIRECTIONAL.
  - Must be unmapped with dma_unmap_single() / dma_unmap_sg() when done.

---

Scatter-Gather DMA (dma_map_sg):

Instead of requiring one large contiguous physical buffer, scatter-gather lets the device DMA to/from a list of physically discontiguous pages. The driver builds a struct scatterlist array describing each chunk, and dma_map_sg() maps them all at once. Modern NICs and storage controllers have hardware scatter-gather engines that process these lists directly.

---

The IOMMU (Intel VT-d, AMD-Vi, ARM SMMU):

An IOMMU sits between PCI devices and physical memory. It translates DMA addresses (sometimes called "bus addresses" or "I/O virtual addresses") to physical addresses using page tables, similar to how the CPU's MMU works.

Benefits:
  - A 32-bit device can DMA to memory above 4GB (the IOMMU remaps it).
  - Devices are isolated — a buggy or malicious device can only access the memory regions mapped for it, not all of RAM.
  - Required for safe SR-IOV VF passthrough to VMs.

Without an IOMMU (or when using the "swiotlb" fallback), the kernel uses bounce buffers: it allocates a buffer in low memory, copies data there, and gives that address to the device.

---

DMA Zones (historical x86 legacy):

ZONE_DMA (0-16MB): For ancient ISA devices that could only address 24 bits.
ZONE_DMA32 (0-4GB): For 32-bit PCI devices that cannot address memory above 4GB.
ZONE_NORMAL (above 4GB): Everything else.

When allocating memory for DMA, the kernel considers the device's addressing capability (set via dma_set_mask()) and picks from the appropriate zone. With a modern IOMMU, these zones are less relevant because the IOMMU can remap any physical address.

    dma_set_mask_and_coherent(dev, DMA_BIT_MASK(64));
    // Tells the kernel: this device can address the full 64-bit range.

---

Bounce Buffers (SWIOTLB):

If a device cannot reach the memory where the kernel buffer lives (e.g., a 32-bit device and the buffer is above 4GB), and there's no IOMMU, the DMA API transparently allocates a "bounce buffer" in addressable memory, copies the data, and gives the bounce buffer's address to the device. On the return path, data is copied back. This is handled by SWIOTLB (Software Translation Lookaside Buffer).

SWIOTLB is also used as a fallback when the IOMMU is in "passthrough" mode or when AMD SEV / Intel TDX confidential VMs need bounce buffers for encrypted memory.

Ref: https://www.kernel.org/doc/html/latest/core-api/dma-api.html
