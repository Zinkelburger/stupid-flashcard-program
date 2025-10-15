Front,Back
"What is UDP hole punching?","A technique for two clients behind separate NATs to establish a direct peer-to-peer connection."
"What problem does UDP hole punching solve?","The problem of NAT traversal, where NAT firewalls block unsolicited incoming connections, preventing direct P2P communication."
"What is the role of the third-party server?","It acts as a 'rendezvous' or 'introduction' server. It discovers the public IP and port of each client and shares that information between them."
"What information does the rendezvous server exchange?","Each client's public IP address and the source port their NAT used for the outbound connection."
"What action creates the 'hole' in the NAT firewall?","An outbound UDP packet from the client to the peer. The NAT temporarily adds an entry to its state table allowing return traffic."
"Why is UDP ideal for this technique?","Because it is connectionless. The NAT's simple state-tracking for outbound UDP packets is what makes hole punching possible."
"Is the rendezvous server needed after the connection is established?","No. Once the direct peer-to-peer link is made, the server is no longer involved in the communication."
"Name a common application of UDP hole punching.","Online gaming, VoIP (like Skype or Discord), or P2P applications."