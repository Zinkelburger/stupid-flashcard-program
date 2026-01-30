Write Memory Barrier (wmb)

Let's assume a serial port, where you have to write bytes to a certain address. The serial chip will then send these bytes over the wires. It is then important that you do not mess up the writes - they have to stay in order or everything is garbled.
But the following is not enough:

   *serial = 'h';
   *serial = 'e';
   *serial = 'l';
   *serial = 'l';
   *serial = 'o';

Because the compiler, the processor, the memory subsystems and the buses in between are allowed to reorder your stores as optimization (believe me, yes they are and yes, they do).

so you'll have to add code that will ensure the stores do not get tangled up. That's what, amongst others, the wmb() macro does: prevent reordering of the stores. 

Note that just making the serial pointer volatile is not enough: while it ensures the compiler will not reorder, the other mechanisms mentioned can still wreak havoc.
https://stackoverflow.com/questions/30236620/what-is-wmb-in-linux-driver