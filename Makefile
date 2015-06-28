build:
	go build -o bin/Gus

build_test:
	go build -o bin/Test ./test

test: clean build build_test
	./bin/Test $(CURDIR)/test/tests

clean:
	-rm bin/*

fmt:
	go fmt github.com/kram/kram/...

all: fmt build