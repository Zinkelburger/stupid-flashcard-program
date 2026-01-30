u16 prev_ntu = rx_ring->next_to_use & ~0x7;
...
val &= ~0x7;

The Mask: 0x7 is 0111 in binary. The tilde ~ flips the bits, so ~0x7 clears the last 3 bits. This effectively rounds the value down to the nearest multiple of 8.

Many Intel NICs ignore the lower 3 bits of the tail pointer anyway. By only writing when the value crosses an 8-descriptor boundary, the driver reduces overhead without slowing down the hardware.

