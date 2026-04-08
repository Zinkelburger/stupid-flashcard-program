Memory Management in the Linux Kernel

The memory management subsystem handles virtual memory, physical page allocation, kernel memory allocation, memory-mapped files, swap, and the OOM killer. It is the most complex subsystem in the kernel.

Ref: Documentation/admin-guide/mm/, Documentation/mm/

---

Virtual Memory and Page Tables:

Every process has its own virtual address space. On x86_64, this is a 48-bit (256 TB) or 57-bit (128 PB, with 5-level paging) address range. The CPU's MMU (Memory Management Unit) translates virtual addresses to physical addresses using page tables.

Page tables are a multi-level tree (4 or 5 levels on x86_64):

    PGD (Page Global Directory)
     └── P4D (Page 4-level Directory, only with 5-level paging)
          └── PUD (Page Upper Directory)
               └── PMD (Page Middle Directory)
                    └── PTE (Page Table Entry) → Physical page frame

Each PTE maps a 4KB virtual page to a 4KB physical page frame. The PTE also stores permission bits (read/write/execute, user/supervisor) and status bits (present, dirty, accessed).

When the CPU tries to access a virtual address that has no valid PTE, a page fault occurs. The kernel's page fault handler decides what to do:
  - Valid mapping, page not in RAM → demand paging: load from disk (swap or file).
  - Valid mapping, copy-on-write → allocate a new page, copy the data.
  - Invalid mapping → SIGSEGV (segmentation fault) to the process.

The TLB (Translation Lookaside Buffer) is a CPU cache of recent virtual→physical translations. TLB misses are expensive (require walking the page table), so the kernel tries to minimize TLB flushes.

---

Physical Page Allocator (the Buddy System):

The kernel manages physical RAM in page frames (4KB each on x86). The buddy allocator organizes free pages into lists of power-of-two sizes: order-0 (4KB), order-1 (8KB), order-2 (16KB), ... up to order-10 (4MB, MAX_ORDER).

    alloc_pages(gfp_mask, order)   — Allocate 2^order contiguous pages.
    __get_free_pages(gfp_mask, order) — Same, but returns a kernel virtual address.
    free_pages(addr, order)        — Return pages to the buddy system.

When the allocator needs an order-N block but only has order-(N+1), it splits the larger block in half (the two halves are "buddies"). When freeing, if a block's buddy is also free, they're merged back into a larger block. This prevents fragmentation.

    cat /proc/buddyinfo    # Shows free pages at each order for each zone

---

Memory Zones:

Physical memory is divided into zones based on addressing constraints:

    ZONE_DMA      (0 - 16MB)     — For legacy ISA DMA (24-bit addressing).
    ZONE_DMA32    (0 - 4GB)      — For 32-bit PCI devices.
    ZONE_NORMAL   (4GB+)         — Everything the kernel can directly address.
    ZONE_MOVABLE                 — Migratable pages (for memory hot-unplug and CMA).

The allocator tries to satisfy requests from the most appropriate zone. GFP flags control this:
  GFP_KERNEL — Normal kernel allocation (can sleep, can reclaim).
  GFP_ATOMIC — Cannot sleep (interrupt context). May fail more easily.
  GFP_DMA    — Must come from ZONE_DMA.
  GFP_DMA32  — Must come from ZONE_DMA32.

---

Slab Allocator (SLUB):

The buddy allocator works in page-sized (4KB) chunks, but the kernel frequently needs much smaller allocations (64 bytes for a struct, 128 bytes for a buffer). Allocating a full page for 64 bytes wastes memory.

The slab allocator sits on top of the buddy allocator. It pre-allocates pages and carves them into fixed-size objects:

    kmalloc(size, GFP_KERNEL)    — Allocate size bytes of kernel memory.
    kfree(ptr)                    — Free it.

kmalloc uses slab caches for common sizes (8, 16, 32, 64, 128, 256, ... 8192 bytes). Each cache holds pages divided into same-sized slots.

Linux has had three slab implementations:
  - SLAB (original, complex, lots of metadata).
  - SLUB (current default since 2.6.23, simpler, better performance, less fragmentation).
  - SLOB (removed in 6.4, was for tiny embedded systems).

kmem_cache_create() lets drivers create custom slab caches for frequently allocated structures (e.g., inode cache, dentry cache, task_struct cache).

    cat /proc/slabinfo              # Shows all slab caches, object sizes, usage
    slabtop                         # Interactive view of slab usage

---

kmalloc vs vmalloc:

