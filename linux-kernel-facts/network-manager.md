NetworkManager

NetworkManager (NM) is the default network configuration daemon on most modern Linux desktops and servers (Fedora, RHEL, Ubuntu, CentOS). It manages all network interfaces — Ethernet, Wi-Fi, VPN, bridges, bonds, VLANs — through a D-Bus API and integrates with systemd.

---

The Concept: Managed vs Unmanaged

When an interface is "managed" by NetworkManager, NM controls its IP configuration, routing, DNS, and state. When "unmanaged," NM completely ignores the device — it won't assign an IP, bring it up, or interfere with it.

An "unmanaged" interface is like a telephone plugged into the wall with no service plan. The hardware exists, but no one has configured it.

Interfaces can become unmanaged because:
  - They are explicitly listed in /etc/NetworkManager/conf.d/ as unmanaged.
  - They were created by a script or container runtime that NM doesn't recognize.
  - The [keyfile] unmanaged-devices setting matches the interface.

---

Profiles (Connections) vs Devices:

This is a key distinction:

Device: A physical or virtual network interface (e.g., enp3s0, wlan0, bond0). It exists in the kernel.

Connection Profile: A named configuration stored on disk that tells NM how to configure a device. A profile includes: IP method (DHCP or static), IP address, gateway, DNS, VLAN, bond settings, 802.1X credentials, etc.

A single device can have multiple profiles (e.g., "Office DHCP" and "Home Static"), but only one profile can be active on a device at a time.

Profiles are stored as keyfiles in /etc/NetworkManager/system-connections/ (modern format) or as ifcfg files in /etc/sysconfig/network-scripts/ (legacy RHEL format, deprecated in RHEL 9+).

---

nmcli — The Command-Line Tool:

Viewing:
    nmcli device status                     # List all devices and their states
    nmcli connection show                   # List all connection profiles
    nmcli connection show "Wired connection 1"  # Show details of a specific profile
    nmcli device show enp3s0                # Show all config for a specific device
    nmcli general status                    # Overall NM status

Managing connections:
    nmcli con up "My Profile"               # Activate a connection profile
    nmcli con down "My Profile"             # Deactivate
    nmcli device reapply enp3s0             # Re-read and apply settings without full down/up

Modifying:
    nmcli con mod "My Profile" ipv4.method manual
    nmcli con mod "My Profile" ipv4.addresses 192.168.1.100/24
    nmcli con mod "My Profile" ipv4.gateway 192.168.1.1
    nmcli con mod "My Profile" ipv4.dns "8.8.8.8,1.1.1.1"
    nmcli con up "My Profile"               # Apply changes

Creating:
    nmcli con add type ethernet con-name "Static Lab" ifname enp3s0 \
        ipv4.method manual ipv4.addresses 10.0.0.5/24 ipv4.gateway 10.0.0.1

Deleting:
    nmcli con delete "My Profile"

---

nmtui — Text User Interface:

    nmtui

Opens a simple ncurses menu for editing connections, activating/deactivating, and setting the hostname. Useful on headless servers over SSH.

---

Dispatcher Scripts:

NM can run scripts in response to network events (interface up, interface down, DHCP lease change, VPN connected, etc.). Scripts are placed in /etc/NetworkManager/dispatcher.d/ and receive the interface name and action as arguments.

Example: /etc/NetworkManager/dispatcher.d/99-custom
    #!/bin/bash
    INTERFACE=$1
    ACTION=$2
    if [ "$ACTION" = "up" ] && [ "$INTERFACE" = "enp3s0" ]; then
        ip route add 10.10.0.0/16 via 10.0.0.1
    fi

Scripts must be executable (chmod +x) and owned by root. They run synchronously, so they block NM until they complete — keep them fast.

---

How NM Interacts with udev:

When a new network device appears in the kernel (hot-plug, driver load, VF creation), udev sends an event. NetworkManager listens for these events and:
  1. Checks if the device matches any unmanaged rules → if yes, ignores it.
  2. If managed, looks for a connection profile that matches the device (by MAC, interface name, or device type).
  3. If a matching profile is found and set to autoconnect=yes, NM activates it automatically.
  4. If no profile matches, NM may create a default "Wired connection N" profile (depending on NM configuration).

---

DNS Management:

NM writes DNS settings to /etc/resolv.conf. On modern systems using systemd-resolved, NM pushes DNS servers to resolved via D-Bus instead.

    nmcli con mod "My Profile" ipv4.dns "8.8.8.8 1.1.1.1"
    nmcli con mod "My Profile" ipv4.dns-search "example.com"
    nmcli con mod "My Profile" ipv4.ignore-auto-dns yes   # Ignore DHCP-provided DNS

---

Restarting NM:

    sudo systemctl restart NetworkManager      # Full restart (brief network disruption)
    sudo nmcli networking off && sudo nmcli networking on  # Softer toggle
