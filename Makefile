SOURCEDIR = src
SOURCES := $(shell find $(SOURCEDIR) -name '*.cpp')

build:
	g++ --std=c++11 -o bin/kram $(SOURCES)

build_test:
	go build -o bin/Test ./test

clean:
	-rm bin/*
	-mkdir -p bin

test: clean build build_test
	./bin/Test $(CURDIR)/test/tests

all: build