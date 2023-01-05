# JPL - A UofU CS4470 Programming Language

![Mandlebrot](/mandlebrot.png)

## History

JPL also known as **J**(ohn Regehr) **P**((avel Pachenkha) **L**(ang), or **J**(ust a) **P**(rogramming) **L**(anguage),
or whatever else you make fit with the letters JPL is a toy programming language to teach the students of the
University of Utah about compilers.

It is an array-based programming language built to process images and videos. To find the entire spec, take a look at SPEC.md

It was initially written with a NASM backend that one would have to assemble with nasm,
and link to the runtime.  The runtime is minimal. It is written in C, and uses LibPNG to do the image processing. Due to constraints of the class
video processing was never finished. 

As a compiler it is pretty simple. It will lex, parse, then do source to source optimizations,
emitting NASM at the end of the process. We were tasked with choosing from a list of optimizations, 
so I implemented dead code elimination, constant propagation, constant folding, and a few peephole optimizations.

All tests were run using Perl, and used our required S-Expression outputs were delt with by Racket to guarantee
our lexing, parsing, type checking, and optimizations were all correct.

## Additions

Since the spec is pretty bare, and missing some "real programming language" features, I've added some new features.

First NASM sucks pretty bad, and there were a lot of steps to get a binary working correctly. Now I've added a LLVM backend
to make things more portable, and I've got the build step down to two commands. One to compile the object file,
and one to link LibPNG.

## Using it

To compile the compiler, run `make compile` and you should get a `jpl` binary outputted in your directory. You may need to
have Tiny Go's llvm-go installed so that LLVM can be invoked while compiling