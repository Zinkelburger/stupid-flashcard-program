FETCH_HEAD is a short-lived reference (a pointer) to the tip of what you just downloaded from a remote repository.

Whenever you run git fetch, Git does two things:

    Downloads the data: It pulls the commits, files, and objects from the remote into your .git folder.

    Updates FETCH_HEAD: It records the SHA-1 hash of the fetched commit in a special file located at .git/FETCH_HEAD.

Think of it as a temporary bookmark. If you fetch a specific commit from a remote, Git doesn't automatically create a local branch for it (like origin/fix-nbd). Instead, it just puts the commit in FETCH_HEAD so you can refer to it immediately.