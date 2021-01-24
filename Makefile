TEST=test.jpl

all: run

compile: cmd/compiler.go
	go build -o jpl cmd/compiler.go

run: compile
	./jpl -l $(TEST)

clean:
	rm -rf jpl