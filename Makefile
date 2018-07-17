CC=gcc
CFLAGS=-O2 -Wall -DCGIMAXTEMPFILESIZE=5242880

# build
all: upload 

upload: upload.o cgic.o
	gcc -o upload upload.o cgic.o 
	strip upload

upload.o: upload.c
cgic.o: cgic.c cgic.h

clean:
	rm -f *.o upload
