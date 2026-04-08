Right after you get your IP from DHCP, you might see one more weird ARP packet where your computer asks, "Who has 192.168.1.50?" Wait, that's your IP. Why are you asking? This is a Gratuitous ARP. Your computer is checking to make sure no one else on the network is accidentally using the same IP address. If someone replies, you’ll get an "IP Address Conflict" error.

---
The "Who has" Request

Your computer broadcasts a message: "Who has 192.168.1.1? Tell 192.168.1.50."

    Target IP: The Router's IP.

    Target MAC: 00:00:00:00:00:00 (This is a placeholder because you don't know it yet).

    Destination MAC: ff:ff:ff:ff:ff:ff (Broadcast to everyone).

The "Is at" Reply

The Router sees its own IP in that broadcast and replies directly to you: "192.168.1.1 is at a4:bb:c8:12:34:56."

    Sender MAC: This is the "golden" info your computer was looking for.

    Your computer now saves this in its ARP Cache so it doesn't have to ask again for a while.

Look for "Opcode" (1 for request, 2 for reply).

---
A. The Server Checks for Squatters (ARP Probe)

Before the Router gives you an IP (like 192.168.1.50), it might send out an ARP request: "Who has 192.168.1.50?" It does this to make sure no one else is already using that address. If nobody answers, it's safe to give to you.

---
Once DHCP is finished, ARP becomes the main character. If you want to go to google.com, your computer realizes Google isn't in your house; it's outside. So, you need to talk to your Default Gateway (your router).

    Your computer looks at the Router IP (which it got from the DHCP ACK packet earlier).

    It checks its ARP Cache. If it's empty, it sends an ARP Request: "Who has 192.168.1.1? Tell 192.168.1.50."

    The Router replies with its MAC address.

    Now you can finally send data.