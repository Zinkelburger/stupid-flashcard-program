Tape drives have two "modes" in Linux:

    /dev/st1 (Rewinding): Every time a program finishes using the tape and closes the file, the drive automatically rewinds the tape to the beginning.

    /dev/nst1 (Non-rewinding): The tape stays exactly where it stopped.