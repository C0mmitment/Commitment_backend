#!/bin/sh
if [ ! -f package.json ]; then
  npm init -y >/dev/null 2>&1
fi
tail -f /dev/null