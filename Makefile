SOURCEDIR = src
SOURCES := $(shell find $(SOURCEDIR) -name '*.cpp')

build: clean
	g++ -Ofast --std=c++11 -o bin/kram -Wall -Wextra $(SOURCES)

clean:
	mkdir -p bin

build_test: clean
	go build -o bin/Test ./test

test: clean build build_test
	./bin/Test $(CURDIR)/bin $(CURDIR)/test/tests

test_only: build_test
	./bin/Test $(CURDIR)/bin $(CURDIR)/test/tests