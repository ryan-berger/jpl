#!/usr/bin/env bash

set -e
set -u

rm -f runtime.o pngstuff.o runtime.a

clang -g -O0 -c runtime.c
clang -g -O0 -c pngstuff.c

ar rcs runtime.a runtime.o pngstuff.o
