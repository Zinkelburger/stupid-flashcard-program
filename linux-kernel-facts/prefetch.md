Modern CPUs are "superscalar," meaning they can do multiple things at once.

    Instruction A: Tell the memory controller to go get rx_buf->page (Prefetch).

    Instruction B: Check if size is zero (The if statement).

The CPU can start the prefetch and, while the memory hardware is busy fetching that data, the CPU's execution unit simultaneously checks the size. By the time the code actually needs to touch the page (later in the function), the data is already arriving in the cache.

