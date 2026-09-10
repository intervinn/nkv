CC=cc
CFLAGS=-Wall -Wextra -I.
TARGET=libkv
SRCS=kv.c

$(TARGET): $(SRCS)
	$(CC) $(CFLAGS) $^ -dynamiclib -o $@.dylib
