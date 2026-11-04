# RHCSA (EX200, RHEL 10) study plan

The flashcards (`rhcsa-flashcards.tsv`) give you **recall**. The exam tests **doing**: about 3 hours on live VMs, everything graded after a reboot. You need both. Cards alone get you maybe halfway there.

## The exam, in short

- Based on **RHEL 10**. Red Hat switched the objectives around Jul–Aug 2025.
  Official page: https://www.redhat.com/en/services/training/ex200-red-hat-certified-system-administrator-rhcsa-exam
- Hands-on only. No internet, no notes. The docs that ship with the system (man, info, `/usr/share/doc`) are what you have.
- **Configs must survive a reboot.** Grading happens on the rebooted machines.
- Length and passing score aren't on the official page. Third parties say 3 h and 210/300. Confirm when you book.
- **Changes from RHEL 9 → 10:**
  - Removed: containers/podman, SGID collaboration directories, MBR partitioning, and "diagnose SELinux policy violations".
  - Added: a "Manage software" section with **Flatpak** remotes and apps, plus **systemd timers** next to at and cron.
  - Stratis, VDO, module streams and LDAP are long gone, so ignore old material that covers them.

## What would stump you if you only did the flashcards

1. **Speed.** There are about 20 tasks in about 3 hours, some with several parts. Knowing `lvextend -r` is different from doing partition → PV → VG → LV → mkfs → fstab → mount -a in 5 minutes without errors.
2. **Tasks that cut across topics.** One task like "serve /srv/web on port 82" needs the package, the service enabled, a firewall rule, an SELinux port label, an fcontext rule plus restorecon, and a test with curl. Each card covers one piece. The exam wants the whole thing working.
3. **The reboot.** A typo in fstab and the node doesn't come back. If you can't fix that from emergency mode quickly, you lose every task on that node.
4. **Things that fail silently.** SELinux denials, firewall rules added without `--permanent`, sshd drop-ins overriding your setting, `nmcli con mod` without `con up`, `usermod -G` without `-a`. The deck covers each one. You only spot them by habitually **checking** your work.
5. **No internet.** When you blank on autofs map or timer syntax, you have to find it in man pages or copy an example from the system (see the `[Docs]` cards). Practise this, because it's slower than you think.
6. **Writing from scratch.** A whole script, a `.timer` + `.service` pair, a `.repo` file or an autofs map, typed in vim under time pressure.
7. **Order of work.** Get networking, hostname and repos working first. If you can't install packages, half the exam stalls. If a node's task is the root password, do that first.
8. **Reading the task literally.** "Between 830 and 865 MiB", exact UIDs, exact group names, exact paths.

None of that comes from flashcards. It comes from repeating the tasks in a lab.

## Resources (ranked)

