#!/bin/python

from subprocess import call
import time
import os

# How many times to execute the script
exec_times = 10

# Will hold the result
times = []

# So that we can send everything to /dev/null
DEVNULL = open(os.devnull, 'wb')

for x in range(0, exec_times):
	start = time.time()
	call(["../bin/kram", "benchmarks/Fib.kr"], stdout = DEVNULL)
	times.append(time.time() - start)

print "Average: ", sum(times) / exec_times
print "    Min: ", min(times)
print "    Max: ", max(times)