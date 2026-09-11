#include "kv.h"

#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>
#include <netinet/in.h>
#include <sys/socket.h>

#define KV_BUF_SIZE 1024

struct _kv_client_t {
    int sock_fd;
    struct sockaddr_in serv_addr;

    const char *addr;
    const char *err;
    int port;
    char buf[KV_BUF_SIZE];
};

static int kv_marshal(char* b, int count, ...);

kv_client_t *kv_new(const char *addr, int port) {
    kv_client_t *c = malloc(sizeof(kv_client_t));
    c->addr = addr;
    c->port = port;
    c->err = "";
    return c;
}

int kv_connect(kv_client_t *c) {
    if ((c->sock_fd = socket(AF_INET, SOCK_STREAM, 0)) <= 0) {
        c->err = "socket creation failed";
        goto error;
    }

    c->serv_addr.sin_family = AF_INET;
    c->serv_addr.sin_port = htons(c->port);

    if (inet_pton(AF_INET, c->addr, &c->serv_addr.sin_addr) <= 0) {
        c->err = "invalid address";
        goto error;
    }

    if (connect(c->sock_fd, (struct sockaddr*)&c->serv_addr, sizeof(c->serv_addr))) {
        c->err = "connection failed";
        goto error;
    }

    return 0;

    error:
        close(c->sock_fd);
        return -1;
}

int kv_close(kv_client_t *c) {
    close(c->sock_fd);
    free(c);
    return 0;
}

const char *kv_err(kv_client_t *c) {
    return c->err;
}

int kv_exists(kv_client_t *c, char *key, bool *value) {
    int n = kv_marshal(c->buf+4, 2, "EXISTS", key);
    uint32_t un = (uint32_t)n;
    memcpy(c->buf, &un, sizeof(un));

    if (send(c->sock_fd, c->buf, n+4, 0) < 0) {
        c->err = "send failed";
        return -1;
    }

    char hd[4] = {0};
    if (recv(c->sock_fd, hd, 4, 0) == -1) {
        c->err = "recv failed";
        return -1;
    }

    uint32_t len;
    memcpy(&len, hd, sizeof(len));

    if (recv(c->sock_fd, c->buf, len, 0) == -1) {
        c->err = "recv failed";
        return -1;
    }

    *value = strcmp(c->buf, "YES") == 0;

    return 0;
}

int kv_put(kv_client_t *c, char *key, char *value) {
    int n = kv_marshal(c->buf+4, 3, "SET", key, value);
    uint32_t un = (uint32_t)n;
    memcpy(c->buf, &un, sizeof(un));

    if (send(c->sock_fd, c->buf, n+4, 0) < 0) {
        c->err = "send failed";
        return -1;
    }

    char hd[4] = {0};
    if (recv(c->sock_fd, hd, 4, 0) == -1) {
        c->err = "recv failed";
        return -1;
    }

    uint32_t len;
    memcpy(&len, hd, sizeof(len));

    if (recv(c->sock_fd, c->buf, len, 0) == -1) {
        c->err = "recv failed";
        return -1;
    }

    if (strncmp(c->buf, "OK", 3) == 0) {
        return 0;
    }

    c->err = "put failed";
    return -1;
}

int kv_get(kv_client_t *c, char *key, char **value) {
    int n = kv_marshal(c->buf+4, 2, "GET", key);
    uint32_t un = (uint32_t)n;
    memcpy(c->buf, &un, sizeof(un));

    if (send(c->sock_fd, c->buf, n+4, 0) < 0) {
        c->err = "send failed";
        return -1;
    }
    
    char hd[4] = {0};
    if (recv(c->sock_fd, hd, 4, 0) == -1) {
        c->err = "recv failed";
        return -1;
    }

    uint32_t len;
    memcpy(&len, hd, sizeof(len));

    if (recv(c->sock_fd, c->buf, len, 0) == -1) {
        c->err = "recv failed";
        return -1;
    }

    *value = strndup(c->buf, len);
    return 0;
}

int kv_del(kv_client_t *c, char *key) {
    int n = kv_marshal(c->buf+4, 2, "DEL", key);
    uint32_t un = (uint32_t)n;
    memcpy(c->buf, &un, sizeof(un));

    if (send(c->sock_fd, c->buf, n+4, 0) < 0) {
        c->err = "send failed";
        return -1;
    }
    
    char hd[4] = {0};
    if (recv(c->sock_fd, hd, 4, 0) == -1) {
        c->err = "recv failed";
        return -1;
    }

    uint32_t len;
    memcpy(&len, hd, sizeof(len));

    if (recv(c->sock_fd, c->buf, len, 0) == -1) {
        c->err = "recv failed";
        return -1;
    }

    if (strncmp(c->buf, "OK", 3) == 0) {
        return 0;
    }

    c->err = "del failed";
    return -1;
}

static int kv_marshal(char *b, int count, ...) {
    va_list args;
    va_start(args, count);

    int n = 0;
    for (int i = 0; i < count; i++) {
        char *str = va_arg(args, char*);
        int l = strlen(str) + 1; // include null terminator
        memcpy(b+n, str, l);
        n += l;
    }

    va_end(args);
    return n;
}