kmalloc(size, flags):
  - Allocates physically contiguous memory.
  - Fast (uses slab caches).
  - Limited to a few MB (practically ~128KB without special flags; up to a few MB with page allocation fallback).
  - Required when hardware needs contiguous physical pages (DMA buffers).

vmalloc(size):
  - Allocates virtually contiguous memory, but pages may be physically scattered.
  - Slower (requires modifying page tables).
  - Can allocate much larger regions (limited only by address space).
  - Cannot be used for DMA (not physically contiguous).
  - Kernel modules' code is loaded into vmalloc space.

kvmalloc(size, flags):
  - Tries kmalloc first; falls back to vmalloc if the contiguous allocation fails.
  - Best of both worlds for code that doesn't care about physical contiguity.

---

The OOM Killer:

When the system runs out of memory and cannot reclaim enough pages (by flushing caches, writing dirty pages to disk, or swapping), the kernel invokes the OOM (Out Of Memory) killer.

The OOM killer selects a process to kill based on a "badness" score:
  - Processes using more memory get higher scores.
  - Root processes get a slight discount.
  - The kernel, init (PID 1), and kernel threads are exempt.
  - Userspace can influence the score via /proc/<pid>/oom_score_adj (-1000 to 1000).
    Setting oom_score_adj to -1000 makes a process immune to the OOM killer.

    cat /proc/<pid>/oom_score      # Current calculated badness score
    echo -1000 > /proc/<pid>/oom_score_adj   # Protect this process from OOM

The OOM killer is a last resort. On servers, hitting OOM usually means misconfigured memory limits or a memory leak. Overcommit settings (/proc/sys/vm/overcommit_memory) control how aggressively the kernel promises memory it may not have:
  0 = Heuristic overcommit (default — the kernel guesses).
  1 = Always overcommit (never refuse malloc, risk OOM later).
  2 = Never overcommit (refuse malloc if commit would exceed RAM + swap * ratio).

---

Hugepages:

Standard pages are 4KB. For workloads with large memory footprints (databases, VMs), walking page tables for millions of 4KB pages wastes TLB entries and causes frequent TLB misses.

Hugepages are larger page sizes:
  x86_64: 2MB (PMD-level) or 1GB (PUD-level).

Fewer pages → fewer TLB entries needed → fewer TLB misses → better performance.

Two types:

Static Hugepages (HugePages, hugetlbfs):
  - Reserved at boot or runtime via /proc/sys/vm/nr_hugepages.
  - Allocated from a dedicated pool, NOT available for normal allocations.
  - Applications must explicitly use them (mmap with MAP_HUGETLB, or mount -t hugetlbfs).
  - Used heavily by databases (Oracle, PostgreSQL) and DPDK.

Transparent Huge Pages (THP):
  - The kernel automatically promotes 4KB pages to 2MB pages when possible.
  - No application changes required.
  - Controlled via /sys/kernel/mm/transparent_hugepage/enabled (always/madvise/never).
  - Can cause latency spikes due to compaction (the kernel rearranging pages to create contiguous 2MB regions). Some latency-sensitive workloads disable THP.

    cat /proc/meminfo | grep -i huge    # Shows hugepage allocation and usage

---

Page Cache:

When you read a file, the kernel doesn't throw away the data after delivering it to your process. It keeps the file's pages in the page cache — a region of RAM used to cache disk data. Subsequent reads of the same file hit RAM instead of disk.

The page cache also handles write buffering: writes go to the page cache first (marking pages "dirty"), and the kernel flushes dirty pages to disk asynchronously (controlled by /proc/sys/vm/dirty_ratio and dirty_expire_centisecs).

    free -h    # Shows "buff/cache" — this is the page cache + buffer cache
    echo 3 > /proc/sys/vm/drop_caches   # Drop page cache (for benchmarking)

The page cache is automatically reclaimed when memory is needed — it's "free" memory that happens to have useful data in it.

---

Swap:

When physical RAM is exhausted, the kernel can move infrequently-used anonymous pages (heap, stack — not file-backed) to a swap device (partition or file) on disk. This frees RAM for active workloads.

    swapon -s                        # Show active swap devices
    /proc/sys/vm/swappiness          # 0-200, how aggressively to swap (default 60)

Swappiness = 0 means "only swap to avoid OOM." Swappiness = 100 means "swap aggressively to keep more file cache in RAM." Values > 100 (with zswap or zram) compress swap pages in RAM before writing to disk.

Ref: Documentation/admin-guide/mm/concepts.rst
Ref: Documentation/admin-guide/mm/hugetlbpage.rst
Ref: Documentation/admin-guide/mm/transhuge.rst
