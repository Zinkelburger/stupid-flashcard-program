MESI (Modified, Exclusive, Shared, Invalid).

Think of it as a "gossip protocol" for CPUs. When Core 0 wants to change a value in its cache, it sends a high-speed signal across the Internal Bus to all other cores saying: "Hey! I’m changing address 0x1000. If you have a copy of that in your cache, mark it as Invalid right now!"

The next time Core 1 tries to read that flag, it sees the "Invalid" tag and is forced to fetch the fresh version from Core 0 or RAM.
