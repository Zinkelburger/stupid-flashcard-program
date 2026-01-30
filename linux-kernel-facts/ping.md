When you run a standard ping command without specifying a size, it defaults to a payload of 56 bytes (which results in a 64-byte ICMP packet after headers).

    FAIL: ping -c 10 -I if0.30 172.172.149.63 (Packet size: 64 bytes)

    PASS: ping -c 10 -I if0.30 -s 1500 172.172.149.63 (Packet size: 1508 bytes)

    PASS: ping -c 10 -I if0.30 -s 65000 172.172.149.63 (Packet size: 65008 bytes)

Why the small packet is failing

In most networking scenarios, if a large packet works, a small one should too. However, in specific virtualized or hardware-accelerated environments (like those using DPDK, specific NIC drivers, or VLAN tagging on sub-interfaces like if0.30), you might be encountering one of the following:

    VLAN/MTU Padding Issues: Some network stacks or middleboxes (firewalls/switches) have a "Minimum Ethernet Frame Size" requirement (usually 64 bytes). If the tagging on if0.30 is adding or stripping bytes incorrectly, the packet might be arriving "undersized" (a "runt" packet) and getting dropped by the receiving interface.

    Firewall ICMP Filtering: It is common for security rules to be configured to allow fragmented traffic or specific packet sizes while blocking "standard" ICMP echo requests to prevent basic discovery scans.

    MTU/Fragmentation Logic: By forcing a size of 1500 or 65000, you are forcing the OS to fragment the packet. The success of the larger pings suggests that the path is perfectly fine with fragmented IP traffic, but something in the pipeline is discarding unfragmented, standard-sized packets.


Check the MTU of the interface: Run ip link show if0.30. If the MTU is unusually high or low, it might explain the behavior.
Run a Packet Capture: In a separate terminal, run tcpdump -i if0.30 icmp while you run the failing ping.
    If you see Request go out but no Reply, the destination or a firewall is dropping it.
    If you see nothing at all, your local routing table or interface driver is dropping it before it hits the wire.