| Resource | Notes |
|---|---|
| **Red Hat RH124 + RH134** (or **RH199** rapid track) via Red Hat Learning Subscription | The official courses, updated for RHEL 10.0. They include real lab environments. As a Red Hat employee you may have internal access, so check the internal learning portal (I haven't confirmed this). |
| **Sander van Vugt, *Red Hat RHCSA 10 Cert Guide*** (Pearson, 2026) + his RHEL 10 video course | The resource r/redhat recommends most. It has end-of-chapter labs and 4 practice exams. https://www.informit.com/store/red-hat-rhcsa-rhel-10-with-exam-labs-pearson-cert-prep-9780135493106 |
| **Asghar Ghori, *RHCSA RHEL 10 Training and Exam Preparation Guide*** (4th ed., Dec 2025) | Lots of labs, plus 4 sample exams. A good second source of practice exams. |
| **Red Hat free skills assessment** | Linked from the EX200 page. Use it to find weak areas. |
| **CertForged** — https://certforged.com/learn | New (domain registered April 2026), a one-person project, targets RHEL 10. The free tier is a *simulated* terminal. Pro is $19/mo, and real VMs cost $5/mo extra. Fine as extra drill; it doesn't replace your own VMs. |
| **chlebik/rhcsa-practice-questions** — https://github.com/chlebik/rhcsa-practice-questions | Written for RHCSA 7/8 and tested on CentOS 7. Last commit was Jan 2024. Still useful for core drills (root reset, LVM, autofs, users, SELinux contexts, cron). **Skip** the LDAP/authconfig, VDO and container questions, and the old yum/NTP-era ones. |
| Udemy "RHCSA 10" practice exams | Some say they're updated for v10; quality varies. Optional. |

**Avoid brain dumps.** That means anything selling "real EX200 questions", "actual exam PDFs", or posts that recall specific exam tasks. You agree to an NDA when you sit the exam, and breaking it can get your result voided and your certification revoked. As a Red Hat employee, that would also be a workplace problem. Practice exams written from the published objectives (all the resources above) are fine.

## Lab setup (do this first)

- **Free RHEL 10:** the Red Hat Developer Subscription for Individuals covers up to 16 systems, from developers.redhat.com. Download the **DVD ISO**, not the boot ISO.
- **Two or three VMs in virt-manager/KVM** (native on Fedora):
  - **node1:** the practice target. Add **2–3 extra virtual disks** (5–10 GB) for partitions, LVM and swap.
  - **node2:** a second target. Also runs as the **NFS server** for autofs practice. Export `/home/guests/*` and `/export/data`.
  - Optional **repo server:** serve the ISO's BaseOS/AppStream over httpd so you can practise http baseurls, and host a local Flatpak repo.
- **Take a snapshot of the clean state.** Revert it before every practice exam.
- **Practise without the network:** disconnect the VMs' internet and use the ISO as a local repo, so it feels like the exam.

## Training plan (about 8 weeks at ~1 h/day; compress if you have more time)

**Every day, the whole way through:** Anki for 10–15 min. Use **20 new cards/day**, which gets through all 505 in about 3.5 weeks, then keep up with reviews.

**Weeks 1–4: learn and do every objective once**
- Work through RH124 → RH134 (or van Vugt chapter by chapter).
- After each chapter, do the labs, then **do them again the next day without looking**.
- Rule: every time a card teaches a command, run it in the lab at least once.

**Weeks 5–6: drills, speed, and breaking things**
- Timed drills by domain, for example:
  - storage: partition + LV + fstab + reboot in under 10 min
  - web server on a custom port with SELinux + firewall in under 10 min
  - user/group/sudo batch in under 5 min
  - a timer unit from scratch in under 5 min
- Break things and fix them:
  - typo in fstab, then recover
  - forget the root password, then reset it
  - mislabel web content, then fix it
  - block a port with SELinux, then fix it
- Practise looking things up: answer 10 "how do I…" questions using only `man -k`, `man 5`, `/usr/share/doc` and `/etc/*` examples.

**Weeks 7–8: full mock exams**
- Run them under exam conditions: 3 hours, no internet, no notes, two VMs, reboot both at the end, then grade yourself honestly.
- Sources: van Vugt's 4 exams, Ghori's 4, the RHLS course review labs, then CertForged or chlebik for variety.
- After each mock, add cards for whatever you got wrong.
- **You're ready when** you pass 2–3 *different* mocks comfortably and finish with about 30 minutes left for checking.

## On exam day

- Read through all the tasks first. Start with networking, hostname and repos, then any root password reset.
- Check every task as you finish it (`systemctl is-enabled`, `findmnt`, `firewall-cmd --list-all`, `ls -Z`, `id`, `swapon --show` …).
- Before any reboot, run `findmnt --verify` and `mount -a`.
- Leave 20–30 min at the end: reboot every node, then check each task again.
- Stuck for more than 10 minutes? Move on and come back later.

## Using the deck

- **Anki:** File → Import → `rhcsa-flashcards.tsv`. The file's header sets the deck ("RHCSA"), the tags and the stable IDs.
  - If the cards are improved later, re-import the file. Changed cards update in place and keep their review history.
  - Tags: `rhcsa::<domain>` for each objective area. Cards marked `beyond-objectives` (ACLs, SGID, SELinux troubleshooting order) aren't in the v10 objective list but help with real tasks.
- The `[Docs]` cards cover where to look things up during the exam. The `[Exam]` cards cover strategy.
