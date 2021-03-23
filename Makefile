TEST=test.jpl

.SILENT: run run-a4t run-a4f
all: run

compile: cmd/compiler.go
	go build -o jpl cmd/compiler.go

run:
	./jpl -s $(TEST)

run-a4t:
	./jpl -t $(TEST)

run-a4f:
	./jpl -f $(TEST)

clean:
	rm -rf jpl
