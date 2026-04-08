RCU (Read-Copy-Update)

RCU is a synchronization mechanism optimized for read-mostly data structures. It allows readers to access shared data with near-zero overhead while writers can update the data concurrently — without locks on the read side.

Ref: Documentation/RCU/whatisRCU.rst, https://www.kernel.org/doc/html/latest/RCU/whatisRCU.html

---

The Core Idea: Split updates into two phases

1. Removal: The writer unlinks the old data from the shared structure (e.g., swaps a pointer using rcu_assign_pointer()). After this, new readers will see the new data, but old readers who were already inside a read-side critical section may still be looking at the old data.

2. Reclamation: The writer waits until all pre-existing readers have finished (a "grace period"), then frees the old data. This wait is done via synchronize_rcu() (blocking) or call_rcu() (async callback).

The key insight: because modern CPUs guarantee that aligned pointer writes are atomic, a reader will always see either the old pointer or the new pointer — never a half-written one. This means readers never need locks.

---

The Core API (5 functions):

rcu_read_lock()      — Marks the beginning of an RCU read-side critical section.
                       The reader promises not to sleep (in classic RCU).
                       Disabling preemption, interrupts, or BHs also implicitly
                       enters an RCU read-side critical section.

rcu_read_unlock()    — Marks the end of the critical section.

synchronize_rcu()    — Called by the writer/reclaimer. Blocks until all
                       pre-existing RCU read-side critical sections have completed.
                       Does NOT wait for sections that begin after it is called.

call_rcu(head, func) — Async version of synchronize_rcu(). Registers a callback
                       that will be invoked after the grace period ends. Used when
                       the writer cannot block (e.g., in interrupt context).

rcu_assign_pointer(p, v) — Writes a new pointer value with the necessary memory
                           barriers so that readers using rcu_dereference() will
                           see a fully initialized structure.

rcu_dereference(p)   — Reads a pointer that was published with rcu_assign_pointer().
                       Includes the required memory barrier (on architectures that
                       need one, like Alpha).

---

What is a "grace period"?

A grace period is the span of time during which the kernel can guarantee that every CPU has passed through at least one "quiescent state."

A quiescent state is any point where the CPU is provably not inside an RCU read-side critical section:
  - A context switch (the CPU ran the scheduler)
  - Returning to userspace
  - The idle loop

Once every CPU has passed through a quiescent state, the kernel knows that no reader could possibly still hold a reference to data that was unlinked before the grace period began. The old data can now be safely freed.

---

Read-side overhead:

In non-preemptible kernels (CONFIG_PREEMPT_RCU=n), rcu_read_lock() compiles down to essentially nothing — it just disables preemption (a single increment to a per-CPU counter). There are no atomic operations, no memory barriers, no cache-line bouncing. This is why RCU is so much faster than rwlocks for readers.

In PREEMPT_RCU kernels, rcu_read_lock() is slightly heavier (it manipulates a per-task counter) but still far cheaper than any lock.

---

RCU vs rwlock vs refcount:

Mechanism        | Read overhead     | Write overhead   | Reader can sleep?
-----------------|-------------------|------------------|------------------
rwlock           | Atomic op + barrier | Exclusive lock  | No
RCU              | Near zero         | Must wait grace period | No (classic RCU)
SRCU             | Counter increment | Must wait grace period | Yes
Reference count  | Atomic inc/dec    | Free when count=0| Yes

RCU's trade-off: readers are extremely cheap, but writers must wait for a full grace period before freeing old data. This is ideal when reads vastly outnumber writes (routing tables, module lists, dcache, etc.).

---

The "Golden Rule": RCU readers must not sleep.

If a reader sleeps, it might hold a reference across a context switch. The kernel would interpret that context switch as a quiescent state and conclude the grace period has ended — but the sleeping reader is still using the old data. Freeing it would cause a use-after-free.

Exception: SRCU (Sleepable RCU) allows readers to sleep by using an explicit counter (srcu_read_lock returns an index that must be passed to srcu_read_unlock). The kernel can't rely on context switches to detect quiescent states, so it uses the counter instead.

---

Typical usage pattern (linked list removal):

    // Writer side:
    spin_lock(&list_lock);
    list_del_rcu(&entry->node);     // Unlink (readers still see old)
    spin_unlock(&list_lock);

    synchronize_rcu();              // Wait for all readers to finish
    kfree(entry);                   // Safe to free now

    // Reader side (no lock needed):
    rcu_read_lock();
    list_for_each_entry_rcu(entry, &head, node) {
        // Use entry safely
    }
    rcu_read_unlock();

---

RCU in practice — where it's used in the kernel:

- Routing table lookups (net/ipv4/fib_trie.c)
- The dcache (filesystem directory entry cache)
- Module list (tracking loaded kernel modules)
- PID lookup (finding tasks by PID)
- Network namespace management
- Device mapper table swaps (dm_get_live_table / dm_put_live_table)

Ref: Documentation/RCU/whatisRCU.rst
Ref: Documentation/RCU/checklist.rst
