SOURCEDIR = src
SOURCES := $(shell find $(SOURCEDIR) -name '*.cpp')

build:
	g++ -Ofast --std=c++11 -o bin/kram -Wall -Wextra $(SOURCES)

build_test:
	go build -o bin/Test ./test

clean:
	-rm bin/*
	-mkdir -p bin

test: clean build build_test
	./bin/Test $(CURDIR)/bin $(CURDIR)/test/tests

test_only: build_test
	./bin/Test $(CURDIR)/bin $(CURDIR)/test/tests

all: build