SELinux (Security-Enhanced Linux)

SELinux is a Mandatory Access Control (MAC) system implemented as a Linux Security Module (LSM). It was developed by the NSA and is the default MAC system on RHEL, CentOS, Fedora, and Amazon Linux.

---

DAC vs MAC:

Traditional Linux uses DAC (Discretionary Access Control): file permissions (rwx) and ownership (user/group). The owner of a file decides who can access it. The problem: if a process running as root is compromised, it has unrestricted access to the entire system.

MAC (Mandatory Access Control): The kernel enforces a security policy that even root cannot override. Every process (subject) and every resource (object) has a security label. The policy defines which labels can interact. Even root cannot read a file if the policy says "no."

SELinux provides MAC. A compromised Apache process running as root can only access files labeled for Apache use — not /etc/shadow, not user home directories, not other services' data.

---

SELinux Modes:

Enforcing: The policy is enforced. Denied actions are blocked AND logged.
    setenforce 1

Permissive: The policy is NOT enforced. Denied actions are logged but ALLOWED. Used for troubleshooting — see what would break without actually breaking it.
    setenforce 0

Disabled: SELinux is completely off. No labels, no logging, no policy. Requires a reboot to re-enable (and a full filesystem relabel).

    getenforce                        # Show current mode
    sestatus                          # Detailed SELinux status

Persistent configuration: /etc/selinux/config (SELINUX=enforcing|permissive|disabled).

---

Security Contexts (Labels):

Every process, file, port, and resource has a security context (label) in the format:

    user:role:type:level

    system_u:system_r:httpd_t:s0

The type field is the most important for policy decisions (this is called "Type Enforcement"):
  - httpd_t: Apache web server process.
  - httpd_sys_content_t: Files Apache is allowed to read.
  - user_home_t: User home directory files.
  - etc_t: Files in /etc/.
  - unconfined_t: No SELinux restrictions (default for user-launched processes on RHEL).

A process with type httpd_t can read files with type httpd_sys_content_t but CANNOT read files with type user_home_t — even if the Unix permissions say rwxrwxrwx.

Viewing labels:

    ls -Z /var/www/html/             # File labels (-Z flag)
    ps -eZ | grep httpd              # Process labels
    id -Z                             # Current user's label

---

Common SELinux Problems and Fixes:

Problem 1: File has the wrong label.

You moved a file (mv) instead of copying it (cp). mv preserves the original label; cp gets the label of the destination directory.

    mv /home/user/index.html /var/www/html/
    # index.html still has user_home_t instead of httpd_sys_content_t

Fix: Restore the correct label:
    restorecon -Rv /var/www/html/    # Recursively reset labels to policy defaults

Or set a label manually:
    chcon -t httpd_sys_content_t /var/www/html/index.html

restorecon uses the "file context database" (/etc/selinux/targeted/contexts/files/file_contexts) to determine the correct label for each path. chcon is temporary — restorecon will override it.

To permanently add a custom path→label mapping:
    semanage fcontext -a -t httpd_sys_content_t "/srv/mysite(/.*)?"
    restorecon -Rv /srv/mysite

Problem 2: A service needs to bind to a non-standard port.

SELinux controls which ports each type can bind to. Apache (httpd_t) is allowed to bind to ports 80, 443, 8080, etc. If you configure Apache to listen on port 9090, SELinux blocks it.

Fix: Add the port to the allowed list:
    semanage port -a -t http_port_t -p tcp 9090
    semanage port -l | grep http     # List allowed ports for HTTP

Problem 3: A boolean needs to be toggled.

Booleans are on/off switches for common policy adjustments:

    getsebool -a                      # List all booleans
    getsebool httpd_can_network_connect  # Check a specific boolean

    setsebool -P httpd_can_network_connect on   # Allow Apache to make outgoing network connections
    # -P makes it persistent across reboots

Common booleans:
    httpd_can_network_connect      — Apache can connect to network services (databases, APIs).
    httpd_enable_homedirs          — Apache can serve files from user home directories.
    samba_enable_home_dirs         — Samba can share home directories.
    nfs_export_all_rw              — NFS can export any directory read/write.
    virt_use_nfs                   — VMs can use NFS storage.

---

Troubleshooting with audit2why and audit2allow:

When SELinux denies an action, it logs an AVC (Access Vector Cache) denial in /var/log/audit/audit.log:

    type=AVC msg=audit(1234567890.123:456): avc:  denied  { read } for
    pid=1234 comm="httpd" name="config.php" dev="sda1" ino=12345
    scontext=system_u:system_r:httpd_t:s0
    tcontext=unconfined_u:object_r:user_home_t:s0 tclass=file

Reading this: Process httpd (type httpd_t) tried to read file config.php (type user_home_t) and was denied.

    audit2why < /var/log/audit/audit.log
    # Explains WHY the denial happened and suggests a fix (boolean, label, or policy module).

    audit2allow -a                    # Generate a policy module to allow all recent denials
    audit2allow -a -M mypolicy       # Compile it into a loadable module
    semodule -i mypolicy.pp          # Load the policy module

Warning: audit2allow can be overly permissive. Always review the generated policy. Prefer restorecon, semanage, and setsebool over custom policy modules when possible.

---

Quick Troubleshooting Checklist:

1. Is SELinux the problem?  → setenforce 0, retry. If it works, SELinux is blocking something.
2. Check the audit log     → ausearch -m avc -ts recent
3. Read the denial          → audit2why < /var/log/audit/audit.log
4. Fix with (in order of preference):
   a. restorecon (fix file labels)
   b. setsebool (toggle a boolean)
   c. semanage (add a port or file context rule)
   d. audit2allow (create a custom policy module — last resort)
5. Re-enable enforcing     → setenforce 1

---

sealert (setroubleshoot):

On RHEL/Fedora with setroubleshoot-server installed, sealert provides human-readable explanations:

    sealert -a /var/log/audit/audit.log
    # Produces detailed explanations with confidence ratings and suggested commands

Ref: man 8 semanage, man 8 restorecon, man 8 audit2why
Ref: Red Hat SELinux documentation
