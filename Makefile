gus:
	go build -o bin/Gus

all: gus

build_test:
	go build -o bin/Test ./test

test: clean gus build_test
	./bin/Test $(CURDIR)/test/tests

clean:
	-rm bin/*