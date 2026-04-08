Kernel Modules

A kernel module is a piece of code that can be loaded into or unloaded from the running kernel without rebooting. Most device drivers, filesystems, and protocol handlers are built as modules (.ko files). This keeps the base kernel image small and lets hardware support be loaded on demand.

---

Module Basics:

A module is a compiled ELF object file (.ko) that exports init and exit functions:

    module_init(my_init);   — Called when the module is loaded (insmod/modprobe).
    module_exit(my_exit);   — Called when the module is removed (rmmod/modprobe -r).

Every module must declare its license:
    MODULE_LICENSE("GPL");    — Required for access to GPL-only kernel symbols.
    MODULE_AUTHOR("Name");
    MODULE_DESCRIPTION("What this module does");

Modules live in /lib/modules/$(uname -r)/ organized by subsystem:
    /lib/modules/6.8.0/kernel/drivers/net/ethernet/intel/ice/ice.ko
    /lib/modules/6.8.0/kernel/fs/ext4/ext4.ko
    /lib/modules/6.8.0/kernel/net/ipv4/netfilter/

---

Loading Modules:

insmod <path/to/module.ko>:
  - Loads a single module. Does NOT resolve dependencies.
  - If the module depends on another module that isn't loaded, insmod fails.
  - Rarely used directly.

modprobe <module_name>:
  - The smart loader. Reads /lib/modules/$(uname -r)/modules.dep to find and load all dependencies first, then loads the requested module.
  - This is the standard way to load modules.

    modprobe ice                     # Load the Intel ice driver + all dependencies
    modprobe bonding mode=4          # Load with a module parameter

Automatic loading:
  When the kernel discovers hardware (PCI device, USB device, etc.), it emits a uevent with the device's identifiers (vendor/device/class). udev catches this event and calls modprobe with the matching module alias. This is why drivers load automatically when you plug in hardware.

    modprobe -c | grep alias         # Show all module aliases

---

Removing Modules:

rmmod <module_name>:
  - Unloads a module. Fails if the module is in use (reference count > 0).
  - Does NOT remove dependent modules.

modprobe -r <module_name>:
  - Removes the module AND any dependencies that are no longer needed.

    rmmod ice                        # Remove the ice driver
    modprobe -r ice                  # Remove ice and unused dependencies

A module's reference count tracks how many things are using it (open devices, active connections, dependent modules). You cannot remove a module with refcount > 0.

    lsmod | grep ice                 # Shows "Used by" column = refcount
    cat /sys/module/ice/refcnt       # Direct refcount check

---

Module Parameters:

Many modules accept parameters that change their behavior. Parameters can be set at load time or via sysfs:

At load time:
    modprobe bonding mode=4 miimon=100

Via sysfs (runtime, if the module supports it):
    echo 100 > /sys/module/bonding/parameters/miimon

Listing available parameters:
    modinfo bonding                  # Shows all parameters with descriptions
    ls /sys/module/bonding/parameters/

Making parameters persistent (load at every boot):
    echo "options bonding mode=4 miimon=100" > /etc/modprobe.d/bonding.conf

---

depmod — Dependency Resolution:

depmod scans all modules in /lib/modules/$(uname -r)/ and builds:
  - modules.dep: Which modules depend on which.
  - modules.alias: Maps hardware IDs (PCI vendor:device, USB IDs) to module names.
  - modules.symbols: Which symbols each module exports.

    depmod -a                        # Rebuild all dependency files
    depmod $(uname -r)               # For a specific kernel version

You must run depmod after installing new modules (e.g., after make modules_install for a custom kernel). Without it, modprobe won't know about new modules or their dependencies.

---

Blacklisting Modules:

To prevent a module from loading (e.g., a buggy driver, or to prefer a different driver):

    echo "blacklist nouveau" > /etc/modprobe.d/blacklist-nouveau.conf

This stops modprobe from auto-loading the module. To also prevent manual loading, add:
    echo "install nouveau /bin/false" >> /etc/modprobe.d/blacklist-nouveau.conf

---

Module Signing (Secure Boot):

When Secure Boot is enabled, the kernel only loads modules that are cryptographically signed with a trusted key. Unsigned or incorrectly signed modules are rejected with "Required key not available" (-ENOKEY).

The signing key is embedded in the kernel at build time (CONFIG_MODULE_SIG_KEY). RHEL/Fedora sign all in-tree modules with Red Hat's key.

For out-of-tree modules (DKMS, custom drivers), you must:
  1. Generate a signing keypair.
  2. Enroll the public key in the UEFI MOK (Machine Owner Key) database:
         mokutil --import my_signing_key.der
  3. Reboot and accept the key enrollment in the MOK manager.
  4. Sign the module:
         /usr/src/kernels/$(uname -r)/scripts/sign-file sha256 my_key.priv my_key.der module.ko

    cat /proc/sys/kernel/tainted     # Non-zero if unsigned/out-of-tree modules are loaded

---

Useful Commands:

    lsmod                            # List all loaded modules
    modinfo <module>                 # Show module info (author, description, params, alias)
    modprobe <module>                # Load module with dependencies
    modprobe -r <module>             # Remove module and unused dependencies
    rmmod <module>                   # Remove module (no dependency handling)
    depmod -a                        # Rebuild module dependency database
    dmesg | grep <module>            # See module init/error messages

Ref: Documentation/kbuild/modules.rst
Ref: Documentation/admin-guide/module-signing.rst
