Message Signaled Interrupts (MSI)

MSI (first defined in PCI 2.2) permits a device to allocate 1, 2, 4, 8, 16 or 32 interrupts. The device is programmed with an address to write to (this address is generally a control register in an interrupt controller), and a 16-bit data word to identify it. The interrupt number is added to the data word to identify the interrupt.

MSI-X (first defined in PCI 3.0) permits a device to allocate up to 2048 interrupts. The single address used by original MSI was found to be restrictive for some architectures. In particular, it made it difficult to target individual interrupts to different processors, which is helpful in some high-speed networking applications. MSI-X allows a larger number of interrupts and gives each one a separate target address and data word. Devices with MSI-X do not necessarily support 2048 interrupts.
