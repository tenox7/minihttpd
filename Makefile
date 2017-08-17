CC=gcc
CFLAGS=-O2 -Wall

# build
all: receive 

receive: receive.o cgic.o
	gcc -o receive receive.o cgic.o 
	strip receive

receive.o:receive.c
cgic.o: cgic.c cgic.h

clean:
	rm -f *.o receive
