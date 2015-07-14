SOURCEDIR = src
SOURCES := $(shell find $(SOURCEDIR) -name '*.cpp')

all:
	g++ -std=c++14 -o kram $(SOURCES)