Front,Back
"Class A IP Address Range","1.0.0.0 to 126.255.255.255"
"Class B IP Address Range","128.0.0.0 to 191.255.255.255"
"Class C IP Address Range","192.0.0.0 to 223.255.255.255"
"Purpose of Class D Addresses","Reserved for Multicasting."
"Purpose of Class E Addresses","Reserved for experimental and research purposes."
"Default Subnet Mask for Class A","255.0.0.0 or /8"
"Default Subnet Mask for Class B","255.255.0.0 or /16"
"Default Subnet Mask for Class C","255.255.255.0 or /24"
"What class is the address 150.10.25.1?","Class B"
"What class is the address 198.162.1.1?","Class C"
"The 127.0.0.0 network is used for...","Loopback testing. 127.0.0.1 is the common address for localhost."
"Total number of IP addresses in a /25 subnet","128 ($2^7$ addresses)"
"Number of USABLE host IPs in a /25 subnet","126 (128 minus the network and broadcast addresses)."
"An address starting with 169.254.x.x indicates...","APIPA (Automatic Private IP Addressing), meaning the device couldn't find a DHCP server."
"Private IP address range in Class A","10.0.0.0 to 10.255.255.255 (/8)"
"Private IP address range in Class B","172.16.0.0 to 172.31.255.255 (/12)"
"Private IP address range in Class C","192.168.0.0 to 192.168.255.255 (/16)"
"IP: 192.168.10.130 with a /25 mask. What is the network address?","192.168.10.128"
"IP: 192.168.10.130 with a /25 mask. What is the broadcast address?","192.168.10.255"
"Are 10.0.0.50/25 and 10.0.0.200/25 on the same subnet?","No. The first is in the 10.0.0.0 subnet, the second is in the 10.0.0.128 subnet."
"IP: 172.20.250.5 with a /22 mask. What is the network address?","172.20.248.0"
"IP: 172.20.250.5 with a /22 mask. What is the broadcast address?","172.20.251.255"
"What is the first usable host address on the 192.168.1.128/26 network?","192.168.1.129"
"What is the last usable host address on the 192.168.1.128/26 network?","192.168.1.190 (The broadcast address is .191)"
"How many subnets and hosts per subnet can you get from the network 10.0.0.0/8 if you use a /12 mask?","16 subnets ($2^4$) and 1,048,574 hosts per subnet ($2^{20} - 2$). You borrowed 4 bits for the subnet."

The Full Localhost Range (127.0.0.0/8)
While 127.0.0.1 is the universally recognized address for "this computer" or localhost, it's not the only one. The entire IP block from 127.0.0.0 to 127.255.255.255 is reserved for loopback purposes.

APIPA - The "I Can't Find a Server" Address (169.254.0.0/16)
APIPA stands for Automatic Private IP Addressing.
If your computer fails to get an IP address from a DHCP server, it will assign itself an address in the range of 169.254.0.1 to 169.254.255.254. This is also known as a link-local address.

255.255.255.255	Broadcast Address	A message sent to this address is delivered to every single device on the local network. It's like shouting to everyone in the room at once.

Front,Back
"Relationship between Class A and Private IPs","The private IP range 10.0.0.0/8 is a reserved *subset* of the much larger Class A address space."
"Relationship between Class B and Private IPs","The private IP range 172.16.0.0/12 is a reserved *subset* of the much larger Class B address space."
"Relationship between Class C and Private IPs","The private IP range 192.168.0.0/16 is a reserved *subset* of the much larger Class C address space."
"What is the purpose of Network Address Translation (NAT)?","NAT allows devices with private, non-unique IP addresses to communicate with the internet by translating them to a single, unique public IP address at the router."
"What does an IP address starting with 169.254 indicate?","This is an Automatic Private IP Addressing (APIPA) address, meaning the device could not reach a DHCP server to get an IP address automatically."
"What is the entire 127.0.0.0/8 block used for?","Loopback. Any packet sent to an address in this range is routed back to the sending device and never leaves the machine."
"While 127.0.0.1 is the standard, can you use other 127.x.x.x addresses?","Yes, the entire range is available for loopback, which is useful for developers testing multiple network services on one machine."
"What is the address 0.0.0.0 used for?","It's an 'unspecified' address. A device often uses it when first connecting to a network, before it has been assigned a proper IP."
"What is the address 255.255.255.255 used for?","This is the 'limited broadcast' address. A packet sent to this address is delivered to every host on the local network segment."
"Can private IP addresses be routed on the public internet?","No. Routers on the public internet are configured to drop any traffic coming from or going to a private IP address range."

