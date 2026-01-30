QinQ (standardized as IEEE 802.1ad) stands for "802.1Q in 802.1Q." It is a technique used to provide multiple VLANs over a single VLAN connection.

Standard VLAN (802.1Q): A normal network packet has one "tag" that identifies which VLAN it belongs to.
QinQ (802.1ad): A second tag is added to the packet. The original tag is called the C-Tag (Customer Tag), and the new, outer tag is called the S-Tag (Service Tag).

Service Providers: An Internet Service Provider (ISP) can use one "Outer" VLAN to carry a specific customer's traffic. Inside that tunnel, the customer can run their own "Inner" VLANs (up to 4096 of them) without interfering with the ISP's other customers.

