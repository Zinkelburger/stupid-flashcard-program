RCU is a synchronization trick used in the Linux kernel to allow one person to "write" to a list while dozens of others "read" it at the exact same time—without using expensive locks.

Removal: The "Writer" unlinks the data from the list. New readers can't find it, but old readers who were already looking at it still can.
Grace Period: The "Writer" waits. It waits until every CPU core has gone through a "quiescent state" (like a context switch), which proves that no one could possibly still be holding a reference to that old data.
Reclamation: Once the wait is over, the memory is finally freed.
