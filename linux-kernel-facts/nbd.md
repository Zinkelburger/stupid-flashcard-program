sock_xmit
noreclaim_flag = memalloc_noreclaim_save();

memalloc_noreclaim_save(): This is critical. NBD is a storage driver. When the system is low on memory, it tries to free memory by writing data to disk. If that "disk" is NBD, and NBD tries to allocate memory to send the data, you get a deadlock (waiting for memory to free memory).
This line sets a flag on the current task saying: "Do not try to reclaim memory while I am running."

sock->sk->sk_allocation = GFP_NOIO | __GFP_MEMALLOC;
__GFP_MEMALLOC will allow the allocation to disregard the watermarks

To fix your SELinux issue, the kernel developers wrap the if (send) block like this:
C

/* Pseudocode of the fix */
old_cred = override_creds(prepare_kernel_cred(NULL)); // Wear the Kernel Hat

if (send)
    result = sock_sendmsg(sock, &msg);
else
    result = sock_recvmsg(sock, &msg, msg.msg_flags);

revert_creds(old_cred); // Take off the hat

By adding those lines, the SELinux check inside sock_sendmsg sees the Kernel as the caller, not your user-space program, and the I/O error disappears.


---
1. The "Registration" Analogy

Think of the struct nbd_device like a rental car and the nbd->pid field like the registration form in the glovebox.

    The Driver (nbd.c): Is the rental car agency. It’s always "running" in the background as part of the kernel.

    The Device (/dev/nbd0): Is the car itself.

    The PID (nbd->pid): Is where the agency writes down who is currently driving the car.

When you run nbd-client, that program (a userspace process) tells the kernel: "I want to use nbd0." The kernel then looks at that process, sees its ID (e.g., PID 500), and writes 500 into nbd->pid.

When the nbd-client finishes or the connection drops, the driver erases the "registration form."

    Setting nbd->pid = 0 doesn't mean the driver stops existing.

    It just means: "No one is currently using this specific NBD device."
    
---