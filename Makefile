FCGI_SOURCES := $(shell find src/fcgi -name '*.cpp')

KRAM_SOURCES := $(shell find . \( -type f -and -path '*/lexer/*'  \
	               -or -path '*/parser/*' \
	               -or -path '*/third_party/*' \
	               -or -path '*/vm/*' \
		\) -name '*.cpp')

build: clean
	g++ -Ofast --std=c++11 -o bin/kram -Wall -Wextra -Wno-unused-parameter src/main.cpp $(KRAM_SOURCES)

fcgi:
	g++ -Ofast --std=c++11 -o bin/fcgi -Wall -Wextra -Wno-unused-parameter -lfcgi++ -lfcgi $(FCGI_SOURCES) $(KRAM_SOURCES)

clean:
	mkdir -p bin

build_test: clean
	go build -o bin/Test ./test

test: clean build build_test
	./bin/Test $(CURDIR)/bin $(CURDIR)/test/tests

test_only: build_test
	./bin/Test $(CURDIR)/bin $(CURDIR)/test/tests