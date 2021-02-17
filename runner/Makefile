CC = gcc
CFLAGS = -std=c11 -Wall -Werror -D_XOPEN_SOURCE=700

.PHONY: all clean install

all: rn run.sh

clean:
	rm -f rn rn.o

install: all
	install -m755 rn /usr/local/bin/rn
	install -m755 run.sh /usr/local/bin/run.sh

rn: rn.o
	$(CC) $(CFLAGS) -o $@ rn.o

%.o: %.c
	$(CC) $(CFLAGS) -o $@ -c $<
