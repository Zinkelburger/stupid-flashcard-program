Feature,"writel(val, addr)","write(fd, buf, count)"
Who uses it?,Device Drivers (Kernel),Applications (User space)
Data Size,Always 32 bits (4 bytes),Variable (Whatever you specify)
Destination,A hardware register (Pointer),A file descriptor (Integer)
Speed,Extremely fast (Single instruction),Slow (Requires a context switch)
End Result,Changes hardware state,Saves data to a file/socket

You might see some old code use a simple pointer dereference like *addr = val; instead of writel. This is dangerous for three reasons:

    Compiler Optimization: The compiler might think, "You just wrote to this address twice, I'll delete the first one." But for hardware, writing twice might be a specific command (like a toggle). writel uses volatile to prevent the compiler from "cleaning up" your code.

    CPU Reordering: Modern CPUs execute instructions out of order to save time. writel forces the CPU to finish its work in the correct sequence.

    Portability: writel works the same way whether you are on an Intel (x86) chip or an ARM chip. If you write your own *addr = val, your driver might break when moved to a different computer architecture.