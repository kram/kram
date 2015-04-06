gus:
	go build -o bin/Gus

build_test:
	go build -o bin/Test ./test

test: clean gus build_test
	./bin/Test $(CURDIR)/test/tests

clean:
	-rm bin/*

fmt:
	go fmt github.com/zegl/Gus/...

all: fmt gus