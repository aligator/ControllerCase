#!/bin/bash

# Kill also openscad if this script gets terminated.
trap "trap - SIGTERM && kill -- -$$" SIGINT SIGTERM EXIT

go get -u github.com/cosmtrek/air
openscad out.scad &
air