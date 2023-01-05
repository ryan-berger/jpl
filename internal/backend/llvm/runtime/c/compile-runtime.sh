#!/usr/bin/env bash

set -e
set -u

rm -f runtime.bc pngstuff.bc runtime.bco

clang -O3 -c -emit-llvm runtime.c
clang -O3 -c -emit-llvm pngstuff.c

llvm-ar-14 rcs runtime.bco runtime.bc pngstuff.bc
