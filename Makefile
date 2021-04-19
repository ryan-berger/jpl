TEST=flatten.jpl

.SILENT: run run-a4t run-a4f run-a5p run-a5t
all: run

compile: cmd/compiler.go
	go build -o jpl cmd/compiler.go

run:
	./jpl -s $(TEST)

run-a4t:
	./jpl -t $(TEST)

run-a4f:
	./jpl -f $(TEST)

run-a5p:
	./jpl -p $(TEST)

run-a5t:
	./jpl -t $(TEST)

generate-asm:
	./jpl -o code.s $(TEST)

assemble: generate-asm
	nasm -felf64 code.s

link: assemble
	clang code.o ./assignment4/runtime.a -lpng -L/usr/local/lib -lm

clean:
	rm -rf jpl
