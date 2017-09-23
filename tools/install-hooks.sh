#!/bin/sh
hooks_dir="$(pwd)/tools/git-hooks/"
for i in $(ls -1 "$hooks_dir")
do
    ln -sfv "$hooks_dir$i" "$(pwd)/.git/hooks"
done
