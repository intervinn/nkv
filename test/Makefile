CC ?= cc
CFLAGS ?= -Wall -Wextra -I../libkv
LDFLAGS = -L../libkv -lkv -Wl,-rpath,../libkv
TARGET = kvtest
SRCS = main.c

all: $(TARGET)

$(TARGET): $(SRCS)
	$(CC) $(CFLAGS) $^ $(LDFLAGS) -o $@

clean:
	rm -f $(TARGET) *.o

.PHONY: all clean