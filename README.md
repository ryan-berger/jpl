CS 4470 Compilers Template
==========================

Please use this repository to hold all of your work for CS 4470
compilers. You will submit your assigments, and they will be graded,
using the contents of this repository only, at the commit
corresponding to the latest possible submission time (including
all granted extensions).

Picking a Language
------------------

This repository supports work in C++ (version 20), Java (15), or
Python (3.9). It correspondingly includes starter files
`compiler.cpp`, `compiler.java`, and `compiler.py`. Each starter file
compiles and runs, doing nothing and returning successfully. For each
language there is also a Makefile: `Makefile_cpp`, `Makefile_java`,
and `Makefile_py`.

To pick your language, delete the starter files for the other two
languages. If you want to use a language other than C++, Java, or
Python, contact the instructors. Keep in mind that using another
language will be more work, and you will not be able to receive the
same level of instructor support. Compilers is a complicated subject.
It is abjectly irresponsible to try to learn a new language at the
same time as you learn compilers.

Compiling your Compiler
-----------------------

A compiler is just a normal computer program. Before running it you
need to compile it. So, once you've picked a language, compile your
compiler by running:

    make compile

This should complete without errors. If you get an error, you likely
need to install one of `clang++`, `javac`/`java`, or `python3`
(depending on the language you chose), or to add those tools to your
system PATH. If you're having trouble with this step, please talk
to your TA or one of your instructors.

Note that if you're using Python, `make compile` will do a bit of
syntax checking but that's about it.

Running your Compiler
---------------------

To test your compiler you need a test program to compile. The file
`test.jpl` is intended to be a quick scratch-pad for such tests. Run
your compiler on it with:

    make run

Of course, longer-term you'll want to save test files and run the
regularly to avoid regressions. You can change the name of the test
file like so:

    make run TEST=something.jpl

We will use this same functionality to grade your assignments.
