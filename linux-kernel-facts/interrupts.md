Interrupts and Bottom Halves

When a hardware device needs attention (a packet arrived, a disk I/O completed, a timer fired), it raises an interrupt. The CPU immediately stops what it's doing, saves state, and runs the interrupt handler. Understanding how Linux processes interrupts is fundamental to understanding kernel performance.

---

Top Half (Hardirq):

The interrupt handler registered with request_irq() runs in "hardirq" context. While it runs:
  - All interrupts on the local CPU may be disabled (depending on the handler flags).
  - The handler cannot sleep, cannot allocate memory with GFP_KERNEL, cannot acquire a mutex.
  - It must complete as fast as possible — every microsecond spent here is a microsecond no other interrupt can be serviced.

The top half should do the absolute minimum: acknowledge the interrupt in the hardware, read any status registers, stash the data somewhere, and schedule a "bottom half" to do the real work.

---

Bottom Half Mechanisms:

The "bottom half" is where the real processing happens — outside of hardirq context, with fewer restrictions. Linux has evolved several mechanisms:

1. Softirqs (highest priority, fixed at compile time):

Softirqs are statically defined in the kernel (you can't add new ones from a module). There are exactly 10 softirq types:

    HI_SOFTIRQ          — High-priority tasklets
    TIMER_SOFTIRQ        — Timer callbacks
    NET_TX_SOFTIRQ       — Network transmit
    NET_RX_SOFTIRQ       — Network receive
    BLOCK_SOFTIRQ        — Block layer completion
    IRQ_POLL_SOFTIRQ     — IRQ polling
    TASKLET_SOFTIRQ      — Normal tasklets
    SCHED_SOFTIRQ        — Scheduler load balancing
    HRTIMER_SOFTIRQ      — High-resolution timers
    RCU_SOFTIRQ          — RCU callbacks

Softirqs run on the CPU that raised them, with interrupts enabled but in atomic context (cannot sleep). The same softirq can run simultaneously on different CPUs — the handler must be fully reentrant and use per-CPU data or locking.

Softirqs are processed at the end of a hardirq or by the ksoftirqd/N kernel thread if there are too many pending.

2. Tasklets (built on softirqs, easier to use):

Tasklets are dynamically allocatable (unlike softirqs) and are guaranteed to run on only one CPU at a time — the same tasklet will never run concurrently on two CPUs. This makes them much easier to write safely.

    DECLARE_TASKLET(name, function);
    tasklet_schedule(&name);    — Schedule the tasklet to run soon.

Tasklets run in softirq context (cannot sleep). They are being phased out in favor of threaded IRQs and workqueues.

3. Workqueues (most flexible, can sleep):

A workqueue schedules work to run in process context — meaning the handler CAN sleep, allocate memory, acquire mutexes, and do I/O. Work items are executed by kernel worker threads (kworker/N:M).

    DECLARE_WORK(name, function);
    schedule_work(&name);                — Queue on the system workqueue.
    schedule_delayed_work(&name, delay); — Queue after a delay.

The system workqueue (system_wq) is shared. For work that may block for a long time, drivers can create dedicated workqueues with alloc_workqueue().

Workqueue flavors:
  - Bound (default): Worker threads are pinned to specific CPUs for cache locality.
  - Unbound: Worker threads can run on any CPU. Better for long-running work that doesn't need cache affinity.

    cat /proc/sys/kernel/workqueue/   # Workqueue tuning knobs

Ref: Documentation/core-api/workqueue.rst

4. Threaded IRQs (the modern approach):

Instead of using a tasklet or workqueue, a driver can request a "threaded IRQ" where the top half is minimal and the bulk of the work runs in a dedicated kernel thread:

    request_threaded_irq(irq, hardirq_handler, thread_fn, flags, name, data);

The hardirq_handler runs in interrupt context (fast acknowledge), then wakes the thread_fn which runs in process context (can sleep). Each threaded IRQ gets its own irq/N-name kernel thread.

This is the preferred model for new drivers on PREEMPT_RT (real-time) kernels because all interrupt handling becomes schedulable and predictable.

---

NAPI (New API) — Network-Specific:

NAPI is a special interrupt mitigation strategy for network drivers. Under heavy traffic, a NIC can generate thousands of interrupts per second — each one context-switching the CPU. NAPI combines interrupts with polling:

1. The NIC raises an interrupt when the first packet arrives.
2. The interrupt handler disables further interrupts from the NIC and schedules a NAPI poll.
3. The NAPI poll function runs in softirq context and processes packets in a loop (up to a budget, typically 64 packets).
4. If more packets keep arriving, the poll keeps running (no interrupts needed).
5. When the queue is empty, NAPI re-enables interrupts and waits for the next packet.

This turns "interrupt per packet" into "interrupt per burst," dramatically reducing CPU overhead under load. Idle networks still use interrupts for low latency on the first packet.

    /sys/class/net/eth0/napi_defer_hard_irqs  — Defer re-enabling interrupts for N more polls
    /sys/class/net/eth0/gro_flush_timeout     — How long to wait before flushing GRO

---

Interrupt Context Restrictions Summary:

Context          | Can sleep? | Can use GFP_KERNEL? | Can acquire mutex? | Runs on
-----------------|-----------|--------------------|--------------------|--------
Hardirq          | No        | No (GFP_ATOMIC)    | No                 | Interrupt CPU
Softirq/Tasklet  | No        | No (GFP_ATOMIC)    | No                 | Same CPU
Workqueue        | Yes       | Yes                | Yes                | kworker thread
Threaded IRQ     | Yes       | Yes                | Yes                | irq/N thread
Process context  | Yes       | Yes                | Yes                | Any

---

Viewing Interrupts:

    cat /proc/interrupts     # Per-CPU interrupt counts for each IRQ line
    cat /proc/softirqs       # Per-CPU softirq counts
    watch -n 1 cat /proc/softirqs   # Watch softirq rates in real time
    mpstat -I ALL 1          # Per-CPU interrupt statistics

Ref: Documentation/core-api/genericirq.rst
