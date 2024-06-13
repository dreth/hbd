#!/bin/zsh
# Change to the backend directory
cd "backend" || exit
# run air
air -c .air.debug.toml -d
