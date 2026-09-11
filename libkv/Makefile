CC ?= cc
CFLAGS ?= -Wall -Wextra -fPIC -I.
SRCS = kv.c

UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S), Darwin)
    EXT = dylib
    LDFLAGS = -dynamiclib
else
    EXT = so
    LDFLAGS = -shared
endif

TARGET = libkv.$(EXT)

all: $(TARGET)

$(TARGET): $(SRCS)
	$(CC) $(CFLAGS) $(LDFLAGS) $^ -o $@

clean:
	rm -f $(TARGET) *.o

.PHONY: all clean