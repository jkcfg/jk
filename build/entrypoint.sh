#!/bin/bash

if [ -z "$SRC_PATH" ]; then
    echo "error: \$SRC_PATH not set."
    exit 1
fi

# If we run make directly, any files created on the bind mount
# will have awkward ownership.  So we switch to a user with the
# same user and group IDs as source directory.  We have to set a
# few things up so that sudo works without complaining later on.
uid=$(stat --format="%u" "$SRC_PATH")
gid=$(stat --format="%g" "$SRC_PATH")
echo "user:x:$uid:$gid::$SRC_PATH:/bin/bash" >>/etc/passwd
echo "user:*:::::::" >>/etc/shadow
echo "user	ALL=(ALL)	NOPASSWD: ALL" >>/etc/sudoers

script="$@"
exec gosu user /bin/bash -c "$script"
