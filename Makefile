gus:
	go build -o Gus ./src

all: gus

build_test:
	go build -o RunTests ./test

test: clean gus build_test
	./RunTests

clean:
	rm Gus RunTests