#!/bin/sh
if [ ! -f go.mod ]; then
  go mod init github.com/86shin/commit_goback
fi
tail -f /dev/null