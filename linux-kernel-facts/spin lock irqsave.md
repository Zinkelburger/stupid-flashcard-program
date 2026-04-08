spin_lock_irqsave(&ndlp->lock, iflags);

A spinlock is a lock where the waiting CPU "spins" (busy-waits) rather than sleeping. If the lock is held by another CPU, the waiting CPU loops until it becomes available.

Early Linux used a simple test-and-set integer (0/1). This was replaced by ticket spinlocks (like a bakery number — each waiter gets a ticket and waits for their number to be called, guaranteeing FIFO fairness). Modern kernels (4.2+) use queued spinlocks (qspinlocks), where each waiting CPU spins on its own per-CPU variable rather than hammering a shared cache line. This dramatically reduces cache-line bouncing on large NUMA systems.

The core idea is the same: the lock holder executes the critical section while everyone else busy-waits. Spinlocks are only appropriate when the critical section is very short (a few instructions) because the waiting CPUs are burning cycles doing nothing useful.

Irqsave:
If the cpu is interrupted:
If that Interrupt handler also tries to grab ndlp->lock, it will see the lock is held and start spinning.
But since the Interrupt is running on the same CPU that was holding the lock, the CPU will never return to the original code to release it.
= Deadlock

---
When you call spin_lock_irqsave, the kernel does three things in order:

Captures the Current State: It reads the CPU's current interrupt state (were interrupts already off, or were they on?).

Saves to iflags: It stuffs that state into your iflags variable.

Disables Interrupts: It executes a command (like cli on x86) to turn off interrupts.
