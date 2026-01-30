"RCU Golden Rule": Readers must be fast, and they must not block.
An RCU reader cannot sleep. You cannot call any function that might wait for a disk, a user input, or a mutex while inside an RCU critical section
If a reader sleeps, it triggers a context switch, and the kernel would mistakenly think it's safe to delete the data the reader was still using.

An SRCU reader can sleep. You can hold an SRCU "lock," go to sleep waiting for a hard drive to respond, and wake up later still holding that lock safely.
Because the kernel can't rely on context switches to know when you're done, it uses an explicit counter. This is why you saw srcu_idx in your code. You are "checking out" a specific version of the table and must "check it back in" manually.

- Slower than an RCU, faster than a traditional Mutex
