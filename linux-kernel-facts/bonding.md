Think of bonding as the "Voltron" of networking. You take multiple physical network cables and combine them so they appear to the computer as a single, super-fast, or ultra-reliable cable.
bonding (also called Link Aggregation) is the process of grouping multiple physical network interfaces (called slaves, like eth0 and eth1) into one single logical interface (called the master, like bond0).

In the Linux kernel, the bonding driver determines which slave to use by calling a hashing function (historically bond_xmit_hash). When you set xmit_hash_policy to layer3+4, the driver attempts to extract:
Layer 3: Source and Destination IP addresses.
Layer 4: Source and Destination Ports.
The Problem: ICMP does not have "ports." For layer3+4 to be effective for ICMP, the kernel must treat the ICMP Identifier (ID) as the "port" equivalent to provide entropy. If it doesn't, every ping from Host A to Host B will result in the exact same hash, pinning all traffic to one interface (as seen in your veth1 capture).

The bonding driver (bonding.ko) is the kernel module that acts as the "brain" for this group. When the kernel wants to send a packet out of bond0, the bonding driver intercepts that packet and makes a split-second decision based on its configuration (the mode):
Failover Mode: "Is the main cable broken? If yes, use the backup."
Load Balancing Mode (XOR/LACP): "I have two cables. I'll use a math formula (the hash) to decide which cable this specific packet should take so I don't overload one of them."

The Socket Layer: Your application opens a socket and sends data to an IP address. The socket doesn't know (or care) that bonding exists.
The Routing Layer: The kernel looks at the destination IP and says, "To get there, I need to send this out of the bond0 interface."
The Bonding Driver: This is where the magic happens. The kernel hands the packet to the bonding driver. The driver looks at the packet headers (IPs, MACs, or in your case, ICMP IDs) and chooses a physical NIC (veth0 or veth1) to actually transmit the bits.

Mode,Name,Switch Support,Best For...
0,balance-rr,Required (Static),Maximum raw throughput (Internal/Private networks).
1,active-backup,None,High availability / Redundancy (Standard).
2,balance-xor,Required (Static),Basic load balancing without LACP.
3,broadcast,Required (Static),"Specific niche cases (e.g., financial/high-stakes data)."
4,802.3ad (LACP),Required (LACP),Enterprise-grade load balancing and redundancy.
5,balance-tlb,None,Outbound load balancing without switch config.
6,balance-alb,None,Full load balancing without switch config.
