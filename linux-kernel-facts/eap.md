EAP (Extensible Authentication Protocol) and 802.1X

EAP is an authentication framework (not a single protocol) defined in RFC 3748. It provides a flexible structure for carrying authentication exchanges, and is used heavily in both wired and wireless network access control.

---

802.1X: Port-Based Network Access Control

802.1X (IEEE standard) uses EAP to control who can send traffic through a network switch port or wireless access point. Until a device authenticates, the port only allows EAP frames through — all other traffic is blocked.

Three roles:

Supplicant: The client device requesting access (your laptop, a server's NIC).
Authenticator: The network switch or wireless AP. It acts as a gatekeeper — it doesn't make authentication decisions itself, it relays EAP messages.
Authentication Server: Almost always a RADIUS server (FreeRADIUS, Microsoft NPS, Cisco ISE). This is where the actual username/password or certificate check happens.

---

The Authentication Flow:

1. The supplicant connects to the switch port. The port starts in an "unauthorized" state — only EAPOL (EAP over LAN) frames are allowed through.

2. The switch sends an EAP-Request/Identity to the supplicant.

3. The supplicant responds with its identity (username, machine name, certificate CN).

4. The switch encapsulates this EAP message inside a RADIUS Access-Request and forwards it to the RADIUS server.

5. The RADIUS server and supplicant exchange EAP challenge/response messages (relayed through the switch). The specific method used (EAP-TLS, PEAP, etc.) determines what this looks like.

6. If authentication succeeds, the RADIUS server sends a RADIUS Access-Accept to the switch. The switch moves the port to "authorized" state and normal traffic flows.

7. If authentication fails, the port stays unauthorized. The device gets no network access.

---

Common EAP Methods:

EAP-TLS (Transport Layer Security):
  - Both the client and server present X.509 certificates.
  - Considered the most secure method — mutual authentication, no passwords.
  - Requires a PKI (Public Key Infrastructure) to issue and manage certificates.
  - Common in enterprise environments with managed devices.

EAP-PEAP (Protected EAP):
  - The server presents a certificate (like HTTPS), creating a TLS tunnel.
  - Inside that tunnel, the client authenticates with a username/password (MSCHAPv2).
  - Simpler to deploy than EAP-TLS because clients don't need individual certificates.
  - The most common method in enterprise Wi-Fi.

EAP-TTLS (Tunneled TLS):
  - Similar to PEAP: server cert creates a tunnel, inner authentication uses password.
  - More flexible than PEAP in the inner methods it supports (PAP, CHAP, MSCHAPv2).
  - Common in eduroam (university networks).

EAP-MD5:
  - Simple challenge-response using MD5 hash of a password.
  - No server authentication, no encryption of credentials.
  - Insecure and deprecated. Should not be used.

---

EAPOL (EAP over LAN):

EAP frames on a wired Ethernet network are carried using EAPOL — a special Ethernet frame type with EtherType 0x888E. These frames use the reserved multicast MAC address 01:80:C2:00:00:03.

By default, 802.1D-compliant bridges (Linux bridge, br0) are supposed to consume (drop) frames destined to the 01:80:C2:xx range and not forward them. This is correct for standard bridging but breaks 802.1X in environments where EAP frames need to pass through a software bridge.

The fix: writing to /sys/class/net/br0/bridge/group_fwd_mask with bit 3 set (value 8) tells the Linux bridge to forward EAPOL frames instead of dropping them.

---

Linux tools:

wpa_supplicant: The standard 802.1X supplicant for both wired and wireless on Linux.
  - For wired 802.1X: wpa_supplicant -i eth0 -D wired -c /etc/wpa_supplicant/wpa_supplicant.conf
  - Config file specifies the EAP method, identity, password/certificate paths.

hostapd: Can act as an 802.1X authenticator on Linux (for building a RADIUS-aware access point or switch).

FreeRADIUS: The most widely deployed open-source RADIUS server.

Ref: IEEE 802.1X-2020, RFC 3748 (EAP), RFC 5216 (EAP-TLS)
