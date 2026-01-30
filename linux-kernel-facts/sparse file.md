truncate -s 256M /tmp/testfile

This creates a sparse file named /tmp/testfile that is 256MB in size. It doesn't actually take up 256MB on your disk yet; it just tells the system "this file is allowed to grow up to this size." This acts as the virtual hard drive.

---
The Unix command

dd of=sparse-file bs=5M seek=1 count=0

will create a file of five mebibytes in size, but with no data stored on the media (only metadata). (GNU dd has this behavior because it calls ftruncate to set the file size; other implementations may merely create an empty file.)

Similarly the truncate command may be used, if available:

truncate -s 5M <filename>

---
The -s option of the ls command shows the occupied space in blocks.

ls -ls sparse-file

Alternatively, the du command prints the occupied space, while ls prints the apparent size. In some non-standard versions of du, the option --block-size=1 prints the occupied space in bytes instead of blocks, so that it can be compared to the ls output:

du --block-size=1 sparse-file
ls -l sparse-file
