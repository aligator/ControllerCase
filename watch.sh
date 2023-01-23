#!/usr/bin/env bash

# Kill also openscad if this script gets terminated.
trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT

go install github.com/cosmtrek/air@latest
openscad out.scad &
air -- $@
