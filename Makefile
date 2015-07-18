SOURCEDIR = src
SOURCES := $(shell find $(SOURCEDIR) -name '*.cpp')

all:
	g++ -std=c++11 -o kram $(SOURCES)