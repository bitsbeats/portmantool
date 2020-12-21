CC = gcc
CFLAGS = -std=c11 -Wall -Werror -D_XOPEN_SOURCE=700

objects = runner.o

.PHONY: all clean

all: runner

clean:
	rm -f runner runner.o

runner: $(objects)
	$(CC) $(CFLAGS) -o $@ $(objects)

%.o: %.c
	$(CC) $(CFLAGS) -o $@ -c $<
