QinQ (802.1ad) — Double VLAN Tagging

QinQ (standardized as IEEE 802.1ad) stands for "802.1Q in 802.1Q." It is a technique for stacking two VLAN tags on a single Ethernet frame, creating a "VLAN tunnel."

---

Standard VLAN (802.1Q): A normal Ethernet frame gets a single 4-byte tag inserted after the source MAC address. This tag contains:
  - TPID (Tag Protocol Identifier): 0x8100 — identifies the frame as 802.1Q tagged.
  - VLAN ID: 12 bits, allowing up to 4094 VLANs (0 and 4095 are reserved).
  - Priority (PCP): 3 bits for QoS priority (0-7).

QinQ (802.1ad): A second VLAN tag is prepended to the existing one. The frame now has two tags:
  - Outer tag (S-Tag / Service Tag): TPID = 0x88A8. Added by the service provider.
  - Inner tag (C-Tag / Customer Tag): TPID = 0x8100. The customer's original VLAN tag.

The outer tag is sometimes called the "provider tag" and the inner tag the "customer tag."

---

Frame Structure:

    | Dst MAC | Src MAC | S-Tag (0x88A8) | C-Tag (0x8100) | EtherType | Payload | FCS |
                        |← Provider VLAN →|← Customer VLAN →|

Each tag adds 4 bytes. A standard 1500-byte MTU frame with double tags becomes 1508 bytes. The infrastructure must support this larger frame size (baby giant frames) or the MTU must be reduced.

---

Why QinQ Exists:

The Problem: 802.1Q only supports 4094 VLANs. An ISP serving thousands of enterprise customers would quickly run out.

The Solution: With QinQ, the ISP uses the outer S-Tag to identify each customer. Each customer can independently use all 4094 inner C-Tag VLANs without conflicting with anyone else. Total theoretical VLAN combinations: 4094 × 4094 ≈ 16.7 million.

Use cases:
  - ISP/carrier networks: Each customer gets one S-Tag. Their internal VLANs (C-Tags) are transparent to the ISP.
  - Data center interconnects: Tunnel one data center's VLAN space over another's infrastructure.
  - Metro Ethernet: Extend customer L2 domains across a provider network.

---

Configuring QinQ in Linux:

First, create the outer (provider) VLAN interface:

    ip link add link eth0 name eth0.100 type vlan id 100 protocol 802.1ad
    # protocol 802.1ad sets the TPID to 0x88A8 (S-Tag)

Then, create the inner (customer) VLAN interface on top of it:

    ip link add link eth0.100 name eth0.100.200 type vlan id 200
    # Default protocol is 802.1Q (0x8100), which is correct for the C-Tag

Bring them up and assign an IP:

    ip link set eth0.100 up
    ip link set eth0.100.200 up
    ip addr add 192.168.1.1/24 dev eth0.100.200

Traffic sent out eth0.100.200 will have two VLAN tags: outer S-Tag 100 (0x88A8) and inner C-Tag 200 (0x8100).

---

MTU Considerations:

Each VLAN tag adds 4 bytes. Double tagging adds 8 bytes total. If the underlying network has a 1500-byte MTU, the maximum payload for QinQ frames is 1492 bytes (1500 - 8).

To avoid fragmentation, increase the physical interface MTU:

    ip link set eth0 mtu 1508    # Or higher, like 9000 for jumbo frames

All switches and routers in the path must also support the larger frame size.

Ref: IEEE 802.1ad, IEEE 802.1Q
