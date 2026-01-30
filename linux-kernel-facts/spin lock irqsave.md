spin_lock_irqsave(&ndlp->lock, iflags);

A spinlock is a simple integer in memory (usually 0 for unlocked, 1 for locked).

When a CPU wants the lock, it performs an Atomic Swap. It tries to write a 1 and read what was there before.

If it reads a 0, it "wins" the lock and continues.

If it reads a 1, it means someone else has it. The CPU literally runs a while(lock == 1); loop, "spinning" at full speed until the other CPU sets it back to 0.

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
