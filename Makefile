CC=gcc
CFLAGS=-O2 -Wall

# build
all: upload 

upload: upload.o cgic.o
	gcc -o upload upload.o cgic.o 
	strip upload

upload.o: upload.c
cgic.o: cgic.c cgic.h

clean:
	rm -f *.o upload
