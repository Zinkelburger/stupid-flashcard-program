Standard Bridge Behavior: Bridges are supposed to consume (drop) these frames and not forward them.
The Fix/Feature: Setting /sys/class/net/br0/bridge/group_fwd_mask to 8 (bit 3) tells the bridge "it is okay to forward these specific EAP packets."


Extensible Authentication Protocol (EAP)
