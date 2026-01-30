PID,Name,Role
0,Idle Task,"Internal kernel ""idle"" loop. Never assigned to an app."
1,init / systemd,"The ""Mother"" of all processes. The first thing that starts."
2,kthreadd,The parent of all other kernel threads.
>2,Userspace Apps,"Your browser, terminal, and nbd-client all live here."

If the code sees nbd->pid == 0, it knows for an absolute fact that no userspace process is currently managing that block device. If it were still using the old struct task_struct *task_recv, it would have to check if the pointer is NULL. Checking if (pid == 0) is functionally the same, but much lighter and avoids the memory "dangling" issues we discussed earlier.

When the nbd-client process starts the device, the kernel saves that process's real ID (like 1422) into that field. When it disconnects, it clears it back to 0.
