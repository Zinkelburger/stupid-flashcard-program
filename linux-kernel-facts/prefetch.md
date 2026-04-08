Prefetch Instructions in the Linux Kernel

Modern CPUs are "superscalar" — they can execute multiple instructions simultaneously. The CPU can start fetching data from memory in the background while doing other work. Prefetch exploits this by issuing a memory load hint before the data is actually needed.

---

How it works in a driver (e.g., a NIC Rx path):

    prefetch(rx_buf->page);          // Start loading the page into L1 cache
    if (unlikely(!size))             // Meanwhile, check if size is zero
        break;
    // By the time we get here, rx_buf->page is likely already in cache

The CPU can start the prefetch and, while the memory controller is busy, the execution unit simultaneously checks the if condition. By the time the code actually dereferences the page pointer (later in the function), the data is already in cache — saving 50-200 cycles of memory latency.

---

The Kernel's prefetch() Macro:

In the kernel, prefetch(addr) is defined in <linux/prefetch.h>. It compiles to an architecture-specific hint instruction:
  - x86: PREFETCHT0 — fetch data into all cache levels (L1, L2, L3).
  - ARM: PLD (Preload Data).
  - If the architecture doesn't support prefetch, it compiles to nothing (no-op).

The key property: prefetch is a hint, not a demand. If the address is invalid (NULL, unmapped), the CPU silently ignores it — no fault, no crash. This makes it safe to prefetch speculatively.

---

GCC's __builtin_prefetch:

    __builtin_prefetch(addr, rw, locality);

    addr     — The address to prefetch.
    rw       — 0 = prepare for read, 1 = prepare for write.
    locality — Temporal locality hint (0-3):
                 0 = No temporal locality (data won't be reused soon; use non-temporal
                     load, goes into L3 or streaming buffer, doesn't pollute L1/L2).
                 1 = Low temporal locality (L3 cache).
                 2 = Moderate temporal locality (L2 cache).
                 3 = High temporal locality (L1 cache — the default).

The kernel's prefetch() macro typically maps to __builtin_prefetch(addr, 0, 3) — "read, high locality."

prefetchw(addr) maps to __builtin_prefetch(addr, 1, 3) — "write, high locality." Used when you know you're about to modify the data (e.g., before acquiring a spinlock, the kernel prefetches the lock variable for writing).

---

When Prefetch Helps:

- Linked list traversal: Before processing a node, prefetch the next node's pointer. This hides the pointer-chasing latency.
- Network drivers: Prefetch the next packet buffer while processing the current one.
- Hash table lookups: After computing the hash, prefetch the bucket before branching.
- Any code with a predictable access pattern and enough "work" between the prefetch and the actual access to hide latency (~100+ cycles of intervening work).

---

When Prefetch Hurts:

- Small, cache-resident datasets: If the data is already in cache, prefetch instructions waste decode bandwidth.
- Unpredictable access patterns: If you prefetch the wrong address, you pollute the cache with useless data and potentially evict data you actually need.
- Tight loops with tiny bodies: Not enough intervening work between the prefetch and the access — the prefetch hasn't completed yet, so you stall anyway.
- Over-prefetching: Flooding the memory controller with prefetch requests can saturate the memory bus and slow down real loads.

As a rule: profile before adding prefetch. Modern CPUs have aggressive hardware prefetchers that detect sequential and strided access patterns automatically. Software prefetch is most valuable for irregular access patterns (pointer chasing) that the hardware prefetcher can't predict.

Ref: Documentation/core-api/cachetlb.rst, include/linux/prefetch.h
