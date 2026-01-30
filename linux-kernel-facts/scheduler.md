EEVDF stands for Earliest Eligible Virtual Deadline First. It combines two concepts to decide which process runs next:
Eligibility (The "Fairness" part): A task is only "eligible" if it has not already used more than its fair share of the CPU. If a task has been hogging the processor, its "lag" becomes negative, and it has to wait until its share catches up.
Virtual Deadline (The "Latency" part): Among all eligible tasks, the scheduler picks the one with the earliest virtual deadline. This deadline is calculated based on the task's requested time slice.

Feature,EEVDF (General Purpose),Deadline Scheduler (SCHED_DEADLINE)
Goal,Proportional fairness + low latency.,Guaranteed execution within a time window.
Failure,"A ""missed"" deadline just causes a tiny stutter.",A missed deadline can be a system failure.
Complexity,Automatic; handles thousands of desktop tasks.,Requires manual parameters (Runtime/Period).
Usage,"Default for Browsers, DEs, Games.","Hard Real-Time (Audio, Kernels, Robotics)."

One property of the EEVDF scheduler that can be seen in the above tables is that the sum of all the lag values in the system is always zero. 

If a background task asks for a 10ms slice, its virtual deadline is 10ms away.
If your Mouse/GUI asks for a 1ms slice, its virtual deadline is much sooner.
The Result: Even if the background task is "fair," the GUI task jumps to the front of the line because its deadline is earlier.

In older kernels, we used "Latency Nice" patches. In modern EEVDF kernels, this is being integrated into the sched_setattr system call. Applications can now tell the kernel: "I don't need much total CPU time, but when I do need it, I need it right now."

Nice Levels: Setting a process to a lower nice value (e.g., -20) gives it more "weight," which indirectly shortens its virtual deadline.

A High "Nice" value (e.g., +19): You are being very nice to other users. You are saying, "Please, go ahead of me. I'll just take the leftover CPU scraps whenever you're done."

A Low "Nice" value (e.g., -20): You are being not nice (selfish). You are saying, "I am more important. Move aside; I'm taking the CPU right now."

Every nice level corresponds to a "weight." For every 1-point drop in nice value, a process gets roughly 10% more CPU time than its peers.

The Debt System: Every time a task runs, it builds up "debt" (called negative lag).
The "You Must Be This Fair" Line: A task is only eligible to run if its lag is ≥0. This means it hasn't already taken more than its fair share of the CPU

The kernel doesn't actually allow preemption to happen at an infinite frequency. There is a setting called sched_min_granularity_ns (usually around 1.5ms to 3ms).

Even if a GUI task wakes up and says, "My deadline is now!", the scheduler will generally let the current "throughput" task finish its minimum window. This ensures that the CPU caches stay "warm" and the system doesn't waste all its cycles on the overhead of switching between tasks.

The core of CFS was simple: Total CPU time should be equal. * Every process had a variable called vruntime (virtual runtime).
If Task A ran for 10ms, its vruntime increased by 10.
The scheduler always picked the task with the lowest vruntime to run next.

Imagine your Mouse Driver. It "sleeps" 99% of the time. When you move the mouse, its vruntime is very low because it hasn't used the CPU in a while. CFS would let it jump to the front of the line.
    The "Icky" Part: To prevent these "sleepers" from jumping too far ahead and hogging the CPU when they wake up, CFS had hundreds of lines of "heuristics" (educated guesses).
    These guesses often failed. Sometimes a game would stutter because the scheduler didn't "guess" correctly that the game needed a tiny slice of CPU right now to render a frame.

EEVDF:  The solution is to decay a sleeping task's lag over virtual run time instead. The implementation of this idea in the patch set is somewhat interesting. When a task sleeps, it is normally removed from the run queue so that the scheduler need not consider it. With the new patch, instead, an ineligible process that goes to sleep will be left on the queue, but marked for "deferred dequeue". Since it is ineligible, it will not be chosen to execute, but its lag will increase according to the virtual run time that passes. Once the lag goes positive, the scheduler will notice the task and remove it from the run queue.
The result of this implementation is that a task that sleeps briefly will not be able to escape a negative lag value, but long-sleeping tasks will eventually have their lag debt forgiven. Interestingly, a positive lag value is, instead, retained indefinitely until the task runs again. 