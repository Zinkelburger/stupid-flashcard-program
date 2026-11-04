# CompTIA Network+ (N10-009) study plan

The deck (`networkplus-flashcards.tsv`, 629 cards) covers every N10-009 objective, 1.1 through 5.5. Network+ is multiple choice plus performance-based questions (PBQs), so recall helps a lot more here than on the RHCSA. Cards still won't teach you to reason through scenario questions, do subnetting quickly, or handle PBQs. That's what practice exams are for.

## The exam

- **N10-009** is the current version. It launched **June 20, 2024**. CompTIA hasn't announced a retirement date; third parties estimate late 2027. N10-008 retired in Dec 2024, so ignore anything older than N10-009.
- **Up to 90 questions in 90 minutes.** A mix of multiple choice and PBQs (simulations, drag-and-drop, matching).
- **Passing score: 720** on a 100–900 scale.
- **Domains and their weights** (the deck is tagged by objective, e.g. `netplus::5.3`):

  | Domain | Weight |
  |---|---|
  | 1.0 Networking Concepts | 23% |
  | 2.0 Network Implementation | 20% |
  | 3.0 Network Operations | 19% |
  | 4.0 Network Security | 14% |
  | 5.0 Network Troubleshooting | **24%** (the biggest) |

- Official page and free objectives PDF: https://www.comptia.org/en-us/certifications/network/

## What the flashcards won't do for you

1. **Scenario wording.** Questions ask for the BEST, FIRST or MOST LIKELY answer, and often two options are technically true. You only get good at picking the intended one by doing lots of practice questions.
2. **Subnetting speed.** You need to work out a network/broadcast address or the right prefix for N hosts in under a minute, on paper. Practice daily (links below).
3. **PBQs.** Typical ones:
   - configure a wireless router or firewall rules
   - match ports to protocols or cables to connectors
   - place devices in a diagram
   - read `ipconfig` / `traceroute` / `show interface` output and spot the fault
4. **Troubleshooting reasoning.** "Users on VLAN 20 get 169.254 addresses" means a DHCP relay or scope problem. The cards in 5.x drill the common symptom → cause pairs, but practice exams throw new combinations at you.
5. **Reading real output.** Run the commands yourself: `ping`, `traceroute`, `dig`, `nslookup`, `ip`, `ss`/`netstat`, `tcpdump`, `nmap`, and Wireshark.

## Was the Google result useful?

| Result | Verdict |
|---|---|
| **PacktPublishing / CompTIA-Network-Certification-N10-009-The-Total-Course** | **Useful.** It's the companion files for Mike Meyers' (Total Seminars) N10-009 course. The episode quiz spreadsheet is good extra drill. |
| PacktPublishing … N10-007 | **Skip.** It's two versions old. |
| certbasepro/comptia-certification-questions | **Be skeptical.** The repo was created June 2026, has 10 stars, and is essentially an ad for certbase.pro claiming 13,645 questions across 52 exams. Mass-produced questions of unknown accuracy. |
| tertiarycourses/C938-CompTIA-NetworkPlus-Training | **Meh.** Courseware advertising a training company (2 stars, created July 2026). The objectives-mapped slides might be OK to skim; the "labs" are theirs. |
| PradyumnaKumarMirdha/CompTIA-Network-Plus | **Skip.** Last updated 2023, so it predates N10-009. |
| stephenschoettler/networkplus-notes | **Skim at most.** It says itself that the notes are *AI-generated* and unverified. |
| craigbu … N10-006 "dumps" / "DumpsBoss" GitHub app | **Avoid.** These are brain dumps: outdated, and against CompTIA's candidate agreement. You can lose your certification over them. |
| The AI overview's topic list | Fine as a list, but it hedges on "N10-008 or N10-009". Only N10-009 matters now. Its AD values (EIGRP 90, OSPF 110) are right. |

## Resources (ranked)

| Resource | Notes |
|---|---|
| **Professor Messer N10-009 course** — https://www.professormesser.com | **Free** full video course in objective order. He also sells course notes and practice exams. This is the default starting point most people use. |
| **Jason Dion** (Udemy / diontraining.com) N10-009 course + practice exams | His practice exams include PBQ-style questions and are widely used as the "am I ready?" check. |
| **Mike Meyers / Total Seminars** N10-009 (Udemy, Packt) | Good video explanations. The Packt repo above has the quiz spreadsheet. |
| **CompTIA official**: objectives PDF (free), CertMaster Learn / Labs / Practice | Expensive but official. Always use the objectives PDF as your checklist. |
| A Sybex/Wiley N10-009 study guide + practice test book | Good if you prefer reading. Extra question bank. |
| **Subnetting practice**: https://subnettingpractice.com , https://subnetipv4.com | 10 minutes a day until it's automatic. |
| **Labs**: Cisco Packet Tracer (free with a Cisco NetAcad account), Wireshark, your own Fedora box | Packet Tracer for VLANs, trunks, STP, routing and DHCP relay. Wireshark to actually see DORA, the TCP handshake, DNS, and ARP. |

## Training plan (about 6–8 weeks at ~1 h/day)

**Every day:** Anki for 15 min at **25 new cards/day**, which gets through all 629 in about 3.5 weeks. Add 10 minutes of subnetting practice.

**Weeks 1–4: learn each objective once**
- Watch Messer's videos in objective order (1.1 → 5.5).
- To study one objective's cards alongside its videos, use an Anki filtered deck with `deck:"CompTIA Network+ N10-009" tag:netplus::1.4`.
- Do the hands-on bits as you go:
  - capture DORA and a TCP handshake in Wireshark
  - `dig` each DNS record type
  - `traceroute` somewhere and read the hops
  - build VLANs + a trunk + router-on-a-stick + DHCP relay in Packet Tracer
  - `nmap` your own machine

**Weeks 4–5: weak spots and PBQs**
- Re-watch whatever your Anki stats and quiz results say you're weakest at. Usually that's ports, wireless, the troubleshooting order, and STP.
- Do PBQ-style practice: firewall rule tables, cable/connector matching, subnet allocation, reading command output.

**Weeks 6–8: practice exams**
- Take full timed practice exams (Dion, Messer, a book's question bank).
- Review **every** question you missed or guessed, and add a card for each real gap.
- **You're ready when** you score 85%+ consistently on exams you *haven't seen before*.

## Exam day

- Read the last sentence of each question first, and look for FIRST / BEST / MOST LIKELY.
- Troubleshooting questions usually follow the 7-step order. "What should the tech do FIRST?" is usually *gather information / question the user*.
- PBQs often come first and can eat time. Flag a long one, finish the multiple choice, then come back.
- Answer everything. A blank is guaranteed wrong.

## Using the deck

- **Anki:** File → Import → `networkplus-flashcards.tsv`. This creates its own deck, "CompTIA Network+ N10-009", separate from the RHCSA deck.
- **Updates:** re-importing an updated file updates cards in place and keeps your review history.
- **Tags:** `netplus::<objective>` (e.g. `netplus::2.2`). 17 cards are tagged `beyond-objectives` (POP3/IMAP, Wi-Fi 7, RIP, phishing variants…). They're commonly asked about but not explicitly in the objectives.
- **Your earlier notes:** the useful bits of `networking-gemini-notes/` (subnetting, spine-leaf, 802.11) are folded into the deck. The 802.11ac speed is corrected to ~6.9 Gbps theoretical. UDP hole punching isn't on the exam, so it was left out.
