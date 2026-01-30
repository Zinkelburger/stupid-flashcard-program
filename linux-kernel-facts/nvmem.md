Historically, every driver did indeed manage its own EEPROM (Electrically Erasable Programmable Read-Only Memory) in its own unique, "wild west" way. Today, modern operating systems (especially Linux) use standardized subsystems to handle this, though some exceptions remain for complex devices like network cards.

NVMEM subsystem

/sys/bus/nvmem/

How it works: The driver acts as a "provider." It simply tells the OS, "I have memory here, it is this big, and here is how you read a byte."

The Benefit: The OS takes over the management. It creates a standardized file (usually in /sys/bus/nvmem/) that looks like a normal file to the user. You can read/write to the EEPROM just like you would edit a text document, without needing to know which specific driver is running.


Driver Type,How it manages EEPROM,Is it unique?
I2C / SPI Sensors,Uses NVMEM Subsystem,"No. It uses a shared, standard kernel driver (like at24)."
Network Cards,Uses ethtool API,"Internally yes, Externally no. The code is unique, but the user interface is standard."
Embedded / Legacy,Direct Register Access,Yes. The driver manually flips bits to talk to the chip.