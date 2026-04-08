LRO (Large Receive Offload) is a performance feature. Instead of the CPU processing every single 1500-byte packet that arrives, the mlx5 hardware waits for a burst of packets from the same stream, stitches them together into one "giant" packet (up to 64KB), and sends that single large buffer to the kernel. This significantly reduces CPU overhead.

1460 bytes: The standard maximum segment size (MSS) for a TCP packet on a 1500 MTU network.
1448 bytes: The MSS when TCP "Timestamps" are enabled (which adds 12 bytes to the header).
1428 bytes: Often seen in environments with extra encapsulation or specific TCP options.

