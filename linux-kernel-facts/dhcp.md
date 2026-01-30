The file /var/lib/dhclient/dhclient.leases is a local database maintained by the DHCP client (dhclient) on Linux systems. It acts as a history and a "memory" of the IP addresses and network configuration details your computer has received from a DHCP server.
What is inside this file?

When your computer connects to a network, it asks a router (the DHCP server) for an IP address. The router "leases" an address to you for a specific amount of time. The dhclient.leases file records the details of that transaction, including:

    IP Address: The specific address assigned to your machine.

    Lease Duration: When the lease was granted and when it expires.

    Gateway & DNS: The addresses of the router and the domain name servers.

    Server ID: The identity of the DHCP server that provided the address.

Why does Linux keep this file?

This file serves two primary purposes:

    Reboot Recovery: If you restart your computer, dhclient looks at this file to see what IP address it had previously. It will then try to request that same IP address again, which makes network reconnection faster and keeps your IP consistent.

    Fallback: If the DHCP server is down or unreachable when you try to connect, your computer can look at the most recent entry in this file and try to use that configuration as a "best guess" to maintain connectivity.

Can I delete it?

Yes, it is safe to delete. If you delete this file, your computer won't "remember" its previous lease. The next time you connect to a network, dhclient will simply start a fresh negotiation with the router and create a new version of the file.

    Note: On some modern distributions (like Ubuntu or Fedora), you might find this information handled by systemd-networkd or NetworkManager instead, located in /var/lib/NetworkManager/ or similar paths. However, the logic remains the same.
  

Getting a new dhcp lease nowadays:
Most modern desktops (Ubuntu, Fedora, CentOS/RHEL) use NetworkManager. It is often better to use nmcli to avoid conflicting with the system's background services:

    Find your connection name: nmcli connection show

    Restart the connection: sudo nmcli con down "Connection Name" && sudo nmcli con up "Connection Name" (Replace "Connection Name" with the name from step 1, e.g., "Wired connection 1" or "enp3s0")

Method 3: The "Soft" Refresh
If you just want to update settings without fully dropping the connection: sudo nmcli device reapply eth0 (Replace eth0 with your interface name)

To set a static IP address using `nmcli`, you essentially modify an existing network profile to change its "method" from `auto` (DHCP) to `manual`.

### Step 1: Identify your Connection

First, you need to find the name of the connection you want to change.

```bash
nmcli connection show

```

*Look for the **NAME** column (e.g., "Wired connection 1" or "enp3s0").*

---

### Step 2: Set the Static Details

Use the following commands to configure your IP, Gateway, and DNS. Replace `"Your Connection Name"` with the name you found in Step 1.

1. **Set the IP address and Subnet Mask:**
`sudo nmcli con mod "Your Connection Name" ipv4.addresses 192.168.1.100/24`
2. **Set the Default Gateway:**
`sudo nmcli con mod "Your Connection Name" ipv4.gateway 192.168.1.1`
3. **Set the DNS Server(s):**
`sudo nmcli con mod "Your Connection Name" ipv4.dns "8.8.8.8,8.8.4.4"`
4. **Change the method to Manual:**
`sudo nmcli con mod "Your Connection Name" ipv4.method manual`

---

### Step 3: Apply the Changes

The changes are saved but won't take effect until you "up" the connection again.

> **Warning:** If you are doing this over SSH, you may lose your connection if the new IP is different from the old one.

```bash
sudo nmcli con up "Your Connection Name"

```

### Quick Verification

To confirm your settings are now active, run:

```bash
ip addr show

```

You should see your new static IP listed under the relevant interface.

### How to go back to DHCP?

If you ever want to revert to automatic settings, run these two commands:

1. `sudo nmcli con mod "Your Connection Name" ipv4.method auto`
2. `sudo nmcli con up "Your Connection Name"`
