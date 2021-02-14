TEST=test.jpl

.SILENT: run
all: run

compile: cmd/compiler.go
	go build -o jpl cmd/compiler.go

run:
	./jpl -p $(TEST)

clean:
	rm -rf jpl