Feature,WARN_ON(condition),WARN_ON_ONCE(condition)
Trigger Frequency,Fires every time the condition is true.,Fires only the first time the condition is true.
Log Spam,High. Can flood dmesg if in a loop or hot path.,Low. Self-throttles to prevent console flooding.
Performance,Constant overhead for logging and stack traces.,Minimal overhead after the first trigger.
Implementation,Checks the condition and prints if true.,Checks condition AND a hidden static bool flag.
Standard Usage,Rarely used in new code.,Highly preferred by kernel maintainers.

Both macros produce a stack trace and a register dump to the kernel log (dmesg). This allows developers to see exactly what code path led to the unexpected state.

WARN_ON_ONCE is implemented using a static bool variable that is unique to that specific line of code.
C

// Simplified logic of WARN_ON_ONCE
static bool __warned;
if (unlikely(condition)) {
    if (!__warned) {
        __warned = true;
        warn_slowpath_fmt(...); // Triggers the stack trace
    }
}