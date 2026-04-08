MESI (Modified, Exclusive, Shared, Invalid)

MESI is the cache coherence protocol used by most modern multi-core CPUs (Intel, AMD, ARM with variants). It ensures that when one CPU core modifies a value in its cache, all other cores see a consistent view of memory. Without it, multi-core systems would read stale data constantly.

---

The Four States:

Every cache line (typically 64 bytes) is tagged with one of four states:

Modified (M):
  - This core has changed the data. The cache line is "dirty."
  - No other core has a copy.
  - Main memory is OUT OF DATE — this cache holds the only correct copy.
  - If another core wants this data, this core must write it back to memory
    (or snoop-forward it) before the other core can read it.

Exclusive (E):
  - This core has the only copy, and it matches main memory exactly.
  - No other core has a copy.
  - The core can silently promote this to Modified (write without telling anyone)
    because it knows nobody else is looking.

Shared (S):
  - Multiple cores may have a copy of this cache line.
  - All copies match main memory.
  - The data is read-only in this state. To write, the core must first
    invalidate all other copies (transition to M).

Invalid (I):
  - The cache line is empty or stale. The data here is garbage.
  - The core must fetch from memory or snoop another core's cache before using it.

---

State Transitions (the important ones):

Read miss (I → E or S):
  Core needs data it doesn't have. It sends a read request on the bus/interconnect.
  If no other core has it → fetched from memory, state = Exclusive.
  If another core has it (Shared or Modified) → gets a copy, state = Shared.
  If another core had it Modified → that core writes back first, both become Shared.

Write hit (E → M):
  Core already has exclusive copy. It writes silently. State becomes Modified.

Write hit (S → M, "RFO — Request For Ownership"):
  Core has a shared copy but wants to write. It broadcasts an invalidation to
  all other cores holding copies. They mark their lines Invalid. Now this core
  is the sole owner and transitions to Modified.

Write miss (I → M):
  Core doesn't have the data at all. It issues a Read-For-Ownership (RFO) that
  fetches the line AND invalidates everyone else. Expensive — a full round trip.

Snoop invalidation (S or E → I):
  Another core wants to write, so it sends an invalidate. This core marks its
  copy as Invalid.

---

Cache Line Bouncing:

When two cores repeatedly write to the same cache line, the line "bounces" between them via constant RFO → Invalidate cycles. Each transition takes ~40-100ns (depending on whether cores share an L3 or cross NUMA boundaries). This kills performance on parallel workloads.

This is why the kernel's queued spinlocks (qspinlocks) have each waiter spin on its own per-CPU variable: to avoid all waiters hammering the same cache line.

---

False Sharing:

False sharing occurs when two unrelated variables happen to live on the same 64-byte cache line. Core 0 writes variable A, Core 1 writes variable B — but since they share a cache line, every write invalidates the other core's copy, even though they're touching completely different data.

The kernel fights this with:
  - ____cacheline_aligned_in_smp: Forces a struct field to start on a cache line boundary.
  - __cacheline_aligned: Same, but unconditional.
  - Per-CPU variables (DEFINE_PER_CPU): Each CPU gets its own copy on its own cache line.

---

MESI and Memory Barriers:

CPUs can reorder loads and stores for performance. MESI guarantees cache coherence (eventually all cores see the same value), but it does NOT guarantee ordering of when those values become visible.

That's where memory barriers come in:
  - wmb() (write memory barrier): All stores before the barrier are visible before stores after it.
  - rmb() (read memory barrier): All loads before the barrier complete before loads after it.
  - mb() (full memory barrier): Both.
  - smp_wmb(), smp_rmb(), smp_mb(): SMP-only variants (compile to nothing on UP).

MESI provides the coherence mechanism; memory barriers provide the ordering guarantees. They are complementary.

Ref: Documentation/memory-barriers.txt
